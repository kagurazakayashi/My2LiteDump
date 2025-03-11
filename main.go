package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	nyamysql "github.com/kagurazakayashi/libNyaruko_Go/nyamysql"
	nyasqlite "github.com/kagurazakayashi/libNyaruko_Go/nyasqlite"
	"gopkg.in/yaml.v3"
)

var (
	NullStr string
	mysql   *nyamysql.NyaMySQL
	sqlite  *nyasqlite.NyaSQLite
)

type My2LiteConfig struct {
	MySQLTable string `json:"mysql_table" yaml:"mysql_table"`
}

func main() {
	log.Println("My2LiteDump v1.0.0")
	var configContent string
	var contentBytes []byte
	var my2LiteConfig My2LiteConfig
	var errInfo string = ""

	var configFile *string = flag.String("c", "", "配置文件路径,默认则用`程序名.yaml||.json`")
	flag.Parse()

	if len(strings.TrimSpace(*configFile)) == 0 {
		execPath, err := os.Executable()
		if err != nil {
			log.Println("获取执行文件路径失败。")
			log.Fatal(err)
		}

		execName := filepath.Base(execPath)
		baseName := execName[:len(execName)-len(filepath.Ext(execName))]

		dir := filepath.Dir(execPath)
		yamlPath := filepath.Join(dir, baseName+".yaml")
		jsonPath := filepath.Join(dir, baseName+".json")

		log.Println("正在加载配置文件: " + yamlPath)
		if bytes, err := os.ReadFile(yamlPath); err == nil {
			contentBytes = bytes
		} else {
			errInfo = err.Error()
		}
		if len(errInfo) > 0 {
			log.Println("正在加载配置文件: " + jsonPath)
			if bytes, err := os.ReadFile(jsonPath); err == nil {
				contentBytes = bytes
				errInfo = ""
			} else {
				errInfo += "\n" + err.Error()
			}
		}
		if len(errInfo) > 0 {
			log.Println("配置文件读取失败。")
			log.Fatal(errInfo)
		}
	} else {
		log.Println("正在加载配置文件: " + *configFile)
		if bytes, err := os.ReadFile(*configFile); err == nil {
			contentBytes = bytes
		} else {
			log.Println("配置文件读取失败：" + *configFile)
			log.Fatal(err)
		}
	}
	configContent = string(contentBytes)
	if err := json.Unmarshal(contentBytes, &my2LiteConfig); err != nil {
		errInfo = err.Error()
	}
	if len(errInfo) > 0 {
		if err := yaml.Unmarshal(contentBytes, &my2LiteConfig); err != nil {
			errInfo += err.Error()
		} else {
			errInfo = ""
		}
	}
	if len(errInfo) > 0 {
		log.Println("配置文件解析失败。")
		log.Fatal(errInfo)
	}
	log.Println("成功。")

	logger := log.New(os.Stdout, "[MySQLDump] ", log.Ldate|log.Ltime|log.Lshortfile)

	log.Println("初始化 MySQL ...")
	mysql = nyamysql.New(configContent, logger)
	if mysql.Error() != nil {
		log.Println("失败！")
		log.Fatal(mysql.Error())
	}
	log.Println("成功。")
	log.Println("初始化 SQLite ...")
	sqlite = nyasqlite.New(configContent, logger)
	if sqlite.Error() != nil {
		log.Println("失败！")
		disconnect()
		log.Fatal(sqlite.Error())
	}
	log.Println("成功。")
	disconnect()
	log.Println("正常退出。")
}

func disconnect() {
	log.Println("断开 MySQL 。")
	mysql.Close()
	log.Println("断开 SQLite 。")
	sqlite.Close()
}
