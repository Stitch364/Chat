# Chat Application Based on Socket.IO

This is a chat application developed using the Go programming language, built on the Gin framework and Socket.IO. The system includes complete chat functionalities such as account management, message transmission, group management, and friend relationship management.

## Technical Architecture
- **Domain**
  - `controller`: Handles HTTP requests
  - `logic`: Processes business logic
  - `model`: Defines data models and request/response structures
  - `dao`: Data access layer
    - `mysql`: MySQL database operations
    - `redis`: Redis cache operations
  - `pkg`: General-purpose utility packages
    - `emailMark`: Email verification code
    - `retry`: Retry mechanism
    - `tool`: Utility functions
  - `setting`: System configuration initialization
  - `manager`: Connection management
  - `middlewares`: Middleware components
    - Authentication and authorization
    - Cross-origin handling
    - Logging

## Core Features

### Account Management
- Create/Delete accounts
- Retrieve account information
- Update account information
- Account authentication

### Chat Features
- Real-time message sending
- Message status updates (read/pinned/revoked)
- Message search
- Persistent message storage

### Group Management
- Create/Disband groups
- Transfer group ownership
- Update group information
- Invite/Leave groups
- Retrieve group member list

### Friend Relationships
- Friend request management
- Accept/Reject requests
- Retrieve friend list
- Remove friends

### File Management
- File upload/download
- File information storage
- Group avatar update
- Account avatar management

### Settings Management
- Notification settings
- Group settings
- Nickname settings
- Conversation pinning settings

## System Configuration
- MySQL database configuration
- Redis cache configuration
- File storage configuration
- Logging configuration
- Token generation configuration
- Email verification code configuration
- Worker pool configuration

## API Documentation
Please refer to the `controller` implementations for each module, including:
- `/api/account`: Account-related APIs
- `/api/user`: User-related APIs
- `/api/group`: Group-related APIs
- `/api/message`: Message-related APIs
- `/api/setting`: Settings-related APIs
- `/api/file`: File-related APIs
- `/api/application`: Application-related APIs
- `/api/email`: Email-related APIs

## Project Build
Build using the Dockerfile:
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

## Data Models
- MySQL Data Models
  - `Account`: Account information
  - `User`: User information
  - `Message`: Message data
  - `Relation`: Relationship management
  - `Setting`: Settings information
  - `File`: File storage

## Version Control
The project uses `.gitignore` for version control, excluding sensitive configurations and temporary files.

## Error Handling
Errors are uniformly handled using `errcode.Err`, including:
- User-related errors
- Account-related errors
- Message-related errors
- Application-related errors
- Settings-related errors
- File-related errors

## Message Queue
Supports RocketMQ message queue, including producer and consumer implementations.

## License
Please check the source repository for specific license information.

## Project Entry Point
`main.go` is the entry file for starting the application.

## Initialization Configuration
System initialization is handled through the `setting` package:
- Database connection initialization
- Redis connection initialization
- Logging system initialization
- File storage initialization
- Socket.IO connection management initialization

## Real-Time Communication
Implemented using Socket.IO:
- Message push notifications
- Online status management
- Real-time communication handling

## Development Tools
- `retry` package: Provides retry mechanism
- `gtype` package: General-purpose type handling
- `tool` package: Error handling utilities
- `wait-for.sh`: Container startup wait script

## Directory Structure
```
/chat
  - Main business logic
/controller
  - API interface definitions
/dao
  - Data access layer
/model
  - Data model definitions
/pkg
  - General-purpose components
/setting
  - System configuration initialization
/manager
  - Connection management
```