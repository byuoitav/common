package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.SugaredLogger
var cfg zap.Config

func init() {

	cfg = zap.NewDevelopmentConfig()
	cfg.OutputPaths = append(cfg.OutputPaths)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, err := cfg.Build()
	if err != nil {
		log.Printf("Couldn't build config for zap logger: %v", err.Error())
		panic(err)
	}
	L = l.Sugar()
	L.Info("Zap Logger Started")

}
