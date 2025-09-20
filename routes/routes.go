package routes

import (
	"database/sql"
	"net/http"
	"reset/controller"
	"reset/middleware"
	"reset/repository"
	"reset/service"

	"github.com/julienschmidt/httprouter"
)

func Routes(db *sql.DB, port string) {
	router := httprouter.New()

	// User
	userRepo := repository.NewUserRepositoryImpl(db)
	userService := service.NewUserServiceImpl(userRepo, db)
	userController := controller.NewUserController(userService)

	router.POST("/api/user/createusers", userController.CreateUser)
	router.POST("/api/user/login", userController.LoginUser)
	router.POST("/api/user/changepassword", userController.ChangePassword)

	// Barang CRUD
	barangRepo := repository.NewBarangRepository()
	barangService := service.NewBarangService(db, barangRepo)
	barangController := controller.NewBarangController(barangService)

	router.POST("/api/barang", barangController.CreateBarang)
	router.GET("/api/barang/:id", barangController.GetBarang)
	router.GET("/api/barang", barangController.GetAllBarang)
	router.PUT("/api/barang/:id", barangController.UpdateBarang)
	router.DELETE("/api/barang/:id", barangController.DeleteBarang)

	// Static file serving
	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads/"))

	handler := middleware.CorsMiddleware(router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}