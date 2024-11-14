package logger

import "go.uber.org/zap"

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func Initialize(l *zap.Logger) {
	logger = l
	logger.Sync()
	sugar = l.Sugar()
	sugar.Sync()
}

func Logger() *zap.Logger {
	if logger != nil {
		return logger
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}

func Sugar() *zap.SugaredLogger {
	if sugar != nil {
		return sugar
	}
	if logger != nil {
		sugar = logger.Sugar()
		return sugar
	}
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar = l.Sugar()
	return sugar
}
