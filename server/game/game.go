package game

import (
	"fmt"

	"github.com/expki/calculator/lib/schema"
	"go.uber.org/zap"
)

type Game struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func New(logger *zap.Logger) *Game {
	var state schema.State
	fmt.Println("game state:", state)
	return &Game{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}
