package application

import (
	"gopkg.in/macaron.v1"
)

func respondWithCode(ctx *macaron.Context, code int) {
	ctx.Resp.WriteHeader(code)
}

func respondWithJSON(ctx *macaron.Context, code int, payload interface{}) {
	ctx.JSON(code, payload)
}

func respondWithError(ctx *macaron.Context, code int, message string) {
	respondWithJSON(ctx, code, map[string]string{"error": message})
}
