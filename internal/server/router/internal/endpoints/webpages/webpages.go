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

	// Retrieve the flashes and save the session
	session, _ := appCtx.Sessions().Get(r)
	flashes := session.Flashes()
	session.Save(r, w)

	completeData := struct {
		User    *models.User
		Flashes map[string][]string
		Data    interface{}
	}{
		Flashes: flashes,
		User:    user,
		Data:    data,
	}

	rendered, err := appCtx.TemplatesRenderer().RenderTemplate(completeData, "base.html", template)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, rendered)
}
