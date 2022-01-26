package main

import (
	"fmt"
	"log"

	go_validator "github.com/go-playground/validator/v10"
	"github.com/marmotedu/iam/pkg/validator"
)

type Data struct {
	Name string `validate:"required,gt=0,lte=4" label:"姓名" label_en:"name"`
	Age  uint   `validate:"required,gte=0,lte=50" label:"年龄" label_en:"age"`
	Sex  uint   `validate:"required,oneof=1 2 3" label:"性别" label_en:"sex"`
}

type Data2 struct {
	Age uint `validate:"required,myValidation,gte=0,lte=50" label:"年龄" label_en:"age"`
}

func main() {
	basic()            // 基础示例
	I18in()            // 国际化
	customValidation() // 自定义验证函数
}

func basic() {
	d := &Data{
		Age: 66, // 不符合验证规则
		Sex: 1,
	}

	opts := &validator.Options{
		Language: "zh",
		Tag:      "label",
	}
	validator.Init(opts)

	if err := validator.Struct(d); err != nil {
		// output: [姓名为必填字段 年龄必须小于或等于50]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())
	}

	d2 := &Data{Name: "tom", Age: 10, Sex: 2}
	if err := validator.Struct(d2); err != nil {
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())
	} else {
		fmt.Println("success") // output: success
	}
}

func I18in() {
	d := &Data{
		Name: "jack",
		Age:  51,
	}
	v := validator.New("en", "label_en")
	if err := v.Struct(d); err != nil {
		// output: map[age:age must be 50 or less sex:sex is a required field]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrsMap())
	}
	// 支持动态切换语言
	if err := v.WithTranslator("zh").Struct(d); err != nil {
		// output: [age必须小于或等于50 sex为必填字段]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())
	}
}

func customValidation() {
	d := &Data2{Age: 66}
	v := validator.New("zh", "label")
	// register custom validation.
	if rErr := v.RegisterValidation("myValidation", "{0}真的不能为66", myValidation); rErr != nil {
		log.Fatalln(rErr.Error())
	}
	if err := v.Struct(d); err != nil {
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrs())    // output: [年龄真的不能为66]
		fmt.Println(err.(*validator.ValidationErrors).TranslateErrsMap()) // output: map[年龄:年龄真的不能为66]
		fmt.Println(err.Error())                                          // output: 年龄真的不能为66
	}

}

func myValidation(fl go_validator.FieldLevel) bool {
	if fl.Field().Uint() == 66 {
		return false
	}

	return true
}
