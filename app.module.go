package main

import (
	"context"
	"log"
	"sync"

	"github.com/google/wire"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg"
)

type App struct {
	CoreObserver *core.Observer
	TgServer     *tg.Server
	Notifier     *tg.Notifier
}

func (this *App) Start(ctx context.Context) {
	log.Println("app: started")
	defer log.Println("app: stopped")

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()

		this.TgServer.Listen(ctx)
	}()

	go func() {
		defer wg.Done()

		this.CoreObserver.Start(ctx)
	}()

	go func() {
		defer wg.Done()

		this.Notifier.Start(ctx)
	}()

	wg.Wait()
}

func ProvideApp(
	observer *core.Observer,
	tgServer *tg.Server,
	notifier *tg.Notifier,
) *App {
	return &App{
		CoreObserver: observer,
		TgServer:     tgServer,
		Notifier:     notifier,
	}
}

var AppModule = wire.NewSet(ProvideApp, tg.TgModule, core.CoreModule)
