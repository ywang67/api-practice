package dbservice

import "database/sql"

type DbServiceInterface interface {
	CreateDbConn() (*sql.DB, error)
}

type DataBaseService struct {
	db       DbServiceInterface
	DbReader *sql.DB
	DbWriter *sql.DB
}
