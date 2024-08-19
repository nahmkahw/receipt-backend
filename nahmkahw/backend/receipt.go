package backend

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type (
	FeeForm struct {
		FeeRole       string `json:"feerole" validate:"required"`
		StatusOperate string `json:"statusoperate" validate:"required"`
	}

	ReceiptForm struct {
		ReceiptId        int64  `json:"receiptid"`
		StatusOperate    string `json:"statusoperate" validate:"required"`
		OrderCode        string `json:"ordercode" validate:"required"`
		UserUpdate       string `json:"userupdate" validate:"required"`
		AdditionDocument string `json:"additiondocument"`
		StatusCase       string `json:"statuscase" validate:"required"`
	}
)

func (h *backendRepoDB) FindReceiptId(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	id := c.Param("id")

	var (
		receipts []Receipt
		receipt  Receipt
		cache = NewRedisCache(h.redis_cache,time.Second*60)
	)

	sql := `select RECEIPT_ID,CODE,r.STD_CODE,AMOUNT,PRICE,STATUS,CREATED,MODIFIED,ORDER_CODE,ORDER_ID,
	decode(r.CODE,'40','FEEPOST',r.STATUS_OPERATE) STATUS_OPERATE,
	decode(USER_UPDATE, null, '-', USER_UPDATE) USER_UPDATE,
	decode(ADDITION_DOCUMENT, null, '-', ADDITION_DOCUMENT) ADDITION_DOCUMENT,
	decode(r.YEAR, null, '-', r.YEAR) YEAR,
	decode(r.SEMESTER, null, '-', r.SEMESTER) SEMESTER,
	decode(sh.FEE_AMOUNT, null, '-', sh.FEE_AMOUNT) FEE_AMOUNT,
	decode(sh.FEE_FORM, null, 'X', sh.FEE_FORM) FEE_FORM,
	decode(sh.FEE_ROLE, null, 'G', sh.FEE_ROLE) FEE_ROLE,
	decode(sh.FEE_DESCRIPTION, null, 'ไม่พบข้อมูล', sh.FEE_DESCRIPTION) FEE_DESCRIPTION,
	decode(sh.FEE_SEND, null, '-', sh.FEE_SEND) FEE_SEND,
	decode(fee.MONEY_NOTE, null, decode(r.CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name), FEE.MONEY_NOTE) description
	from (select * from fees_receipt where ORDER_CODE = :1) r 
	left join dbeng000.vm_feesem_money_web fee on (r.std_code = fee.std_code and r.code = fee.fee_no)
	left join fees_sheet sh on r.code = sh.fee_no 
	order by 1`

	key := fmt.Sprintf("receipt-%s", id) 

	receiptcache := cache.GetReceiptAll(key)

	if receiptcache == nil {
		fmt.Println("receipt database")

		rows, err := h.oracle_db.Query(sql, id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount,
				&receipt.Price, &receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode, &receipt.OrderId,
				&receipt.StatusOperate, &receipt.UserUpdate, &receipt.AdditionDocument,
				&receipt.Year, &receipt.Semester,
				&receipt.FeeAmount, &receipt.FeeForm, &receipt.FeeRole, &receipt.FeeDescription,
				&receipt.FeeSend, &receipt.Description)
			receipts = append(receipts, receipt)
		}

		if err = rows.Err(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		cache.SetReceiptAll(key, &receipts)

		return c.JSON(http.StatusOK, receipts)
	}

	fmt.Println("receipt redis")

	return c.JSON(http.StatusOK, receiptcache)

}

