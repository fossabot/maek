package auth

import (
	"net/http"

	"github.com/karngyan/maek/routers/base"
)

func Info(ctx *base.WebContext) {
	base.Respond(ctx, map[string]any{
		"user":       ctx.User,
		"workspaces": ctx.User.Workspaces,
	}, http.StatusOK)
}
