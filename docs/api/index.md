# API 文档

DooPush 提供完整的 REST API 接口，支持通过程序化方式实现推送功能。

## 🔑 认证方式

DooPush API 按调用身份提供三种认证方式：

- **App Key**（`X-App-Key`）：客户端 SDK 注册设备、连接 Gateway 和上报统计，可以包含在应用中。
- **App Secret**（`Authorization: Bearer`）：客户服务端根据 Scope 调用推送发送接口。
- **JWT Token**（`Authorization: Bearer`）：Web 控制台用户操作。

App Key 永远不能发送推送。App Secret 与具体应用绑定，并使用最小权限 Scope。

## 📚 API 文档目录

### 🔐 认证相关
- [**API 认证**](./authentication.md) - App Key 和 App Secret 的使用方法

### 📨 推送相关
- [**推送接口**](./push-apis.md) - 单推、批量、广播推送 API

### 📱 设备相关
- [**设备接口**](./device-apis.md) - 设备注册、查询、分组、标签 API

### 📊 数据相关
- [**数据接口**](./data-apis.md) - 推送统计、日志、审计数据查询 API

## 🌐 API 基础信息

- **Base URL**: `https://doopush.com/api/v1`
- **认证方式**: App Key / App Secret / JWT Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 🛠 快速开始

1. 首先阅读 [API 认证](./authentication.md) 了解 App Key 与 App Secret 的边界
2. 查看 [推送接口](./push-apis.md) 了解如何发送推送
3. 参考具体接口文档中的代码示例

---

*所有 API 文档基于当前生产环境的实际接口，确保内容准确可用。*
