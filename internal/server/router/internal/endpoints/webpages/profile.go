package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//ProfileGetHandler returns a handler that renders the profile page
func ProfileGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		pgpKeys, err := app.GetUserPGPKeys(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			User                    *models.User
			PGPKeys                 []*models.PGPKey
			ExpectedPGPKeyProofText string
		}{
			User:                    user,
			PGPKeys:                 pgpKeys,
			ExpectedPGPKeyProofText: app.ExpectedPGPKeyProofText(user.ID),
		}

		rendered, err := app.TemplatesRenderer().RenderTemplate(data, "base.html", "profile.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler := auth.WithUser(handlerFunc, app)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

//AddPGPKeyPostHandler returns a handler that adds PGP key to the user's profile
func AddPGPKeyPostHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		key := r.Form.Get("key")
		proof := r.Form.Get("proof")

		if err := app.AddUserKey(user.ID, key, proof); err != nil {
			app.Logger().Error(err)
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

	handler := auth.WithUser(handlerFunc, app)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
