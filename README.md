# Slack MCP Server (Allegro Security Fork)

Hardened fork of [korotovsky/slack-mcp-server](https://github.com/korotovsky/slack-mcp-server) — a Go-based [Model Context Protocol](https://modelcontextprotocol.io/) server for Slack workspaces. Supports stdio, SSE, and HTTP transports.

This fork removes browser session token support (xoxc/xoxd), uTLS fingerprinting, and other unnecessary attack surface. It uses only official Slack Bot API tokens (`xoxb-*`) or User OAuth tokens (`xoxp-*`).

## Quick Start

### 1. Download the binary

Grab the latest binary for your platform from the [Releases](https://github.com/PiotrMackowski/slack-mcp-server/releases) page.

```bash
# macOS Apple Silicon (most likely)
curl -L -o slack-mcp-server \
  https://github.com/PiotrMackowski/slack-mcp-server/releases/latest/download/slack-mcp-server-darwin-arm64
chmod +x slack-mcp-server
```

Or build from source:

```bash
git clone https://github.com/PiotrMackowski/slack-mcp-server.git
cd slack-mcp-server
go build -o slack-mcp-server ./cmd/slack-mcp-server/
```

### 2. Get a Slack token

You need a Bot token (`xoxb-*`) or User OAuth token (`xoxp-*`). See [Authentication Setup](docs/01-authentication-setup.md) for details on creating a Slack app and obtaining tokens.

Bot tokens are recommended for most use cases — they have a limited, well-defined scope.

### 3. Configure your MCP client

#### OpenCode (`opencode.jsonc`)

```jsonc
{
  "mcp": {
    "slack": {
      "type": "stdio",
      "command": "/path/to/slack-mcp-server",
      "args": ["--transport", "stdio"],
      "env": {
        "SLACK_MCP_XOXB_TOKEN": "xoxb-..."
      }
    }
  }
}
```

#### Claude Desktop (`claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "slack": {
      "command": "/path/to/slack-mcp-server",
      "args": ["--transport", "stdio"],
      "env": {
        "SLACK_MCP_XOXB_TOKEN": "xoxb-..."
      }
    }
  }
}
```

## Tools

All write tools are disabled by default and must be explicitly enabled via environment variables.

### Read-only (always registered)

| Tool | Description |
|------|-------------|
| `conversations_history` | Get messages from a channel or DM, with pagination by date or count |
| `conversations_replies` | Get a thread of messages by channel and `thread_ts` |
| `conversations_search_messages` | Search messages with filters (date, user, channel). Not available with bot tokens. |
| `channels_list` | List channels with optional sorting and filtering |
| `users_search` | Search users by name, email, or display name |
| `usergroups_list` | List user groups in the workspace |
| `usergroups_me` | Manage your own user group membership (list/join/leave) |
| `conversations_unreads` | Get unread messages across channels. Not available with bot tokens. |

### Write (opt-in via env vars)

| Tool | Env Var to Enable | Description |
|------|-------------------|-------------|
| `conversations_add_message` | `SLACK_MCP_ADD_MESSAGE_TOOL` | Post messages. Set to `true`, channel IDs, or `!`-prefixed blocklist. |
| `conversations_edit_message` | `SLACK_MCP_EDIT_MESSAGE_TOOL` | Edit messages. Same channel restriction syntax. |
| `conversations_delete_message` | `SLACK_MCP_DELETE_MESSAGE_TOOL` | Delete messages. Same channel restriction syntax. |
| `reactions_add` | `SLACK_MCP_ADD_MESSAGE_TOOL` | Add emoji reactions. Shares gate with add_message. |
| `reactions_remove` | `SLACK_MCP_ADD_MESSAGE_TOOL` | Remove emoji reactions. Shares gate with add_message. |
| `conversations_mark` | `SLACK_MCP_MARK_TOOL` | Mark channel as read. Set to `true` or `1`. |
| `usergroups_create` | `SLACK_MCP_USERGROUP_WRITE_TOOL` | Create user groups. Set to `true` or `1`. |
| `usergroups_update` | `SLACK_MCP_USERGROUP_WRITE_TOOL` | Update user group metadata. |
| `usergroups_users_update` | `SLACK_MCP_USERGROUP_WRITE_TOOL` | Replace user group members. |

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SLACK_MCP_XOXP_TOKEN` | Yes* | — | User OAuth token (`xoxp-...`) |
| `SLACK_MCP_XOXB_TOKEN` | Yes* | — | Bot token (`xoxb-...`). Bot must be invited to channels. |
| `SLACK_MCP_PORT` | No | `13080` | Port for SSE/HTTP transport |
| `SLACK_MCP_HOST` | No | `127.0.0.1` | Bind address for SSE/HTTP transport |
| `SLACK_MCP_API_KEY` | No | — | Bearer token for SSE/HTTP auth. **Required for SSE/HTTP modes.** |
| `SLACK_MCP_TLS_CERT` | No | — | Path to TLS certificate (SSE/HTTP) |
| `SLACK_MCP_TLS_KEY` | No | — | Path to TLS private key (SSE/HTTP) |
| `SLACK_MCP_CORS_ORIGIN` | No | — | Allowed CORS origin (SSE/HTTP). If unset, no CORS headers (same-origin only). |
| `SLACK_MCP_PROXY` | No | — | Proxy URL for outgoing Slack API requests |
| `SLACK_MCP_SERVER_CA` | No | — | Path to CA certificate for Slack API TLS verification |
| `SLACK_MCP_ADD_MESSAGE_TOOL` | No | — | Enable posting/reactions. `true`, channel IDs, or `!`-prefixed blocklist. |
| `SLACK_MCP_ADD_MESSAGE_MARK` | No | — | Auto-mark sent messages as read when set to `true` |
| `SLACK_MCP_ADD_MESSAGE_UNFURLING` | No | — | Enable link unfurling. `true` or comma-separated domain whitelist. |
| `SLACK_MCP_EDIT_MESSAGE_TOOL` | No | — | Enable editing. Same syntax as add_message. |
| `SLACK_MCP_DELETE_MESSAGE_TOOL` | No | — | Enable deletion. Same syntax as add_message. |
| `SLACK_MCP_USERGROUP_WRITE_TOOL` | No | — | Enable usergroup write tools. `true` or `1`. |
| `SLACK_MCP_MARK_TOOL` | No | — | Enable mark-as-read. `true` or `1`. |
| `SLACK_MCP_ENABLED_TOOLS` | No | — | Comma-separated tool whitelist. Overrides default registration. |
| `SLACK_MCP_USERS_CACHE` | No | OS default | Path to users cache JSON file |
| `SLACK_MCP_CHANNELS_CACHE` | No | OS default | Path to channels cache JSON file |
| `SLACK_MCP_LOG_LEVEL` | No | `info` | Log level: `debug`, `info`, `warn`, `error`, `panic`, `fatal` |
| `SLACK_MCP_GOVSLACK` | No | — | Set to `true` for FedRAMP GovSlack endpoints |

