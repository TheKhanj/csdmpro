package core

import (
	"errors"
	"sort"
	"sync"

	rd "github.com/Pallinder/go-randomdata"
)

var ERR_FAKE_PLAYER_NOT_FOUND error = errors.New("fake player not found")

type internalPlayer struct {
	isOnline bool
	player   Player
}

type StubCrawler struct {
	mutex   sync.Mutex
	players map[string]*internalPlayer
}

func (this *StubCrawler) MakeOnline(name string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	p, ok := this.players[name]
	if !ok {
		return ERR_FAKE_PLAYER_NOT_FOUND
	}

	p.isOnline = true

	return nil
}

func (this *StubCrawler) MakeOffline(name string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	p, ok := this.players[name]
	if !ok {
		return ERR_FAKE_PLAYER_NOT_FOUND
	}

	p.isOnline = false

	return nil
}

func (this *StubCrawler) RemovePlayer(name string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.players[name]; !ok {
		return ERR_FAKE_PLAYER_NOT_FOUND
	}

	delete(this.players, name)

	this.rerankPlayers()

	return nil
}

func (this *StubCrawler) AddPlayer() string {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var rank int

	p := Player{
		Name:     rd.SillyName(),
		Country:  rd.Country(rd.TwoCharCountry),
		Rank:     &rank,
		Score:    rd.Number(50000),
		Kills:    rd.Number(50000),
		Deaths:   rd.Number(50000),
		Accuracy: rd.Number(100),
	}

	ip := &internalPlayer{
		isOnline: false,
		player:   p,
	}

	this.players[p.Name] = ip

	this.rerankPlayers()

	return p.Name
}

func (this *StubCrawler) Stats(page int) ([]Player, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	ret := make([]Player, 0)

	for _, p := range this.players {
		ret = append(ret, p.player)
	}

	sort.Slice(ret, func(i, j int) bool {
		return *ret[i].Rank < *ret[j].Rank
	})

	l := (page - 1) * 50
	r := page * 50

	if r > len(ret) {
		return ret[:], nil
	}

	return ret[l:r], nil
}

func (this *StubCrawler) Online() ([]Player, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	ret := make([]Player, 0)

	for _, p := range this.players {
		if p.isOnline {
			ret = append(ret, p.player)
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return *ret[i].Rank < *ret[j].Rank
	})

	return ret, nil
}

func (this *StubCrawler) rerankPlayers() {
	rank := 1

	for _, p := range this.players {
		*p.player.Rank = rank
		rank++
	}
}

var _ Crawler = (*StubCrawler)(nil)

func NewStubCrawler() *StubCrawler {
	return &StubCrawler{
		players: make(map[string]*internalPlayer),
	}
}