func (h *backendRepoDB) FindReceipt(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		receipts []Receipt
		receipt  Receipt
		cache = NewRedisCache(h.redis_cache,time.Second*60)
	)

	feeform := new(FeeForm)

	if err := c.Bind(feeform); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(feeform); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	/* Formatted on 23/9/2021 11:43:37 (QP5 v5.185.11230.41888) */
	sql := `select receipt_id, 
		   code,
		   r.std_code,
		   amount,
		   price,
		   status,
		   r.created,
		   r.modified,
		   r.order_code,
		   r.order_id,
		   status_operate,
		   decode (user_update, null, '-', user_update) user_update,
		   decode (addition_document, null, '-', addition_document)
			  addition_document,
		   decode (r.year, null, '-', r.year) year,
		   decode (r.semester, null, '-', r.semester) semester,
		   decode (sh.fee_amount, null, '-', sh.fee_amount) fee_amount,
		   decode (sh.fee_form, null, 'X', sh.fee_form) fee_form,
		   decode (sh.fee_role, null, 'G', sh.fee_role) fee_role,
		   decode (sh.fee_description,
				   null, 'ไม่พบข้อมูล.',
				   sh.fee_description)
			  fee_description,
		   decode (sh.fee_send, null, '-', sh.fee_send) fee_send,
		   decode(fee.MONEY_NOTE, null, decode(CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name), FEE.MONEY_NOTE) description,
		   case
          when     o.status_payment is not null
               and o.status_confirm is null
               and o.status_approve is null
               and o.status_process is null
               and o.status_success is null
          then
             o.status_payment
          when     o.status_payment is not null
               and o.status_confirm is not null
			   and o.status_process is null
               and o.status_approve is null               
               and o.status_success is null
          then
             o.status_confirm
		  when     o.status_payment is not null
               and o.status_confirm is not null
               and o.status_process is not null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_process
          when     o.status_payment is not null
               and o.status_confirm is not null
			   and o.status_process is not null             
               and o.status_approve is not null
               and o.status_success is null
          then
             o.status_approve
          when     o.status_payment is not null
               and o.status_confirm is not null
               and o.status_approve is not null
               and o.status_process is not null
               and o.status_success is not null
          then
             o.status_success
          else
             'ERORR'
       end status,
	   decode(note.note,null,'-',note.note) NOTE
	  from fees_receipt r
		   inner join fees_order o
			  on (    r.order_id = o.order_id
				  and o.ORDER_ID > 789
				  and o.std_code not in ('6299999991','6299999992') 
				  and o.status_approve = 'APPROVE'
				  and o.status_success is null)
		   left join dbeng000.vm_feesem_money_web fee
			  on (r.std_code = fee.std_code and r.code = fee.fee_no)
		   left join fees_sheet sh
			  on r.code = sh.fee_no 
		   left join (select distinct std_code, document_code, decode(NOTE,null,'-',note) NOTE from regis000.VM_MNY_BANK_MACTH_RU WHERE MATCH_RECEIPT IS NOT NULL and NOTE = '-') note
			  on o.std_code = note.std_code and o.document_code = note.document_code `

	switch feeform.FeeRole {
	case "A":
		sql += ` where r.CODE != 40 and sh.FEE_ROLE = 'A' `
	case "B":
		sql += ` where r.CODE != 40 and sh.FEE_ROLE = 'B' `
	case "C":
		sql += ` where r.CODE != 40 and sh.FEE_ROLE = 'C' `
	default:
		sql += ` where r.CODE != 40 `
	}

	switch feeform.StatusOperate {
	case "CANCEL":
		sql += ` and r.STATUS_OPERATE = 'CANCEL' `
	case "PENDING":
		sql += ` and r.STATUS_OPERATE = 'PENDING' `
	case "OPERATE":
		sql += ` and r.STATUS_OPERATE = 'OPERATE' `
	case "SUCCESS":
		sql += ` and r.STATUS_OPERATE = 'SUCCESS' `
	}

	sql += ` order by 1 desc`

	key := fmt.Sprintf("receipt-%s", feeform.StatusOperate)

	receiptcache := cache.GetReceiptAll(key)

	if receiptcache == nil {
		fmt.Println("receipt database :",key)

		rows, err := h.oracle_db.Query(sql)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount, &receipt.Price,
				&receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode, &receipt.OrderId,
				&receipt.StatusOperate, &receipt.UserUpdate, &receipt.AdditionDocument, &receipt.Year, &receipt.Semester,
				&receipt.FeeAmount, &receipt.FeeForm,
				&receipt.FeeRole, &receipt.FeeDescription, &receipt.FeeSend, &receipt.Description,
				&receipt.OrderStatus, &receipt.Note)
			receipts = append(receipts, receipt)
		}

		if err = rows.Err(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if len(receipts) < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูลรายการ"})
		}

		cache.SetReceiptAll(key, &receipts)

		return c.JSON(http.StatusOK, receipts)
	}

	fmt.Println("receipt redis :",key)

	return c.JSON(http.StatusOK, receiptcache)

}

func SqlUpdateStatus(receiptForm *ReceiptForm) string {
	sql := ``
	switch receiptForm.StatusOperate {
	case "CANCEL":
		sql = `UPDATE FEES_RECEIPT SET STATUS_OPERATE = 'CANCEL' , 
				USER_UPDATE = :1, MODIFIED = sysdate WHERE ORDER_CODE = :2 and RECEIPT_ID = :3 `
	case "PENDING":
		sql = `UPDATE FEES_RECEIPT SET STATUS_OPERATE = 'PENDING' , 
				USER_UPDATE = :1, MODIFIED = sysdate WHERE ORDER_CODE = :2 and RECEIPT_ID = :3 `
	case "OPERATE":
		sql = `UPDATE FEES_RECEIPT SET STATUS_OPERATE = 'OPERATE' , 
				USER_UPDATE = :1, MODIFIED = sysdate WHERE ORDER_CODE = :2 and RECEIPT_ID = :3 `
	case "SUCCESS":
		sql = `UPDATE FEES_RECEIPT SET STATUS_OPERATE = 'SUCCESS' , 
				USER_UPDATE = :1, MODIFIED = sysdate WHERE ORDER_CODE = :2 and RECEIPT_ID = :3 `
	}
	return sql
}

