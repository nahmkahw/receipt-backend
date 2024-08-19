package frontend

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *frontendRepoDB) FindReceipt(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		receipts []Receipt
		receipt  Receipt
	)

	sql := `select RECEIPT_ID,CODE,r.STD_CODE,AMOUNT,PRICE,STATUS,CREATED,MODIFIED,ORDER_CODE,ORDER_ID,
			STATUS_OPERATE,decode(USER_UPDATE, null, '-', USER_UPDATE) USER_UPDATE,
			decode(ADDITION_DOCUMENT, null, '-', ADDITION_DOCUMENT) ADDITION_DOCUMENT,
			sh.FEE_AMOUNT,sh.FEE_FORM,sh.FEE_ROLE,sh.FEE_DESCRIPTION,sh.FEE_SEND,
			decode(fee.MONEY_NOTE, null, decode(CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name), FEE.MONEY_NOTE) description
			from fees_receipt r 
			inner join dbeng000.vm_feesem_money_web fee on (r.std_code = fee.std_code and r.code = fee.fee_no)
			left join fees_sheet sh on r.code = sh.fee_no 
			order by 1 desc`

	rows, err := h.oracle_db.Query(sql)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount, &receipt.Price,
			&receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode, &receipt.OrderId,
			&receipt.StatusOperate, &receipt.UserUpdate, &receipt.AdditionDocument,
			&receipt.FeeAmount, &receipt.FeeForm, &receipt.FeeRole, &receipt.FeeDescription,
			&receipt.FeeSend, &receipt.Description)
		receipts = append(receipts, receipt)
	}

	if err = rows.Err(); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(receipts) < 1 {
		c.Logger().Error("ไม่พบข้อมูล")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูลของ : "})
	}

	c.Logger().Info("frontend-receipt")

	return c.JSON(http.StatusOK, receipts)

}
