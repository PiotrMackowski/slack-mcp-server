### 1. Authentication Setup

You need a Slack Bot token (`xoxb-*`) or User OAuth token (`xoxp-*`).

> **Note**: If both are provided, priority is `xoxp` > `xoxb`.

#### Option 1: Using `SLACK_MCP_XOXB_TOKEN` (Bot Token) — Recommended

Bot tokens have limited, well-defined scope and are the safer option.

1. Go to [api.slack.com/apps](https://api.slack.com/apps) and create a new app
2. Under "OAuth & Permissions", add Bot Token Scopes:
    - `channels:history` - View messages in public channels
    - `channels:read` - View basic information about public channels
    - `groups:history` - View messages in private channels
    - `groups:read` - View basic information about private channels
    - `im:history` - View messages in direct messages
    - `im:read` - View basic information about direct messages
    - `im:write` - Start direct messages with people
    - `mpim:history` - View messages in group direct messages
    - `mpim:read` - View basic information about group direct messages
    - `mpim:write` - Start group direct messages with people
    - `users:read` - View people in a workspace
    - `chat:write` - Send messages
    - `usergroups:read` - View user groups
    - `usergroups:write` - Create and manage user groups (only needed if using usergroup write tools)
3. Install the app to your workspace
4. Copy the "Bot User OAuth Token" (starts with `xoxb-`)
5. **Important**: The bot must be invited to channels it needs to access

> **Note**: Bot tokens cannot use `search.messages` API, so `conversations_search_messages` and `conversations_unreads` will not be available.

#### Option 2: Using `SLACK_MCP_XOXP_TOKEN` (User OAuth)

User tokens have broader access — all channels the user can see, plus search.

1. Go to [api.slack.com/apps](https://api.slack.com/apps) and create a new app
2. Under "OAuth & Permissions", add User Token Scopes:
    - `channels:history`, `channels:read`
    - `groups:history`, `groups:read`
    - `im:history`, `im:read`, `im:write`
    - `mpim:history`, `mpim:read`, `mpim:write`
    - `users:read`
    - `chat:write`
    - `search:read` - Search workspace content
    - `usergroups:read`, `usergroups:write`
3. Install the app to your workspace
4. Copy the "User OAuth Token" (starts with `xoxp-`)

##### App manifest (preconfigured scopes)

```json
{
    "display_information": {
        "name": "Slack MCP"
    },
    "oauth_config": {
        "scopes": {
            "user": [
                "channels:history",
                "channels:read",
                "groups:history",
                "groups:read",
                "im:history",
                "im:read",
                "im:write",
                "mpim:history",
                "mpim:read",
                "mpim:write",
                "users:read",
                "chat:write",
                "search:read",
                "usergroups:read",
                "usergroups:write"
            ]
        }
    },
    "settings": {
        "org_deploy_enabled": false,
        "socket_mode_enabled": false,
        "token_rotation_enabled": false
    }
}
```

See next: [Installation](02-installation.md)
