// Package webpages contains handlers that serve autodeb-server's web pages
package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func renderWithBase(
	r *http.Request,
	w http.ResponseWriter,
	appCtx *appctx.Context,
	user *models.User,
	template string,
	data interface{}) {

	completeData := struct {
		User *models.User
		Data interface{}
	}{
		User: user,
		Data: data,
	}

	rendered, err := appCtx.TemplatesRenderer().RenderTemplate(completeData, "base.html", template)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, rendered)
}
