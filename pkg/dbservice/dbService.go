package dbservice

import (
	"api-project/pkg/db/postgres"
	"database/sql"
)

var (
	rds       = &postgres.Rds_Postgres{}
	DbService = &DataBaseService{
		db: rds,
	}
)

func Init() {

	db, err := DbService.FetchDbConn()
	if err != nil {
		panic("postgresDB connection initialization with error failed")
	}
	DbService.DbReader = db
	DbService.DbWriter = db
}

func (dbs *DataBaseService) FetchDbConn() (*sql.DB, error) {
	return dbs.db.CreateDbConn()
}
