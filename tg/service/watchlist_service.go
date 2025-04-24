package service

import (
	"github.com/thekhanj/csdmpro/core"
	"github.com/thekhanj/csdmpro/tg/repo"
)

type WatchlistService struct {
	PlayerRepo    *core.PlayerRepo
	WatchlistRepo *repo.WatchlistRepo
}

func (this *WatchlistService) GetTracking(chatId int64) ([]TrackingPlayer, error) {
	ids, err := this.WatchlistRepo.List(chatId)
	if err != nil {
		return nil, err
	}
	ret := []TrackingPlayer{}
	for _, id := range ids {
		p, err := this.PlayerRepo.GetPlayer(id)
		if err != nil {
			return nil, err
		}

		isOnline, err := this.PlayerRepo.IsOnline(id)
		if err != nil {
			return nil, err
		}
		ret = append(ret, TrackingPlayer{
			DbPlayer: p,
			IsOnline: isOnline,
		})
	}

	return ret, nil
}

type TrackingPlayer struct {
	DbPlayer core.DbPlayer
	IsOnline bool
}
