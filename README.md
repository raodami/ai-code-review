# AI Code Review Assistant

## 功能特性
- GitHub PR URL输入 → 自动获取变更信息
- AI代码审查（支持DeepSeek API）
- 问题分级（高/中/低优先级）
- 质量评分（0-100分）
- 导出报告（Markdown/JSON/Text）
- 历史记录管理

## 技术栈
- 后端：Go + Gin + SQLite
- 前端：Next.js 14 + Tailwind CSS
- AI：DeepSeek API

## 安装运行

### 后端
```bash
cd D:/ai-code-review
go run cmd/server/main.go
```

### 前端
```bash
cd D:/ai-code-review/web
npm install
npm run dev
```

## 环境变量
创建 `.env` 文件：
```
GITHUB_TOKEN=your_github_token
DEEPSEEK_API_KEY=your_deepseek_key
PORT=8081
```

## API端点
- POST /api/review - 提交PR审查
- GET /api/history - 获取历史记录
- GET /api/stats - 获取统计数据
- GET /api/export/:id/:format - 导出报告
