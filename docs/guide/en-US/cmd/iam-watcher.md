## iam-watcher

IAM watcher server

### Synopsis

IAM Watcher is a pluggable watcher service used to do some periodic work like cron job. 
But the difference with cron job is iam-watcher also support sleep some duration after previous job done.

Find more iam-pump information at:
    https://github.com/marmotedu/iam/blob/master/docs/guide/en-US/cmd/iam-watcher.md

```
iam-watcher [flags]
```

### Options

```
      --alsologtostderr                           log to standard error as well as files
  -c, --config FILE                               Read configuration from specified FILE, support JSON, TOML, YAML, HCL, or Java properties formats.
      --health-check-address string               Specifies liveness health check bind address. (default "0.0.0.0:6060")
      --health-check-path string                  Specifies liveness health check request path. (default "healthz")
  -h, --help                                      help for iam-watcher
      --log-backtrace-at traceLocation            when logging hits line file:N, emit a stack trace (default :0)
      --log-dir string                            If non-empty, write log files in this directory
      --log.development                           Development puts the logger in development mode, which changes the behavior of DPanicLevel and takes stacktraces more liberally.
      --log.disable-caller                        Disable output of caller information in the log.
      --log.disable-stacktrace                    Disable the log to record a stack trace for all messages at or above panic level.
      --log.enable-color                          Enable output ansi colors in plain format logs.
      --log.error-output-paths strings            Error output paths of log. (default [stderr])
      --log.format FORMAT                         Log output FORMAT, support plain or json format. (default "console")
      --log.level LEVEL                           Minimum log output LEVEL. (default "info")
      --log.name string                           The name of the logger.
      --log.output-paths strings                  Output paths of log. (default [stdout])
      --logtostderr                               log to standard error instead of files
      --mysql.database string                     Database name for the server to use.
      --mysql.host string                         MySQL service host address. If left blank, the following related mysql options will be ignored. (default "127.0.0.1:3306")
      --mysql.log-mode int                        Specify gorm log level. (default 1)
      --mysql.max-connection-life-time duration   Maximum connection life time allowed to connecto to mysql. (default 10s)
      --mysql.max-idle-connections int            Maximum idle connections allowed to connect to mysql. (default 100)
      --mysql.max-open-connections int            Maximum open connections allowed to connect to mysql. (default 100)
      --mysql.password string                     Password for access to mysql, should be used pair with password.
      --mysql.username string                     Username for access to mysql service.
      --redis.addrs strings                       A set of redis address(format: 127.0.0.1:6379).
      --redis.database int                        By default, the database is 0. Setting the database is not supported with redis cluster. As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.
      --redis.enable-cluster                      If you are using Redis cluster, enable it here to enable the slots mode.
      --redis.host string                         Hostname of your Redis server. (default "127.0.0.1")
      --redis.master-name string                  The name of master redis instance.
      --redis.optimisation-max-active int         In order to not over commit connections to the Redis server, we may limit the total number of active connections to Redis. We recommend for production use to set this to around 4000. (default 4000)
      --redis.optimisation-max-idle int           This setting will configure how many connections are maintained in the pool when idle (no traffic). Set the --redis.optimisation-max-active to something large, we usually leave it at around 2000 for HA deployments. (default 2000)
      --redis.password string                     Optional auth password for Redis db.
      --redis.port int                            The port the Redis server is listening on. (default 6379)
      --redis.ssl-insecure-skip-verify            Allows usage of self-signed certificates when connecting to an encrypted Redis database.
      --redis.timeout int                         Timeout (in seconds) when connecting to redis service.
      --redis.use-ssl                             If set, IAM will assume the connection to Redis is encrypted. (use with Redis providers that support in-transit encryption).
      --redis.username string                     Username for access to redis service.
      --stderrthreshold severity                  logs at or above this threshold go to stderr (default 2)
  -v, --v Level                                   log level for V logs
      --version version[=true]                    Print version information and quit.
      --vmodule moduleSpec                        comma-separated list of pattern=N settings for file-filtered logging
      --watcher.counter.max-reserve-days int      Policy audit log maximum retention days. (default 180)
      --watcher.task.max-inactive-days int        Maximum user inactivity time. Otherwise the account will be disabled.
```

###### Auto generated by spf13/cobra on 17-Nov-2021
