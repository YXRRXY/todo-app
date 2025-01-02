# Todo App API 文档

## 基础信息
- 基础URL: `http://localhost:8888`
- 认证方式: Bearer Token (在 Header 中添加 `Authorization: Bearer {token}`)
- 响应格式: JSON

## API 接口列表

### 1. 用户认证

#### 1.1 用户注册
- **路径**: `/user/register`
- **方法**: `POST`
- **认证**: 不需要
- **请求体**:
```json
{
    "username": "string",
    "password": "string"
}
```
- **响应**:
```json
{
    "status": 200,
    "msg": "注册成功",
    "data": {
        "id": 1,
        "username": "string",
        "token": "jwt_token_string"
    }
}
```

#### 1.2 用户登录
- **路径**: `/user/login`
- **方法**: `POST`
- **认证**: 不需要
- **请求体**:
```json
{
    "username": "string",
    "password": "string"
}
```
- **响应**:
```json
{
    "status": 200,
    "msg": "登录成功",
    "data": {
        "id": 1,
        "username": "string",
        "token": "jwt_token_string"
    }
}
```

### 2. 待办事项管理

#### 2.1 创建待办事项
- **路径**: `/todo/add`
- **方法**: `POST`
- **认证**: 需要
- **请求体**:
```json
{
    "title": "string",
    "content": "string",
    "start_time": 1647302400,  // Unix时间戳
    "end_time": 1647388800
}
```
- **响应**:
```json
{
    "status": 200,
    "msg": "创建成功",
    "data": {
        "id": 1,
        "title": "string",
        "content": "string",
        "status": 0,
        "start_time": "2024-03-15 10:00:00",
        "end_time": "2024-03-15 12:00:00"
    }
}
```

#### 2.2 获取待办事项列表
- **路径**: `/todo/list`
- **方法**: `GET`
- **认证**: 需要
- **查询参数**:
  - `page`: 页码 (默认: 1)
  - `page_size`: 每页数量 (默认: 10)
  - `status`: 状态过滤 (可选: 0-未完成, 1-已完成)
- **响应**:
```json
{
    "status": 200,
    "msg": "ok",
    "data": {
        "items": [
            {
                "id": 1,
                "title": "string",
                "content": "string",
                "status": 0,
                "start_time": "2024-03-15 10:00:00",
                "end_time": "2024-03-15 12:00:00"
            }
        ],
        "total": 100
    }
}
```

#### 2.3 搜索待办事项
- **路径**: `/todo/search`
- **方法**: `GET`
- **认证**: 需要
- **查询参数**:
  - `keyword`: 搜索关键词
  - `page`: 页码
  - `page_size`: 每页数量
  - `status`: 状态过滤 (可选)
- **响应**: 同获取列表

#### 2.4 更新待办事项状态
- **路径**: `/todo/:id/status/:status`
- **方法**: `PUT`
- **认证**: 需要
- **路径参数**:
  - `id`: 待办事项ID
  - `status`: 新状态 (0或1)
- **响应**:
```json
{
    "status": 200,
    "msg": "更新成功"
}
```

#### 2.5 批量更新状态
- **路径**: `/todo/status/batch`
- **方法**: `PUT`
- **认证**: 需要
- **请求体**:
```json
{
    "status": 1,
    "current_status": 0,  // 可选，当前状态
    "ids": [1, 2, 3]     // 可选，指定ID列表
}
```
- **响应**:
```json
{
    "status": 200,
    "msg": "更新成功",
    "data": {
        "updated_count": 3
    }
}
```

#### 2.6 删除待办事项
- **路径**: `/todo/:id`
- **方法**: `DELETE`
- **认证**: 需要
- **路径参数**:
  - `id`: 待办事项ID
- **响应**:
```json
{
    "status": 200,
    "msg": "删除成功"
}
```

#### 2.7 批量删除
- **路径**: `/todo/batch`
- **方法**: `DELETE`
- **认证**: 需要
- **请求体**:
```json
{
    "status": 1,     // 可选，按状态删除
    "ids": [1, 2, 3] // 可选，指定ID列表
}
```
- **响应**:
```json
{
    "status": 200,
    "msg": "删除成功",
    "data": {
        "deleted_count": 3
    }
}
```

## 错误响应
所有接口在发生错误时会返回统一格式：
```json
{
    "status": 400,  // 或其他错误码
    "error": "错误信息描述"
}
```

## 常见错误码
- 400: 请求参数错误
- 401: 未授权或token无效
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误