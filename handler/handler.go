package handler

import (
	"net/http"

	"github.com/techiehe/goweb/context"
)

func Home(ctx *context.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"message": "Welcome Home!"})
}
