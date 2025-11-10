package router

import (
	"crypto/rsa"
	"os"

	"github.com/go-gorote/auth/controller"
	"github.com/go-gorote/auth/goroteadmin"
	"github.com/go-gorote/gorote"
	"github.com/go-gorote/gorote/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/swagger"
)

type AppRouter struct {
	*fiber.App
	PublicKey  *rsa.PublicKey
	Storage    storage.StorageProvider
	Controller controller.Controller
}

func (r *AppRouter) RegisterBaseRouter(router fiber.Router, docSwagger bool) {
	if r.Storage == nil {
		if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
			os.Mkdir("./uploads", 0755)
		}
		r.registerStaticRouter(router.Group("/uploads"))
	}
	if docSwagger {
		r.registerSwagger(router)
	}
	r.Health(router.Group("/health"))
	r.UpdateLogo(router.Group("/logo"))
	// Route Group auth
	r.Login(router.Group("/auth", gorote.Limited(60)))
	r.Logout(router.Group("/auth"))
	r.Refresh(router.Group("/auth", gorote.Limited(60)))
	// Route Group users
	r.ListUser(router.Group("/users"))
	r.RecieveUser(router.Group("/users"))
	r.CreateUser(router.Group("/users"))
	r.UpdateUser(router.Group("/users"))
	r.ChangePassword(router.Group("/users"))
	// Route Group roles
	r.ListRole(router.Group("/roles"))
	r.CreateRole(router.Group("/roles"))
	r.UpdateRole(router.Group("/roles"))
	// Route Group permissions
	r.ListPermission(router.Group("/permissions"))
	r.CreatePermission(router.Group("/permissions"))
	r.UpdatePermission(router.Group("/permissions"))
	// Route Group tenant
	r.ListTenant(router.Group("/tenants"))
	r.CreateTenant(router.Group("/tenants"))
	r.UpdateTenant(router.Group("/tenants"))
}

func (r *AppRouter) registerStaticRouter(router fiber.Router) {
	router.Use(helmet.New(helmet.Config{
		XSSProtection:             "1; mode=block",
		CrossOriginEmbedderPolicy: "unsafe-none",
		CrossOriginOpenerPolicy:   "unsafe-none",
		CrossOriginResourcePolicy: "cross-origin",
	})).Static("/", "./uploads")
}

func (r *AppRouter) RegisterAdminRouter() {
	r.App.Get("/login", r.Controller.SetCookiePainelAdminHandler)
	r.App.Get("/", r.Controller.SetCookiePainelAdminHandler)

	r.App.Use("*", filesystem.New(filesystem.Config{
		Root:         goroteadmin.BuildHTTPFS(),
		PathPrefix:   "",
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))
}

func (r *AppRouter) registerSwagger(router fiber.Router) {
	router.Get("/swagger/*", swagger.HandlerDefault)
}
