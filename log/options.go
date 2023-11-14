package log

/**
* created by mengqi on 2023/11/14
* 这里使用函数选项模式
 */

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	FormatConsole = "console"
	FormatJson    = "json"
)

type Options struct {
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	Level             string   `json:"level"              mapstructure:"level"`
	Format            string   `json:"format"             mapstructure:"format"`
	Name              string   `json:"name"               mapstructure:"name"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool     `json:"development"        mapstructure:"development"`
	EnableTraceID     bool     `json:"enable-trace-id"    mapstructure:"enable-trace-id"`    //是否开启traceID
	EnableTraceStack  bool     `json:"enable-trace-stack" mapstructure:"enable-trace-stack"` //是否开启traceStack
}

func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            FormatConsole,
		EnableColor:       false,
		Development:       false,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != FormatConsole && format != FormatJson {
		errs = append(errs, fmt.Errorf("invalid log format: %q", o.Format))
	}

	return errs
}

func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == FormatConsole && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
	}
	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

type Option func(l *Logger)

func WithMinLevel(lvl zapcore.Level) Option {
	return func(l *Logger) {
		l.minLevel = lvl
	}
}

func WithErrorStatusLevel(lvl zapcore.Level) Option {
	//默认zap.ErrorLevel.
	return func(l *Logger) {
		l.errorStatusLevel = lvl
	}
}

func WithCaller(on bool) Option {
	return func(l *Logger) {
		l.caller = on
	}
}

func WithStackTrace(on bool) Option {
	return func(l *Logger) {
		l.stackTrace = on
	}
}

func WithTraceIDField(on bool) Option {
	return func(l *Logger) {
		l.withTraceID = on
	}
}
