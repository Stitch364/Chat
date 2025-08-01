

# 基于Socket.IO的实时聊天应用

这是一个使用Go语言开发的聊天应用，基于Gin框架和Socket.IO实现。系统包含账户管理、消息传递、群组管理、好友关系管理等完整聊天功能。

## 技术架构
- 域
  - `controller`：处理HTTP请求
  - `logic`：业务逻辑处理
  - `model`：数据模型和请求/响应结构
  - `dao`：数据访问层
    - `mysql`：MySQL数据库操作
    - `redis`：Redis缓存操作
  - `pkg`：通用功能包
    - `emailMark`：邮箱验证码
    - `retry`：重试机制
    - `tool`：工具函数
  - `setting`：系统配置初始化
  - `manager`：连接管理
  - `middlewares`：中间件
    - 认证授权
    - 跨域处理
    - 日志记录

## 核心功能

### 账户管理
- 创建/删除账户
- 获取账户信息
- 更新账户信息
- 账户登录认证

### 聊天功能
- 实时消息发送
- 消息状态更新（已读/置顶/撤销）
- 消息搜索
- 消息持久化存储

### 群组管理
- 创建/解散群组
- 转移群主
- 更新群组信息
- 邀请/退出群组
- 获取群成员列表

### 好友关系
- 好友申请管理
- 接受/拒绝申请
- 获取好友列表
- 删除好友

### 文件管理
- 文件上传/下载
- 文件信息存储
- 群头像更新
- 账户头像管理

### 设置管理
- 消息提醒设置
- 群设置
- 昵称设置
- 会话置顶设置

## 系统配置
- MySQL数据库配置
- Redis缓存配置
- 文件存储配置
- 日志配置
- Token生成配置
- 邮箱验证码配置
- 工作池配置

## API文档
请参考各模块的`controller`实现，包含：
- `/api/account`：账户相关API
- `/api/user`：用户相关API
- `/api/group`：群组相关API
- `/api/message`：消息相关API
- `/api/setting`：设置相关API
- `/api/file`：文件相关API
- `/api/application`：申请相关API
- `/api/email`：邮箱相关API

## 项目构建
使用Dockerfile进行构建：
```Dockerfile
FROM golang:alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY .env .
COPY go.mod .
COPY go.sum .
COPY . .
...
```

## 数据模型
- MySQL数据模型
  - `Account`：账户信息
  - `User`：用户信息
  - `Message`：消息数据
  - `Relation`：关系管理
  - `Setting`：设置信息
  - `File`：文件存储

## 版本控制
项目使用`.gitignore`进行版本控制，忽略敏感配置和临时文件。

## 错误处理
统一使用`errcode.Err`进行错误处理，包含：
- 用户相关错误
- 账户相关错误
- 消息相关错误
- 申请相关错误
- 设置相关错误
- 文件相关错误

## 消息队列
支持RocketMQ消息队列，包含生产者和消费者实现。

## 许可协议
请查看源仓库获取具体的许可协议信息

## 项目入口
`main.go`为项目启动入口文件，初始化整个应用。

## 初始化配置
通过`setting`包完成系统初始化：
- 数据库连接初始化
- Redis连接初始化
- 日志系统初始化
- 文件存储初始化
- Socket.IO连接管理初始化

## 实时通信
使用Socket.IO实现：
- 消息推送
- 在线状态管理
- 实时通信处理

## 开发工具
- `retry`包：提供重试机制
- `gtype`包：通用类型处理
- `tool`包：错误处理工具
- `wait-for.sh`：容器启动等待脚本

## 目录结构
```
/chat
  - 主要业务逻辑
/controller
  - API接口定义
/dao
  - 数据访问层
/model
  - 数据模型定义
/pkg
  - 通用组件
/setting
  - 系统配置初始化
/manager
  - 连接管理
```