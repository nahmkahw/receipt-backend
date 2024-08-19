package backend

import (
	"runtime"
	"fmt"
	"strings"
	"receipt-backend/nahmkahw/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/godror/godror"
	"github.com/labstack/echo"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"

)

type (

	backendRepoDB struct {
		oracle_db *sqlx.DB
		redis_cache *redis.Client
		logger *logrus.Logger
		discordURL string
	}

	BackendRepoInterface interface {
		GetUserSignIn(c echo.Context) error
		GetPhoto(c echo.Context) error
		AuthenRole(c echo.Context) error

		FindOrder(c echo.Context) error 
		FindOrderId(c echo.Context) error
		FindOrderDate(c echo.Context) error
		UpdateOrder(c echo.Context) error

		FindReceiptId(c echo.Context) error
		FindReceipt(c echo.Context) error
		UpdateReceipt(c echo.Context) error

		CreateLogs(c echo.Context) error
		FindLogs(c echo.Context) error

		FindPayment(c echo.Context) error

		FindStudentId(c echo.Context) error
		
		UpdateUser(c echo.Context) error

	}
)

func NewBackendRepo(oracle_db *sqlx.DB,redis_cache *redis.Client, logger *logrus.Logger) BackendRepoInterface {
	return &backendRepoDB{oracle_db: oracle_db, redis_cache: redis_cache ,logger: logger}
}

func (r *backendRepoDB) logAndNotifyError(err error,param string) {
    oraCode := extractORACode(err.Error())

    pc, file, line, ok := runtime.Caller(1)
    if !ok {
        r.logger.Error("Failed to retrieve caller information")
    }
    funcName := runtime.FuncForPC(pc).Name()

    r.logger.WithFields(logrus.Fields{
        "func_name": funcName,
        "file":      file,
        "line":      line,
        "error":     err.Error(),
        "ORA_CODE":  oraCode,
    }).Error("SQL Error")

    message := fmt.Sprintf("SQL Error in %s File: %s Line: %d ORA_CODE: %s Parameter: %s", funcName, file, line, oraCode, param)
    if err := util.SendToDiscord(message); err != nil {
        r.logger.Error("Failed to send message to Discord: ", err)
    }

	if err := util.SendToTeams(message); err != nil {
        r.logger.Error("Failed to send message to Discord: ", err)
    }

	
}

func extractORACode(errorMessage string) string {
    parts := strings.Split(errorMessage, ":")
    if len(parts) > 1 && strings.Contains(parts[0], "ORA-") {
        return strings.TrimSpace(parts[0])
    }
    return "Unknown ORA Code"
}
