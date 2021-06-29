module github.com/marmotedu/iam

go 1.16

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/buger/jsonparser v1.1.1
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0
	github.com/dgraph-io/ristretto v0.0.3
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v7 v7.4.0
	github.com/golang/mock v1.5.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/gosuri/uitable v0.0.4
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/influxdata/influxdb v1.8.4
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.11.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.9.0 // indirect
	github.com/likexian/host-stat-go v0.0.0-20190516151207-c9cf36dd6ce9
	github.com/marmotedu/api v1.0.1
	github.com/marmotedu/component-base v1.0.0
	github.com/marmotedu/errors v1.0.2
	github.com/marmotedu/marmotedu-sdk-go v1.0.2-0.20210528170801-2c91b80cb4cf
	github.com/mattn/go-isatty v0.0.13
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/novalagung/gubrak v1.0.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/olivere/elastic/v7 v7.0.22
	github.com/onsi/ginkgo v1.16.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/ory/ladon v1.2.0
	github.com/parnurzeal/gorequest v0.2.16
	github.com/prometheus/client_golang v1.10.0
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/russross/blackfriday v1.6.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.12
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/tpkeeper/gin-dump v1.0.0
	github.com/vmihailenco/msgpack v3.3.3+incompatible
	github.com/zsais/go-gin-prometheus v0.1.0
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	golang.org/x/tools v0.1.4
	google.golang.org/grpc v1.38.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.0.5
	gorm.io/gorm v1.21.4
	k8s.io/klog v1.0.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.27.0
