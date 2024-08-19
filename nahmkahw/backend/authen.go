package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	ExtExpiresIn string `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	RefreshToken string `json:"refresh_token"`
	Resource     string `json:"resource"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type TokenError struct {
	CorrelationID    string  `json:"correlation_id"`
	Error            string  `json:"error"`
	ErrorCodes       []int64 `json:"error_codes"`
	ErrorDescription string  `json:"error_description"`
	ErrorURI         string  `json:"error_uri"`
	Timestamp        string  `json:"timestamp"`
	TraceID          string  `json:"trace_id"`
}

type Role struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	RoleToken    string `json:"role_token"`
	DisplayName  string `json:"displayName"`
	Email        string `json:"email"`
	ExpiresTime  int64  `json:"expiretime"`
	role         string `json:"role"`
	org          string `json:"org"`
}

type User struct {
	Odata_context     string      `json:"@odata.context"`
	BusinessPhones    []string    `json:"businessPhones"`
	DisplayName       string      `json:"displayName"`
	GivenName         string      `json:"givenName"`
	ID                string      `json:"id"`
	JobTitle          string      `json:"jobTitle"`
	Mail              string      `json:"mail"`
	MobilePhone       string      `json:"mobilePhone"`
	OfficeLocation    string      `json:"officeLocation"`
	PreferredLanguage interface{} `json:"preferredLanguage"`
	Surname           string      `json:"surname"`
	UserPrincipalName string      `json:"userPrincipalName"`
	AccessToken       string      `json:"access_token"`
	RefreshToken      string      `json:"refresh_token"`
	RoleToken         string      `json:"role_token"`
	ExpiresTime       int64       `json:"expiretime"`
}

func getAuth(c echo.Context) error {
	contentType := c.Request().Header.Get("Authorization")
	return c.String(http.StatusOK, contentType)
}

// e.GET("/users/:id", getUser)
func (h *backendRepoDB) GetUserSignIn(c echo.Context) error {
	// User ID from path `users/:id`
	var (
		token        Token
		user         User
		errormessage ErrorMessage
		err          error
	)
	username := c.FormValue("username")
	password := c.FormValue("password")
	token, err = token.getToken(username, password)
	if err != nil {
		errormessage.Message = "Error validating credentials due to invalid username or password."
		return c.JSON(http.StatusBadRequest, errormessage)
	}
	user, err = token.getUser()
	if err != nil {
		errormessage.Message = "ไม่พบสิทธิการใช้งานระบบรับเงินประจำวัน."
		return c.JSON(http.StatusBadRequest, errormessage)
	}
	user, err = h.getRole(user,token, c)
	if err != nil {
		errormessage.Message = "คุณไม่ได้สิทธิเข้าใช้งานระบบรับเงินประจำวัน."
		return c.JSON(http.StatusBadRequest, errormessage)
	}

	return c.JSON(http.StatusOK, user)

}

func (h *backendRepoDB) getRole(user User,token Token, c echo.Context) (User, error) {
	var (
		role string
	)

	sql := `select role from fees_userrole where username = :1 `

	row := h.oracle_db.QueryRow(sql, user.Mail) 

	err := row.Scan(&role)

	if err != nil {
		return user, err
	}

	// Create token
	tokenJwt := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := tokenJwt.Claims.(jwt.MapClaims)
	claims["email"] = user.Mail
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 4).Unix()

	// Generate encoded token and send it as response.
	t, err := tokenJwt.SignedString([]byte(viper.GetString("token.client_secret")))
	if err != nil {
		return user, err
	}
	user.RoleToken = t
	user.RefreshToken = token.RefreshToken
	user.ExpiresTime = 12 * 60 * 60 * 1000
	user.AccessToken = token.AccessToken

	return user, nil
}

func (h *backendRepoDB) AuthenRole(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	fmt.Println(role)
	return c.JSON(http.StatusOK, role)
}

func (token Token) getToken(Username string, Password string) (Token, error) {
	timeout := time.Duration(22 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	data := url.Values{}
	data.Add("client_secret", viper.GetString("token.client_secret"))
	data.Add("client_id", viper.GetString("token.client_id"))
	data.Add("grant_type", "password")
	data.Add("resource", "https://graph.microsoft.com")
	data.Add("username", Username)
	data.Add("password", Password)

	surl := "https://login.microsoftonline.com/" + viper.GetString("token.tenant_id") + "/oauth2/token"

	req, err := http.NewRequest("POST", surl, bytes.NewBufferString(data.Encode()))
	//req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
	if err != nil {
		return token, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return token, err
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		var tokenerror TokenError
		err = json.Unmarshal(f, &tokenerror)

		if err != nil {
			return token, err
		}

		return token, errors.New("Error validating credentials due to invalid username or password." + tokenerror.ErrorDescription)
	}

	err = json.Unmarshal(f, &token)

	if err != nil {
		return token, err
	}

	return token, nil
}

func (token *Token) getUser() (User, error) {

	var user User

	timeout := time.Duration(5 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	surl := "https://graph.microsoft.com/v1.0/me"

	req, err := http.NewRequest("GET", surl, nil)
	req.Header.Set("Authorization", token.AccessToken)

	if err != nil {
		return user, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return user, err
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	resp.Body.Close()
	if err != nil {
		return user, err
	}

	if resp.StatusCode != 200 {
		var tokenerror TokenError
		err = json.Unmarshal(f, &tokenerror)
		if err != nil {
			return user, err
		}
		return user, err
	}

	err = json.Unmarshal(f, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (h *backendRepoDB) GetPhoto(c echo.Context) error {

	timeout := time.Duration(10 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	surl := "https://graph.microsoft.com/v1.0/me/photos/48x48/$value"

	req, err := http.NewRequest("GET", surl, nil)
	AccessToken := c.Request().Header.Get("Authorization")
	fmt.Println(AccessToken)
	req.Header.Set("Authorization", AccessToken)

	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		var tokenerror TokenError
		/*err = json.Unmarshal(f, &tokenerror)
		if err != nil {
			return err
		}*/

		return errors.New(tokenerror.ErrorDescription)
	}

	return c.Blob(http.StatusOK, "image/jpeg", f)
}
