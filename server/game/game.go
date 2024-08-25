package game

import (
	"github.com/expki/calculator/lib/schema"
	"go.uber.org/zap"
)

type Game struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	state  schema.Global
}

func New(logger *zap.Logger) *Game {
	return &Game{
		logger: logger,
		sugar:  logger.Sugar(),
		state: schema.Global{
			Calculator: schema.Calculator{},
			Members:    make(map[string]*schema.Member),
		},
	}
}
