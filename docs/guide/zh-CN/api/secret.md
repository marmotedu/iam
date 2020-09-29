# 密钥相关接口

## 1. 创建密钥

### 1.1 接口描述

创建密钥。

### 1.2 请求方法

POST /v1/secrets

### 1.3 输入参数

**Body 参数**

| 参数名称 | 必选 | 类型                      | 描述               |
| -------- | ---- | ------------------------- | ------------------ |
| metadata | 是   | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| expires | 否   | Int64                    | 过期时间               |
| description | 否   | String                    | 密钥描述               |

### 1.4 输出参数

| 参数名称    | 类型                                 | 描述                |
| ----------- | ------------------------------------ | ------------------- |
| metadata    | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| username    | String                               | 用户名              |
| secretID    | String                               | 密钥 ID              |
| secretKey   | String                               | 密钥 Key             |
| expires     | Int64                                | 过期时间            |
| description | String                               | 密钥描述            |

### 1.5 请求示例

**输入示例**

```bash
 curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'{
  "metadata": {
    "name": "secret"
  },
  "expires": 0,
  "description": "admin secret"
}' http://marmotedu.io:8080/v1/secrets
```
**输出示例**

```json
{
  "metadata": {
    "id": 28,
    "name": "secret",
    "createdAt": "2020-09-23T11:03:43.189962859+08:00",
    "updatedAt": "2020-09-23T11:03:43.189962859+08:00"
  },
  "username": "admin",
  "secretID": "lXirSIJV5tA34V8hffffFYq7CnDhfc4gDxrz",
  "secretKey": "PK8NMhHnapVdNHAoPxhrN5Beg0C5fcmT",
  "expires": 0,
  "description": "admin secret"
}
```

## 2. 删除密钥

### 2.1 接口描述

删除密钥。

### 2.2 请求方法

DELETE /v1/secrets/:name

### 2.3 输入参数

**Path 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（密钥名） |

### 2.4 输出参数

Null

### 2.5 请求示例

**输入示例**

```bash
curl -XDELETE -H'Content-Type: application/json' -H'Authorization: Bearer $Token' http://marmotedu.io:8080/v1/secrets/foo
```

**输出示例**

```json
null
```

## 3. 修改密钥属性

### 3.1 接口描述

修改密钥属性。

### 3.2 请求方法

PUT /v1/secrets/:name

### 3.3 输入参数

**Body 参数**

| 参数名称 | 必选 | 类型                      | 描述               |
| -------- | ---- | ------------------------- | ------------------ |
| metadata | 是   | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| expires | 否   | Int64                    | 过期时间               |
| description | 否   | String                    | 密钥描述               |

### 3.4 输出参数

| 参数名称    | 类型                                 | 描述                |
| ----------- | ------------------------------------ | ------------------- |
| metadata    | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| username    | String                               | 用户名              |
| secretID    | String                               | 密钥 ID              |
| secretKey   | String                               | 密钥 Key             |
| expires     | Int64                                | 过期时间            |
| description | String                               | 密钥描述            |

### 3.5 请求示例

**输入示例**

```bash
 curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'{
  "metadata": {
    "name": "secret"
  },
  "expires": 0,
  "description": "admin secret(modify)"
}' http://marmotedu.io:8080/v1/secrets/secret
```
**输出示例**

```json
{
  "metadata": {
    "id": 28,
    "name": "secret",
    "createdAt": "2020-09-23T11:03:43+08:00",
    "updatedAt": "2020-09-23T11:26:01.798471148+08:00"
  },
  "username": "admin",
  "secretID": "lXirSIJV5tA34V8hffffFYq7CnDhfc4gDxrz",
  "secretKey": "PK8NMhHnapVdNHAoPxhrN5Beg0C5fcmT",
  "expires": 0,
  "description": "admin secret(modify)"
}
```

## 4. 查询密钥信息

### 4.1 接口描述

查询密钥信息。

### 4.2 请求方法

GET /v1/secrets/:name

### 4.3 输入参数

**Path 参数**

| 参数名称 | 必选 | 类型   | 描述     |
| -------- | ---- | ------ | -------- |
| name | 是   | String | 资源名称（密钥名） |

### 4.4 输出参数

| 参数名称    | 类型                                 | 描述                |
| ----------- | ------------------------------------ | ------------------- |
| metadata    | [ObjectMeta](./struct.md#ObjectMeta) | REST 资源的功能属性 |
| username    | String                               | 用户名              |
| secretID    | String                               | 密钥 ID              |
| secretKey   | String                               | 密钥 Key             |
| expires     | Int64                                | 过期时间            |
| description | String                               | 密钥描述            |

### 4.5 请求示例

**输入示例**

```bash
curl -XGET -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'' http://marmotedu.io:8080/v1/secrets/secret
```

**输出示例**

```json
{
  "metadata": {
    "id": 28,
    "name": "secret",
    "createdAt": "2020-09-23T11:03:43+08:00",
    "updatedAt": "2020-09-23T11:26:02+08:00"
  },
  "username": "admin",
  "secretID": "lXirSIJV5tA34V8hffffFYq7CnDhfc4gDxrz",
  "secretKey": "PK8NMhHnapVdNHAoPxhrN5Beg0C5fcmT",
  "expires": 0,
  "description": "admin secret(modify)"
}
```

## 5. 查询密钥列表

### 5.1 接口描述

查询密钥列表。

### 5.2 请求方法

GET /v1/secrets

### 5.3 输入参数

**Query 参数**

| 参数名称      | 必选 | 类型   | 描述                                                           |
| ------------- | ---- | ------ | -------------------------------------------------------------- |
| fieldSelector | 否   | String | 字段选择器，格式为 `name=foo,phone=181`,当前只支持 name 字段过滤 |

### 5.4 输出参数

| 参数名称   | 类型     | 描述               |
| ---------- | -------- | ------------------ |
| totalCount | Uint64     | 资源总个数         |
| items      | Array of [Secret](./struct.md#Secret) | 符合条件的密钥列表 |

### 5.5 请求示例

**输入示例**

```bash
curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer $Token' -d'' http://marmotedu.io:8080/v1/secrets?offset=0&limit=10&fieldSelector=name=secret1
```

**输出示例**

```json
{
  "totalCount": 1,
  "items": [
    {
      "metadata": {
        "id": 22,
        "name": "secret1",
        "createdAt": "2020-09-20T10:09:09+08:00",
        "updatedAt": "2020-09-20T10:09:09+08:00"
      },
      "username": "admin",
      "secretID": "Uh5xpXBI5BCivVUU7kyejMvMhvRv5jcDeGYb",
      "secretKey": "D4tMymjnAKAD5w44Zf648smpK8PGw5Gf",
      "expires": 0,
      "description": "admin secret"
    }
  ]
}
```
