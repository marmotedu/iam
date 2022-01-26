package validator

import "github.com/spf13/pflag"

const (
	flagLanguage = "validator.language" // 验证器使用的语言 zh, en...
	flagTag      = "validator.tag"
)

type Options struct {
	Language string `json:"language" mapstructure:"language"`
	// Struct field tag
	Tag string `json:"tag" mapstructure:"tag"`
}

func NewOptions() *Options {
	return &Options{
		Language: "zh",
		Tag:      "label",
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Language, flagLanguage, o.Language, "Validator use language.")
	fs.StringVar(&o.Tag, flagTag, o.Tag, "Validator struct field tag.")
}
