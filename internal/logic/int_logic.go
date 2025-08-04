package logic

import (
	"context"
	"random/internal/utils"
	"strconv"

	"random/internal/svc"
	"random/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IntLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIntLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IntLogic {
	return &IntLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IntLogic) Int(req *types.IntRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	seed, secret, err := utils.GetSeedAndSecret(l.ctx, l.svcCtx.Rds, req.Secret)
	if err != nil {
		return nil, err
	}

	random := utils.Random(req.End, seed, secret)

	return &types.Response{
		Result: strconv.FormatInt(random, 10),
	}, nil
}
