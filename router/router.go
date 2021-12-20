package router

import (
	"freq/handlers"
	"freq/middleware"
	"freq/repository"
	"freq/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repository.NewAuthRepoImpl())}
	ch := handlers.CouponHandler{CouponService: services.NewCouponService(repository.NewCouponRepoImpl())}
	ih := handlers.LoginIpHandler{LoginIpService: services.NewLoginIpService(repository.NewLoginIpRepoImpl())}
	crh := handlers.CustomerHandler{CustomerService: services.NewCustomerService(repository.NewCustomerRepoImpl())}
	ph := handlers.PurchaseHandler{PurchaseService: services.NewPurchaseService(repository.NewPurchaseRepoImpl())}
	prh := handlers.ProductHandler{ProductService: services.NewProductService(repository.NewProductRepoImpl())}

	app.Use(recover.New())

	api := app.Group("", logger.New())

	auth := api.Group("/iriguchi/auth")
	auth.Post("/login", ah.Login)

	product := api.Group("/products")
	product.Get("/:id", prh.FindByProductId)
	product.Get("", prh.FindAll)

	purchase := api.Group("/iriguchi/purchase")
	purchase.Get("/:id", middleware.IsLoggedIn, ph.FindByPurchaseById)
	purchase.Get("", middleware.IsLoggedIn, ph.FindAll)

	items := api.Group("/iriguchi/items")
	items.Post("", middleware.IsLoggedIn, prh.Create)
	items.Put("", middleware.IsLoggedIn, prh.UpdateById)
	items.Delete("/:id", middleware.IsLoggedIn, prh.DeleteById)

	_ = api.Group("/iriguchi/email")

	coupon := api.Group("/iriguchi/coupon")
	coupon.Post("", middleware.IsLoggedIn, ch.Create)
	coupon.Get("", middleware.IsLoggedIn, ch.FindAll)
	coupon.Post("", middleware.IsLoggedIn, ch.FindByCode)
	coupon.Delete("", middleware.IsLoggedIn, ch.DeleteByCode)

	ip := api.Group("/iriguchi/ip")
	ip.Post("", middleware.IsLoggedIn, ih.FindByIp)
	ip.Get("", middleware.IsLoggedIn, ih.FindAll)

	customer := api.Group("/iriguchi/customer")
	customer.Post("", middleware.IsLoggedIn, crh.FindAllByFullName)
	customer.Get("", middleware.IsLoggedIn, crh.FindAll)
}

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	SetupRoutes(app)

	return app
}
