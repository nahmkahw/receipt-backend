package frontend

import (
	"context"
	"database/sql"
	"net/http"

	//"strconv"
	"fmt"
	"time"

	//"encoding/json"

	"github.com/labstack/echo"
)

type (
	Receipt struct {
		ReceiptId        int64   `json:"receiptid"`
		Code             string  `json:"code" validate:"required"`
		StudentCode      string  `json:"studentcode" validate:"required"`
		Amount           int64   `json:"amount" validate:"gt=0"`
		Price            float64 `json:"price" validate:"required"`
		Status           string  `json:"status" validate:"required"`
		Created          string  `json:"created"`
		Modified         string  `json:"modified"`
		OrderCode        string  `json:"ordercode"`
		OrderId          int64   `json:"orderid"`
		StatusOperate    string  `json:"statusoperate"`
		FeeAmount        string  `json:"feeamount"`
		FeeForm          string  `json:"feeform"`
		FeeRole          string  `json:"feerole"`
		FeeDescription   string  `json:feedescription`
		FeeSend          string  `json:"feesend"`
		Description      string  `json:"description"`
		UserUpdate       string  `json:"userupdate"`
		AdditionDocument string  `json:"additiondocument"`
		Year             string  `json:"year"`
		Semester         string  `json:"semester"`
	}

	Item struct {
		StudentCode   string    `json:"studentcode" validate:"required"`
		StatusPayment string    `json:"statuspayment" validate:"required"`
		FlagAddress   string    `json:"flagaddress" validate:"required"`
		Total         float64   `json:"total" validate:"required"`
		Cart          []Receipt `json:"cart" validate:"required"`
	}

	Order struct {
		OrderId       int64          `json:"OrderId"`
		OrderCode     string         `json:"OrderCode" validate:"required"`
		StudentCode   string         `json:"StudentCode"`
		Created       string         `json:"Created"`
		Modified      string         `json:"Modified"`
		StatusVerify string         `json:"StatusVerify"`
		DateVerify   string         `json:"DateVerify"`
		StatusPayment sql.NullString         `json:"StatusPayment"`
		DatePayment   sql.NullString         `json:"DatePayment"`
		StatusConfirm sql.NullString `json:"StatusConfirm"`
		DateConfirm   sql.NullString `json:"DateConfirm"`
		StatusApprove sql.NullString `json:"StatusApprove"`
		DateApprove   sql.NullString `json:"DateApprove"`
		StatusProcess sql.NullString `json:"StatusProcess"`
		DateProcess   sql.NullString `json:"DateProcess"`
		StatusSuccess sql.NullString `json:"StatusSuccess"`
		DateSuccess   sql.NullString `json:"DateSuccess"`
		Status        string         `json:"Status" validate:"required"`
		OrderSlip     string         `json:"OrderSlip"`
		FlagAddress   string         `json:"FlagAddress" validate:"required"`
		SendCode      string         `json:"SendCode"`
		SendDate      string         `json:"SendDate"`
		Total         float64        `json:"Total"` 
		Cart          []Receipt      `json:"cart"`
		FiscalYear    string         `json:"FiscalYear"`
		CounterNo     string         `json:"CounterNo"`
		ReceiptNo     string         `json:"ReceiptNo"`
		CheckBill     string         `json:"CheckBill"`
		SumTotal string         `json:"SumTotal"`
		StatusPackage string         `json:"StatusPackage"`
	}

	StudentForm struct {
		StudentCode string `json:"StudentCode" validate:"required"`
		Status      string `json:"Status" validate:"required"`
	}

	OrderMessage struct {
		StudentCode string `json:"StudentCode" validate:"required"`
	}
)

var (
	ctx context.Context
)

