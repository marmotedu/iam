package rollinglog

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel            = "log.level"
	flagEnableCaller     = "log.enable_caller"
	flagFormat           = "log.format"
	flagEnableColor      = "log.enable_color"
	flagOutputPaths      = "log.output_paths"
	flagErrorOutputPaths = "log.error_output_paths"
	flagName             = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"

	flagRolling           = "log.Rolling"
	flagRollingMaxSize    = "log.Rolling_max_size"
	flagRollingMaxAge     = "log.Rolling_max_age"
	flagRollingMaxBackups = "log.Rolling_max_backups"
	flagRollingLocalTime  = "log.Rolling_local_time"
	flagRollingCompress   = "log.Rolling_compress"
)

type Options struct {
	OutputPaths      []string `json:"output_paths"       mapstructure:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths" mapstructure:"error_output_paths"`
	Level            string   `json:"level"              mapstructure:"level"`
	Format           string   `json:"format"             mapstructure:"format"`
	EnableCaller     bool     `json:"enable_caller" mapstructure:"enable_caller"`
	EnableColor      bool     `json:"enable_color"       mapstructure:"enable_color"`
	Name             string   `json:"name"               mapstructure:"name"`

	Rolling           bool `json:"rolling" mapstructure:"rolling"`
	RollingMaxSize    int  `json:"rolling_max_size" mapstructure:"rolling_max_size"`
	RollingMaxAge     int  `json:"rolling_max_age" mapstructure:"rolling_max_age"`
	RollingMaxBackups int  `json:"rolling_max_backups" mapstructure:"rolling_max_backups"`
	RollingLocalTime  bool `json:"rolling_local_time" mapstructure:"rolling_local_time"`
	RollingCompress   bool `json:"rolling_compress" mapstructure:"rolling_compress"`
}

// NewOptions creates an Options object with default parameters.
func NewOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel.String(),
		EnableCaller:     false,
		Format:           consoleFormat,
		EnableColor:      false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Rolling:          true,
		RollingLocalTime: true,
	}
}

// Validate the options fields.
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds flags for log to the specified FlagSet object.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.EnableCaller, flagEnableCaller, o.EnableCaller, "Enable output of caller information in the log.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")

	fs.BoolVar(&o.Rolling, flagRolling, o.Rolling, "Enable log rolling.")
	fs.IntVar(&o.RollingMaxAge, flagRollingMaxAge, o.RollingMaxAge, "Maximum number of days to retain old log files.")
	fs.IntVar(&o.RollingMaxSize, flagRollingMaxSize, o.RollingMaxSize, " Maximum size in megabytes of the log file.")
	fs.IntVar(&o.RollingMaxBackups, flagRollingMaxBackups, o.RollingMaxBackups, "Maximum number of old log files to retain.")
	fs.BoolVar(&o.Rolling, flagRollingLocalTime, o.Rolling, "Determines if the time used for formatting the timestamps in backup files is the computer's local time. The defaults is to use UTC time.")
	fs.BoolVar(&o.RollingCompress, flagRollingCompress, o.RollingCompress, "Determines if the rotated log files should be compressed using gzip. The default is not to perform compression.")

}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
