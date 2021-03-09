# marmotedu/log

`marmotedu/log` 是一个生产可用的日志包，基于 `zap` 包封装。除了实现 `Go` 日志包的基本功能外，还实现了很多高级功能，`marmotedu/log`具有如下特性：

- 支持日志级别：`Debug`、`Info`、`Warn`、`Error`、`Panic`、`Fatal`。
- 支持自定义配置。
- 支持文件名和行号。
- 支持输出掉标准输出和文件，可以同时输出到多个地方。
- 支持 `JSON` 和 `Text` 两种日志格式。
- 支持颜色输出。
- 兼容标准的 `log` 包。
- 高性能。
- 支持结构化日志记录。
- **兼容标准库 `log` 包和 `glog`**。
- **支持Context（业务定制）**

## 使用方法


### 一个简单的示例

创建一个 `example2.go` 文件，内容如下：

```go
package main

import (
    "github.com/marmotedu/iam/pkg/log"
)

func main() {
    defer log.Flush()

    // Debug、Info(with field)、Warnf、Errorw使用
    log.Debug("This is a debug message")
    log.Info("This is a info message", log.Int32("int_key", 10))
    log.Warnf("This is a formatted %s message", "warn")
}
```

执行代码：

```bash
$ go run example2.go 
2020-12-05 07:56:37.154	info	example/example2.go:12	This is a info message	{"int_key": 10}
2020-12-05 07:56:37.154	warn	example/example2.go:13	This is a formatted warn message
```

上述代码使用 `marmotedu/log` 包默认的全局 `logger`，分别在 `Debug` 、`Info` 和 `Warn` 级别打印了一条日志。

### 初始化日志包

可以使用 `Init` 来初始化一个日志包，如下：

```go
// logger配置    
opts := &log.Options{
    Level:            "debug",
    Format:           "console",
    EnableColor:      true,
    EnableCaller:     true,
    OutputPaths:      []string{"test.log", "stdout"},
    ErrorOutputPaths: []string{"error.log"},
}
// 初始化全局logger    
log.Init(opts)
```

Format 支持 `console` 和 `json` 2 种格式：
- console：输出为 text 格式。例如：`2020-12-05 08:12:02.324	DEBUG	example/example.go:43	This is a debug message`
- json：输出为 json 格式，例如：`{"level":"debug","time":"2020-12-05 08:12:54.113","caller":"example/example.go:43","msg":"This is a debug message"}`

OutputPaths，可以设置日志输出：
- stdout：输出到标准输出。
- stderr：输出到标准错误输出。
- /var/log/test.log：输出到文件。

支持同时输出到多个输出。

EnableColor 为 `true` 开启颜色输出，为 `false` 关闭颜色输出。

### 结构化日志输出

`marmotedu/log` 也支持结构化日志打印，例如：

```go
log.Info("This is a info message", log.Int32("int_key", 10))
log.Infow("Message printed with Errorw", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508") 
```
对应的输出结果为：

```
2020-12-05 08:16:18.749	INFO	example/example.go:44	This is a info message	{"int_key": 10}
2020-12-05 08:16:18.749	ERROR	example/example.go:46	Message printed with Errorw	{"X-Request-ID": "fbf54504-64da-4088-9b86-67824a7fb508"}
```

log.Info 这类函数需要指定具体的类型，以最大化的 提高日志的性能。log.Infow 这类函数，不用指定具体的类型，底层使用了反射，性能会差些。建议用在低频调用的函数中。

## 支持V level

创建 `v_level.go`，内容如下：

```go
package main

import (
    "github.com/marmotedu/iam/pkg/log"
)

func main() {
    defer log.Flush()

    log.V(0).Info("This is a V level message")
    log.V(0).Infow("This is a V level message with fields", "X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")
}
```

执行如上代码：

```bash
$ go run v_level.go 
2020-12-05 08:20:37.763	info	example/v_level.go:10	This is a V level message
2020-12-05 08:20:37.763	info	example/v_level.go:11	This is a V level message with fields	{"X-Request-ID": "7a7b9f24-4cae-4b2a-9464-69088b45b904"}
```

## 完整的示例

一个完整的示例请参考[example.go](./example/example.go)。
