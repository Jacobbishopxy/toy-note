// logger package origin from
// https://gist.github.com/rnyrnyrny/a6dc926ae11951b753ecd66c00695397#file-logger-go-L8
package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	// DPanic, Panic and Fatal level can not be set by user
	DebugLevelStr   string = "debug"
	InfoLevelStr    string = "info"
	WarningLevelStr string = "warning"
	ErrorLevelStr   string = "error"
)

type ToyNoteLogger struct {
	globalLogger *zap.Logger
	devMode      bool
}

var TNLogger *ToyNoteLogger

// Init logger
func Init(logLevel string, logFile string, dev bool) error {
	var level zapcore.Level
	switch logLevel {
	case DebugLevelStr:
		level = zap.DebugLevel
	case InfoLevelStr:
		level = zap.InfoLevel
	case WarningLevelStr:
		level = zap.WarnLevel
	case ErrorLevelStr:
		level = zap.ErrorLevel
	default:
		return fmt.Errorf("unknown log level %s", logLevel)
	}

	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, //MB
		MaxBackups: 30,
		MaxAge:     30, //days
		Compress:   false,
	})

	// encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		// use NewConsoleEncoder for human readable output
		zapcore.NewJSONEncoder(encoderConfig),
		// write to stdout as well as log files
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), ws),
		// shorten of create AtomicLevel and set level
		zap.NewAtomicLevelAt(level),
	)

	var _globalLogger *zap.Logger
	if dev {
		_globalLogger = zap.New(core, zap.AddCaller(), zap.Development())
	} else {
		_globalLogger = zap.New(core)
	}
	zap.ReplaceGlobals(_globalLogger)

	TNLogger = &ToyNoteLogger{
		globalLogger: _globalLogger,
		devMode:      dev,
	}

	return nil
}

// Call it in defer
func (l *ToyNoteLogger) Sync() error {
	return l.globalLogger.Sync()
}

// Each package can have its own logger
func (l *ToyNoteLogger) NewSugar(name string) *zap.SugaredLogger {
	return l.globalLogger.Named(name).Sugar()
}
