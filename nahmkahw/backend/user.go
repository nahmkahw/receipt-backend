package backend

import (
	"context"
	"net/http"

	//"strconv"
	"time"
	//"encoding/json"

	"github.com/labstack/echo"
)

type (
	UserLogin struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		Key      string `json:"key"`
		Created  string `json:"created"`
		Modified string `json:"modified"`
	}

	GraphUser struct {
		displayName       string
		givenName         string
		Id                string
		jobTitle          string
		Mail              string
		mobilePhone       string
		officeLocation    string
		preferredLanguage string
		surname           string
		userPrincipalName string
		access_token      string
		refresh_token     string
		role_token        string
		expiretime        string
	}
)

func (h *backendRepoDB) UpdateUser(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		user UserLogin
	)

	GraphUser := new(GraphUser)

	if err := c.Bind(GraphUser); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(GraphUser); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = tx.ExecContext(ctx, "UPDATE fees_userrole SET KEY = :1 , MODIFIED = sysdate WHERE USERNAME = :2", GraphUser.Id, GraphUser.Mail)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	//s := fmt.Sprintf("ปรับสถานะ %s คำสั่งซื้อเลขที่ : %s ", Order.Status , Order.OrderCode)
	row := h.oracle_db.QueryRow("select username,role,key from fees_userrole where key = :1", GraphUser.Id)

	err = row.Scan(&user.Username, &user.Role, &user.Key)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, user)

}
