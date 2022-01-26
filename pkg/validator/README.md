# validator

基于 `go-playground/validator/v10` 封装的验证包，特性如下：

- 支持结构体、值验证等必备功能；
- 支持国际化功能，同时可动态切换语言；
- 错误信息可根据需要返回切片错误或 `map` 类型；
- 支持注册自定义验证函数。


# 使用方法


## 基础示例

```go
package main

import (
	"fmt"
	
	go_validator "github.com/go-playground/validator/v10"
	"github.com/marmotedu/iam/pkg/validator"
)

type Data struct {
	Name string `validate:"required,gt=0,lte=4" label:"姓名" label_en:"name"`
	Age  uint   `validate:"required,gte=0,lte=50" label:"年龄" label_en:"age"`
}

func main() {
	opts := &validator.Options{
		Language: "zh",
		Tag:      "label",
	}
	validator.Init(opts)

	d := &Data{
		Age: 66, // 不符合验证规则
	}
	
	if err := validator.Struct(d); err != nil {
		// output: [姓名为必填字段 年龄必须小于或等于50]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())
	}

	d2 := &Data{Name: "tom", Age: 10}
	if err := validator.Struct(d2); err != nil {
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())
	} else {
		fmt.Println("success") // output: success
	}
}
```

## 国际化

```go
package main

import (
	"fmt"

	go_validator "github.com/go-playground/validator/v10"
	"github.com/marmotedu/iam/pkg/validator"
)

type Data struct {
	Name string `validate:"required,gt=0,lte=4" label:"姓名" label_en:"name"`
	Age  uint   `validate:"required,gte=0,lte=50" label:"年龄" label_en:"age"`
}

func main() {
	d := &Data{
		Name: "jack",
		Age:  51,
	}
	v1 := validator.New("en", "label_en")
	if err := v1.Struct(d); err != nil {
		// output: map[age:age must be 50 or less sex:sex is a required field]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrsMap())
	}
}
```

支持动态切换语言

```go
package main

import (
	"fmt"

	go_validator "github.com/go-playground/validator/v10"
	"github.com/marmotedu/iam/pkg/validator"
)

type Data struct {
	Name string `validate:"required,gt=0,lte=4" label:"姓名" label_en:"name"`
	Age  uint   `validate:"required,gte=0,lte=50" label:"年龄" label_en:"age"`
}

func main() {
	d := &Data{
		Name: "jack",
		Age:  51,
	}
	v1 := validator.New("en", "label_en")
	if err := v1.Struct(d); err != nil {
		// output: map[age:age must be 50 or less sex:sex is a required field]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrsMap())
	}

	if err := v1.WithTranslator("zh").Struct(d); err != nil {
		// output: [age必须小于或等于50 sex为必填字段]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())
	}
}
```

## 注册自义验证函数

```go
package main

import (
	"fmt"
	"log"

	go_validator "github.com/go-playground/validator/v10"
	"github.com/marmotedu/iam/pkg/validator"
)


type Data struct {
	Age uint `validate:"required,myValidation,gte=0,lte=50" label:"年龄"`
}

func main() {
	d := &Data{Age: 66}
	v := validator.New("zh", "label")
	// register custom validation.
	if rErr := v.RegisterValidation("myValidation", "{0}真的不能为66", myValidation); rErr != nil {
		log.Fatalln(rErr.Error())
	}
	if err := v.Struct(d); err != nil {
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())    // output: [年龄真的不能为66]
		fmt.Println(err.Error())                                          // output: 年龄真的不能为66
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrsMap()) // output: map[年龄:年龄真的不能为66]
	}
}

func myValidation(fl go_validator.FieldLevel) bool {
	if fl.Field().Uint() == 66 {
		return false
	}

	return true
}
```