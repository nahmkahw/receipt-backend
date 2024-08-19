package backend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type (
	Payment struct {
		DocumentCode     string  `json:"documentcode"`
		NameReport    string  `json:"namereport"`
		Ref1        string  `json:"ref1"`
		Ref2        string  `json:"ref2"`
		StudentCode string  `json:"studentcode"`
		NameThai    string  `json:"namethai"`
		Mobile      string  `json:"mobile"`
		DateBank	string  `json:"datebank"`
		TimeBank	string  `json:"timebank"`
		MatchReceipt string  `json:"matchreceipt"`
		AmountBank      float64 `json:"amountbank"`
		AmountTotal      float64 `json:"amounttotal"`
		Note      string `json:"note"`
	}
)

func (h *backendRepoDB) FindPayment(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		payments []Payment
		payment  Payment
		cache  = NewRedisCache(h.redis_cache,time.Second*60)
	)

	order := new(Order)

	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(order); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sql := `SELECT DOCUMENT_CODE,
	GROUP_NAME_REPORT,
	REF1,
	REF2,
	STD_CODE,
	NAME_THAI NAMETHAI,
	SUBSTR (REF2, 0, 10) MOBILE,
	DATEBANK,
	TIMEBANK,
	decode(MATCH_RECEIPT,null,'-',MATCH_RECEIPT) MATCH_RECEIPT,
	AMOUNTBANK,
	TOTAL_AMOUNT,
	decode(NOTE,null,'-',note) NOTE
FROM REGIS000.VM_MNY_BANK_MACTH_RU QR
WHERE MATCH_RECEIPT IS NOT NULL AND STD_CODE = :1 AND DOCUMENT_CODE = :2 `

	key := fmt.Sprintf("%s-%s", order.StudentCode, order.OrderSlip)

	paymentcache := cache.GetPaymentAll(key)

	if paymentcache == nil {
		fmt.Println("payment database :",key)
		rows, err := h.oracle_db.Query(sql, order.StudentCode, order.OrderSlip)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&payment.DocumentCode, &payment.NameReport, &payment.Ref1, &payment.Ref2,
				&payment.StudentCode, &payment.NameThai, &payment.Mobile, 
				&payment.DateBank,&payment.TimeBank,&payment.MatchReceipt,
				&payment.AmountBank, &payment.AmountTotal,
				&payment.Note)
			payments = append(payments, payment)
		}

		if err = rows.Err(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if len(payments) < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูลรายการ"})
		}

		cache.SetPaymentAll(key, &payments)

		return c.JSON(http.StatusOK, payments)
	}

	fmt.Println("payment redis:",key)

	return c.JSON(http.StatusOK, paymentcache)
}
