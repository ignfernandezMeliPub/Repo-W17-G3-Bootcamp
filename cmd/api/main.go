package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/melisource/fury_go-platform/pkg/fury"

	"app/internal"
	"app/internal/handler"
	"app/internal/logger"
	"app/internal/repository/repositories/buyer_repository"
	"app/internal/repository/repositories/carries_repository"
	"app/internal/repository/repositories/employee_repository"
	"app/internal/repository/repositories/inbound_order_repository"
	"app/internal/repository/repositories/locality_repository"
	"app/internal/repository/repositories/product_batch_repository"
	"app/internal/repository/repositories/product_record_repository"
	"app/internal/repository/repositories/product_repository"
	"app/internal/repository/repositories/product_type_repository"
	"app/internal/repository/repositories/purchase_order_repository"
	"app/internal/repository/repositories/sections_repository"
	"app/internal/repository/repositories/seller_repository"
	"app/internal/repository/repositories/warehouse_repository"
	"app/internal/service"
)

func main() {
	// - app
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app, err := fury.NewWebApplication()
	if err != nil {
		return err
	}

	// Obtener configuración
	config := internal.GetConfig()

	// Log de configuración para debugging
	log.Printf("Configuración de BD: User=%s, Addr=%s, DBName=%s", config.DB.User, config.DB.Address, config.DB.Name)

	dbConf := &mysql.Config{
		User:                 config.DB.User,
		Passwd:               config.DB.Password,
		Net:                  config.DB.Protocol,
		Addr:                 config.DB.Address,
		DBName:               config.DB.Name,
		AllowNativePasswords: config.DB.AllowNativePasswords,
	}

	db, err := sql.Open(config.DB.Driver, dbConf.FormatDSN())
	if err != nil {
		return err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	logger.SetLogDb(db)
	logger.SetLogLevel(logger.LogLevelDebug)

	// Seller
	sellerRepo := seller_repository.NewSellerRepositorySql(db)
	sellerService := service.NewSellerServiceImpl(&sellerRepo)
	sellerHandler := handler.NewSellerHandler(&sellerService)

	// Locality
	localityRepo := locality_repository.NewLocalityRepositorySql(db)
	localityService := service.NewLocalityServiceImpl(&localityRepo)
	localityHandler := handler.NewLocalityHandler(&localityService)

	// Carries
	carriesRepo := carries_repository.NewCarriesSql(db)
	carriesService := service.NewCarriesServiceDefault(carriesRepo)
	carriesHandler := handler.NewCarriesHandler(carriesService)

	// Buyer
	buyerRp := buyer_repository.NewBuyerSQL(db)
	buyerSv := service.NewBuyerDefault(buyerRp)
	buyerHd := handler.NewBuyerDefault(buyerSv)

	// PurchaseOrders
	purchaseOrderRp := purchase_order_repository.NewPurchaseOrderRepositorySQL(db)
	purchaseOrderSv := service.NewPurchaseOrderDefault(purchaseOrderRp)
	purchaseOrderHd := handler.NewPurchaseOrderDefault(purchaseOrderSv)

	// Product
	productRpSQL := product_repository.NewProductRepositoryMySQL(db)
	productTypeRpSQL := product_type_repository.NewProductTypeRepositoryMySQL(db)
	productTypeSv := service.NewProductTypeService(productTypeRpSQL)
	productSv := service.NewProductService(productRpSQL, productTypeSv, &sellerService)
	productHd := handler.NewProductController(&productSv)

	// warehouse
	warehouseRp := warehouse_repository.NewWarehouseSql(db)
	warehouseSv := service.NewWarehouseDefault(warehouseRp)
	warehouseHd := handler.NewWarehouseDefault(warehouseSv)

	// sections
	sectionsRp := sections_repository.NewSectionsRepositorySQL(db)
	sectionsSv := service.NewSectionsService(sectionsRp)
	sectionsHd := handler.NewSectionsController(sectionsSv)

	// product-batch
	productBatchRP := product_batch_repository.NewProductBatchRepositorySQL(db)
	productBatchSv := service.NewProductBatchService(productBatchRP)
	productBatchHd := handler.NewProductBatchController(productBatchSv)

	// Employee
	rpEmployee := employee_repository.NewEmployeeDb(db)
	svEmployee := service.NewEmployeeService(rpEmployee)
	hdEmployee := handler.NewEmployeeController(svEmployee)

	// Product Record
	productRecordRp := product_record_repository.NewProductRecordRepositorySQL(db)
	productRecordSv := service.NewProductRecordService(productRecordRp)
	productRecordHd := handler.NewProductRecordHandler(productRecordSv)

	// InboundOrders
	rpInbounOrders := inbound_order_repository.NewInboundOrderDb(db)
	svInbounOrders := service.NewInboundOrderService(rpInbounOrders)
	hdInbounOrders := handler.NewInboundOrderController(svInbounOrders)

	// - router
	apiV1 := app.Group("/api/v1")

	// 1. Sellers
	sellersGroup := apiV1.Group("/sellers")
	sellersGroup.Get("/", sellerHandler.GetAllSellers)
	sellersGroup.Get("/{id}", sellerHandler.GetSellerById)
	sellersGroup.Post("/", sellerHandler.CreateSeller)
	sellersGroup.Patch("/{id}", sellerHandler.PatchSeller)
	sellersGroup.Delete("/{id}", sellerHandler.DeleteSeller)

	// 2. Warehouses
	warehousesGroup := apiV1.Group("/warehouses")
	warehousesGroup.Get("/", warehouseHd.GetWarehouse)
	warehousesGroup.Get("/{id}", warehouseHd.GetWarehouseById)
	warehousesGroup.Post("/", warehouseHd.CreateWarehouse)
	warehousesGroup.Patch("/{id}", warehouseHd.PatchWarehouse)
	warehousesGroup.Delete("/{id}", warehouseHd.DeleteWarehouse)

	// 3. Sections
	sectionsGroup := apiV1.Group("/sections")
	sectionsGroup.Get("/", sectionsHd.GetSections)
	sectionsGroup.Get("/{id}", sectionsHd.GetSectionById)
	sectionsGroup.Post("/", sectionsHd.CreateSection)
	sectionsGroup.Patch("/{id}", sectionsHd.PatchSection)
	sectionsGroup.Delete("/{id}", sectionsHd.DeleteSection)
	sectionsGroup.Get("/reportProducts", sectionsHd.GetAllProductBatchesBySection)

	// 4. Products
	productsGroup := apiV1.Group("/products")
	productsGroup.Get("/", productHd.GetAllProducts)
	productsGroup.Get("/{id}", productHd.GetProductById)
	productsGroup.Post("/", productHd.CreateProduct)
	productsGroup.Patch("/{id}", productHd.PatchProduct)
	productsGroup.Delete("/{id}", productHd.DeleteProduct)
	productsGroup.Get("/reportRecords", productHd.GetReportRecords)

	// 5. Employees
	employeesGroup := apiV1.Group("/employees")
	employeesGroup.Get("/", hdEmployee.GetAllEmployees)
	employeesGroup.Get("/{id}", hdEmployee.GetEmployeeById)
	employeesGroup.Get("/reportInboundOrders", hdEmployee.GetReportInboundOrders)
	employeesGroup.Post("/", hdEmployee.CreateEmployee)
	employeesGroup.Patch("/{id}", hdEmployee.PatchEmployee)
	employeesGroup.Delete("/{id}", hdEmployee.DeleteEmployee)

	// 6. Buyers
	buyersGroup := apiV1.Group("/buyers")
	buyersGroup.Get("/", buyerHd.GetAllBuyers)
	buyersGroup.Get("/{id}", buyerHd.GetBuyerById)
	buyersGroup.Post("/", buyerHd.CreateBuyer)
	buyersGroup.Patch("/{id}", buyerHd.PatchBuyer)
	buyersGroup.Delete("/{id}", buyerHd.DeleteBuyer)
	buyersGroup.Get("/reportPurchaseOrders", buyerHd.GetBuyersPurchaseOrdersCount)

	// 7. Localities

	localitiesGroup := apiV1.Group("/localities")
	localitiesGroup.Post("/", localityHandler.CreateLocality)
	localitiesGroup.Get("/reportCarries", localityHandler.GetCarriesReport)
	localitiesGroup.Get("/reportSellers", localityHandler.GetLocalitySellerCount)

	// 8. Carries
	carriesGroup := apiV1.Group("/carries")
	carriesGroup.Post("/", carriesHandler.CreateCarrie)

	// 9. Product Records
	productRecordsGroup := apiV1.Group("/productRecords")
	productRecordsGroup.Get("/", productRecordHd.GetAllProductRecords)
	productRecordsGroup.Post("/", productRecordHd.CreateProductRecord)

	// 10. ProductBatches
	productBatchesGroup := apiV1.Group("/productBatches")
	productBatchesGroup.Post("/", productBatchHd.CreateProductBatch)

	// 11. InboundOrders
	inboundOrdersGroup := apiV1.Group("/inboundOrders")
	inboundOrdersGroup.Post("/", hdInbounOrders.CreateInboundOrder)

	// 12. Purchase Orders
	purchaseOrdersGroup := apiV1.Group("/purchaseOrders")
	purchaseOrdersGroup.Post("/", purchaseOrderHd.CreatePurchaseOrder)

	return app.Run()
}
