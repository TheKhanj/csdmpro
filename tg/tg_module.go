package tg

import (
	"log"
	"os"

	"github.com/google/wire"
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/db"
	"github.com/thekhanj/csdmpro/tg/controllers"
	"github.com/thekhanj/csdmpro/tg/repo"
	"github.com/thekhanj/csdmpro/tg/service"
	"github.com/thekhanj/tgool"
)

type TgControllers []tgool.Controller

func ProvideWatchlistRepo(db db.Database) *repo.WatchlistRepo {
	repo, err := repo.CreateWatchlistRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func ProvideBilakhRepo(db db.Database) *repo.BilakhRepo {
	repo, err := repo.CreateBilakhRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	return repo
}

func ProvideWatchlistService(
	playerRepo *core.PlayerRepo,
	watchlistRepo *repo.WatchlistRepo,
) *service.WatchlistService {
	return &service.WatchlistService{
		PlayerRepo:    playerRepo,
		WatchlistRepo: watchlistRepo,
	}
}

func ProvideControllers(
	playerRepo *core.PlayerRepo,
	watchlistRepo *repo.WatchlistRepo,
	service *service.WatchlistService,
) TgControllers {
	start := &controllers.StartController{}
	watchlist := &controllers.WatchlistController{
		Service:       service,
		PlayerRepo:    playerRepo,
		WatchlistRepo: watchlistRepo,
	}
	stats := &controllers.StatsController{PlayerRepo: playerRepo}
	onlines := &controllers.OnlinesController{PlayerRepo: playerRepo}

	return TgControllers{
		start,
		watchlist,
		stats,
		onlines,
	}
}

func ProvideTg(
	controllers TgControllers,
	bilakhRepo *repo.BilakhRepo,
) *Server {
	serverBuilder := ServerBuilder{}

	serverBuilder.
		WithToken(os.Getenv("API_TOKEN")).
		WithControllers(controllers...).
		WithBilakhRepo(bilakhRepo)

	socks_proxy := os.Getenv("http_proxy")
	if socks_proxy != "" {
		serverBuilder.WithProxy(socks_proxy)
	}

	s, err := serverBuilder.Build()
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func ProvideNotifier(
	observer *core.Observer,
	watchlistRepo *repo.WatchlistRepo,
	playerRepo *core.PlayerRepo,
	server *Server,
) *Notifier {
	return &Notifier{
		observer:      observer,
		watchlistRepo: watchlistRepo,
		playerRepo:    playerRepo,
		bot:           server.bot,
	}
}

var TgModule = wire.NewSet(
	ProvideTg, ProvideControllers,
	ProvideWatchlistRepo, ProvideBilakhRepo,
	ProvideWatchlistService, ProvideNotifier,
)
