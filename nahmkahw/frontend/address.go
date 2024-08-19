package frontend

import (
	"context"
	"net/http"
	//"strconv"
	"time"
	//"fmt"

	"github.com/labstack/echo"
)

type (
	Address struct {
		StudentCode string  `json:"studentcode" validate:"required"`
		Detail      string  `json:"detail" validate:"required"`
		Amphoe      string  `json:"amphoe" validate:"required"`
		District    string  `json:"district" validate:"required"`
		ProvinceNo  float64 `json:"provinceno" validate:"required"`
		Province    string  `json:"province" validate:"required"`
		Postcode    string  `json:"zipcode" validate:"required"`
		Mobile      string  `json:"mobile" validate:"required"`
		Email       string  `json:"email" validate:"required"`
		Code        string  `json:"code"`
		Created     string  `json:"created"`
		Modified    string  `json:"modified"`
	}
)

func (h *frontendRepoDB) FindAddressDefaultId(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	id := c.Param("id")

	var (
		address Address
	)

	sql := `select std_code, addr_number||' '||road detail ,decode(area,null,'-',area) amphoe,decode(district,null,'-',district) district, decode(s.province_no,null,0,s.province_no) provinceno
	, decode(s.province_no,null,'-',p.province_name_thai) province ,decode(postcode,null,'-',postcode) zipcode,decode(mobile_telephone,null,'-',mobile_telephone) mobile, decode(email_address,null,'-',email_address) email
	from DBBACH00.ugb_student_address s
	left join fees_province p on s.province_no = p.province_no
	where std_code like :1`

	row := h.oracle_db.QueryRow(sql, id)

	err := row.Scan(&address.StudentCode, &address.Detail, &address.Amphoe, &address.District, &address.ProvinceNo, &address.Province, &address.Postcode, &address.Mobile, &address.Email)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}



	return c.JSON(http.StatusOK, address)

}

func (h *frontendRepoDB) CreateAddress(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	r := new(Address)

	if err := c.Bind(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		code    string
		address Address
	)

	stmt, err := h.oracle_db.Prepare("select LOWER(SYS_GUID()) AS code from dual")
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&code)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = tx.ExecContext(ctx, "insert into FEES_ADDRESS (STD_CODE,DETAIL,AMPHOE,DISTRICT,PROVINCE_NO,PROVINCE,POSTCODE,MOBILE,EMAIL,CREATED,MODIFIED,CODE) values (:1,:2,:3,:4,:5,:6,:7,:8,:9,sysdate,sysdate,:10)", r.StudentCode, r.Detail, r.Amphoe, r.District, r.ProvinceNo, r.Province, r.Postcode, r.Mobile, r.Email, code)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	row := h.oracle_db.QueryRow("select * from fees_address where std_code = : 1", r.StudentCode)

	err = row.Scan(&address.StudentCode, &address.Detail, &address.Amphoe, &address.District, &address.ProvinceNo, &address.Province, &address.Postcode, &address.Mobile, &address.Email, &address.Created, &address.Modified, &address.Code)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Logger().Info("frontend-address-create")

	return c.JSON(http.StatusOK, address)

}

func (h *frontendRepoDB) FindAddressId(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	id := c.Param("id")

	var (
		address Address
	)

	sql := `select * from fees_address where STD_CODE like :1`

	row := h.oracle_db.QueryRow(sql, id)

	err := row.Scan(&address.StudentCode, &address.Detail, &address.Amphoe, &address.District, &address.ProvinceNo, &address.Province, &address.Postcode, &address.Mobile, &address.Email, &address.Created, &address.Modified, &address.Code)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Logger().Info("frontend-address-id")

	return c.JSON(http.StatusOK, address)

}

func (h *frontendRepoDB) UpdateAddress(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "PUT")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		address Address
	)

	Address := new(Address)

	if err := c.Bind(Address); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(Address); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = tx.ExecContext(ctx, "UPDATE FEES_ADDRESS SET DETAIL = :1, AMPHOE = :2, DISTRICT = :3, PROVINCE_NO = :4, PROVINCE  = :5, POSTCODE  = :6, MOBILE  = :7, EMAIL= :8, MODIFIED = sysdate WHERE STD_CODE = :9", Address.Detail, Address.Amphoe, Address.District, Address.ProvinceNo, Address.Province, Address.Postcode, Address.Mobile, Address.Email, Address.StudentCode)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sql := `select * from fees_address where STD_CODE like :1`

	row := h.oracle_db.QueryRow(sql, Address.StudentCode)

	err = row.Scan(&address.StudentCode, &address.Detail, &address.Amphoe, &address.District, &address.ProvinceNo, &address.Province, &address.Postcode, &address.Mobile, &address.Email, &address.Created, &address.Modified, &address.Code)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Logger().Info("frontend-address-update")

	return c.JSON(http.StatusOK, address)

}
