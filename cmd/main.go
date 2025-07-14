package main

import (
	"app/internal/application"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// - env

	// - config
	cfg := &application.ConfigServerChi{
		ServerAddress:        ":8080",
		EmployeesFilePath:    "docs/db/employees.json",
		BuyerLoaderFilePath:  "docs/db/buyers.json",
		WarehouseFilePath:    "docs/db/warehouse.json",
		ProductTypesFilePath: "docs/db/product_types.json",
		ProductsFilePath:     "docs/db/products.json",

		DbConf: &mysql.Config{
			User:   "root",
			Passwd: "",
			Net:    "tcp",
			Addr:   "localhost:3306",
			DBName: "fresh_db",
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
