package main

import (
	"app/internal/application"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// - config
	cfg := &application.ConfigServerChi{
		ServerAddress: application.ServerAddress,
		DbConf: &mysql.Config{
			User:   application.DbUser,
			Passwd: application.DbPassword,
			Net:    application.DbProtocol,
			Addr:   application.DbAddress,
			DBName: application.DbName,
		},
	}

	// - app
	app := application.NewServerChi(cfg)

	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
