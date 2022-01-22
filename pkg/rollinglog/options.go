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
	flagDisableCaller    = "log.disable-caller"
	flagFormat           = "log.format"
	flagEnableColor      = "log.enable-color"
	flagOutputPaths      = "log.output-paths"
	flagErrorOutputPaths = "log.error-output-paths"
	flagDevelopment      = "log.development"
	flagName             = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"

	flagRolling           = "log.rolling"
	flagRollingMaxSize    = "log.rolling-max-size"
	flagRollingMaxAge     = "log.rolling-max-age"
	flagRollingMaxBackups = "log.rolling-max-backups"
	flagRollingLocalTime  = "log.rolling-local-time"
	flagRollingCompress   = "log.rolling-compress"
)

type Options struct {
	OutputPaths      []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level            string   `json:"level"              mapstructure:"level"`
	Format           string   `json:"format"             mapstructure:"format"`
	DisableCaller    bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	EnableColor      bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development      bool     `json:"development"        mapstructure:"development"`
	Name             string   `json:"name"               mapstructure:"name"`

	Rolling           bool `json:"rolling" mapstructure:"rolling"`
	RollingMaxSize    int  `json:"rolling-max-size" mapstructure:"rolling-max-size"`
	RollingMaxAge     int  `json:"rolling-max-age" mapstructure:"rolling-max-age"`
	RollingMaxBackups int  `json:"rolling-max-backups" mapstructure:"rolling-max-backups"`
	RollingLocalTime  bool `json:"rolling-local-time" mapstructure:"rolling-local-time"`
	RollingCompress   bool `json:"rolling-compress" mapstructure:"rolling-compress"`
}

// NewOptions creates an Options object with default parameters.
func NewOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel.String(),
		Format:           jsonFormat,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
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
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")

	fs.BoolVar(&o.Rolling, flagRolling, o.Rolling, "Enable log rolling.")
	fs.IntVar(&o.RollingMaxAge, flagRollingMaxAge, o.RollingMaxAge, "Maximum number of days to retain old log files.")
	fs.IntVar(&o.RollingMaxSize, flagRollingMaxSize, o.RollingMaxSize, " Maximum size in megabytes of the log file.")
	fs.IntVar(&o.RollingMaxBackups, flagRollingMaxBackups, o.RollingMaxBackups, "Maximum number of old log files to retain.")
	fs.BoolVar(&o.RollingLocalTime, flagRollingLocalTime, o.RollingLocalTime, "Determines if the time used for formatting the timestamps in backup files is the computer's local time. The defaults is to use UTC time.")
	fs.BoolVar(&o.RollingCompress, flagRollingCompress, o.RollingCompress, "Determines if the rotated log files should be compressed using gzip. The default is not to perform compression.")
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
