package logic

import (
	"context"
	"random/internal/utils"
	"strconv"

	"random/internal/svc"
	"random/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArrRetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArrRetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArrRetLogic {
	return &ArrRetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArrRetLogic) ArrRet(req *types.ArrRetRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	seed, secret, err := utils.GetSeedAndSecret(l.ctx, l.svcCtx.Rds, req.Secret)
	if err != nil {
		return nil, err
	}

	ret := utils.RandomArrRet(req.Arr, secret, seed)

	return &types.Response{
		Result: strconv.FormatInt(ret, 10),
	}, nil
}
