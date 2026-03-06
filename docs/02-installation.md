### 2. Installation

#### Pre-built binary (recommended)

Download the binary for your platform from the [Releases](https://github.com/PiotrMackowski/slack-mcp-server/releases) page.

```bash
# macOS Apple Silicon
curl -L -o slack-mcp-server \
  https://github.com/PiotrMackowski/slack-mcp-server/releases/latest/download/slack-mcp-server-darwin-arm64

# macOS Intel
curl -L -o slack-mcp-server \
  https://github.com/PiotrMackowski/slack-mcp-server/releases/latest/download/slack-mcp-server-darwin-amd64

# Linux x86_64
curl -L -o slack-mcp-server \
  https://github.com/PiotrMackowski/slack-mcp-server/releases/latest/download/slack-mcp-server-linux-amd64

chmod +x slack-mcp-server
```

Move it somewhere on your PATH or reference the full path in your MCP client config.

#### Build from source

Requires Go 1.24+.

```bash
git clone https://github.com/PiotrMackowski/slack-mcp-server.git
cd slack-mcp-server
go build -o slack-mcp-server ./cmd/slack-mcp-server/
```

See next: [Configuration and Usage](03-configuration-and-usage.md)
