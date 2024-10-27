package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/DavidHuie/gomigrate"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mfsyahrz/image_feed_api/internal/config"
)

type PostgreDB struct {
	Conn *sqlx.DB
}

func NewPostgres(cfg config.Postgres) (PostgreDB, error) {
	db := PostgreDB{}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Name)
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	db.Conn.SetConnMaxIdleTime(time.Minute * time.Duration(cfg.MaxIdleLifetime))
	db.Conn.SetMaxOpenConns(cfg.MaxOpenConns)
	db.Conn.SetConnMaxLifetime(time.Minute * time.Duration(cfg.MaxConnLifetime))

	return db, startMigration(db.Conn.DB)
}

func startMigration(db *sql.DB) error {
	migrator, err := gomigrate.NewMigrator(db, gomigrate.Postgres{}, "./migration/postgres")
	if err != nil {
		return err
	}

	err = migrator.Migrate()
	if err != nil {
		_ = migrator.Rollback()
		return err
	}
	return nil
}
