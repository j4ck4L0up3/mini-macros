package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strings"
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

func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fname := r.FormValue("first-name")
	lname := r.FormValue("last-name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	rePassword := r.FormValue("reenter-password")

	if strings.Compare(password, rePassword) != 0 {
		w.WriteHeader(http.StatusTeapot)
		c := templates.PasswordMatchError()
		c.Render(r.Context(), w)
		return
	}

	user, _ := h.userStore.GetUser(email)
	if user != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.EmailInUseError()
		c.Render(r.Context(), w)
		return
	}

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
