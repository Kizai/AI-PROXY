# AI-PROXY

AI-PROXY 是一个支持多种AI服务商（如OpenAI、Claude、Gemini等）的API统一代理、管理与监控系统。它为开发者和企业提供高效、可视化的API配置、请求转发、日志统计和健康监控能力。

## 主要功能
- **多API统一管理**：支持多家AI服务API的集中配置、启用/禁用、健康检测
- **智能路由转发**：自动将请求转发到最优API
- **实时监控统计**：请求量、成功率、响应时间等多维度统计与可视化
- **API配置管理**：可视化界面，支持API测试、编辑、删除
- **请求日志**：详细记录所有请求，支持筛选、导出、清空
- **权限认证**：基于Token的后台管理权限

## 技术架构
- **后端**：Go 1.18+ / Gin / GORM / MySQL
- **前端**：HTML5 / CSS3 / JavaScript / Bootstrap 5 / Chart.js
- **数据库**：MySQL 8.0+

## 目录结构
```
AI-PROXY/
├── main.go                # 程序入口
├── config.json            # 主配置文件
├── model/                 # 数据模型
├── repository/            # 数据访问层
├── service/               # 业务逻辑层
├── controller/            # 控制器
├── router/                # 路由配置
├── util/                  # 工具函数
├── web/                   # 前端页面与静态资源
│   ├── pages/             # 主要页面（dashboard, api-config, ...）
│   ├── js/                # 前端JS逻辑
│   ├── css/               # 样式文件
│   └── assets/            # 图片等资源
└── ...
```

## 安装与启动
1. **环境准备**：Go 1.18+、MySQL 8.0+
2. **克隆项目**：
   ```bash
   git clone <your-repo-url>
   cd AI-PROXY
   ```
3. **配置数据库**：
   - 创建数据库 `proxy`
   - 修改 `config.json` 填写数据库连接信息
4. **安装依赖**：
   ```bash
   go mod tidy
   ```
5. **启动后端**：
   ```bash
   go run main.go
   ```
6. **访问前端**：
   - 浏览器打开 `http://localhost:8080` 进入管理后台

## 配置说明
- `config.json` 包含 server、database、log、apis 等配置项
- 支持多API配置，示例：
```json
"apis": {
  "openai": {
    "base_url": "https://api.openai.com",
    "headers": {"Content-Type": "application/json"},
    "auth_type": "bearer",
    "auth_value": "your_openai_api_key",
    "timeout": 60,
    "rate_limit": 0,
    "description": "OpenAI API"
  }
}
```

## 主要API接口
- **代理转发**：`/{apiName}/{path}` 支持所有HTTP方法
- **管理接口**（需Token认证）：
  - `GET /admin/api-config` 获取API配置列表
  - `POST /admin/api-config` 新增API配置
  - `PUT /admin/api-config/{name}` 更新API配置
  - `DELETE /admin/api-config/{name}` 删除API配置
  - `POST /admin/api-config/test` 测试API配置
  - `GET /admin/logs` 查询请求日志
  - `DELETE /admin/logs` 批量删除日志
  - `GET /admin/stats` 获取统计数据

## 前端用法
- 管理后台支持：API配置、测试、日志、统计、健康监控等全部可视化操作
- 支持Token登录，默认Token可在 `config.json` 或环境变量中设置

## 常见问题
- **数据库连接失败**：请检查 `config.json` 数据库配置与MySQL服务状态
- **端口被占用**：修改 `config.json` 中的 `port` 或释放端口
- **页面空白/报错**：请检查后端服务是否正常启动，浏览器控制台查看错误

## 贡献与许可证
- 欢迎PR和Issue，详见项目贡献指南
- 本项目采用 MIT 许可证 