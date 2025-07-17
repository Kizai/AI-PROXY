# AI-PROXY 快速上手与部署手册

## 项目简介
AI-PROXY 是一个开箱即用的 AI API 代理平台，支持多家主流 AI 厂商（如 OpenAI、Claude、Gemini、Grok 等），让你可以用自己的服务器统一转发和管理所有 AI API 请求。

<img width="1931" height="778" alt="QQ_1752766224004" src="https://github.com/user-attachments/assets/6fd2cae0-d715-4671-9ef0-a44b62c486b7" />

<img width="1955" height="777" alt="QQ_1752766255971" src="https://github.com/user-attachments/assets/62bafa27-2cd6-4f61-9854-eefc1bf7c97f" />

<img width="1964" height="555" alt="QQ_1752766115130" src="https://github.com/user-attachments/assets/c2eff4ef-2551-441e-b1f2-854e978351ab" />



- 支持多厂商、多模型
- 无需注册、无需登录，直接用
- 管理后台可动态添加/删除/禁用 API
- 适合个人、团队、企业自建安全代理

---

## 环境要求
- 操作系统：Windows、Linux、MacOS 均可
- 运行环境：Go 1.18+（推荐 1.20 及以上）
- 数据库：MySQL 5.7+ 或 MariaDB 10.2+
- 推荐有公网服务器和域名（如需对外服务）

---

## 快速部署步骤

### 1. 下载源码

```bash
git clone https://github.com/你的仓库/AI-PROXY.git
cd AI-PROXY
```

### 2. 配置数据库

- 新建一个 MySQL 数据库（如：proxy）
- 修改 `config.json` 里的数据库连接信息（用户名、密码、库名等）

### 3. 编译并启动服务

#### Windows：
```bash
go build -o ai-proxy.exe main.go
./ai-proxy.exe
```

#### Linux/Mac：
```bash
go build -o ai-proxy main.go
./ai-proxy
```

启动后，终端会显示 `HTTP服务器启动,监听地址: 0.0.0.0:8080`

### 4. 访问前端页面

- 用户首页：http://你的服务器IP:8080/
- 管理后台：http://你的服务器IP:8080/admin

### 5. 登录管理后台
- 默认需要输入 `config.json` 里配置的 `auth.token` 作为访问令牌
- 登录后可添加/编辑/删除 API 配置

### 6. 添加你的 API 配置
- 在管理后台“API配置”页面，点击“添加API配置”
- 填写 API 名称、基址URL（如 https://api.openai.com）、描述，勾选启用
- 保存即可

### 7. 开始使用代理
- 直接用 http://你的服务器IP:8080/厂商名/xxx 作为 API 地址
- 例如：
  - OpenAI: http://你的服务器IP:8080/openai/v1/chat/completions
  - Claude: http://你的服务器IP:8080/claude/v1/messages
  - Gemini: http://你的服务器IP:8080/gemini/v1beta/models/gemini-pro:generateContent

---

## 常见问题

**Q: 为什么访问管理后台需要令牌？**
A: 为了安全，防止他人随意修改你的API配置。

**Q: 如何让别人也能用我的代理？**
A: 只需把你的服务器IP或域名告诉对方，对方用你的代理地址即可。

**Q: 如何支持 HTTPS？**
A: 推荐用 Nginx/Caddy 配置 SSL 证书做反向代理，详见官方教程或联系你的服务器运维。

**Q: 如何彻底删除某个API配置？**
A: 在管理后台删除即可，后台会物理删除数据库记录。

**Q: 数据库出错/端口被占用怎么办？**
A: 检查数据库配置和端口占用情况，或换一个端口。

---

## 联系方式
如遇到无法解决的问题，可联系开发者邮箱：support@gptocean.com

---

祝你用得顺利！如需定制开发或企业支持，请联系作者。 
