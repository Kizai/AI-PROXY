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

## 环境要求与安装

### 必需软件
1. **Go 1.18+**
   ```bash
   # 检查版本
   go version
   
   # 下载地址：https://golang.org/dl/
   # Windows: 下载 .msi 安装包
   # Linux: sudo apt install golang-go 或下载官方包
   # Mac: brew install go
   ```

2. **MySQL 8.0+**
   ```bash
   # 检查版本
   mysql --version
   
   # Windows: 下载 MySQL Installer
   # Linux: sudo apt install mysql-server
   # Mac: brew install mysql
   ```

3. **Git**
   ```bash
   # 检查版本
   git --version
   
   # Windows: 下载 Git for Windows
   # Linux: sudo apt install git
   # Mac: brew install git
   ```

### 项目获取
```bash
# 克隆项目
git clone <your-repo-url>
cd AI-PROXY

# 检查项目结构
ls -la
```

## 数据库配置

### 1. 创建数据库和用户
```sql
-- 登录MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE proxy_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户（替换 'your_password' 为你的密码）
CREATE USER 'proxy_user'@'localhost' IDENTIFIED BY 'your_password';

-- 授权
GRANT ALL PRIVILEGES ON proxy_db.* TO 'proxy_user'@'localhost';
FLUSH PRIVILEGES;

-- 验证
SHOW DATABASES;
SELECT User, Host FROM mysql.user WHERE User = 'proxy_user';
```

### 2. 配置数据库连接
编辑 `config.json` 文件：
```json
{
  "server": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "database": {
    "host": "localhost",
    "port": 3306,
    "username": "proxy_user",
    "password": "your_password",
    "database": "proxy_db",
    "charset": "utf8mb4"
  },
  "log": {
    "level": "info",
    "file": "log/app.log"
  },
  "auth": {
    "token": "your_admin_token_here"
  },
  "apis": {}
}
```

## 安装与启动

### 1. 安装Go依赖
```bash
# 在项目根目录下执行
go mod tidy

# 检查依赖
go mod verify
```

### 2. 启动后端服务
```bash
# 开发模式启动
go run main.go

# 或编译后运行
go build -o ai-proxy main.go
./ai-proxy
```

### 3. 访问前端
- 浏览器打开：`http://localhost:8080`
- 使用 `config.json` 中配置的 `auth.token` 登录管理后台

## 配置文件详解

### config.json 字段说明
```json
{
  "server": {
    "port": 8080,           // 服务端口
    "host": "0.0.0.0"      // 监听地址，0.0.0.0表示所有地址
  },
  "database": {
    "host": "localhost",    // 数据库主机
    "port": 3306,          // 数据库端口
    "username": "proxy_user", // 数据库用户名
    "password": "your_password", // 数据库密码
    "database": "proxy_db", // 数据库名
    "charset": "utf8mb4"   // 字符集
  },
  "log": {
    "level": "info",        // 日志级别：debug, info, warn, error
    "file": "log/app.log"   // 日志文件路径
  },
  "auth": {
    "token": "your_admin_token_here" // 管理后台访问Token
  },
  "apis": {                 // API配置（初始为空）
    "openai": {
      "base_url": "https://api.openai.com",
      "headers": {
        "Content-Type": "application/json"
      },
      "auth_type": "bearer",
      "auth_value": "sk-your-openai-key",
      "timeout": 60,
      "rate_limit": 0,
      "description": "OpenAI API"
    }
  }
}
```

## 主要API接口

### 代理转发接口
- **格式**：`/{apiName}/{path}`
- **方法**：支持所有HTTP方法（GET、POST、PUT、DELETE等）
- **示例**：
  ```
  POST /openai/v1/chat/completions
  GET /claude/v1/messages
  ```

### 管理接口（需Token认证）
```bash
# API配置管理
GET    /admin/api-config          # 获取API配置列表
POST   /admin/api-config          # 新增API配置
PUT    /admin/api-config/{name}   # 更新API配置
DELETE /admin/api-config/{name}   # 删除API配置
POST   /admin/api-config/test     # 测试API配置

# 日志管理
GET    /admin/logs                # 查询请求日志
DELETE /admin/logs                # 批量删除日志

# 统计信息
GET    /admin/stats               # 获取统计数据
```

## 前端使用指南

### 登录管理后台
1. 访问 `http://localhost:8080`
2. 输入 `config.json` 中配置的 `auth.token`
3. 点击登录进入管理界面

