package auth

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// authKey is a custom context key for storing the auth token.
type authKey struct{}

// withAuthKey adds an auth key to the context.
func withAuthKey(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

// resolveAPIKey reads the API key from environment variables once.
// Returns the resolved key (may be empty if no key is configured).
func resolveAPIKey(logger *zap.Logger) string {
	key := os.Getenv("SLACK_MCP_API_KEY")
	if key == "" {
		key = os.Getenv("SLACK_MCP_SSE_API_KEY")
		if key != "" {
			logger.Warn("SLACK_MCP_SSE_API_KEY is deprecated, please use SLACK_MCP_API_KEY")
		}
	}

	if key == "" {
		logger.Warn("No API key configured — all HTTP/SSE requests will be allowed without authentication. Set SLACK_MCP_API_KEY to enable auth.",
			zap.String("context", "http"),
		)
	}

	return key
}

// validateToken checks if the request context contains a valid auth token.
// The expected API key is passed in to avoid reading env vars on every request.
func validateToken(ctx context.Context, expectedKey string, logger *zap.Logger) (bool, error) {
	if expectedKey == "" {
		return true, nil
	}

	keyB, ok := ctx.Value(authKey{}).(string)
	if !ok {
		logger.Warn("Missing auth token in context",
			zap.String("context", "http"),
		)
		return false, fmt.Errorf("missing auth")
	}

	logger.Debug("Validating auth token",
		zap.String("context", "http"),
		zap.Bool("has_bearer_prefix", strings.HasPrefix(keyB, "Bearer ")),
	)

	if strings.HasPrefix(keyB, "Bearer ") {
		keyB = strings.TrimPrefix(keyB, "Bearer ")
	}

	if subtle.ConstantTimeCompare([]byte(expectedKey), []byte(keyB)) != 1 {
		logger.Warn("Invalid auth token provided",
			zap.String("context", "http"),
		)
		return false, fmt.Errorf("invalid auth token")
	}

	logger.Debug("Auth token validated successfully",
		zap.String("context", "http"),
	)
	return true, nil
}

// AuthFromRequest extracts the auth token from the request headers.
func AuthFromRequest(logger *zap.Logger) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		authHeader := r.Header.Get("Authorization")
		return withAuthKey(ctx, authHeader)
	}
}

// BuildMiddleware creates a middleware function that ensures authentication based on the provided transport type.
// The API key is resolved once at middleware creation time, not on every request.
func BuildMiddleware(transport string, logger *zap.Logger) server.ToolHandlerMiddleware {
	apiKey := resolveAPIKey(logger)
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			logger.Debug("Auth middleware invoked",
				zap.String("context", "http"),
				zap.String("transport", transport),
				zap.String("tool", req.Params.Name),
			)

			if authenticated, err := isAuthenticated(ctx, transport, apiKey, logger); !authenticated {
				logger.Error("Authentication failed",
					zap.String("context", "http"),
					zap.String("transport", transport),
					zap.String("tool", req.Params.Name),
					zap.Error(err),
				)
				return nil, err
			}

			logger.Debug("Authentication successful",
				zap.String("context", "http"),
				zap.String("transport", transport),
				zap.String("tool", req.Params.Name),
			)

			return next(ctx, req)
		}
	}
}

// IsAuthenticated public api
func IsAuthenticated(ctx context.Context, transport string, logger *zap.Logger) (bool, error) {
	apiKey := resolveAPIKey(logger)
	return isAuthenticated(ctx, transport, apiKey, logger)
}

// isAuthenticated is the internal implementation that accepts a pre-resolved API key.
func isAuthenticated(ctx context.Context, transport string, apiKey string, logger *zap.Logger) (bool, error) {
	switch transport {
	case "stdio":
		return true, nil

	case "sse", "http":
		authenticated, err := validateToken(ctx, apiKey, logger)

		if err != nil {
			logger.Error("HTTP/SSE authentication error",
				zap.String("context", "http"),
				zap.Error(err),
			)
			return false, fmt.Errorf("authentication error: %w", err)
		}

		if !authenticated {
			logger.Warn("HTTP/SSE unauthorized request",
				zap.String("context", "http"),
			)
			return false, fmt.Errorf("unauthorized request")
		}

		return true, nil

	default:
		logger.Error("Unknown transport type",
			zap.String("context", "http"),
			zap.String("transport", transport),
		)
		return false, fmt.Errorf("unknown transport type: %s", transport)
	}
}
