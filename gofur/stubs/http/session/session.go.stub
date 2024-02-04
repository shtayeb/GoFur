package session

import (
	"database/sql"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

var SessionManager *scs.SessionManager

func InitSessionManager(db *sql.DB) {
	SessionManager = scs.New()
	SessionManager.Lifetime = 240 * time.Minute
	SessionManager.Store = postgresstore.New(db)
}
