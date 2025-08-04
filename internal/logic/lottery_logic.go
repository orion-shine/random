package logic

import (
	"context"
	"random/internal/utils"

	"random/internal/svc"
	"random/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LotteryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLotteryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LotteryLogic {
	return &LotteryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LotteryLogic) Lottery(req *types.LotteryRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	seed, secret, err := utils.GetSeedAndSecret(l.ctx, l.svcCtx.Rds, req.Secret)
	if err != nil {
		return nil, err
	}

	var result []int64
	if req.Type == "pc28" {
		result = utils.RandomPC28(secret, seed)
	} else if req.Type == "markSix" {
		result = utils.RandomMarkSix(secret, seed)
	}

	return &types.Response{
		Result: utils.Int64SliceToString(result),
	}, nil
}
