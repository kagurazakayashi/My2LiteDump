package main

import (
	"github.com/kagurazakayashi/libNyaruko_Go/nyamysql"
)

// ExportConfig 結構體，封裝匯出配置
type ExportConfig struct {
	MySQLDSN    string // MySQL 連線字串
	SQLiteFile  string // SQLite 資料庫檔案
	MySQLTable  string // MySQL 表名
	SQLiteTable string // SQLite 目標表名
	KeyColumn   string // 用於分頁查詢的列（通常是主鍵）
	StartValue  int64  // 查詢起始值
	EndValue    int64  // 查詢結束值
	BatchSize   int64  // 每批次匯出多少條資料
}

// type MySQLDBConfig struct {
// 	User     string `json:"mysql_user"`
// 	Password string `json:"mysql_pwd"`
// 	Address  string `json:"mysql_addr"`
// 	Port     string `json:"mysql_port"`
// 	DbName   string `json:"mysql_db"`
// 	Limit    string `json:"mysql_limit"`
// }

func BatchExportMySQLToSQLite(config ExportConfig, mysqlConfig nyamysql.MySQLDBConfig) error {
	return linkMySQL(mysqlConfig)
}
