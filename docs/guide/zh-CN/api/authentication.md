# 认证相关接口

## 1. 用户登录

### 1.1 接口描述

用户登录。

### 1.2 请求方法

POST /login

### 1.3 输入参数

**Header 参数**

| 参数名称      | 必选 | 类型   | 描述                                   |
| ------------- | ---- | ------ | -------------------------------------- |
| Authorization | 是   | String | "Basic" + <`username:password` Base64格式的字符串> |

### 1.4 输出参数

| 参数名称 | 类型   | 描述              |
| -------- | ------ | ----------------- |
| expire   | String | JWT Token过期时间 |
| token    | String | JWT Token         |

### 1.5 请求示例

**输入示例**

```bash
$ curl -XPOST -H"Authorization: Basic `echo -n 'admin:Admin@2021'|base64`" http://iam.api.marmotedu.com:8080/login
```

**输出示例**

```json
{
  "expire": "2021-10-05T14:34:07+08:00",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXBpLm1hcm1vdGVkdS5jb20iLCJleHAiOjE2MzM0MTU2NDcsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE2MjU2Mzk2NDcsInN1YiI6ImFkbWluIn0.NjQ2Q5fQ1hIEPtZBtLnOZcVywYBNWTqxysZOJSYnSGM"
}
```

## 2. 用户登出

### 2.1 接口描述

用户登出。如果通过前端登录iam-apiserver系统，可以调用此接口登出。登出操作会清理浏览器Cookie。

### 2.2 请求方法

POST /logout

### 2.3 输入参数

Null

### 2.4 输出参数

Null

### 2.5 请求示例

**输入示例**

```bash
$ curl -XPOST http://iam.api.marmotedu.com:8080/logout
```

**输出示例**

```json
null
```

## 3. 刷新Token

### 3.1 接口描述

刷新Token。

### 3.2 请求方法

POST /refresh

### 3.3 输入参数

**Header 参数**

| 参数名称      | 必选 | 类型   | 描述                                   |
| ------------- | ---- | ------ | -------------------------------------- |
| Authorization | 是   | String | "Bearer " + <JWT Token> |

### 3.4 输出参数

| 参数名称 | 类型   | 描述              |
| -------- | ------ | ----------------- |
| expire   | String | JWT Token过期时间 |
| token    | String | JWT Token         |

### 3.5 请求示例

**输入示例**

```bash
$ curl -XPOST -H"Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXBpLm1hcm1vdGVkdS5jb20iLCJleHAiOjE2MzM0MTY1MzUsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE2MjU2NDA1MzUsInN1YiI6ImFkbWluIn0.tFJC5ZO2UGy-3NI_FLiGrQF-DztRmBSDP4C5gazQYW4" http://iam.api.marmotedu.com:8080/refresh
```

**输出示例**

```json
{
  "expire": "2021-10-05T14:49:20+08:00",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpYW0uYXBpLm1hcm1vdGVkdS5jb20iLCJleHAiOjE2MzM0MTY1NjAsImlkZW50aXR5IjoiYWRtaW4iLCJpc3MiOiJpYW0tYXBpc2VydmVyIiwib3JpZ19pYXQiOjE2MjU2NDA1NjAsInN1YiI6ImFkbWluIn0.yq4jMqP338CkhM1DQWGd9K7v_9L_J4tNe76qNk3U10A"
}
```
