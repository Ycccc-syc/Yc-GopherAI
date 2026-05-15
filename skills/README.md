# Skills 使用指南 — AI 全栈挑战赛

## 工作流

```
接到需求
    │
    ▼
/planner  ───────────── 输出技术方案 + 分步计划
    │
    ▼
/fullstack-builder ──── 按计划逐步实现
    │
    ▼
/simplify  ──────────── review 代码质量
    │
    ▼
/loop /simplify ─────── 持续监控（可选）
    │
    ▼  (出 Bug)
/debug    ──────────── 分析根因并修复
```

## 快速启动新项目

1. 复制 `skills/` 文件夹到新项目
2. 编辑项目根目录 `CLAUDE.md`（基于模板修改）
3. 将 `skills-config-template.json` 的内容添加到 VSCode `settings.json`
4. 将 `mcp-servers.md` 的 MCP 配置添加到 VSCode `settings.json`
5. 安装依赖：Node.js、Playwright、Docker

## Skills 一览

| Skill | 功能 | 何时用 |
|---|---|---|
| `/planner` | 架构设计 + 实施计划 | 接到需求第一步 |
| `/fullstack-builder` | 全栈功能实现 | 有明确任务后 |
| `/debug` | Bug 根因分析 | 出问题后 |
| `/simplify` | 代码质量 review | 实现完后 |
| `/review` | PR 审查 | 提 PR 前 |
| `/fewer-permission-prompts` | 减少权限弹窗 | 新项目初始化 |
| `/init` | 生成 CLAUDE.md | 新项目初始化 |

## 技能目录

```
skills/
├── README.md                  ← 本文件
├── planner.md                 ← Planner Skill 定义
├── fullstack-builder.md       ← Fullstack Builder Skill 定义
├── debug.md                   ← Debug Skill 定义
├── mcp-servers.md             ← MCP 服务器配置模板
└── skills-config-template.json ← Skills VSCode 配置模板
```
