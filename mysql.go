package main

import (
	nyamysql "github.com/kagurazakayashi/libNyaruko_Go/nyamysql"
)

func linkMySQL(mysqlConfig nyamysql.MySQLDBConfig) error {
	var pool *nyamysql.MySQLPool = nyamysql.NewPoolC(mysqlConfig, 1)
	mysqlLinkID, err := pool.MysqlIsRun(true)
	if err != nil {
		println("MySQL Link failed:", err.Error())
		return err
	}
	pool.MysqlClose(mysqlLinkID, true)
	return nil
}
