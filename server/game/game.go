package game

import (
	"context"
	"sync"

	"github.com/expki/calculator/lib/schema"
)

type Game struct {
	appCtx     context.Context
	stateLock  sync.RWMutex
	state      schema.StateState
	clientLock sync.RWMutex
	clientMap  map[int]*Session
}

func New(appCtx context.Context) *Game {
	game := &Game{
		appCtx:    appCtx,
		clientMap: make(map[int]*Session),
		state: schema.StateState{
			Calculator: schema.Calculator{},
			Members:    make([]schema.MemberState, 0),
		},
	}
	go game.gameLoop()
	return game
}
