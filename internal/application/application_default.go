package application

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository/buyer_repository"
	employee_repository "app/internal/repository/employee_repository"
	"app/internal/repository/product_repository"
	"app/internal/repository/product_type_repository"
	"app/internal/repository/sections_repository"
	"app/internal/repository/seller_repository"
	warehouse_repository "app/internal/repository/warehouse_repository"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ProductTypesFilePath string
	ProductsFilePath     string
	ServerAddress        string
	EmployeesFilePath    string
	BuyerLoaderFilePath  string
	WarehouseFilePath    string
	SectionsFilePath     string
}

// NewServerChi is a function that returns a new instance of ServerChi
func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	// default values
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
	}
	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}
		if cfg.EmployeesFilePath != "" {
			defaultConfig.EmployeesFilePath = cfg.EmployeesFilePath
		}
		if cfg.BuyerLoaderFilePath != "" {
			defaultConfig.BuyerLoaderFilePath = cfg.BuyerLoaderFilePath
		}
		if cfg.WarehouseFilePath != "" {
			defaultConfig.WarehouseFilePath = cfg.WarehouseFilePath
		}
		if cfg.ProductTypesFilePath != "" {
			defaultConfig.ProductTypesFilePath = cfg.ProductTypesFilePath
		}
		if cfg.ProductsFilePath != "" {
			defaultConfig.ProductsFilePath = cfg.ProductsFilePath
		}
		if cfg.SectionsFilePath != "" {
			defaultConfig.SectionsFilePath = cfg.SectionsFilePath
		}
	}

	return &ServerChi{
		serverAddress:       defaultConfig.ServerAddress,
		employeesFilePath:   defaultConfig.EmployeesFilePath,
		buyerLoaderFilePath: defaultConfig.BuyerLoaderFilePath,
		warehouseFilePath:   defaultConfig.WarehouseFilePath,
		productTypeFilePath: defaultConfig.ProductTypesFilePath,
		productsFilePath:    defaultConfig.ProductsFilePath,
		sectionsFilePath:    defaultConfig.SectionsFilePath,
	}
}

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress       string
	employeesFilePath   string
	buyerLoaderFilePath string
	warehouseFilePath   string
	productTypeFilePath string
	productsFilePath    string
	sectionsFilePath    string
}