func SqlUpdateAdditionDocument() string {
	sql := `UPDATE FEES_RECEIPT SET ADDITION_DOCUMENT = :1, USER_UPDATE = :2, MODIFIED = sysdate WHERE ORDER_CODE = :3 and RECEIPT_ID = :4 `
	return sql
}

func SqlUpdateVerify() string {
	sql := `UPDATE FEES_RECEIPT SET STATUS_OPERATE = :1 , ADDITION_DOCUMENT = :2, USER_UPDATE = :3, MODIFIED = sysdate WHERE ORDER_CODE = :4 and RECEIPT_ID = :5 `
	return sql
}

func (h *backendRepoDB) UpdateReceipt(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		StringError string
		result      sql.Result
		receipt     Receipt
	)

	receiptForm := new(ReceiptForm)

	if err := c.Bind(receiptForm); err != nil {
		StringError = fmt.Sprintf("\n s%) \n ....Form Error!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	if err := c.Validate(receiptForm); err != nil {
		StringError = fmt.Sprintf("\n (s%) \n ....Validate Form Error!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		StringError = fmt.Sprintf("\n (s%) \n .... Connection Error!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	switch receiptForm.StatusCase {
	case "StatusOperate":
		sql := SqlUpdateStatus(receiptForm)
		result, err = tx.ExecContext(ctx, sql,
			receiptForm.UserUpdate, receiptForm.OrderCode, receiptForm.ReceiptId)
	case "AdditionDocument":
		sql := SqlUpdateAdditionDocument()
		result, err = tx.ExecContext(ctx, sql,
			receiptForm.AdditionDocument, receiptForm.UserUpdate, receiptForm.OrderCode, receiptForm.ReceiptId)
	case "VerifyUpdate":
		sql := SqlUpdateVerify()
		result, err = tx.ExecContext(ctx, sql,
			receiptForm.StatusOperate, receiptForm.AdditionDocument, receiptForm.UserUpdate, receiptForm.OrderCode, receiptForm.ReceiptId)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "โปรดระบุสถานะให้ถูกต้อง. (PENDING,OPERATE,SUCCESS)"})
	}

	if err != nil {
		tx.Rollback()
		StringError = fmt.Sprintf("\n (s%) \n ....Transaction rollback!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	rows, err := result.RowsAffected()

	switch {
	case rows != 1:
		StringError = fmt.Sprintf("\n ....Receipt affect %d with id %d!\n", rows, receiptForm.ReceiptId)
		return c.JSON(http.StatusBadRequest, StringError)
	case err != nil:
		tx.Rollback()
		StringError = fmt.Sprintf("\n (s%) \n ....Transaction rollback!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	err = tx.Commit()

	if err != nil {
		StringError = fmt.Sprintf("\n (s%) \n ....Transaction not commit!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	row := h.oracle_db.QueryRow(`select r.RECEIPT_ID,r.CODE,r.STD_CODE,r.AMOUNT,
		r.PRICE,r.STATUS,r.CREATED,r.MODIFIED,r.ORDER_CODE,r.ORDER_ID,
		STATUS_OPERATE,decode(USER_UPDATE, null, '-', USER_UPDATE) USER_UPDATE,decode(ADDITION_DOCUMENT, null, '-', ADDITION_DOCUMENT) ADDITION_DOCUMENT,
		sh.FEE_AMOUNT,sh.FEE_FORM,sh.FEE_ROLE,sh.FEE_DESCRIPTION,sh.FEE_SEND,
		decode(fee.MONEY_NOTE, null, decode(CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name), FEE.MONEY_NOTE) description
		from (select * from fees_receipt where receipt_id = :1) r
		left join dbeng000.vm_feesem_money_web fee on (r.std_code = fee.std_code and r.code = fee.fee_no)
		left join fees_sheet sh on r.code = sh.fee_no 
		where r.ORDER_CODE = :2 order by 1`, receiptForm.ReceiptId, receiptForm.OrderCode)

	err = row.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount,
		&receipt.Price, &receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode, &receipt.OrderId,
		&receipt.StatusOperate, &receipt.UserUpdate, &receipt.AdditionDocument,
		&receipt.FeeAmount, &receipt.FeeForm, &receipt.FeeRole,
		&receipt.FeeDescription, &receipt.FeeSend,
		&receipt.Description)

	if err != nil {
		StringError = fmt.Sprintf("\n (s%) \n .Can not scan row data!\n", err.Error())
		return c.JSON(http.StatusBadRequest, StringError)
	}

	return c.JSON(http.StatusOK, receipt)

}
