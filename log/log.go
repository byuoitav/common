package log

import (
	"log"

	"github.com/byuoitav/common/nerr"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.SugaredLogger
var cfg zap.Config
var atom *zap.AtomicLevel

func init() {

	atom = zap.NewAtomicLevel()

	cfg = zap.ProductionConfig()
	cfg.OutputPaths = append(cfg.OutputPaths)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = atom
	l, err := cfg.Build()
	if err != nil {
		log.Printf("Couldn't build config for zap logger: %v", err.Error())
		panic(err)
	}

	//we're gonna default to info for now
	atom.SetLevel(zapcore.InfoLevel)

	L = l.Sugar()
	L.Info("Zap Logger Started")

}

func SetLevel(level string) *nerr.E {
	switch level {
	case "debug":
		atom.SetLevel(zapcore.DebugLevel)
	case "info":
		atom.SetLevel(zapcore.InfoLevel)
	case "warn":
		atom.SetLevel(zapcore.WarnLevel)
	case "error":
		atom.SetLevel(zapcore.ErrorLevel)
	case "panic":
		atom.SetLevel(zapcore.PanicLevel)
	case "fatal":
		atom.SetLevel(zapcore.FatalLevel)
	default:
		return nerr.Create("Invalid level", "invalid_args")
	}

	return nil
}