// Handlers
func (h *frontendRepoDB) CreateOrder(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")

	r := new(Item)

	if err := c.Bind(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(r); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		code         string
		orderid      int64
		documentcode string
		order        Order
	)

	stmt, err := h.oracle_db.Prepare(`select LOWER(SYS_GUID()) AS code , FEES_ORDER_SEQ.NEXTVAL AS orderid, TO_CHAR(SYSDATE+3, 'YYMMDDHH24MISS') AS documentcode  from dual`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error: Get Code "+err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&code, &orderid, &documentcode)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, "Error: Scan Code "+err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, "Error: Begin Transaction "+err.Error())
	}

	_, err = tx.ExecContext(ctx, `insert into FEES_ORDER (ORDER_ID,ORDER_CODE,STD_CODE,CREATED,MODIFIED,STATUS_VERIFY,DATE_VERIFY,FLAG_ADDRESS,DOCUMENT_CODE,TOTAL) values (:1,:2,:3,sysdate,sysdate,:4,sysdate,:5,:6,:7)`, orderid, code, r.StudentCode, r.StatusPayment, r.FlagAddress, documentcode, r.Total)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, "Error: Insert Order "+err.Error())
	}

	for _, item := range r.Cart {

		_, err := tx.ExecContext(ctx, `insert into FEES_RECEIPT (RECEIPT_ID,CODE,STD_CODE,AMOUNT,PRICE,STATUS,CREATED,MODIFIED,ORDER_CODE,ORDER_ID,YEAR,SEMESTER) values (FEES_RECEIPT_SEQ.NEXTVAL,:1,:2,:3,:4,:5,sysdate,sysdate,:6,:7,:8,:9)`, item.Code, item.StudentCode, item.Amount, item.Price, item.Status, code, orderid, item.Year, item.Semester)

		if err != nil {
			c.Logger().Error(err.Error())
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, "Error: Insert Receipt "+err.Error())
		}

	}

	err = tx.Commit()

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, "Error: Commit Transaction "+err.Error())
	}

	row := h.oracle_db.QueryRow(`select o.ORDER_ID,o.ORDER_CODE,o.STD_CODE,o.CREATED,o.MODIFIED,
	o.STATUS_VERIFY,o.DATE_VERIFY,o.STATUS_PAYMENT,o.DATE_PAYMENT,o.STATUS_CONFIRM,o.DATE_CONFIRM,
	o.STATUS_PROCESS,o.DATE_PROCESS,o.STATUS_SUCCESS,o.DATE_SUCCESS,o.FLAG_ADDRESS,o.SEND_CODE,o.DOCUMENT_CODE,o.TOTAL,o.STATUS_APPROVE,
	o.DATE_APPROVE,o.SEND_DATE,
	decode(o.FISCAL_YEAR,null,'-',o.FISCAL_YEAR) FISCAL_YEAR,
	decode(o.COUNTER_NO,null,'-',o.COUNTER_NO) COUNTER_NO,
	decode(o.RECEIPT_NO,null,'-',o.RECEIPT_NO) RECEIPT_NO,
	decode(o.CHECK_BILL,null,'-',o.CHECK_BILL) CHECK_BILL, 
	case
		   when     o.status_verify is not null
               and o.status_payment is null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_verify
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_payment
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_confirm
          
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null
			   and o.status_process is not null
			   and o.status_approve is null
               and o.status_success is null
          then
			 o.status_process
		 when     o.status_verify is not null
             and o.status_payment is not null
			 and o.status_confirm is not null  
			 and o.status_process is not null             
			 and o.status_approve is not null
			 and o.status_success is null
		then
		   o.status_approve
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null               
			   and o.status_process is not null
			   and o.status_approve is not null
               and o.status_success is not null
          then
             o.status_success
          else
             'ERORR'
	   end status
	   from fees_order o where o.ORDER_ID = :1 order by 1`, orderid)

	err = row.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
		&order.StatusVerify, &order.DateVerify,&order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm,
		&order.StatusProcess, &order.DateProcess, &order.StatusSuccess, &order.DateSuccess,
		&order.FlagAddress, &order.SendCode, &order.OrderSlip, &order.Total, &order.StatusApprove, &order.DateApprove, &order.SendDate,
		&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
		&order.Status)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, "Error: Scan Order "+err.Error())
	}

	c.Logger().Info("frontend-order-create")

	return c.JSON(http.StatusOK, order)

}

