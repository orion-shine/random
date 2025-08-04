package logic

import (
	"context"
	"random/internal/svc"
	"random/internal/types"
	"random/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type PukeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPukeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PukeLogic {
	return &PukeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 13 * 花色
// 黑桃 (♠)：0–12
// 红桃 (♥)：13–25
// 梅花 (♣)：26–38
// 方块 (♦)：39–51
// 小王：52
// 大王：53

// Puke
func (l *PukeLogic) Puke(req *types.PukeRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	seed, secret, err := utils.GetSeedAndSecret(l.ctx, l.svcCtx.Rds, req.Secret)
	if err != nil {
		return nil, err
	}

	pukes := initPukesWithoutJokers()
	if req.Type == "with" {
		pukes = utils.RandomSort(pukes, secret, seed)
	} else if req.Type == "without" {
		pukes = append(pukes, 53, 54)
		pukes = utils.RandomSort(pukes, secret, seed)
	}

	pukeStr := utils.Int64SliceToString(pukes)

	return &types.Response{
		Result: pukeStr,
	}, nil
}

func initPukesWithoutJokers() []int64 {
	pukes := make([]int64, 0, 52)
	for i := int64(0); i < 52; i++ {
		pukes = append(pukes, i)
	}
	return pukes
}
