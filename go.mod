module github.com/marmotedu/iam

go 1.15

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/appleboy/gin-jwt/v2 v2.6.4
	github.com/arnaud-deprez/gsemver v0.5.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/buger/jsonparser v1.0.0
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/cpuguy83/go-md2man/v2 v2.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-contrib/requestid v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v7 v7.4.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/influxdata/influxdb v1.8.3
	github.com/jinzhu/gorm v1.9.16
	github.com/jinzhu/now v1.1.1
	github.com/json-iterator/go v1.1.10
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/likexian/host-stat-go v0.0.0-20190516151207-c9cf36dd6ce9
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/marmotedu/api v1.0.0
	github.com/marmotedu/component-base v1.0.1
	github.com/marmotedu/errors v1.0.0
	github.com/marmotedu/marmotedu-sdk-go v1.0.1
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/mapstructure v1.3.3
	github.com/moby/term v0.0.0-20200915141129-7f0af18e79f2
	github.com/novalagung/gubrak v1.0.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/olivere/elastic/v7 v7.0.21
	github.com/ory/ladon v1.2.0
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/prometheus/client_golang v1.8.0
	github.com/russross/blackfriday v1.5.2
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.8
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tpkeeper/gin-dump v1.0.0
	github.com/zsais/go-gin-prometheus v0.1.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20201027140754-0fcbb8f4928c // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	golang.org/x/tools v0.0.0-20201028025901-8cd080b735b3
	google.golang.org/grpc v1.33.1
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.5
	k8s.io/klog v1.0.0
)

replace (
	github.com/marmotedu/api => /home/colin/workspace/golang/src/github.com/marmotedu/api
	github.com/marmotedu/component-base => /home/colin/workspace/golang/src/github.com/marmotedu/component-base
  google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