func (h *frontendRepoDB) FindOrder(c echo.Context) error {

	student := new(StudentForm)

	if err := c.Bind(student); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(student); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		orders  []Order
		order   Order
		receipt Receipt
		cache = NewRedisCacheFrontEnd(h.redis_cache,time.Second*60)
	)

	sql := `select o.ORDER_ID,o.ORDER_CODE,o.STD_CODE,o.CREATED,o.MODIFIED,o.STATUS_VERIFY,o.DATE_VERIFY,o.STATUS_PAYMENT,o.DATE_PAYMENT,o.STATUS_CONFIRM,o.DATE_CONFIRM,
	o.STATUS_PROCESS,o.DATE_PROCESS,o.STATUS_SUCCESS,o.DATE_SUCCESS,o.FLAG_ADDRESS,o.SEND_CODE,o.DOCUMENT_CODE,o.TOTAL,o.STATUS_APPROVE,
	o.DATE_APPROVE,o.SEND_DATE,
	decode(o.FISCAL_YEAR,null,'-',o.FISCAL_YEAR) FISCAL_YEAR,
	decode(o.COUNTER_NO,null,'-',o.COUNTER_NO) COUNTER_NO,
	decode(o.RECEIPT_NO,null,'-',o.RECEIPT_NO) RECEIPT_NO,
	decode(o.CHECK_BILL,null,'-',o.CHECK_BILL) CHECK_BILL,
	case
		   when     o.status_verify is not null
               and o.status_payment is null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_verify
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_payment
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_confirm
          
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null
			   and o.status_process is not null
			   and o.status_approve is null
               and o.status_success is null
          then
			 o.status_process
		 when     o.status_verify is not null
             and o.status_payment is not null
			 and o.status_confirm is not null  
			 and o.status_process is not null             
			 and o.status_approve is not null
			 and o.status_success is null
		then
		   o.status_approve
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null               
			   and o.status_process is not null
			   and o.status_approve is not null
               and o.status_success is not null
          then
             o.status_success
          else
             'ERORR'
	   end status,
	   decode(o.STATUS_PACKAGE,null,'-',o.STATUS_PACKAGE) STATUS_PACKAGE,  
	   r.sumtotal
	from fees_order o 
	left join (select r.order_id , sum(r.amount*r.price) sumtotal from fees_receipt r
	where R.STATUS_OPERATE <> 'CANCEL'
	group by r.order_id) r on r.order_id = o.order_id
	where o.ORDER_ID > 789 and o.STD_CODE = :1 `

	switch student.Status {
	case "VERIFY":
		sql += ` and o.STATUS_VERIFY = 'VERIFY' and o.STATUS_PAYMENT is null and o.STATUS_CONFIRM is null and o.STATUS_APPROVE is null and o.STATUS_PROCESS is null and o.STATUS_SUCCESS is null `
	case "QR":
		sql += ` and o.STATUS_PAYMENT = 'QR' and o.STATUS_CONFIRM is null and o.STATUS_APPROVE is null and o.STATUS_PROCESS is null and o.STATUS_SUCCESS is null `
	case "CONFIRM":
		sql += ` and o.STATUS_CONFIRM = 'CONFIRM' and o.STATUS_APPROVE is null and o.STATUS_PROCESS is null and o.STATUS_SUCCESS is null `
	case "PROCESS":
		sql += ` and o.STATUS_PROCESS = 'PROCESS' and o.STATUS_APPROVE is null and o.STATUS_SUCCESS is null `
	case "APPROVE":
		sql += ` and o.STATUS_APPROVE = 'APPROVE' and o.STATUS_SUCCESS is null `
	case "SUCCESS":
		sql += ` and o.STATUS_SUCCESS = 'SUCCESS' `
	}

	sql += ` order by 1 desc `

	key := fmt.Sprintf("order-%s-%s", student.Status, student.StudentCode)

	fmt.Println(key)

	ordercache := cache.GetOrderAll(key)

	if ordercache == nil {

		fmt.Println("order frontend oracle")

		rows, err := h.oracle_db.Query(sql, student.StudentCode)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
				&order.StatusVerify, &order.DateVerify,&order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm,
				&order.StatusProcess, &order.DateProcess, &order.StatusSuccess, &order.DateSuccess,
				&order.FlagAddress, &order.SendCode, &order.OrderSlip, &order.Total, &order.StatusApprove, &order.DateApprove, &order.SendDate,
				&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
				&order.Status , &order.StatusPackage, &order.SumTotal)

			sql = `select r.RECEIPT_ID,r.CODE,r.STD_CODE,r.AMOUNT,r.PRICE,r.STATUS,r.CREATED,r.MODIFIED,r.ORDER_CODE,r.ORDER_ID,
			STATUS_OPERATE,
			decode(r.USER_UPDATE, null, '-', r.USER_UPDATE) USER_UPDATE,
			decode(r.ADDITION_DOCUMENT, null, '-', r.ADDITION_DOCUMENT) ADDITION_DOCUMENT,
			decode(r.YEAR, null, '-', r.YEAR) YEAR,
			decode(r.SEMESTER, null, '-', r.SEMESTER) SEMESTER,
			decode(sh.FEE_AMOUNT, null, '-', sh.FEE_AMOUNT) FEE_AMOUNT,
			decode(sh.FEE_FORM, null, 'X', sh.FEE_FORM) FEE_FORM,
			decode(sh.FEE_ROLE, null, 'G', sh.FEE_ROLE) FEE_ROLE,
			decode(sh.FEE_DESCRIPTION, null, 'ไม่พบข้อมูล', sh.FEE_DESCRIPTION) FEE_DESCRIPTION,
			decode(sh.FEE_SEND, null, '-', sh.FEE_SEND) FEE_SEND,
			decode(fee.MONEY_NOTE, null, decode(r.CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name) , FEE.MONEY_NOTE) description
			from (select RECEIPT_ID ,CODE ,STD_CODE ,AMOUNT ,PRICE ,STATUS , CREATED ,MODIFIED ,ORDER_CODE ,ORDER_ID ,STATUS_OPERATE ,USER_UPDATE ,ADDITION_DOCUMENT,YEAR,SEMESTER  from fees_receipt where ORDER_CODE = :1 ) r 
			left join dbeng000.vm_feesem_money_web fee on ( r.std_code = fee.std_code and r.code = fee.fee_no )
			left join fees_sheet sh on r.code = sh.fee_no 
			order by 1`

			rowRecripts, err := h.oracle_db.Query(sql, order.OrderCode)

			order.Cart = nil

			if err != nil {
				c.Logger().Error(err.Error())
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			defer rowRecripts.Close()

			for rowRecripts.Next() {
				rowRecripts.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount,
					&receipt.Price, &receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode,
					&receipt.OrderId, &receipt.StatusOperate, &receipt.UserUpdate, &receipt.AdditionDocument,
					&receipt.Year, &receipt.Semester,
					&receipt.FeeAmount, &receipt.FeeForm, &receipt.FeeRole,
					&receipt.FeeDescription, &receipt.FeeSend, &receipt.Description)
				//fmt.Println(order.OrderId, receipt)
				order.Cart = append(order.Cart, receipt)
			}

			if err = rowRecripts.Err(); err != nil {
				c.Logger().Error(err.Error())
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			orders = append(orders, order)
		}

		if err = rows.Err(); err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if len(orders) < 1 {
			c.Logger().Error("ไม่พบข้อมูลของ : " + student.StudentCode)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูลของ : " + student.StudentCode})
		}

		cache.SetOrderAll(key, &orders)

		c.Logger().Info("frontend-order-database")

		return c.JSON(http.StatusOK, orders)

	}

	fmt.Println("order frontend redis")

	c.Logger().Info("frontend-order-redis")

	return c.JSON(http.StatusOK, ordercache)

}

