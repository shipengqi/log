package log

import (
	"encoding/json"
	"errors"

	"github.com/spf13/pflag"
)

// Options Configuration for logging.
type Options struct {
	// DisableConsole whether to log to console
	DisableConsole bool `json:"disable-console" mapstructure:"disable-console"`
	// DisableConsoleColor force disabling colors.
	DisableConsoleColor bool `json:"disable-console-color" mapstructure:"disable-console-color"`
	// DisableConsoleTime whether to add a time
	DisableConsoleTime bool `json:"disable-console-time" mapstructure:"disable-console-time"`
	// DisableConsoleLevel whether to add a level
	DisableConsoleLevel bool `json:"disable-console-level" mapstructure:"disable-console-level"`
	// DisableConsoleCaller whether to log caller info
	DisableConsoleCaller bool `json:"disable-console-caller" mapstructure:"disable-console-caller"`

	// DisableFile whether to log to file
	DisableFile bool `json:"disable-file" mapstructure:"disable-file"`
	// DisableFileJson whether to enable json format for log file
	DisableFileJson bool `json:"disable-file-json" mapstructure:"disable-file-json"`
	// DisableFileTime whether to add a time
	DisableFileTime bool `json:"disable-file-time" mapstructure:"disable-file-time"`
	// DisableFileCaller whether to log caller info
	DisableFileCaller bool `json:"disable-file-caller" mapstructure:"disable-file-caller"`

	// DisableRotate whether to enable log file rotate
	DisableRotate bool `json:"disable-rotate" mapstructure:"disable-rotate"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `json:"max-size" mapstructure:"max-size"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `json:"max-backups" mapstructure:"max-backups"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `json:"max-age" mapstructure:"max-age"`

	// CallerSkip the max age in days to keep a logfile
	CallerSkip int `json:"caller-skip" mapstructure:"caller-skip"`

	// ConsoleLevel sets the standard logger level
	ConsoleLevel string `json:"console-level" mapstructure:"console-level"`
	// FileLevel sets the file logger level.
	FileLevel string `json:"file-level" mapstructure:"file-level"`

	// Output directory for logging when DisableFile is false
	Output string `json:"output" mapstructure:"output"`
}

// NewOptions creates an Options with default parameters.
func NewOptions() *Options {
	return &Options{
		DisableFile:          true,
		DisableRotate:        true,
		DisableConsoleCaller: true,
		DisableFileCaller:    true,
		ConsoleLevel:         InfoLevel.String(),
		CallerSkip:           DefaultCallerSkip,
	}
}

// AddFlags adds flags related to logger to the specified FlagSet.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConsoleLevel, "log.console-level", o.ConsoleLevel,
		"Sets the standard logger level.")

	fs.StringVar(&o.FileLevel, "log.file-level", o.FileLevel,
		"Sets the file logger level.")

	fs.BoolVar(&o.DisableConsole, "log.disable-console", o.DisableConsole,
		"Whether to log to console.")

	fs.BoolVar(&o.DisableConsoleColor, "log.disable-console-color", o.DisableConsoleColor,
		"Force disabling colors.")

	fs.BoolVar(&o.DisableConsoleTime, "log.disable-console-time", o.DisableConsoleTime,
		"Whether to add a time.")

	fs.BoolVar(&o.DisableConsoleLevel, "log.disable-console-time", o.DisableConsoleLevel,
		"Whether to add a level.")

	fs.BoolVar(&o.DisableConsoleCaller, "log.disable-console-time", o.DisableConsoleCaller,
		"Whether to add caller info.")

	fs.BoolVar(&o.DisableFile, "log.disable-file", o.DisableFile,
		"Whether to log to file.")

	fs.BoolVar(&o.DisableFileJson, "log.disable-file-json", o.DisableFileJson,
		"Whether to enable json format for log file.")

	fs.BoolVar(&o.DisableFileTime, "log.disable-file-time", o.DisableFileTime,
		"Whether to add a time.")

	fs.BoolVar(&o.DisableFileCaller, "log.disable-file-caller", o.DisableFileCaller,
		"Whether to add caller info.")

	fs.BoolVar(&o.DisableRotate, "log.disable-file-caller", o.DisableRotate,
		"Whether to enable log file rotate.")

	fs.IntVar(&o.MaxSize, "log.max-size", o.MaxSize,
		"Sets the max size in MB of the logfile before it's rolled.")

	fs.IntVar(&o.MaxBackups, "log.max-backups", o.MaxBackups,
		"Sets the max number of rolled files to keep.")

	fs.IntVar(&o.MaxAge, "log.max-backups", o.MaxAge,
		"Sets the max age in days to keep a logfile.")

	fs.StringVar(&o.Output, "log.output", o.Output,
		"Sets the directory for logging when DisableFile is false.")
}

// Validate validates the options fields.
func (o *Options) Validate() []error {
	var errs []error
	var level Level

	if o.ConsoleLevel != "" {
		if err := level.UnmarshalText([]byte(o.ConsoleLevel)); err != nil {
			errs = append(errs, err)
		}
	}

	if o.FileLevel != "" {
		if err := level.UnmarshalText([]byte(o.FileLevel)); err != nil {
			errs = append(errs, err)
		}
	}

	if o.DisableConsole && o.DisableFile {
		errs = append(errs, errors.New("no enabled logger, one or more of "+
			"(DisableConsole, DisableFile) must be set to false"))
	}

	if !o.DisableFile && o.Output == "" {
		errs = append(errs, errors.New("no log output, 'Output' must be set"))
	}
	return errs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
