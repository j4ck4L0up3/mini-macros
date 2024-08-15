package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)

type PostRegisterHandler struct {
	userStore store.UserStore
}

// TODO: update RegisterParams to include SessionStore and SessionCookieName
type PostRegisterHandlerParams struct {
	UserStore store.UserStore
}

func NewPostRegisterHandler(params PostRegisterHandlerParams) *PostRegisterHandler {
	return &PostRegisterHandler{
		userStore: params.UserStore,
	}
}

// TODO: setup cookie and redirect to dashboard after registering
func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("first-name")
	lname := r.FormValue("last-name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.userStore.CreateUser(fname, lname, email, password)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}
