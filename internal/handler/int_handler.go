package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"random/internal/logic"
	"random/internal/svc"
	"random/internal/types"
)

func IntHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IntRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewIntLogic(r.Context(), svcCtx)
		resp, err := l.Int(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
