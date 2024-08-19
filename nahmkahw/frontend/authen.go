package frontend

import ( 
	"context"
	"net/http"

	//"strconv"
	"fmt"
	"time"

	//"encoding/json"

	"github.com/labstack/echo"
)

type(
	StudentLogin struct {
		Code   string    `json:"code" validate:"required"`
		Expire string
	}

	Loginstatus struct {
		EncCode string 
		InsertDate string
	}
)


func (h *frontendRepoDB) CreateLoginStatus(c echo.Context) error {

	var loginstatus Loginstatus


	r := new(StudentLogin)

	if err := c.Bind(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sql := `select * from fees_login_status	where ENC_CODE = :1 `

	err := h.oracle_db.QueryRow(sql, r.Code).Scan(&loginstatus.EncCode, &loginstatus.InsertDate)

	if err != nil {
		c.Logger().Error(err.Error())
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		tx, err := h.oracle_db.BeginTx(ctx, nil)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Begin Transaction "+err.Error())
		}

		_, err = tx.ExecContext(ctx, `insert into FEES_LOGIN_STATUS (ENC_CODE,INSERT_DATE) values (:1,sysdate)`, r.Code)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Insert Order "+err.Error())
		}
		err = tx.Commit()

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Commit Transaction "+err.Error())
		}

		sql := `select * from fees_login_status	where ENC_CODE = :1 `

		err = h.oracle_db.QueryRow(sql, r.Code).Scan(&loginstatus.EncCode, &loginstatus.InsertDate)

		if err != nil{
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Select Login "+err.Error())
		}

		c.Logger().Info("frontend-login-status")

		return c.JSON(http.StatusOK, loginstatus)

	}	
	
	c.Logger().Info("frontend-login-status")

	return c.JSON(http.StatusOK, loginstatus)

}

func (h *frontendRepoDB) Logout(c echo.Context) error {

	var (
		cache = NewRedisCacheFrontEnd(h.redis_cache,time.Second*20)
	)

	r := new(StudentLogin)


	if err := c.Bind(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	key := fmt.Sprintf("login-%s", r.Code)

	cache.DeleteLoginStatus(key)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		tx, err := h.oracle_db.BeginTx(ctx, nil)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Begin Transaction "+err.Error())
		}

		chkdel, err := tx.ExecContext(ctx, `delete  from fees_login_status	where ENC_CODE = :1 `, r.Code)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Insert Order "+err.Error())
		}
		err = tx.Commit()

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: Commit Transaction "+err.Error())
		}

		orderID, err := chkdel.RowsAffected()
		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, "Error: RowsAffected Transaction "+err.Error())
		}
		
		c.Logger().Info("frontend-logout")

		return c.JSON(http.StatusOK, orderID)

}

func (h *frontendRepoDB) CheckLogin(c echo.Context) error {
	
	var (
		loginstatus Loginstatus
		cache = NewRedisCacheFrontEnd(h.redis_cache,time.Second*60*60)
	)

	student := new(StudentLogin)

	if err := c.Bind(student); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(student); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	key := fmt.Sprintf("login-%s", student.Code)

	fmt.Println(key)

	logincache := cache.GetLoginStatus(key)

	if logincache == nil {

		fmt.Println("loginstatus frontend oracle")

		sql := `select * from fees_login_status	where ENC_CODE = :1 `

		err := h.oracle_db.QueryRow(sql, student.Code).Scan(&loginstatus.EncCode, &loginstatus.InsertDate)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest,err.Error())
		}

		cache.SetLoginStatus(key, &loginstatus)

		c.Logger().Info("frontend-check-login-database")

		return c.JSON(http.StatusOK, true)

	}

	fmt.Println("loginstatus frontend redis")
	c.Logger().Info("frontend-check-login-redis")

	return c.JSON(http.StatusOK, true)

}
