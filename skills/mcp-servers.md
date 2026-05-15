# MCP 服务器配置模板

将此配置添加到 VSCode `settings.json` 的 `claudeCode.mcpServers` 中。

```json
{
    "claudeCode.mcpServers": {
        "playwright": {
            "command": "npx",
            "args": ["@playwright/mcp"]
        },
        "docker": {
            "command": "npx",
            "args": ["-y", "mcp-docker-server"]
        },
        "sqlite": {
            "command": "npx",
            "args": ["-y", "mcp-server-sqlite", "--db-path", "${workspaceFolder}\\data.db"]
        }
    }
}
```

## 前置依赖

- **Node.js** (LTS) — 所有 MCP 依赖 npx 运行
- **Playwright** — `npx playwright install chromium` (首次使用前运行)
- **Docker** — Docker MCP 需要 Docker 引擎
- **uv** (可选) — Python 生态 MCP 使用: `powershell -c "irm https://astral.sh/uv/install.ps1 | iex"`
