package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	V1Usecase "github.com/kelvinator07/go-rest-template/internal/business/usecases/v1"
	"github.com/kelvinator07/go-rest-template/internal/datasources/caches"
	V1PostgresRepository "github.com/kelvinator07/go-rest-template/internal/datasources/repositories/postgres/v1"
	V1Handler "github.com/kelvinator07/go-rest-template/internal/http/handlers/v1"
	"github.com/kelvinator07/go-rest-template/pkg/jwt"
	"github.com/kelvinator07/go-rest-template/pkg/mailer"
)

type usersRoutes struct {
	V1Handler      V1Handler.UserHandler
	router         *gin.RouterGroup
	db             *sqlx.DB
	authMiddleware gin.HandlerFunc
}

func NewUsersRoute(router *gin.RouterGroup, db *sqlx.DB, jwtService jwt.JWTService, redisCache caches.RedisCache, ristrettoCache caches.RistrettoCache, authMiddleware gin.HandlerFunc, mailer mailer.OTPMailer) *usersRoutes {
	V1UserRepository := V1PostgresRepository.NewUserRepository(db)
	V1UserUsecase := V1Usecase.NewUserUsecase(V1UserRepository, jwtService, mailer)
	V1UserHandler := V1Handler.NewUserHandler(V1UserUsecase, redisCache, ristrettoCache)

	return &usersRoutes{V1Handler: V1UserHandler, router: router, db: db, authMiddleware: authMiddleware}
}

func (r *usersRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		// auth
		V1AuhtRoute := V1Route.Group("/auth")
		V1AuhtRoute.POST("/register", r.V1Handler.Register)
		V1AuhtRoute.POST("/login", r.V1Handler.Login)
		V1AuhtRoute.POST("/send-otp", r.V1Handler.SendOTP)
		V1AuhtRoute.POST("/verify-otp", r.V1Handler.VerifyOTP)

		// users
		userRoute := V1Route.Group("/users")
		userRoute.Use(r.authMiddleware)
		{
			userRoute.GET("/", r.V1Handler.GetAllUsers)
			userRoute.GET("/me", r.V1Handler.GetUserData)
			// ...
		}
	}

}
