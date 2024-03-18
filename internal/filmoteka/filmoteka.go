package filmoteka

import (
	"TestVK/internal/config"
	"database/sql"
	"log/slog"
)

type Filmoteka struct {
	Db     *sql.DB
	Config *config.AppConfig
	Logger *slog.Logger
}

func NewFilmoteka(db *sql.DB, appConfig *config.AppConfig, logger *slog.Logger) *Filmoteka {
	return &Filmoteka{
		Db:     db,
		Config: appConfig,
		Logger: logger,
	}
}
