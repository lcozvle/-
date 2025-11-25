# 📺 DY-DanMu: 斗鱼弹幕实时监控系统

![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![Redis](https://img.shields.io/badge/Redis-5.0+-DC382D?style=flat&logo=redis)
![MySQL](https://img.shields.io/badge/MySQL-5.7+-4479A1?style=flat&logo=mysql)
![License](https://img.shields.io/badge/License-MIT-green.svg)

**DY-DanMu** 是一个基于 Go 语言开发的高性能、微服务架构的斗鱼直播弹幕监控系统。它能够实时抓取指定直播间的弹幕，进行持久化存储，并提供 Web 界面进行实时展示和管理。

> ✨ **最新功能**：支持网页端动态切换直播间、一键启动/关闭系统、以及导出弹幕词频数据供 AI 分析。

## 🌟 功能特性

*   **微服务架构**：
    *   🕷️ **Spider**: 基于 WebSocket 协议，实时连接斗鱼弹幕服务器抓取数据。
    *   💾 **PersistServer**: 负责数据清洗、去重，并持久化到 MySQL 和 ElasticSearch（可选）。
    *   🌐 **Web Server**: 提供 RESTful API 和前端监控页面。
*   **实时监控**：Web 页面实时显示弹幕流，支持无限滚动和自动跟随。
*   **动态控制**：无需重启服务，直接在网页上输入房间号即可切换监控目标。
*   **数据隔离**：支持多房间数据存储，查看和导出时自动按房间隔离数据。
*   **AI 数据支持**：提供专门的 API 导出弹幕词频统计数据（JSON格式），方便对接 LLM 进行情感分析或热词统计。
*   **一键管理**：提供 Windows 批处理脚本，一键启动所有服务，网页端一键优雅关闭。

## 🏗️ 技术栈

*   **语言**: Golang
*   **Web 框架**: Gin
*   **数据库**: MySQL (持久化), Redis (缓存/消息/配置), ElasticSearch (全文检索/可选)
*   **通信**: gRPC (服务间通信), WebSocket (弹幕抓取)
*   **前端**: 原生 HTML/CSS/JS (轻量级监控页)

## 🚀 快速开始

### 1. 环境准备

确保你的环境已安装：
*   Go (1.18+)
*   MySQL
*   Redis
*   (可选) ElasticSearch

### 2. 数据库配置

1.  创建 MySQL 数据库 `dybarrage`。
2.  导入表结构（系统会自动创建，但需确保数据库存在）。
3.  **重要**：确保 `barrage` 表包含 `rid` 字段（如果使用旧版数据库，请执行 `ALTER TABLE barrage ADD COLUMN rid VARCHAR(50);`）。

### 3. 修改配置

编辑 `start_all.bat` (Windows) 或环境变量，设置你的数据库密码：

```bat
set MYSQLPWD=你的MySQL密码
set REDISPWD=你的Redis密码(如果有)
```

### 4. 一键启动 (Windows)

双击项目根目录下的 **`start_all.bat`**。

脚本将自动：
1.  启动 PersistServer
2.  启动 Spider
3.  启动 Web Server
4.  自动打开浏览器访问监控页面 (`http://localhost:8080/monitor.html`)

## 🖥️ 使用指南

### 监控页面
启动后，浏览器会自动打开监控页。
*   **状态栏**：显示当前连接状态和弹幕总数。
*   **控制面板**：
    *   输入框：输入新的斗鱼房间号（如 `9999`）。
    *   **切换房间**：点击后 Spider 会自动重连到新房间，无需重启。
    *   **导出 AI 数据**：下载当前房间的弹幕统计 JSON 文件。
    *   **关闭系统**：一键关闭所有后台服务进程。

### API 接口

| 方法 | 路径 | 描述 |
| :--- | :--- | :--- |
| POST | `/config/rid` | 修改当前监控的直播间 ID |
| GET | `/config/current_rid` | 获取当前监控的直播间 ID |
| GET | `/export/ai` | 导出当前房间的弹幕词频统计 (JSON) |
| POST | `/search/all` | 查询弹幕列表 (支持分页和 rid 过滤) |
| POST | `/system/shutdown` | 关闭所有服务进程 |

## 📂 项目结构

```
DY-DanMu/
├── DMconfig/       # 全局配置文件
├── dbConn/         # 数据库连接初始化
├── persistServer/  # [微服务] 数据持久化服务 (RPC Server)
├── spider/         # [微服务] 弹幕抓取服务 (WebSocket Client)
├── web/            # [微服务] Web API 和静态资源
├── monitor.html    # 前端监控页面
├── start_all.bat   # Windows 一键启动脚本
└── README.md       # 项目文档
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 开源协议

MIT License