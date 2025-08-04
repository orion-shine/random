package logic

import (
	"context"
	"random/internal/utils"

	"random/internal/svc"
	"random/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArrLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArrLogic {
	return &ArrLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArrLogic) Arr(req *types.ArrRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	seed, secret, err := utils.GetSeedAndSecret(l.ctx, l.svcCtx.Rds, req.Secret)
	if err != nil {
		return nil, err
	}

	sort := utils.RandomSort(req.Arr, seed, secret)

	return &types.Response{
		Result: utils.Int64SliceToString(sort),
	}, nil
}
