package config

import (
	"os"
	"strconv"
)

func GetApplicationPort() int {
	var port string
	if port = os.Getenv("APPLICATION_PORT"); port == "" {
		return 50053
	}

	portInt, _ := strconv.Atoi(port)
	return portInt
}

func GetRedisAddr() string {
	var port string
	if port = os.Getenv("REDIS_ADDR"); port == "" {
		return "redis:6379"
	}

	return port
}

func GetDatabaseDSN() string {
	var dsn string
	if dsn = os.Getenv("DATA_SOURCE_URL"); dsn == "" {
		return "root:pa55word@tcp(mysql:3306)/customers?charset=utf8&parseTime=True&loc=Local"
	}

	return dsn
}