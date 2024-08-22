package main

import (
	"receipt-backend/nahmkahw/backend"
	"receipt-backend/nahmkahw/frontend"
	"receipt-backend/nahmkahw/environments"
	"receipt-backend/nahmkahw/databases"
	"receipt-backend/nahmkahw/repositories"
	"receipt-backend/nahmkahw/services"
	"receipt-backend/nahmkahw/handlers"
	"receipt-backend/nahmkahw/util" 
	
	"os"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/godror/godror"
	"gopkg.in/go-playground/validator.v9"
	//"github.com/labstack/gommon/log"
	"github.com/jmoiron/sqlx"
	"github.com/go-redis/redis/v7"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var oracle_db *sqlx.DB
var redis_cache *redis.Client


func init() {
	//logger.LoggerInit()
	environments.TimeZoneInit()
	environments.EnvironmentInit()
	oracle_db = databases.NewDatabases().OracleInit()
	redis_cache = databases.NewDatabases().RedisInint()
}

func setupLogger() *logrus.Logger {
    logger := logrus.New()
	file, err := os.OpenFile("/tmp/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        logger.Fatalf("Failed to open log file: %v", err)
    }
    logger.SetOutput(file)
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetLevel(logrus.InfoLevel)
    return logger
}

func main() {
	logger := setupLogger()

	defer oracle_db.Close()
	defer redis_cache.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	backend := backend.NewBackendRepo(oracle_db,redis_cache, logger)
	frontend := frontend.NewFrontendRepo(oracle_db,redis_cache)

	logfile, err := os.OpenFile("/tmp/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer logfile.Close()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
  		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
				`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
				`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
				`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
   		Output: logfile,
	}))
	e.Use(middleware.Logger())
	e.Logger.SetOutput(logfile)

    //e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//config validate

	e.Validator = &CustomValidator{validator: validator.New()}

	testing := e.Group("/testing")

	upload := testing.Group("/upload")
	{
		// Setup fileupload package
		err := util.CreateUploadsDir("/app/fileuploads")
		if err != nil {
			e.Logger.Fatal(err)
		}

		uploadService := services.NewUploadServices("/app/fileuploads")
		uploadHandler := handlers.NewUploadtHandlers(uploadService)
		upload.POST("/", uploadHandler.UploadFileImage)
		upload.POST("/image", uploadHandler.GetFileImage)
	}

	//private group frontend
	private := testing.Group("/frontend")
	{
		//frontend api
		private.GET("/checklogin", frontend.CheckLogin)
		private.POST("/login", frontend.CreateLoginStatus)
		private.POST("/logout", frontend.Logout)

		private.GET("/student/:id", frontend.FindStudentId)

		private.GET("/receipt", frontend.FindReceipt)

		private.POST("/order", frontend.CreateOrder)
		private.GET("/order", frontend.FindOrder)
		private.GET("/order/:id", frontend.FindOrderId)
		private.PUT("/order", frontend.UpdateOrder)

		private.GET("/fees/:id", frontend.FindFeesId)
	
		private.GET("/address-default/:id", frontend.FindAddressDefaultId)

		private.POST("/address", frontend.CreateAddress)
		private.GET("/address/:id", frontend.FindAddressId)
		private.PUT("/address", frontend.UpdateAddress)

		private.GET("/yearsemster", frontend.FindYearSemester)
	}

	//private group backend
	privateBackend := testing.Group("/backend")
	{
		//backend api
		privateBackend.POST("/login", backend.GetUserSignIn)
		privateBackend.GET("/photo", backend.GetPhoto)
	}


	restrict := privateBackend.Group("/restricted")
	restrict.Use(util.ErrorHandlingMiddleware(logger))
	restrict.Use(middleware.JWT([]byte(viper.GetString("token.client_secret"))))
	{
		restrict.GET("/role", backend.AuthenRole)

		//backend api
		restrict.GET("/order", backend.FindOrder)
		restrict.GET("/orderdate", backend.FindOrderDate)
		restrict.GET("/order/:id", backend.FindOrderId) 
		restrict.PUT("/order", backend.UpdateOrder)

		//backend api
		restrict.GET("/receipt", backend.FindReceipt)
		restrict.GET("/receipt/:id", backend.FindReceiptId)
		restrict.PUT("/receipt", backend.UpdateReceipt)

		//backend api
		restrict.PUT("/user", backend.UpdateUser)

		//backend api logs
		restrict.POST("/logs", backend.CreateLogs)
		restrict.GET("/logs", backend.FindLogs)

		//backend api
		restrict.GET("/payment", backend.FindPayment)

		//backend api
		restrict.GET("/student/:id", backend.FindStudentId)
	}

	report := restrict.Group("/report") 
	{

		reportRepo := repositories.NewReportRepo(oracle_db, logger)
		reportService := services.NewReportServices(reportRepo, redis_cache)
		reportHandler := handlers.NewReportHandlers(reportService)
		
		report.POST("/fees", reportHandler.GetReportFees)
		report.POST("/", reportHandler.GetReport)
	}

	report.Use(middleware.Recover())

	//start server
	PORT := viper.GetString("ruipay.port")
	e.Logger.Fatal(e.Start(PORT))

}
