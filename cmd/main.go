package main

import (
	"app/internal/application"
	"fmt"
)

func main() {
	// env
	// ...

	// app
	// - config
	cfg := &application.ConfigServerChi{
		ServerAddress:       ":8080",
		EmployeesFilePath:   "docs/db/employees.json",
		BuyerLoaderFilePath: "docs/db/buyers_10.json",
		WarehouseFilePath:   "docs/db/warehouse.json",
		SectionsFilePath:    "docs/db/sections.json",
	}
	app := application.NewServerChi(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
