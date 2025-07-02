package application

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository/buyer_repository"
	employee_repository "app/internal/repository/employee_repository"
	"app/internal/repository/product_repository"
	"app/internal/repository/product_type_repository"
	"app/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ServerAddress        string
	EmployeesFilePath    string
	BuyerLoaderFilePath  string
	WarehouseFilePath    string
	ProductTypesFilePath string
	ProductsFilePath     string
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
	}

	return &ServerChi{
		serverAddress:       defaultConfig.ServerAddress,
		employeesFilePath:   defaultConfig.EmployeesFilePath,
		buyerLoaderFilePath: defaultConfig.BuyerLoaderFilePath,
		warehouseFilePath:   defaultConfig.WarehouseFilePath,
		productTypeFilePath: defaultConfig.ProductTypesFilePath,
		productsFilePath:    defaultConfig.ProductsFilePath,
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
}

// Run is a method that runs the server
func (a *ServerChi) Run() (err error) {

	ldEmployee := loader.NewEmployeeJSONFile(a.employeesFilePath)

	dbEmployee, err := ldEmployee.Load()

	if err != nil {
		return
	}

	buyer_ld := loader.NewBuyerLoaderJSONFile(a.buyerLoaderFilePath)
	buyer_db, err := buyer_ld.Load()
	if err != nil {
		return
	}

	//load products_type
	product_ld := loader.NewProductLoaderJSONFile(a.productsFilePath)
	product_db, err := product_ld.Load()
	if err != nil {
		return
	}

	product_type_ld := loader.NewProductTypeLoaderJSONFile(a.productTypeFilePath)
	product_type_db, err := product_type_ld.Load()
	if err != nil {
		return
	}

	//Employee - repository
	rpEmployee := employee_repository.NewEmployeeMap(dbEmployee)
	//Employee - service
	svEmployee := service.NewEmployeeService(rpEmployee)
	//Employee - handler
	hdEmployee := handler.NewEmployeeController(svEmployee)

	buyer_rp := buyer_repository.NewBuyerMap(buyer_db)
	buyer_sv := service.NewBuyerDefault(buyer_rp)
	buyer_hd := handler.NewBuyerDefault(buyer_sv)

	// Product - repository
	product_rp := product_repository.NewProductRepositoryMap(product_db)
	product_type_rp := product_type_repository.NewProductTypeRepositoryMap(product_type_db)

	// Product - service
	product_type_sv := service.NewProductTypeService(&product_type_rp)
	product_sv := service.NewProductService(product_rp, &product_type_sv, nil) //agregar service seller

	// Product - handler
	product_hd := handler.NewProductController(&product_sv)

	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// - endpoints
	rt.Route("/api/v1", func(rt chi.Router) {

		//1
		rt.Route("/sellers", func(rt chi.Router) {
			rt.Get("/EXAMPLE", nil)
		})
		//2
		rt.Route("/warehouse", func(rt chi.Router) {
			rt.Get("/EXAMPLE", nil)
		})
		//3
		rt.Route("/sections", func(rt chi.Router) {
			rt.Get("/EXAMPLE", nil)
		})
		//4
		rt.Route("/products", func(rt chi.Router) {
			rt.Get("/", product_hd.GetAllProducts())
			rt.Get("/{id}", product_hd.GetProductById())
			rt.Post("/", product_hd.CreateProduct())
			rt.Patch("/{id}", product_hd.UpdateProduct())
			rt.Delete("/{id}", product_hd.DeleteProduct())
		})
		//5
		rt.Route("/employees", func(rt chi.Router) {
			rt.Get("/", hdEmployee.GetEmployeesList)
			rt.Get("/{id}", hdEmployee.GetEmployeeById)
			rt.Post("/", hdEmployee.SaveEmployee)
			rt.Patch("/{id}", hdEmployee.UpdateEmployee)
			rt.Delete("/{id}", hdEmployee.DeleteEmployee)
		})
		//6
		rt.Route("/buyers", func(rt chi.Router) {
			rt.Get("/", buyer_hd.GetAllBuyers())
			rt.Get("/{id}", buyer_hd.GetBuyerByID())
			rt.Post("/", buyer_hd.CreateBuyer())
			rt.Patch("/{id}", buyer_hd.PatchBuyer())
			rt.Delete("/{id}", buyer_hd.DeleteBuyer())
		})
	})
	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}
