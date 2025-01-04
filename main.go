package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"

	nyamysql "github.com/kagurazakayashi/libNyaruko_Go/nyamysql"
	"golang.org/x/term"
)

var NullStr string

func main() {
	// 定義命令列引數: MySQL
	var host *string = flag.String("h", "127.0.0.1", "MySQL 服务器地址")
	var port *int = flag.Int("P", 3306, "MySQL 端口号")
	var dbname *string = flag.String("D", "", "MySQL 数据库名称")
	var mytable *string = flag.String("T", "", "MySQL 表名")
	var user *string = flag.String("u", "root", "MySQL 用户名")
	var password *string = flag.String("p", "", "MySQL 密码")
	var exKey *string = flag.String("K", "id", "要导出内容的 key")
	var exVal *string = flag.String("V", "", "要导出的内容区间（整数-整数）")
	// 定義命令列引數: SQLite
	var dbfile *string = flag.String("f", "", "SQLite 数据库文件路径")
	var litetable *string = flag.String("t", "", "SQLite 表名")
	flag.Parse()

	// _, err := os.Stat(*dbfile)
	// if os.IsNotExist(err) {
	// 	log.Fatalf("SQLite 数据库文件 %s 不存在!\n", *dbfile)
	// 	return
	// }

	// 如果 -p 選項提供但沒有指定密碼，則要求使用者輸入密碼（隱藏輸入）
	if len(*password) == 0 {
		fmt.Print("请输入 MySQL 密码 (不回显): ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("读取密码失败: %v", err)
		}
		fmt.Println()
		*password = string(bytePassword)
	}

	// 構造 MySQL 連線字串
	var dsn string
	if *dbname != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *user, *password, *host, *port, *dbname)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/", *user, *password, *host, *port)
	}
	fmt.Println("正在连接到数据库:", MaskPassword(dsn))

	var vall []string = strings.Split(*exVal, "-")

	startValue, err := strconv.ParseInt(vall[0], 10, 64)
	if err != nil {
		log.Fatalf("无效的 StartValue: %v", err)
	}
	endValue, err := strconv.ParseInt(vall[1], 10, 64)
	if err != nil {
		log.Fatalf("无效的 EndValue: %v", err)
	}
	config := ExportConfig{
		MySQLDSN:    dsn,
		SQLiteFile:  *dbfile,
		MySQLTable:  *mytable,
		SQLiteTable: *litetable,
		KeyColumn:   *exKey,
		StartValue:  startValue,
		EndValue:    endValue,
		BatchSize:   50,
	}
	mysqlConfig := nyamysql.MySQLDBConfig{
		User:     *user,
		Password: *password,
		Address:  *host,
		Port:     strconv.Itoa(*port),
		DbName:   *dbname,
		MaxLimit: "10000",
	}
	if err := BatchExportMySQLToSQLite(config, mysqlConfig); err != nil {
		return
	}
}
