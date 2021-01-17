package unilog

// TODO: support logfmt format

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// LogFormatJSON means json log
	LogFormatJSON = "json"
	// LogFormatConsole means console log
	LogFormatConsole = "console"

	// KeyTime is default TimeKey in zapcore.EncoderConfig
	KeyTime = "time"
	// KeyLevel is default LevelKey in zapcore.EncoderConfig
	KeyLevel = "level"
	// KeyName is default NameKey in zapcore.EncoderConfig
	KeyName = "logger"
	// KeyCaller is default CallerKey in zapcore.EncoderConfig
	KeyCaller = "caller"
	// KeyMessage is default MessageKey in zapcore.EncoderConfig
	KeyMessage = "msg"
	// KeyStackTrace is default StacktraceKey in zapcore.EncoderConfig
	KeyStackTrace = "stacktrace"

	// MaxSizeMB is default max size of log file in MB
	MaxSizeMB = 1
	// MaxBackups is default max count of backup log files
	MaxBackups = 5
	// MaxAge is default max age of log file before delete
	MaxAge = 7
)

// L return *zap.Logger
func L() *zap.Logger {
	return zap.L()
}

// LogOption is option struct
type LogOption struct {
	Format     string
	MaxSizeMB  int
	MaxBackups int
	MaxAge     int
	Level      zapcore.Level
	FileName   string
}

// DefaultLogOption return default LogOption
func DefaultLogOption() *LogOption {
	return &LogOption{
		Format:     LogFormatJSON,
		MaxSizeMB:  MaxSizeMB,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Level:      zap.DebugLevel,
		FileName:   "/dev/stdout",
	}
}

// SetMaxSizeMB set max size in MB
func (l *LogOption) SetMaxSizeMB(size int) *LogOption {
	l.MaxSizeMB = size
	return l
}

// SetMaxBackups set max backup count
func (l *LogOption) SetMaxBackups(count int) *LogOption {
	l.MaxBackups = count
	return l
}

// SetMaxAge set max age in days
func (l *LogOption) SetMaxAge(days int) *LogOption {
	l.MaxAge = days
	return l
}

// SetFormat set format name of log
func (l *LogOption) SetFormat(f string) *LogOption {
	l.Format = f
	return l
}

// SetFileName set filename of log file
func (l *LogOption) SetFileName(f string) *LogOption {
	l.FileName = f
	return l
}

// SetLevel set log level
func (l *LogOption) SetLevel(level zapcore.Level) *LogOption {
	l.Level = level
	return l
}

// SetSimpleLogger set a simple zap logger
func SetSimpleLogger(format string, filename string, level zapcore.Level) {
	option := DefaultLogOption()
	option.SetFormat(format).SetFileName(filename).SetLevel(level)
	SetLogger(option)
}

// SetLogger create default zap logger
func SetLogger(option *LogOption) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        KeyTime,
		LevelKey:       KeyLevel,
		NameKey:        KeyName,
		CallerKey:      KeyCaller,
		MessageKey:     KeyMessage,
		StacktraceKey:  KeyStackTrace,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	switch option.Format {
	case LogFormatJSON:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	stdoutSyncer := zapcore.AddSync(os.Stdout)
	syncers := []zapcore.WriteSyncer{stdoutSyncer}

	if option.FileName != "" && option.FileName != "/dev/stdout" && option.FileName != "/dev/stderr" {
		hook := lumberjack.Logger{
			Filename:   option.FileName,
			MaxSize:    MaxSizeMB,
			MaxBackups: MaxBackups,
			MaxAge:     MaxAge,
			Compress:   true,
		}
		fileSyncer := zapcore.AddSync(&hook)
		syncers = append(syncers, fileSyncer)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(syncers...),
		zap.NewAtomicLevelAt(option.Level),
	)

	caller := zap.AddCaller()
	development := zap.Development()
	logger := zap.New(core, caller, development)

	zap.ReplaceGlobals(logger)
}
