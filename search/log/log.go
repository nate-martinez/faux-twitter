package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	atomic zap.AtomicLevel
	Log    *zap.SugaredLogger
)

func init() {
	atomic = zap.NewAtomicLevel()
	Log = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atomic,
	)).Sugar()
	Log = Log.WithOptions(zap.AddStacktrace(zapcore.WarnLevel))
}

func SetLogLevel(level string) {
	var lvl zapcore.Level
	switch level {
	case "Error":
		lvl = zap.ErrorLevel
	case "Warn":
		lvl = zap.WarnLevel
	case "Info":
		lvl = zap.InfoLevel
	case "Debug":
		lvl = zap.DebugLevel
	default:
		Log.Infof("unrecognized level '%s'", level)
		return
	}
	atomic.SetLevel(lvl)
	Log.Infof("set log level to %s", level)
}