// Run is a method that runs the server
func (a *ServerChi) Run() (err error) {

	ldEmployee := loader.NewEmployeeJSONFile(a.employeesFilePath)

	dbEmployee, err := ldEmployee.Load()

	if err != nil {
		return
	}

	// Seller
	sellerRepo := seller_repository.NewSellerRepositoryMap(make(map[int]models.Seller))
	sellerService := service.NewSellerService(&sellerRepo)
	sellerHandler := handler.NewSellerHandler(&sellerService)

	buyerLd := loader.NewBuyerLoaderJSONFile(a.buyerLoaderFilePath)
	buyerDb, err := buyerLd.Load()
	if err != nil {
		return
	}

	// load products_type
	productLd := loader.NewProductLoaderJSONFile(a.productsFilePath)
	productDb, err := productLd.Load()
	if err != nil {
		return
	}

	productTypeLd := loader.NewProductTypeLoaderJSONFile(a.productTypeFilePath)
	productTypeDb, err := productTypeLd.Load()

	if err != nil {
		return
	}

	warehouseLb := loader.NewWarehouseJSONFile(a.warehouseFilePath)
	warehouseDb, err := warehouseLb.Load()
	if err != nil {
		return
	}

	buyerRp := buyer_repository.NewBuyerMap(buyerDb)
	buyerSv := service.NewBuyerDefault(buyerRp)
	buyerHd := handler.NewBuyerDefault(buyerSv)

	// Product - repository
	productRp := product_repository.NewProductRepositoryMap(productDb)
	productTypeRp := product_type_repository.NewProductTypeRepositoryMap(productTypeDb)

	// Product - service
	productTypeSv := service.NewProductTypeService(productTypeRp)
	productSv := service.NewProductService(productRp, productTypeSv, nil) // agregar service seller

	// Product - handler
	productHd := handler.NewProductController(&productSv)

	// warehouse
	warehouseRp := warehouse_repository.NewWarehouseMap(warehouseDb)
	warehouseSv := service.NewWarehouseDefault(warehouseRp)
	warehouseHd := handler.NewWarehouseDefault(warehouseSv)

	// sections
	sectionsRp := sections_repository.NewSectionsRepositoryMap()
	sectionsDb, err := loader.LoadDataFromFile[models.Section](a.sectionsFilePath)
	if err != nil {
		return err
	}
	err = sectionsRp.PoblateSectionsRepo(sectionsDb)
	if err != nil {
		return err
	}
	sectionsSv := service.NewSectionsService(sectionsRp, warehouseSv, productTypeSv)
	sectionsHd := handler.NewSectionsController(sectionsSv)

	//Employee - repository
	rpEmployee := employee_repository.NewEmployeeMap(dbEmployee)
	//Employee - service
	svEmployee := service.NewEmployeeService(rpEmployee, *warehouse_sv)
	//svEmployee := service.NewEmployeeService(rpEmployee, svWarehouse)
	//Employee - handler
	hdEmployee := handler.NewEmployeeController(svEmployee)
  
	// - router

	rt := chi.NewRouter()

	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// - endpoints
	rt.Route("/api/v1", func(rt chi.Router) {
		// 1
		rt.Route("/sellers", func(rt chi.Router) {
			rt.Get("/", sellerHandler.GetAll)
			rt.Get("/{id}", sellerHandler.GetById)
			rt.Post("/", sellerHandler.Create)
			rt.Patch("/", sellerHandler.Patch)
			rt.Delete("/{id}", sellerHandler.Delete)
		})
		// 2
		rt.Route("/warehouses", func(rt chi.Router) {
			rt.Get("/", warehouseHd.FindWarehouse())
			rt.Get("/{id}", warehouseHd.FindWarehouseById())
			rt.Post("/", warehouseHd.CreateWarehouse())
			rt.Patch("/{id}", warehouseHd.UpdateWarehouse())
			rt.Delete("/{id}", warehouseHd.DeleteWarehouse())

		})
		// 3
		rt.Route("/sections", func(rt chi.Router) {
			rt.Get("/", sectionsHd.GetSections)
			rt.Get("/{id}", sectionsHd.GetSection)
			rt.Post("/", sectionsHd.CreateSection)
			rt.Patch("/{id}", sectionsHd.UpdateSection)
			rt.Delete("/{id}", sectionsHd.DeleteSection)
		})
		// 4
		rt.Route("/products", func(rt chi.Router) {
			rt.Get("/", productHd.GetAllProducts())
			rt.Get("/{id}", productHd.GetProductById())
			rt.Post("/", productHd.CreateProduct())
			rt.Patch("/{id}", productHd.UpdateProduct())
			rt.Delete("/{id}", productHd.DeleteProduct())
		})
		// 5
		rt.Route("/employees", func(rt chi.Router) {
			rt.Get("/", hdEmployee.GetEmployeesList)
			rt.Get("/{id}", hdEmployee.GetEmployeeById)
			rt.Post("/", hdEmployee.SaveEmployee)
			rt.Patch("/{id}", hdEmployee.UpdateEmployee)
			rt.Delete("/{id}", hdEmployee.DeleteEmployee)
		})
		// 6
		rt.Route("/buyers", func(rt chi.Router) {
			rt.Get("/", buyerHd.GetAllBuyers())
			rt.Get("/{id}", buyerHd.GetBuyerByID())
			rt.Post("/", buyerHd.CreateBuyer())
			rt.Patch("/{id}", buyerHd.PatchBuyer())
			rt.Delete("/{id}", buyerHd.DeleteBuyer())
		})
	})

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}
