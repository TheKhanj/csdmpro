package main

import (
	"github.com/google/wire"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg"
)

type App struct {
	CoreObserver *core.Observer
	TgServer     *tg.Server
}

func ProvideApp(
	observer *core.Observer,
	tgServer *tg.Server,
) *App {
	return &App{
		CoreObserver: observer,
		TgServer:     tgServer,
	}
}

var AppModule = wire.NewSet(ProvideApp, tg.TgModule, core.CoreModule)
