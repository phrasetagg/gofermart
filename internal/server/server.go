package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/phrasetagg/gofermart/internal/app/config"
	"github.com/phrasetagg/gofermart/internal/app/db"
	"github.com/phrasetagg/gofermart/internal/app/handlers/auth"
	"github.com/phrasetagg/gofermart/internal/app/handlers/order"
	"github.com/phrasetagg/gofermart/internal/app/handlers/user"
	"github.com/phrasetagg/gofermart/internal/app/middlewares"
	"github.com/phrasetagg/gofermart/internal/app/repositories"
	"github.com/phrasetagg/gofermart/internal/app/services"
	"log"
	"net/http"
)

var cfg = config.PrepareCfg()

func StartServer() {
	r := chi.NewRouter()

	DB := db.NewDB(cfg.DBDsn)
	err := DB.CreateTables()
	if err != nil {
		panic(err)
	}

	userRepository := repositories.NewUserRepository(DB)
	balanceRepository := repositories.NewBalanceRepository(DB)
	orderRepository := repositories.NewOrderRepository(DB)

	userService := services.NewUserService(userRepository, balanceRepository)
	authService := services.NewAuthService()
	orderService := services.NewOrderService(orderRepository)
	accrualService := services.NewAccrualService(cfg.AccrualAddr, orderRepository)

	go accrualService.StartOrderStatusesUpdating()

	authMiddleware := middlewares.NewAuthMiddleware(authService, userRepository)

	r.Route("/api/user", func(r chi.Router) {
		// Auth routes
		r.Group(func(r chi.Router) {
			r.Post("/register", auth.Register(userService, authService))
			r.Post("/login", auth.Login(userService, authService))
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.CheckAuth())
			r.Get("/orders", order.Get(orderService))
			r.Post("/orders", order.Upload(orderService))

			r.Get("/balance", user.GetBalance(userService))
			r.Post("/balance/withdraw", user.RegisterWithDraw(userService, orderService))

			r.Get("/withdrawals", user.GetWithdrawals(userService))
		})
	})

	log.Fatal(http.ListenAndServe(cfg.ServerAddr, r))
}
