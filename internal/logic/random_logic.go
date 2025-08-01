package logic

import (
	"context"

	"random/internal/svc"
	"random/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RandomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRandomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RandomLogic {
	return &RandomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RandomLogic) Random(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
