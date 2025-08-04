package logic

import (
	"context"
	"random/internal/utils"
	"strconv"

	"random/internal/svc"
	"random/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FloatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFloatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FloatLogic {
	return &FloatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FloatLogic) Float(req *types.FloatRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	seed, secret, err := utils.GetSeedAndSecret(l.ctx, l.svcCtx.Rds, req.Secret)
	if err != nil {
		return nil, err
	}

	float := utils.RandomFloat(secret, seed)
	return &types.Response{
		Result: strconv.FormatFloat(float, 'f', 2, 64),
	}, nil
}