func (h *frontendRepoDB) FindOrderId(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")

	id := c.Param("id")

	var (
		order   Order
		receipt Receipt
	)

	sql := `select o.ORDER_ID,o.ORDER_CODE,o.STD_CODE,o.CREATED,o.MODIFIED,o.STATUS_VERIFY,o.DATE_VERIFY,o.STATUS_PAYMENT,o.DATE_PAYMENT,o.STATUS_CONFIRM,o.DATE_CONFIRM,
	o.STATUS_PROCESS,o.DATE_PROCESS,o.STATUS_SUCCESS,o.DATE_SUCCESS,o.FLAG_ADDRESS,o.SEND_CODE,o.DOCUMENT_CODE,o.TOTAL,o.STATUS_APPROVE,
	o.DATE_APPROVE,o.SEND_DATE,
	decode(o.FISCAL_YEAR,null,'-',o.FISCAL_YEAR) FISCAL_YEAR,
	decode(o.COUNTER_NO,null,'-',o.COUNTER_NO) COUNTER_NO,
	decode(o.RECEIPT_NO,null,'-',o.RECEIPT_NO) RECEIPT_NO,
	decode(o.CHECK_BILL,null,'-',o.CHECK_BILL) CHECK_BILL,
	case
		   when     o.status_verify is not null
               and o.status_payment is null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_verify
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_payment
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null               
			   and o.status_process is null
			   and o.status_approve is null
               and o.status_success is null
          then
             o.status_confirm
          
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null
			   and o.status_process is not null
			   and o.status_approve is null
               and o.status_success is null
          then
			 o.status_process
		 when     o.status_verify is not null
             and o.status_payment is not null
			 and o.status_confirm is not null  
			 and o.status_process is not null             
			 and o.status_approve is not null
			 and o.status_success is null
		then
		   o.status_approve
          when     o.status_verify is not null
               and o.status_payment is not null
               and o.status_confirm is not null               
			   and o.status_process is not null
			   and o.status_approve is not null
               and o.status_success is not null
          then
             o.status_success
          else
             'ERORR'
	   end status,r.sumtotal
	from fees_order o 
	left join (select r.order_id , sum(r.amount*r.price) sumtotal from fees_receipt r
	where R.STATUS_OPERATE <> 'CANCEL'
	group by r.order_id) r on r.order_id = o.order_id
	where o.ORDER_CODE = :1 `

	row := h.oracle_db.QueryRow(sql, id)

	err := row.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
		&order.StatusVerify, &order.DateVerify,&order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm, &order.StatusProcess,
		&order.DateProcess, &order.StatusSuccess, &order.DateSuccess, &order.FlagAddress,
		&order.SendCode, &order.OrderSlip, &order.Total, &order.StatusApprove, &order.DateApprove, &order.SendDate,
		&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
		&order.Status,&order.SumTotal)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sql = `select RECEIPT_ID,CODE,r.STD_CODE,AMOUNT,PRICE,STATUS,CREATED,MODIFIED,ORDER_CODE,ORDER_ID,STATUS_OPERATE,
	sh.FEE_AMOUNT,sh.FEE_FORM,sh.FEE_ROLE,sh.FEE_DESCRIPTION,sh.FEE_SEND,
	decode(fee.MONEY_NOTE, null, decode(r.CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name) , FEE.MONEY_NOTE) description
	from fees_receipt r 
	left join dbeng000.vm_feesem_money_web fee on (r.std_code = fee.std_code and r.code = fee.fee_no)
	left join fees_sheet sh on r.code = sh.fee_no 
	where r.ORDER_CODE = :1 order by 1`

	rows, err := h.oracle_db.Query(sql, id)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount, &receipt.Price, &receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode, &receipt.OrderId, &receipt.StatusOperate, &receipt.FeeAmount, &receipt.FeeForm, &receipt.FeeRole, &receipt.FeeDescription, &receipt.FeeSend, &receipt.Description)
		order.Cart = append(order.Cart, receipt)
	}

	if err = rows.Err(); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Logger().Info("frontend-order-id")

	return c.JSON(http.StatusOK, order)

}

