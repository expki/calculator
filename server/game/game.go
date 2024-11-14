package game

import (
	"context"
	"sync"

	"github.com/expki/calculator/lib/schema"
	"github.com/google/uuid"
)

type Game struct {
	appCtx     context.Context
	stateLock  sync.RWMutex
	state      schema.Global
	clientLock sync.RWMutex
	clientMap  map[uuid.UUID]*Session
}

func New(appCtx context.Context) *Game {
	game := &Game{
		appCtx:    appCtx,
		clientMap: make(map[uuid.UUID]*Session),
		state: schema.Global{
			Calculator: schema.Calculator{},
			Members:    make(map[string]*schema.Member),
		},
	}
	go game.gameLoop()
	return game
}
