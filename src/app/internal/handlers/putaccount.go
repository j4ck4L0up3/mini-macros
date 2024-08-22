package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"goth/internal/hash"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strings"
)

type PutFirstNameHandler struct {
	sessionStore      store.SessionStore
	userStore         store.UserStore
	sessionCookieName string
}

type PutLastNameHandler struct {
	sessionStore      store.SessionStore
	userStore         store.UserStore
	sessionCookieName string
}

type PutEmailHandler struct {
	sessionStore      store.SessionStore
	userStore         store.UserStore
	sessionCookieName string
}

type PutPasswordHandler struct {
	sessionStore      store.SessionStore
	userStore         store.UserStore
	passwordhash      hash.PasswordHash
	sessionCookieName string
}

type PutFirstNameHandlerParams struct {
	SessionStore      store.SessionStore
	UserStore         store.UserStore
	SessionCookieName string
}

type PutLastNameHandlerParams struct {
	SessionStore      store.SessionStore
	UserStore         store.UserStore
	SessionCookieName string
}

type PutEmailHandlerParams struct {
	SessionStore      store.SessionStore
	UserStore         store.UserStore
	SessionCookieName string
}

type PutPasswordHandlerParams struct {
	SessionStore      store.SessionStore
	UserStore         store.UserStore
	Passwordhash      hash.PasswordHash
	SessionCookieName string
}

func NewPutFirstNameHandler(params PutFirstNameHandlerParams) *PutFirstNameHandler {
	return &PutFirstNameHandler{
		sessionStore:      params.SessionStore,
		userStore:         params.UserStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func NewPutLastNameHandler(params PutLastNameHandlerParams) *PutLastNameHandler {
	return &PutLastNameHandler{
		sessionStore:      params.SessionStore,
		userStore:         params.UserStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func NewPutEmailHandler(params PutEmailHandlerParams) *PutEmailHandler {
	return &PutEmailHandler{
		sessionStore:      params.SessionStore,
		userStore:         params.UserStore,
		sessionCookieName: params.SessionCookieName,
	}
}

func NewPutPasswordHandler(params PutPasswordHandlerParams) *PutPasswordHandler {
	return &PutPasswordHandler{
		sessionStore:      params.SessionStore,
		userStore:         params.UserStore,
		passwordhash:      params.Passwordhash,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PutFirstNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	newFirstName1 := r.FormValue("new-first-name")
	newFirstName2 := r.FormValue("reenter-new-first-name")

	if strings.Compare(newFirstName1, newFirstName2) != 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.FirstNameMatchError()
		c.Render(r.Context(), w)
		return
	}

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.FirstNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve sessionCookieName: %v\n", err)
		return
	}

	sessionBytes, err := b64.RawStdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.FirstNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not decode cookie value into bytes: %v\n", err)
		return
	}

	sessionID := string(sessionBytes)
	user, err := h.sessionStore.GetUserFromSession(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.FirstNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve user from session: %v\n", err)
		return
	}

	err = h.userStore.UpdateUserFirstName(user.ID, newFirstName1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.FirstNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not update first name in database: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	c := templates.FirstNameChangeSuccess()
	c.Render(r.Context(), w)
	w.Header().Set("HX-Redirect", "/account")
}

func (h *PutLastNameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	newLastName1 := r.FormValue("new-last-name")
	newLastName2 := r.FormValue("reenter-new-last-name")

	if strings.Compare(newLastName1, newLastName2) != 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.LastNameMatchError()
		c.Render(r.Context(), w)
		return
	}

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.LastNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve sessionCookieName: %v\n", err)
		return
	}

	sessionBytes, err := b64.RawStdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.LastNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not decode cookie value into bytes: %v\n", err)
		return
	}

	sessionID := string(sessionBytes)

	user, err := h.sessionStore.GetUserFromSession(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.LastNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve user from session: %v\n", err)
	}

	err = h.userStore.UpdateUserLastName(user.ID, newLastName1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.LastNameUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not update last name in database: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	c := templates.LastNameChangeSuccess()
	c.Render(r.Context(), w)
	w.Header().Set("HX-Redirect", "/account")
}

func (h *PutEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	newEmail1 := r.FormValue("new-email")
	newEmail2 := r.FormValue("reenter-new-email")

	if strings.Compare(newEmail1, newEmail2) != 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.EmailMatchError()
		c.Render(r.Context(), w)
		return
	}

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.EmailUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve sessionCookieName: %v\n", err)
		return
	}

	sessionBytes, err := b64.RawStdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.EmailUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not decode cookie value into bytes: %v\n", err)
		return
	}

	sessionID := string(sessionBytes)

	user, err := h.sessionStore.GetUserFromSession(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.EmailUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve user from session: %v\n", err)
	}

	err = h.userStore.UpdateUserEmail(user.ID, newEmail1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.EmailUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not update email in database: %v\n", err)
		return
	}

	// TODO: send verification email, set a verify email button if email not verified

	w.WriteHeader(http.StatusAccepted)
	c := templates.EmailChangeSuccess()
	c.Render(r.Context(), w)
	w.Header().Set("HX-Redirect", "/account")
}

func (h *PutPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	currPassword := r.FormValue("curr-password")
	newPassword1 := r.FormValue("new-password")
	newPassword2 := r.FormValue("reenter-new-password")

	if strings.Compare(newPassword1, newPassword2) != 0 {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.PasswordMatchError()
		c.Render(r.Context(), w)
		return
	}

	currCookie, err := r.Cookie(h.sessionCookieName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.PasswordUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve sessionCookieName: %v\n", err)
		return
	}

	sessionBytes, err := b64.RawStdEncoding.DecodeString(currCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.PasswordUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not decode cookie value into bytes: %v\n", err)
		return
	}

	sessionID := string(sessionBytes)

	sessionUser, err := h.sessionStore.GetUserFromSession(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.PasswordUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve user from session: %v\n", err)
		return
	}

	user, err := h.userStore.GetUser(sessionUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.PasswordUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not retrieve user: %v\n", err)
		return
	}

	fmt.Printf("User: %v\n", user)
	passwordIsValid, err := h.passwordhash.ComparePasswordAndHash(currPassword, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.PasswordUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not succesfully validate password: %v\n", err)
		return
	}

	if !passwordIsValid {
		w.WriteHeader(http.StatusNotAcceptable)
		c := templates.CurrentPasswordError()
		c.Render(r.Context(), w)
		return
	}

	err = h.userStore.UpdateUserPassword(user.ID, newPassword1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.PasswordUpdateError()
		c.Render(r.Context(), w)
		fmt.Printf("Could not update password in database: %v\n")
	}

	w.WriteHeader(http.StatusAccepted)
	c := templates.PasswordChangeSuccess()
	c.Render(r.Context(), w)
	w.Header().Set("HX-Redirect", "/account")
}
