package main

import (
	"context"
	"errors"
	"goth/internal/config"
	"goth/internal/handlers"
	"goth/internal/hash/passwordhash"
	database "goth/internal/store/db"
	"goth/internal/store/dbstore"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	m "goth/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
* Set to production at build time
* used to determine what assets to load
 */
var Environment string

func init() {
	Environment = os.Getenv("ENV")
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r := chi.NewRouter()

	cfg := config.MustLoadConfig()

	dsn := cfg.Dsn

	db := database.MustOpen(dsn)
	passwordhash := passwordhash.NewHPasswordHash()

	userStore := dbstore.NewUserStore(
		dbstore.NewUserStoreParams{
			DB:           db,
			PasswordHash: passwordhash,
		},
	)

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
			DB: db,
		},
	)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	authMiddleware := m.NewAuthMiddleware(userStore, sessionStore, cfg.SessionCookieName)

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Logger,
			m.TextHTMLMiddleware,
			m.CSPMiddleware,
			authMiddleware.AddUserToContext,
		)

		r.NotFound(handlers.NewNotFoundHandler().ServeHTTP)

		r.Get("/", handlers.NewHomeHandler().ServeHTTP)

		r.Get("/about", handlers.NewAboutHandler().ServeHTTP)

		r.Get("/register", handlers.NewGetRegisterHandler().ServeHTTP)

		r.Post("/register", handlers.NewPostRegisterHandler(handlers.PostRegisterHandlerParams{
			UserStore: userStore,
		}).ServeHTTP)

		r.Get("/login", handlers.NewGetLoginHandler().ServeHTTP)

		r.Post("/login", handlers.NewPostLoginHandler(handlers.PostLoginHandlerParams{
			UserStore:         userStore,
			SessionStore:      sessionStore,
			PasswordHash:      passwordhash,
			SessionCookieName: cfg.SessionCookieName,
		}).ServeHTTP)

		r.Post("/logout", handlers.NewPostLogoutHandler(handlers.PostLogoutHandlerParams{
			SessionCookieName: cfg.SessionCookieName,
			SessionStore:      sessionStore,
			UserStore:         userStore,
		}).ServeHTTP)

		r.Get("/account", handlers.NewGetAccountHandler(handlers.GetAccountHandlerParams{
			SessionStore:      sessionStore,
			SessionCookieName: cfg.SessionCookieName,
		}).ServeHTTP)

		r.Put(
			"/account/first-name",
			handlers.NewPutFirstNameHandler(handlers.PutFirstNameHandlerParams{
				SessionStore:      sessionStore,
				UserStore:         userStore,
				SessionCookieName: cfg.SessionCookieName,
			}).ServeHTTP,
		)

		r.Put(
			"/account/last-name",
			handlers.NewPutLastNameHandler(handlers.PutLastNameHandlerParams{
				SessionStore:      sessionStore,
				UserStore:         userStore,
				SessionCookieName: cfg.SessionCookieName,
			}).ServeHTTP,
		)

		r.Put(
			"/account/email",
			handlers.NewPutEmailHandler(handlers.PutEmailHandlerParams{
				SessionStore:      sessionStore,
				UserStore:         userStore,
				SessionCookieName: cfg.SessionCookieName,
			}).ServeHTTP,
		)

		r.Put(
			"/account/password",
			handlers.NewPutPasswordHandler(handlers.PutPasswordHandlerParams{
				SessionStore:      sessionStore,
				UserStore:         userStore,
				Passwordhash:      passwordhash,
				SessionCookieName: cfg.SessionCookieName,
			}).ServeHTTP,
		)

		r.Delete(
			"/account/delete-account",
			handlers.NewDeleteAccountHandler(handlers.DeleteAccountHandlerParams{
				UserStore:         userStore,
				SessionStore:      sessionStore,
				SessionCookieName: cfg.SessionCookieName,
			}).ServeHTTP,
		)

		// FIXME: edit as completed
		r.Get("/dashboard", handlers.NewGetDashboardHandler(handlers.GetDashboardHandlerParams{
			SessionCookieName: cfg.SessionCookieName,
			SessionStore:      sessionStore,
		}).ServeHTTP)
	})

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info(
		"Server started",
		slog.String("server address", cfg.ServerAddr),
		slog.String("env", Environment),
	)
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
