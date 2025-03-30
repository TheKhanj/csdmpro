//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/thekhanj/csdmpro/tg"
)

func WireBuild() *tg.Server {
	wire.Build(tg.TgModule)

	return &tg.Server{}
}
