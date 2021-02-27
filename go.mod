module github.com/marmotedu/iam

go 1.15

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/buger/jsonparser v1.1.1
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/cpuguy83/go-md2man/v2 v2.0.0
	github.com/dgraph-io/ristretto v0.0.3
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/fatih/color v1.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/mock v1.4.4
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/influxdata/influxdb v1.8.3
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.11.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.9.0 // indirect
	github.com/likexian/host-stat-go v0.0.0-20190516151207-c9cf36dd6ce9
	github.com/magefile/mage v1.11.0 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/marmotedu/api v1.0.0
	github.com/marmotedu/component-base v1.0.0
	github.com/marmotedu/errors v1.0.0
	github.com/marmotedu/marmotedu-sdk-go v1.0.0
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635
	github.com/novalagung/gubrak v1.0.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/olivere/elastic/v7 v7.0.22
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/ory/ladon v1.2.0
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/prometheus/client_golang v1.9.0
	github.com/russross/blackfriday v1.6.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.8
	github.com/sirupsen/logrus v1.8.0
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tpkeeper/gin-dump v1.0.0
	github.com/ugorji/go v1.1.13 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210226181700-f36f78243c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	golang.org/x/tools v0.1.0
	google.golang.org/grpc v1.35.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.11
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	k8s.io/klog v1.0.0
)

replace (
	github.com/golang/mock => /home/colin/workspace/golang/src/github.com/golang/mock
	github.com/marmotedu/api => /home/colin/workspace/golang/src/github.com/marmotedu/api
	github.com/marmotedu/component-base => /home/colin/workspace/golang/src/github.com/marmotedu/component-base
	github.com/marmotedu/errors => /home/colin/workspace/golang/src/github.com/marmotedu/errors
	github.com/marmotedu/marmotedu-sdk-go => /home/colin/workspace/golang/src/github.com/marmotedu/marmotedu-sdk-go
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
