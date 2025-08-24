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

	// VerifJwt := func(h httprouter.Handle) httprouter.Handle {
	// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 		mw := middleware.JwtVerifyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			h(w, r, ps)
	// 		}))
	// 		mw.ServeHTTP(w, r)
	// 	}
	// }

	// User
	userRepo := repository.NewUserRepositoryImpl(db)
	userService := service.NewUserServiceImpl(userRepo, db)
	userController := controller.NewUserController(userService)

	router.POST("/api/user/createusers", userController.CreateUser)
	router.POST("/api/user/login", userController.LoginUser)
	router.POST("/api/user/changepassword", userController.ChangePassword)


	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads/"))

	// Inventaris
	inventarisRepo := repository.NewItemRepositoryImpl(db)
	inventarisService := service.NewInventarisService(inventarisRepo)
	inventarisController := controller.NewInventarisController(inventarisService)

	router.POST("/api/inventaris/add", inventarisController.CreateInventaris)
	router.GET("/api/inventaris/getbyid/:id", inventarisController.GetByIDInventaris)
	router.GET("/api/inventaris/get", inventarisController.GetAllInventaris)
	router.PUT("/api/inventaris/update/:id", inventarisController.UpdateInventaris)
	router.DELETE("/api/inventaris/delete/:id", inventarisController.DeleteInventaris)
	router.GET("/api/inventaris/search", inventarisController.SearchInventaris)

	// Peminjaman / Barang
	barangRepo := repository.NewBarangRepositoryImpl(db)
	peminjamanRepo := repository.NewPeminjamanRepositoryImpl(db)
	peminjamanService := service.NewPeminjamanServiceImpl(db, barangRepo, peminjamanRepo)
	peminjamanController := controller.NewPeminjamanController(peminjamanService)

	router.GET("/api/barang/tersedia", peminjamanController.GetBarangTersedia)
	router.POST("/api/peminjaman", peminjamanController.CreatePeminjaman)
	router.PUT("/api/peminjaman/kembali/:id", peminjamanController.ReturnPeminjaman)
	router.GET("/api/peminjaman", peminjamanController.ListPeminjaman)
	router.DELETE("/api/peminjaman/delete/:id", peminjamanController.DeletePeminjaman)

	// Stats dashboard
	statsRepo := repository.NewReportRepository(db)
	statsService := service.NewReportService(statsRepo)
	statsController := controller.NewReportController(statsService)

	router.GET("/api/stats/barang", statsController.CountAllBarang)
	router.GET("/api/stats/barang/dipinjam", statsController.CountBarangDipinjam)
	router.GET("/api/stats/barang/rusakberat", statsController.CountBarangRusakBerat)

	reportRepo := repository.NewRepository(db)
	reportService := service.NewService(reportRepo)
	reportController := controller.NewController(reportService)

	// Report endpoints
	router.POST("/api/report/start", reportController.StartReport)
	router.GET("/api/report", reportController.GetReports)
	router.GET("/api/report/detail/:id", reportController.GetReportDetail)
	router.PUT("/api/report/finalize/:id", reportController.FinalizeReport)
	router.GET("/api/report/export/:id", reportController.ExportPDF)

	// Check endpoints
	router.POST("/api/check/:id/add", reportController.AddCheck)
	router.PUT("/api/check/update/:id", reportController.UpdateCheck)
	router.DELETE("/api/check/delete/:id", reportController.DeleteCheck)

	handler := middleware.CorsMiddleware(router)

	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
