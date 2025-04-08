package tg

import (
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/thekhanj/csdmpro/tg/middlewares"
	"github.com/thekhanj/csdmpro/tg/repo"
	"github.com/thekhanj/tgool"
	"golang.org/x/net/proxy"
)

type ServerBuilder struct {
	err         error
	http_client *http.Client
	token       string
	controllers []tgool.Controller
	bilakhRepo  *repo.BilakhRepo
}

// WithProxy uses a socks5 proxy for connecting to telegram'a api.
//
//	Example: socks5://127.0.0.1:9050
func (this *ServerBuilder) WithProxy(socks_proxy string) *ServerBuilder {
	p := strings.ReplaceAll(socks_proxy, "socks5://", "")

	dialer, err := proxy.SOCKS5("tcp", p, nil, proxy.Direct)
	if err != nil {
		this.err = err

		return this
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	this.http_client = &http.Client{
		Transport: transport,
	}

	return this
}

func (this *ServerBuilder) WithToken(token string) *ServerBuilder {
	this.token = token

	return this
}

func (this *ServerBuilder) WithControllers(controllers ...tgool.Controller) *ServerBuilder {
	this.controllers = controllers
	return this
}

func (this *ServerBuilder) WithBilakhRepo(repo *repo.BilakhRepo) *ServerBuilder {
	this.bilakhRepo = repo
	return this
}

func (this *ServerBuilder) Build() (*Server, error) {
	if this.err != nil {
		return nil, this.err
	}

	if this.http_client == nil {
		this.http_client = &http.Client{}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(
		this.token, tgbotapi.APIEndpoint, this.http_client,
	)
	if err != nil {
		return nil, err
	}

	ms := make([]tgool.Middleware, 0)

	if this.bilakhRepo != nil {
		ms = append(ms, middlewares.NewBilakhMiddleware(this.bilakhRepo))
	}

	if this.controllers != nil {
		m := tgool.NewControllerMiddleware(this.controllers...)
		ms = append(ms, m)
	}

	router := tgool.NewRouter(ms...)

	return &Server{bot, router}, nil
}
