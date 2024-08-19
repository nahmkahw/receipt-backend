package frontend

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/godror/godror"
	"github.com/go-redis/redis/v7"
)

type (

	frontendRepoDB struct {
		oracle_db *sqlx.DB
		redis_cache *redis.Client
	}

	FrontendRepoInterface interface {
		//address
		CreateAddress(c echo.Context) error
		FindAddressDefaultId(c echo.Context) error
		FindAddressId(c echo.Context) error
		UpdateAddress(c echo.Context) error

		//authen
		CreateLoginStatus(c echo.Context) error
		Logout(c echo.Context) error
		CheckLogin(c echo.Context) error

		//fee
		FindFeesId(c echo.Context) error
		FindYearSemester(c echo.Context) error

		//order
		CreateOrder(c echo.Context) error 
		UpdateOrder(c echo.Context) error 
		FindOrder(c echo.Context) error
		FindOrderId(c echo.Context) error
		
		//receipt
		FindReceipt(c echo.Context) error

		//student
		FindStudentId(c echo.Context) error
	}
)

func NewFrontendRepo(oracle_db *sqlx.DB,redis_cache *redis.Client) FrontendRepoInterface {
	return &frontendRepoDB{oracle_db: oracle_db, redis_cache: redis_cache}
}