func (h *frontendRepoDB) UpdateOrder(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	var (
		order Order
	)

	Order := new(Order)

	if err := c.Bind(Order); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(Order); err != nil {
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

	switch Order.Status {
	case "QR":
		_, err = tx.ExecContext(ctx, "UPDATE FEES_ORDER SET FLAG_ADDRESS = :1 , MODIFIED = sysdate WHERE ORDER_CODE = :2 and STATUS_CONFIRM is null ", Order.FlagAddress, Order.OrderCode)
	case "CONFIRM":
		_, err = tx.ExecContext(ctx, "UPDATE FEES_ORDER SET STATUS_CONFIRM = 'CONFIRM' ,DATE_CONFIRM = sysdate , MODIFIED = sysdate WHERE ORDER_CODE = :1 and STATUS_CONFIRM is null ", Order.OrderCode)
	default:
		s := fmt.Sprintf("ไม่พบสถานะ %s ของ คำสั่งซื้อเลขที่ : %s ", Order.Status, Order.OrderCode)
		return c.JSON(http.StatusOK, map[string]string{"message": s})
	}

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	//s := fmt.Sprintf("ปรับสถานะ %s คำสั่งซื้อเลขที่ : %s ", Order.Status , Order.OrderCode)
	row := h.oracle_db.QueryRow(`select o.ORDER_ID,o.ORDER_CODE,o.STD_CODE,o.CREATED,o.MODIFIED,o.STATUS_PAYMENT,o.DATE_PAYMENT,o.STATUS_CONFIRM,o.DATE_CONFIRM,
	o.STATUS_PROCESS,o.DATE_PROCESS,o.STATUS_SUCCESS,o.DATE_SUCCESS,o.FLAG_ADDRESS,o.SEND_CODE,o.DOCUMENT_CODE,o.TOTAL,o.STATUS_APPROVE,
	o.DATE_APPROVE,o.SEND_DATE,
	decode(o.FISCAL_YEAR,null,'-',o.FISCAL_YEAR) FISCAL_YEAR,
	decode(o.COUNTER_NO,null,'-',o.COUNTER_NO) COUNTER_NO,
	decode(o.RECEIPT_NO,null,'-',o.RECEIPT_NO) RECEIPT_NO,
	decode(o.CHECK_BILL,null,'-',o.CHECK_BILL) CHECK_BILL, 
	case
          when     o.status_payment is not null
               and o.status_confirm is null               
			   and o.status_process is null
			   and o.status_approve is null
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
			   and o.status_process is not null
			   and o.status_approve is not null
               and o.status_success is not null
          then
             o.status_success
          else
             'ERORR'
       end status 
	from fees_order o where o.ORDER_CODE = :1 order by 1`, Order.OrderCode)

	err = row.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
		&order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm, &order.StatusProcess,
		&order.DateProcess, &order.StatusSuccess, &order.DateSuccess, &order.FlagAddress,
		&order.SendCode, &order.OrderSlip, &order.Total, &order.StatusApprove, &order.DateApprove, &order.SendDate,
		&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
		&order.Status)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Logger().Info("frontend-order-update")

	return c.JSON(http.StatusOK, order)

}

func DeleteOrder(c echo.Context) error {
	return c.String(http.StatusOK, "delete receipt.")
}
