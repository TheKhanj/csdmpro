//go:build wireinject

package main

import (
	"github.com/google/wire"
)

func WireBuild() *App {
	wire.Build(AppModule)

	return &App{}
}
