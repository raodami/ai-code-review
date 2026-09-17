# AI Code Review Assistant

## 项目概述
GitHub PR自动化审查工具 - 输入PR链接，AI自动生成代码审查报告

## 核心功能
- GitHub PR信息获取（diff、文件变更、评论）
- AI代码审查（DeepSeek API）
- 代码质量分析（复杂度、安全性、性能）
- 自动生成审查报告
- 支持多种导出格式
- 历史记录管理

## 技术栈
- 后端：Go + Gin + SQLite
- 前端：Next.js 14 + Tailwind CSS
- AI：DeepSeek API
- 集成：GitHub API

## 项目结构
```
D:/ai-code-review/
├── cmd/server/main.go          # 启动入口
├── internal/
│   ├── api/router.go           # API路由
│   ├── github/client.go        # GitHub API客户端
│   ├── review/review.go        # AI代码审查逻辑
│   ├── export/export.go        # 导出功能
│   └── store/store.go          # SQLite数据层
├── web/                        # Next.js前端
└── docs/PRD.md                 # 产品需求文档
```

## API端点
- POST /api/review - 提交PR进行审查
- GET /api/review/{id} - 获取审查结果
- GET /api/history - 获取历史记录
- GET /api/stats - 获取统计数据
- GET /api/export/{id}/{format} - 导出审查报告

## 环境配置
创建 `.env` 文件：
```
GITHUB_TOKEN=***
DEEPSEEK_API_KEY=***
PORT=8081
```
