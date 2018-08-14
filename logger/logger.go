package logger

import "go.uber.org/zap"

// ZapLogger logging everything
var ZapLogger, _ = zap.NewProduction()
