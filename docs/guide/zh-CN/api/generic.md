# 通用说明

IAM 系统 API 严格遵循 REST 标准进行设计，采用 JSON 格式进行数据传输，使用 JWT Token 进行 API 认证。

## 1. 公共参数

每个参数都属于不同的类型，根据参数位置不同，参数有如下类型：
- 请求头参数 (Header)：例如 `Content-Type: application/json`。
- 路径参数 (Path)：例如 `/user/:id` 中的 id 参数就位于 path 中。
- 查询参数 (Query)：例如 `users?username=colin&username=james&value=`。
- 请求体参数 (Body)：例如 `{"metadata":{"name":"secretdemo"},"expires":0,"description":"admin secret"}`。

这些公共参数，是每个接口都需要传入的，在每个接口文档中，不再一一说明。

IAM API 接口公共参数如下：

| 参数名称          | 位置     | 类型     | 必选  | 描述                                          |
|---------------|--------|--------|-----|---------------------------------------------|
| Content-Type  | Header | String | 是   | 固定值：application/json，本书示例项目统一使用 JSON 数据传输格式 |
| Authorization | Header | String | 是   | JWT Token，值以 `Bearer` 开头                    |

## 2. 返回结果

一个 API 接口调用返回结果只会有 2 种结果：成功和失败，2 种结果返回的内容会有所差异。

- 成功时，返回结果中，包含以下内容：
    1. X-Request-Id：位于 HTTP 返回的请求头中，调用的请求 ID，用来唯一标识一次请求。
    2. HTTP 状态码：HTTP 状态码，成功的请求，状态码永远为 200。
    3. 接口请求的数据：位于 HTTP 返回的 Body 中，API 请求需要的返回数据，JSON 格式。
- 失败时，返回的结果中，包含以下内容：
    1. X-Request-Id：位于 HTTP 返回的请求头中，调用的请求 ID，用来唯一标识一次请求。
    2. HTTP 状态码：HTTP 状态码，不同的错误类型返回的 HTTP 状态码不同，可能的状态码为：400、401、403、404、500。
    3. 返回的错误信息：返回格式为：`{"code":100101,"message":"Database error","reference":"https://github.com/marmotedu/iam/tree/master/docs/guide/zh-CN/faq"}`， `code` 表示错误码，`message` 表示该错误的具体信息，`reference` 表示参考文档（可选）。

成功和失败返回结果的 header 中，还有一些其它返回信息，比如：`Cache-Control`、`Content-Type`、`Access-Control-Allow-Origin` 等，这些在非排障场景下，可以不用关注，这里不再一一说明。

### 2.1 成功返回结果

成功时返回的 HTTP 状态码是 200，在 Body 中返回数据，如下是创建密钥 API 接口的返回结果：

```json
{
  "metadata": {
    "id": 24,
    "name": "secretdemo",
    "createdAt": "2020-09-20T10:17:58.108812081+08:00",
    "updatedAt": "2020-09-20T10:17:58.108812081+08:00"
  },
  "username": "admin",
  "secretID": "k5jZYMJCAk4jGH1nqgszTn6hPaZ8aZbKO0ZO",
  "secretKey": "cKdfmDJlTELfumu3SpLPf0k0SXQDqvdJ",
  "expires": 0,
  "description": "admin secret"
}
```

### 2.2 失败返回结果

失败时返回的 HTTP 状态码是 400、401、403、404、500 中的一个，以下是创建重复密钥时，API 接口返回的错误结果：

```json
{
  "code": 100101,
  "message": "Database error",
  "reference": "https://github.com/marmotedu/iam/tree/master/docs/guide/zh-CN/faq"
}
```

## 3. 返回参数类型

本书的数据传输格式为 JSON 格式，所以支持的数据类型就是 JSON 所支持的数据类型。在 JSON 中，有如下数据类型：string、number、array、boolean、null、object。JSON 中的 number 是数字类型的统称，但是在实际的 Go 项目开发中，我们需要知道更精确的 number 类型，来将 JSON 格式的数据解码（unmarshal）为 Go 的结构体类型。同时，Object 类型在 Go 中也可以直接用结构体名替代。

所以在 IAM 系统中，返回参数支持的数据类型为：
`String、Int、Uint、Int8、Uint8、Int16、Uint16、Int32、Uint32、Int64、Uint64、Float、Float64、Array、Boolean、Null、Struct`。

另外请求的编码格式均为：`UTF-8` 格式。

## 4. 认证

IAM 采用 JWT Token 进行认证，具体操作步骤如下：

1. 获取在 IAM 系统创建的 secretKey 和s ecretID
2. 通过 secretKey 和 secretID 生成 JWT Token，以下是一个可以生成 JWT Token 的 Go 源码（main.go）：

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/pflag"
)

var (
	cliAlgorithm = pflag.StringP("algorithm", "", "HS256", "Signing algorithm - possible values are HS256, HS384, HS512")
	cliTimeout   = pflag.DurationP("timeout", "", 2*time.Hour, "JWT token expires time")
	help         = pflag.BoolP("help", "h", false, "Print this help message")
)

func main() {
	pflag.Usage = func() {
		fmt.Println(`Usage: gentoken [OPTIONS] SECRETID SECRETKEY`)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	if pflag.NArg() != 2 {
		pflag.Usage()
		os.Exit(1)
	}

	token, err := createJWTToken(*cliAlgorithm, *cliTimeout, os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println(token)
}

func createJWTToken(algorithm string, timeout time.Duration, secretID, secretKey string) (string, error) {
	expire := time.Now().Add(timeout)

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
		"kid": secretID,
		"exp": expire.Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString([]byte(secretKey))
}
```

在命令行执行如下命令，即可生成 JWT Token：

```bash
$  go run main.go <secretID> <secretKey>
```

默认会生成 HS256 算法签名、2 小时后过期的 Token。可以通过 `--algorithm` 指定签名算法，通过 `--timeout` 指定 token 过期时间。

3. 携带 Token，发送 HTTP 请求：

```bash
curl -XPOST -H'Content-Type: application/json' -H'Authorization: Bearer <Token>' -d'{"metadata":{"name":"secretdemo"},"expires":0,"description":"admin secret"}' http://marmotedu.io:8080/v1/secrets
```
## 5. 请求方法

本 API 接口文档中请求方法格式为：`HTTP方法 请求路径`，例如请求方法为：`GET /v1/users`，请求地址为：`marmotedu.io`，请求协议为：`HTTP`，则实际的请求格式为：`curl -XGET http://marmotedu.io/v1/users`

## 6. 错误码

IAM 系统同时返回 2 类错误码：HTTP 状态码和业务错误码。IAM会返回 3 类 HTTP 状态码：
- 200：代表成功响应。
- 4xx：响应失败，说明客户端发生错误。
- 500：响应失败，说明服务端发生错误。

**HTTP 状态码说明：**

| 状态码 | 说明                                       |
| ------ | ------------------------------------------ |
| 200    | 成功响应                                   |
| 400    | 客户端发生错误，比如参数不合法、格式错误等 |
| 401    | 认证失败                                   |
| 403    | 授权失败                                   |
| 404    | 页面或者资源不存在                         |
| 500    | 响应失败，说明服务端发生了错误             |

**业务错误码说明**

业务错误码请参考：[错误码](./error_code_generated.md)

## 7. 其它说明

无
