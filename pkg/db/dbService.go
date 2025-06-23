package dbservice

import "database/sql"

type DbServiceInterface interface {
	CreateDbConn() (*sql.DB, error)
}

type DataBaseService struct {
	DbConn DbServiceInterface
}

func (dbs *DataBaseService) RegisterDbConn() (*sql.DB, error) {
	return dbs.DbConn.CreateDbConn()
}
