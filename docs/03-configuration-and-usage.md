## 3. Configuration and Usage

Configure the MCP server using command-line arguments and environment variables.

### Stdio transport (most common)

Stdio is the simplest setup — the MCP client starts the server as a subprocess.

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

### SSE / HTTP transport

For running the server as a standalone service. Requires `SLACK_MCP_API_KEY` for authentication.

```bash
export SLACK_MCP_XOXB_TOKEN="xoxb-..."
export SLACK_MCP_API_KEY="your-secret-key"
./slack-mcp-server --transport sse
```

The server listens on `127.0.0.1:13080` by default. Configure with `SLACK_MCP_HOST` and `SLACK_MCP_PORT`.

For TLS:

```bash
export SLACK_MCP_TLS_CERT="/path/to/cert.pem"
export SLACK_MCP_TLS_KEY="/path/to/key.pem"
```

### Console Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `--transport` / `-t` | Yes | Transport: `stdio`, `sse`, or `http` |
| `--enabled-tools` / `-e` | No | Comma-separated tool whitelist |

### Environment Variables

See the [full environment variables table](../README.md#environment-variables) in the README.

### Tool Registration and Permissions

Tools are controlled at two levels:
- **Registration** (`SLACK_MCP_ENABLED_TOOLS`) — determines which tools are visible to MCP clients
- **Runtime permissions** (tool-specific env vars) — channel restrictions for write tools

Write tools are **not registered by default**. To enable them, either:
1. Set their specific environment variable (e.g., `SLACK_MCP_ADD_MESSAGE_TOOL=true`), or
2. Explicitly list them in `SLACK_MCP_ENABLED_TOOLS`

#### Examples

**Read-only mode (default):**

```json
{
  "env": {
    "SLACK_MCP_XOXB_TOKEN": "xoxb-..."
  }
}
```

**Enable messaging to specific channels:**

```json
{
  "env": {
    "SLACK_MCP_XOXB_TOKEN": "xoxb-...",
    "SLACK_MCP_ADD_MESSAGE_TOOL": "C123456789,C987654321"
  }
}
```

**Enable messaging to all channels except specific ones:**

```json
{
  "env": {
    "SLACK_MCP_XOXB_TOKEN": "xoxb-...",
    "SLACK_MCP_ADD_MESSAGE_TOOL": "!C123456789"
  }
}
```

**Minimal read-only setup with specific tools:**

```json
{
  "env": {
    "SLACK_MCP_XOXB_TOKEN": "xoxb-...",
    "SLACK_MCP_ENABLED_TOOLS": "channels_list,conversations_history"
  }
}
```

#### Behavior Matrix

| `ENABLED_TOOLS` | Tool-specific env var | Write tool registered? | Channel restrictions |
|-----------------|----------------------|------------------------|---------------------|
| empty/not set   | not set              | No                     | N/A                 |
| empty/not set   | `true`               | Yes                    | None                |
| empty/not set   | `C123,C456`          | Yes                    | Only listed channels |
| includes tool   | not set              | Yes                    | None                |
| includes tool   | `C123,C456`          | Yes                    | Only listed channels |
| excludes tool   | any                  | No                     | N/A                 |

### Debugging

```bash
# Run with debug logging
SLACK_MCP_LOG_LEVEL=debug ./slack-mcp-server --transport stdio

# Run the MCP inspector
npx @modelcontextprotocol/inspector ./slack-mcp-server --transport stdio
```
