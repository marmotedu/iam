## iam-authz-server

IAM Authorization Server

### Synopsis

Authorization server to run ladon policies which can protecting your resources.
It is written inspired by AWS IAM policiis.

Find more iam-authz-server information at:
    https://github.com/marmotedu/iam/blob/master/docs/guide/en-US/cmd/iam-authz-server.md,

Find more ladon information at:
    https://github.com/ory/ladon

```
iam-authz-server [flags]
```

### Options

```
      --alsologtostderr                               log to standard error as well as files
      --analytics.enable                              This sets the iam-authz-server to record analytics data. (default true)
      --analytics.enable-detailed-recording           Enable detailed analytics at the key level. (default true)
      --analytics.pool-size int                       Specify number of pool workers. (default 50)
      --analytics.records-buffer-size uint            Specifies buffer size for pool workers (size of each pipeline operation). (default 1000)
      --analytics.storage-expiration-time duration    Set to a value larger than the Pump's purge_delay. This allows the analytics data to exist long enough in Redis to be processed by the Pump. (default 24h0m0s)
      --client-ca-file string                         If set, any request presenting a client certificate signed by one of the authorities in the client-ca-file is authenticated with an identity corresponding to the CommonName of the client certificate.
  -c, --config FILE                                   Read configuration from specified FILE, support JSON, TOML, YAML, HCL, or Java properties formats.
      --feature.enable-metrics                        Enables metrics on the apiserver at /metrics (default true)
      --feature.profiling                             Enable profiling via web interface host:port/debug/pprof/ (default true)
  -h, --help                                          help for iam-authz-server
      --insecure.bind-address string                  The IP address on which to serve the --insecure.bind-port (set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces). (default "127.0.0.1")
      --insecure.bind-port int                        The port on which to serve unsecured, unauthenticated access. It is assumed that firewall rules are set up such that this port is not reachable from outside of the deployed machine and that port 443 on the iam public address is proxied to this port. This is performed by nginx in the default setup. Set to zero to disable. (default 8080)
      --log-backtrace-at traceLocation                when logging hits line file:N, emit a stack trace (default :0)
      --log-dir string                                If non-empty, write log files in this directory
      --log.development                               Development puts the logger in development mode, which changes the behavior of DPanicLevel and takes stacktraces more liberally.
      --log.disable-caller                            Disable output of caller information in the log.
      --log.disable-stacktrace                        Disable the log to record a stack trace for all messages at or above panic level.
      --log.enable-color                              Enable output ansi colors in plain format logs.
      --log.error-output-paths strings                Error output paths of log. (default [stderr])
      --log.format FORMAT                             Log output FORMAT, support plain or json format. (default "console")
      --log.level LEVEL                               Minimum log output LEVEL. (default "info")
      --log.name string                               The name of the logger.
      --log.output-paths strings                      Output paths of log. (default [stdout])
      --logtostderr                                   log to standard error instead of files
      --redis.addrs strings                           A set of redis address(format: 127.0.0.1:6379).
      --redis.database int                            By default, the database is 0. Setting the database is not supported with redis cluster. As such, if you have --redis.enable-cluster=true, then this value should be omitted or explicitly set to 0.
      --redis.enable-cluster                          If you are using Redis cluster, enable it here to enable the slots mode.
      --redis.host string                             Hostname of your Redis server. (default "127.0.0.1")
      --redis.master-name string                      The name of master redis instance.
      --redis.optimisation-max-active int             In order to not over commit connections to the Redis server, we may limit the total number of active connections to Redis. We recommend for production use to set this to around 4000. (default 4000)
      --redis.optimisation-max-idle int               This setting will configure how many connections are maintained in the pool when idle (no traffic). Set the --redis.optimisation-max-active to something large, we usually leave it at around 2000 for HA deployments. (default 2000)
      --redis.password string                         Optional auth password for Redis db.
      --redis.port int                                The port the Redis server is listening on. (default 6379)
      --redis.ssl-insecure-skip-verify                Allows usage of self-signed certificates when connecting to an encrypted Redis database.
      --redis.timeout int                             Timeout (in seconds) when connecting to redis service.
      --redis.use-ssl                                 If set, IAM will assume the connection to Redis is encrypted. (use with Redis providers that support in-transit encryption).
      --redis.username string                         Username for access to redis service.
      --rpcserver string                              The address of iam rpc server. The rpc server can provide all the secrets and policies to use. (default "127.0.0.1:8081")
      --secure.bind-address string                    The IP address on which to listen for the --secure.bind-port port. The associated interface(s) must be reachable by the rest of the engine, and by CLI/web clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces). (default "0.0.0.0")
      --secure.bind-port int                          The port on which to serve HTTPS with authentication and authorization. It cannot be switched off with 0. (default 8443)
      --secure.tls.cert-dir string                    The directory where the TLS certs are located. If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, this flag will be ignored. (default "/var/run/iam")
      --secure.tls.cert-key.cert-file string          File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated after server cert).
      --secure.tls.cert-key.private-key-file string   File containing the default x509 private key matching --secure.tls.cert-key.cert-file.
      --secure.tls.pair-name string                   The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key (default "iam")
      --server.healthz                                Add self readiness check and install /healthz router. (default true)
      --server.middlewares strings                    List of allowed middlewares for server, comma separated. If this list is empty default middlewares will be used.
      --server.mode string                            Start the server in a specified server mode. Supported server mode: debug, test, release. (default "release")
      --stderrthreshold severity                      logs at or above this threshold go to stderr (default 2)
  -v, --v Level                                       log level for V logs
      --version version[=true]                        Print version information and quit.
      --vmodule moduleSpec                            comma-separated list of pattern=N settings for file-filtered logging
```

###### Auto generated by spf13/cobra on 17-Nov-2021
