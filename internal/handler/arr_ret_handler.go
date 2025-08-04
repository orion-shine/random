package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"random/internal/logic"
	"random/internal/svc"
	"random/internal/types"
)

func ArrRetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArrRetRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewArrRetLogic(r.Context(), svcCtx)
		resp, err := l.ArrRet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
