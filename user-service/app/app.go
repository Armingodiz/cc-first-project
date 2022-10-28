package app

import (
	"database/sql"
	"log"
	"net/url"

	_ "github.com/lib/pq"

	"cc-first-project/user-service/controllers/health"
	"cc-first-project/user-service/controllers/userController"
	"cc-first-project/user-service/middlewares"
	"cc-first-project/user-service/store"
	"cc-first-project/user-service/services/brokerService"
	"cc-first-project/user-service/services/userService"

	"github.com/gin-gonic/gin"
)

type App struct {
	route *gin.Engine
}

func NewApp() *App {
	r := gin.Default()
	routing(r)
	return &App{
		route: r,
	}
}

func (a *App) Start(addr string) error {
	return a.route.Run(addr)
}

func routing(r *gin.Engine) {
	r.Use(middlewares.CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	Store := store.NewStore(getPostgresConnection())
	if Store == nil {
		log.Fatal("Failed to setup database")
	}
	UserService := userService.NewUserService(Store)
	BrokerService := brokerService.NewBrokerService()
	UserController := userController.UserController{UserService: UserService, BrokerService: BrokerService}
	healthCheckController := health.NewHealthCheckController()
	//unprotected routes
	r.GET("/health", healthCheckController.GetStatus())
	r.POST("/advertisement", UserController.AddAdvertisement())
	r.GET("/advertisement/:ad_id", UserController.GetAdvertisement())
}


func getPostgresConnection() *sql.DB {
	serviceURI := "postgres://avnadmin:AVNS_8OSIiwW4r_EQh5LHb-u@pg-9240e5b-armingodarzi1380-0360.aivencloud.com:12488/defaultdb?sslmode=require"
	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

	db, err := sql.Open("postgres", conn.String())

	if err != nil {
		log.Fatal(err)
	}
	return db
}