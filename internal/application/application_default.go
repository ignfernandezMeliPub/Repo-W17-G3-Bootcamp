package application

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository/buyer_repository"
	employee_repository "app/internal/repository/employee_repository"
	warehouse_repository "app/internal/repository/warehouse_repository"
	"app/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ServerAddress       string
	EmployeesFilePath   string
	BuyerLoaderFilePath string
	WarehouseFilePath   string
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
	}

	return &ServerChi{
		serverAddress:       defaultConfig.ServerAddress,
		employeesFilePath:   defaultConfig.EmployeesFilePath,
		buyerLoaderFilePath: defaultConfig.BuyerLoaderFilePath,
		warehouseFilePath:   defaultConfig.WarehouseFilePath,
	}
}

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress       string
	employeesFilePath   string
	buyerLoaderFilePath string
	warehouseFilePath   string
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

	warehouse_lb := loader.NewWarehouseJSONFile(a.warehouseFilePath)
	warehouse_db, err := warehouse_lb.Load()
	if err != nil {
		return
	}

	//Employee - repository
	rpEmployee := employee_repository.NewEmployeeMap(dbEmployee)
	//Employee - service
	svEmployee := service.NewEmployeeService(rpEmployee, *svWarehouse)
	//svEmployee := service.NewEmployeeService(rpEmployee, svWarehouse)
	//Employee - handler
	hdEmployee := handler.NewEmployeeController(svEmployee)

	buyer_rp := buyer_repository.NewBuyerMap(buyer_db)
	buyer_sv := service.NewBuyerDefault(buyer_rp)
	buyer_hd := handler.NewBuyerDefault(buyer_sv)

	//warehouse
	warehouse_rp := warehouse_repository.NewWarehouseMap(warehouse_db)
	warehouse_sv := service.NewWarehouseDefault(warehouse_rp)
	warehouse_hd := handler.NewWarehouseDefault(warehouse_sv)

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
		rt.Route("/warehouses", func(rt chi.Router) {
			rt.Get("/", warehouse_hd.FindWarehouse())
			rt.Get("/{id}", warehouse_hd.FindWarehouseById())
			rt.Post("/", warehouse_hd.CreateWarehouse())
			rt.Patch("/{id}", warehouse_hd.UpdateWarehouse())
			rt.Delete("/{id}", warehouse_hd.DeleteWarehouse())

		})
		//3
		rt.Route("/sections", func(rt chi.Router) {
			rt.Get("/EXAMPLE", nil)
		})
		//4
		rt.Route("/products", func(rt chi.Router) {
			rt.Get("/EXAMPLE", nil)
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
