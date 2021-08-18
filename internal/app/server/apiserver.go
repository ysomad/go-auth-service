package server

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/ysomad/go-auth-service/internal/app/store/sqlstore"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return nil
	}
	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, s)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
