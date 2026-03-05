FROM golang:1.26 AS build

ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=local
ENV GOCACHE=/go/pkg/mod

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . /app

RUN --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w" -o /go/bin/mcp-server ./cmd/slack-mcp-server

FROM build AS dev

RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/go-delve/delve/cmd/dlv@v1.25.0 && cp /go/bin/dlv /dlv

RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

WORKDIR /app/mcp-server
RUN chown appuser:appgroup /app/mcp-server

USER appuser

EXPOSE 3001

ENTRYPOINT ["mcp-server"]
CMD ["--transport", "sse"]

FROM alpine:3.23 AS production

RUN apk add --no-cache ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=build /go/bin/mcp-server /usr/local/bin/mcp-server

WORKDIR /app
RUN chown appuser:appgroup /app

USER appuser

EXPOSE 3001

ENTRYPOINT ["mcp-server"]
CMD ["--transport", "sse"]
