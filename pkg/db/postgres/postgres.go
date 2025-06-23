package postgres

import (
	"api-project/pkg/envvar"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

var (
	isDebug = envvar.GetBool("DEBUG", false)
)

func (rdsPgs *Rds_Postgres) CreateDbConn() (*sql.DB, error) {
	if isDebug {
		return rdsPgs.newDefaultClient()
	}
	return nil, errors.New("pending for RDS configuration in cloud")
}

func (rdsPgs *Rds_Postgres) newDefaultClient() (*sql.DB, error) {
	var (
		username = "postgres"
		password = ""
		host     = "localhost"
		dbPort   = 5432
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		host, dbPort, username, password, username,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Error().Caller().Err(err).Msg("db open connection failed")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Error().Caller().Err(err).Msg("ping db failed")
		return nil, err
	}

	return db, nil
}