*One of `SLACK_MCP_XOXP_TOKEN` or `SLACK_MCP_XOXB_TOKEN` is required.

### Console Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `--transport` / `-t` | Yes | Transport: `stdio`, `sse`, or `http` |
| `--enabled-tools` / `-e` | No | Comma-separated tool whitelist (same as env var) |

## Resources

The server exposes two MCP resources for workspace metadata:

| URI | Format | Description |
|-----|--------|-------------|
| `slack://<workspace>/channels` | CSV | Directory of all channels (id, name, topic, purpose, memberCount) |
| `slack://<workspace>/users` | CSV | Directory of all users (userID, userName, realName) |

## Security Changes in This Fork

Compared to upstream `korotovsky/slack-mcp-server`:

- **Removed**: Browser session tokens (xoxc/xoxd), edge client, uTLS fingerprinting, Chrome UA spoofing
- **Removed**: `SLACK_MCP_SERVER_CA_INSECURE` (TLS verification bypass)
- **Removed**: Demo mode, npm packaging, DXT extension, Docker distribution
- **Removed**: slackdump/slackauth dependencies (playwright, go-rod, TUI libs)
- **Added**: Mandatory auth for SSE/HTTP transport (`SLACK_MCP_API_KEY`)
- **Added**: TLS support for SSE/HTTP (`SLACK_MCP_TLS_CERT` / `SLACK_MCP_TLS_KEY`)
- **Added**: CORS restriction (`SLACK_MCP_CORS_ORIGIN`)
- **Hardened**: Write tools gated behind env vars (usergroup writes, mark-as-read, edit, delete)
- **Hardened**: Error messages sanitized, search limits capped, bind address respects `SLACK_MCP_HOST`
- **Fixed**: Stdio transport no longer requires `SLACK_MCP_API_KEY`

## License

Licensed under MIT — see [LICENSE](LICENSE). This is not an official Slack product.

Based on [korotovsky/slack-mcp-server](https://github.com/korotovsky/slack-mcp-server).
