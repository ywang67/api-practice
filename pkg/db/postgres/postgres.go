package postgres

import (
	dbservice "api-project/pkg/db"
	"api-project/pkg/envvar"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

var (
	isDebug             = envvar.GetBool("DEBUG", false)
	DbConnectionService *dbservice.DataBaseService
)

type Rds_Postgres struct {
	DbReader *sql.DB
	DbWriter *sql.DB
}

func Init() {
	rdsPostgres := &Rds_Postgres{}
	dbInit, err := rdsPostgres.CreateDbConn()
	if err != nil {
		panic("postgresDB connection initialization with error failed")
	}
	rdsPostgres.DbReader = dbInit
	rdsPostgres.DbWriter = dbInit
	DbConnectionService = &dbservice.DataBaseService{
		DbConn: rdsPostgres,
	}
}

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
