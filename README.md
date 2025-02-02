# 待办事项应用

一个简单的待办事项管理系统，支持用户注册、登录，以及待办事项的增删改查等功能。

## 项目结构

```
todo-app/
├── config/             # 配置文件
├── controller/         # 控制器层，处理HTTP请求
│   ├── todo.go        # 待办事项控制器
│   └── user.go        # 用户控制器
├── model/             # 数据模型层
│   ├── todo.go        # 待办事项模型
│   └── user.go        # 用户模型
├── repository/        # 数据访问层
│   ├── todo_repo.go   # 待办事项数据库操作
│   └── user_repo.go   # 用户数据库操作
├── service/           # 业务逻辑层
│   ├── todo.go        # 待办事项服务
│   └── user.go        # 用户服务
├── static/            # 静态资源
│   ├── css/          # 样式文件
│   └── js/           # JavaScript文件
├── main.go           # 主程序入口
└── index.html        # 前端页面
```

## 功能特性

### 用户管理
- [x] 用户注册
- [x] 用户登录
- [x] JWT认证
- [x] 登出功能

### 待办事项管理
- [x] 添加待办事项
- [x] 设置开始和结束时间
- [x] 查看待办事项列表
- [x] 更新待办事项状态（完成/未完成）
- [x] 删除待办事项
- [x] 批量操作（完成/删除）
- [x] 按状态筛选（全部/未完成/已完成）
- [x] 搜索待办事项
- [x] 分页显示

## 技术栈

### 后端
- Go
- Hertz (Web框架)
- GORM (ORM框架)
- MySQL (数据库)
- JWT (认证)

### 前端(此部分是为了了解前后端是如何对接的，用cursor生成的。用于理解前端的基本工作和对接流程)
- HTML5
- CSS3
- JavaScript (原生)

## API文档

详细的API文档请参考 `docs/api.md`。

## 使用说明

1. 在当前目录下运行docker-compose up -d
2. 点击index.html即可访问

## 注意事项

- 所有API请求需要JWT认证（除了登录和注册）
- 每个用户只能看到和操作自己的待办事项
- 支持批量操作时要谨慎使用