### 主要功能页面
- **Dashboard**：总览统计、健康监控
- **API配置**：添加、编辑、测试、删除API配置
- **请求日志**：查看、筛选、清空请求记录
- **统计图表**：成功率、响应时间等可视化数据

## 服务器部署

### 1. 环境准备
```bash
# 安装必要软件
sudo apt update
sudo apt install golang-go mysql-server nginx git

# 启动MySQL
sudo systemctl start mysql
sudo systemctl enable mysql
```

### 2. 部署项目
```bash
# 克隆项目
git clone <your-repo-url>
cd AI-PROXY

# 配置数据库（参考本地配置步骤）
# 配置config.json

# 编译
go build -o ai-proxy main.go
```

### 3. 配置systemd服务
创建 `/etc/systemd/system/ai-proxy.service`：
```ini
[Unit]
Description=AI-PROXY Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/AI-PROXY
ExecStart=/path/to/AI-PROXY/ai-proxy
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable ai-proxy
sudo systemctl start ai-proxy
sudo systemctl status ai-proxy
```

### 4. 配置Nginx反向代理
创建 `/etc/nginx/sites-available/ai-proxy`：
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/ai-proxy /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 常见问题与排查

### 1. 启动问题
**问题**：`package not in std` 或找不到包
```bash
# 解决：检查go.mod和import路径一致性
cat go.mod
# 确保module名称与import路径匹配
# 重新整理依赖
go mod tidy
```

**问题**：端口被占用
```bash
# 查看端口占用
netstat -tulpn | grep 8080
# 或
lsof -i :8080

# 停止占用进程或修改config.json中的端口
```

### 2. 数据库问题
**问题**：数据库连接失败
```bash
# 检查MySQL服务状态
sudo systemctl status mysql

# 测试连接
mysql -u proxy_user -p proxy_db

# 检查config.json配置
cat config.json
```

**问题**：表结构问题
```bash
# 检查表是否存在
mysql -u proxy_user -p proxy_db -e "SHOW TABLES;"

# 如需重建表，删除数据库重新创建
mysql -u root -p -e "DROP DATABASE proxy_db; CREATE DATABASE proxy_db;"
```

### 3. 前端问题
**问题**：页面空白或"Failed to fetch"
```bash
# 检查后端服务状态
curl http://localhost:8080/admin/stats

# 检查浏览器控制台错误
# 确认API_BASE_URL配置正确
```

**问题**：登录后页面变灰
- 检查登录模态框是否正确关闭
- 确认Token配置正确

### 4. API配置问题
**问题**：保存API配置报400错误
- 检查必填字段是否完整
- 确认API密钥格式正确
- 检查网络连接

**问题**：清空日志报500错误
```bash
# 检查GORM软删除配置
# 确认数据库权限
mysql -u proxy_user -p proxy_db -e "SHOW GRANTS;"
```

### 5. 性能问题
**问题**：响应慢或超时
- 检查网络连接
- 调整config.json中的timeout配置
- 检查数据库性能

## 开发与维护

### 代码结构
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
│   ├── pages/             # 主要页面
│   ├── js/                # 前端JS逻辑
│   ├── css/               # 样式文件
│   └── assets/            # 图片等资源
└── doc/                   # 文档目录
    └── CHANGELOG.md       # 开发日志
```

### 开发日志维护
- 重要修改记录在 `doc/CHANGELOG.md`
- 格式：UTC时间戳 + 修改说明
- 包含问题排查和解决方案

### 版本更新
1. 更新代码
2. 重新编译：`go build -o ai-proxy main.go`
3. 重启服务：`sudo systemctl restart ai-proxy`
4. 更新CHANGELOG.md

## FAQ

**Q: 如何添加新的AI服务商？**
A: 在管理后台的"API配置"页面添加新的API配置，填写base_url、认证信息等。

**Q: 如何修改管理Token？**
A: 修改config.json中的auth.token字段，重启服务生效。

**Q: 如何备份数据？**
A: 导出数据库：`mysqldump -u proxy_user -p proxy_db > backup.sql`

**Q: 如何查看详细日志？**
A: 查看log/app.log文件，或修改config.json中的log.level为debug。

**Q: 支持哪些AI服务商？**
A: 理论上支持所有提供HTTP API的AI服务，包括OpenAI、Claude、Gemini、百度文心等。

## 贡献与许可证
- 欢迎提交PR和Issue
- 本项目采用 MIT 许可证
- 详细贡献指南请参考项目文档 

## 重要变更

- 2024-06-XX：应产品需求，已彻底移除“请求日志”功能，包括前端页面、后端接口、数据库表等，其他功能不受影响。 