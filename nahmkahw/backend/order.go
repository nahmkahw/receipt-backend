package backend

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
		OrderStatus      string  `json:"orderstatus"`
		AdditionDocument string  `json:"additiondocument"`
		Year             string  `json:"year"`
		Semester         string  `json:"semester"`
		Note             string  `json:"note"`
	}

	Item struct {
		StudentCode   string    `json:"studentcode" validate:"required"`
		StatusPayment string    `json:"statuspayment" validate:"required"`
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
		OrderSlip     string         `json:"orderslip"`
		FlagAddress   string         `json:"FlagAddress" validate:"required"`
		Total         float64        `json:"Total"`
		CountPayment  float64        `json:"CountPayment"`
		CountReceipt  float64        `json:"CountReceipt"`
		Pending       float64        `json:"Pending"`
		Operate       float64        `json:"Operate"`
		Success       float64        `json:"Success"`
		Cancel       float64        `json:"Cancel"`
		CountNone          float64         `json:"None"`
		Cart          []Receipt      `json:"cart"`
		SendCode      string         `json:"SendCode"`
		SendDate      string         `json:"SendDate"`
		Payment       []Payment      `json:"payment"`
		RowNum        float64        `json:"RowNum"`
		Note          string         `json:"Note"`
		CounterA          string         `json:"CounterA"`
		CounterB          string         `json:"CounterB"`
				CounterC          string         `json:"CounterC"`
						CounterD          string         `json:"CounterD"`
								CounterE          string         `json:"CounterE"`
										CounterF          string         `json:"CounterF"`
		FiscalYear		string         `json:"FiscalYear"`
		CounterNo	string         `json:"CounterNo"`
		ReceiptNo string         `json:"ReceiptNo"`
		CheckBill string         `json:"CheckBill"`
		StatusPackage string         `json:"StatusPackage"`
		NameThai string         `json:"NameThai"`
		ZipCode string         `json:"ZipCode"`
		MobileTelephone string         `json:"MobileTelephone"`
		Mobile string         `json:"Mobile"`
		SumTotal float64         `json:"SumTotal"`
	}

	StudentForm struct {
		Status    string `json:"Status" validate:"required"`
		StartDate string `json:"StartDate"`
		EndDate   string `json:"EndDate"`
	}
)

var (
	ctx   context.Context
	cache Cache
)

func (h *backendRepoDB) FindOrder(c echo.Context) error {

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
		orders []Order
		order  Order
		cache  = NewRedisCache(h.redis_cache,time.Second*60)
	)

	sql := `select o.ORDER_ID,o.ORDER_CODE,o.STD_CODE,o.CREATED,o.MODIFIED,
	o.STATUS_VERIFY,o.DATE_VERIFY,o.STATUS_PAYMENT,o.DATE_PAYMENT,o.STATUS_CONFIRM,o.DATE_CONFIRM,
	o.STATUS_PROCESS,o.DATE_PROCESS,o.STATUS_APPROVE,o.DATE_APPROVE,o.STATUS_SUCCESS,o.DATE_SUCCESS,
	o.FLAG_ADDRESS,decode(o.SEND_CODE,null,'-',o.SEND_CODE) SEND_CODE,o.DOCUMENT_CODE,o.TOTAL,decode(o.SEND_DATE,null,'0000-00-00',o.SEND_DATE) SEND_DATE,
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
	decode(p.countpayment,null,0,p.countpayment) countpayment,
	NVL(rc.countreceipt, 0) AS countreceipt,
	NVL(rc.pending, 0) AS pending,
    NVL(rc.operate, 0) AS operate,
    NVL(rc.success, 0) AS success,
	NVL(rc.cancel, 0) AS cancel,
	NVL(note.note, '-') AS NOTE,
	NVL(rc.none, 0) AS none,
	s.NAME_THAI,
    case
        when O.FLAG_ADDRESS  = '1'  then decode(a.postcode,null,'-',a.postcode)
        when O.FLAG_ADDRESS  = '2'  then decode(FA.POSTCODE,null,'-',FA.POSTCODE)
        else decode(a.postcode,null,'-',a.postcode)
    end zipcode,st.sumtotal, NVL(A.MOBILE_TELEPHONE, '-') MOBILE_TELEPHONE, NVL(FA.MOBILE,'-') MOBILE,
	NVL(fr.counterA, 0) AS counterA, 
	NVL(fr.counterB, 0) AS counterB,
	NVL(fr.counterC, 0) AS counterC,
	NVL(fr.counterD, 0) AS counterD,
	NVL(fr.counterE, 0) AS counterE,
	NVL(fr.counterF, 0) AS counterF
	FROM
    fees_order o
LEFT JOIN
    (
        SELECT
            r.order_id,
            SUM(r.amount * r.price) AS sumtotal
        FROM
            fees_receipt r
        WHERE
            R.STATUS_OPERATE <> 'CANCEL'
        GROUP BY
            r.order_id
    ) st ON st.order_id = o.order_id
LEFT JOIN (SELECT
            r.order_id,
            COUNT(r.order_id) AS countreceipt,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'PENDING' THEN 1 END) AS PENDING,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'OPERATE' THEN 1 END) AS OPERATE,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'SUCCESS' THEN 1 END) AS SUCCESS,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'CANCEL' THEN 1 END) AS CANCEL,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'NONE' THEN 1 END) AS NONE
        FROM
            fees_receipt r
        WHERE R.STATUS_OPERATE IN ('PENDING', 'OPERATE','SUCCESS', 'CANCEL','NONE')
        and code != 40 

        GROUP BY
            r.order_id) rc ON rc.order_id = o.order_id
LEFT JOIN
    (
        SELECT
            std_code,
            document_code,
            COUNT(*) AS countpayment
        FROM
            REGIS000.VM_MNY_BANK_MACTH_RU
        WHERE
            MATCH_RECEIPT IS NOT NULL
        GROUP BY
            std_code,
            document_code
    ) p ON o.std_code = p.std_code AND o.document_code = p.document_code
LEFT JOIN
    (
        SELECT
            DISTINCT std_code,
            document_code,
            NVL(NOTE, '-') AS NOTE
        FROM
            REGIS000.VM_MNY_BANK_MACTH_RU
        WHERE
            MATCH_RECEIPT IS NOT NULL
    ) note ON o.std_code = note.std_code AND o.document_code = note.document_code
LEFT JOIN
    DBBACH00.VM_STUDENT_MOBILE s ON o.std_code = s.std_code
LEFT JOIN
    DBBACH00.UGB_STUDENT_ADDRESS a ON o.std_code = a.std_code
LEFT JOIN
    fees_address fa ON o.std_code = fa.std_code
LEFT JOIN
    (
        SELECT
            r.order_id,
			COUNT(CASE WHEN s.fee_role = 'A' THEN 1 END) AS counterA,
            COUNT(CASE WHEN s.fee_role = 'B' THEN 1 END) AS counterB,
			COUNT(CASE WHEN s.fee_role = 'C' THEN 1 END) AS counterC,         
			COUNT(CASE WHEN s.fee_role = 'D' THEN 1 END) AS counterD,
			COUNT(CASE WHEN s.fee_role = 'E' THEN 1 END) AS counterE,
			COUNT(CASE WHEN s.fee_role = 'F' THEN 1 END) AS counterF
        FROM
            fees_receipt r
        LEFT JOIN
            fees_sheet s ON s.FEE_NO = r.CODE
        WHERE
            s.fee_role IN ('A', 'B','C','D','E','F')
        GROUP BY
            r.order_id
    ) fr ON o.order_id = fr.order_id
	where o.ORDER_ID > 789 and 1=1 `

	switch student.Status {
	case "VERIFY":
		sql += ` and o.STATUS_VERIFY = 'VERIFY' and o.STATUS_CONFIRM is null`
	case "QR":
		sql += ` and o.STATUS_PAYMENT = 'QR'`
	case "CONFIRM":
		sql += ` and o.STATUS_CONFIRM = 'CONFIRM' `
	case "PROCESS":
		sql += ` and o.STATUS_PROCESS = 'PROCESS'`
	case "APPROVE":
		sql += ` and o.STATUS_APPROVE = 'APPROVE' and o.STATUS_SUCCESS is null`
	case "SUCCESS":
		sql += ` and o.STATUS_SUCCESS = 'SUCCESS'`
	}

	sql += ` order by 1 desc `

	fmt.Println(student.Status)

	key := fmt.Sprintf("KEY-%s", student.Status)

	fmt.Println(key)

	ordercache := cache.GetOrderAll(key)

	if ordercache == nil {

		fmt.Println("order database :",key)

		rows, err := h.oracle_db.Query(sql)

		if err != nil {
			c.Logger().Error(err.Error())
			param := fmt.Sprintf("query %s",key)
			h.logAndNotifyError(err,param)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
				&order.StatusVerify, &order.DateVerify, &order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm,
				&order.StatusProcess, &order.DateProcess, &order.StatusApprove, &order.DateApprove, &order.StatusSuccess, &order.DateSuccess,
				&order.FlagAddress, &order.SendCode, &order.OrderSlip, &order.Total, &order.SendDate, 
				&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
				&order.Status,&order.StatusPackage,
				&order.CountPayment, &order.CountReceipt, &order.Pending, &order.Operate, &order.Success,&order.Cancel, &order.Note,&order.CountNone,
				&order.NameThai, &order.ZipCode, &order.SumTotal,&order.MobileTelephone,&order.Mobile,&order.CounterA,&order.CounterB,&order.CounterC,&order.CounterD,&order.CounterE,&order.CounterF)
			//fmt.Println(order)
			orders = append(orders, order)

		}

		if err = rows.Err(); err != nil {
			param := fmt.Sprintf("scan %s",key)
			h.logAndNotifyError(err,param)
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if len(orders) < 1 {
			c.Logger().Error("ไม่พบข้อมูล")
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูล."})
		}

		cache.SetOrderAll(key, &orders)

		c.Logger().Info("order backend database")

		return c.JSON(http.StatusOK, orders)

	}

	fmt.Println("order redis :",key)

	return c.JSON(http.StatusOK, ordercache)

}

func (h *backendRepoDB) FindOrderId(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")

	if c.Request().Method != "GET" {
		return c.JSON(http.StatusMethodNotAllowed, "Error Status Method Not Allowed")
	}

	id := c.Param("id")

	var (
		order   Order
		receipt Receipt
		cache  = NewRedisCache(h.redis_cache,time.Second*60)
	)

	key := fmt.Sprintf("KEY-ID-%v", id)

	ordercache := cache.GetOrder(key)

	if ordercache == nil {
		fmt.Println("order database :",key)

		sql := `select o.ORDER_ID,o.ORDER_CODE,o.STD_CODE,o.CREATED,o.MODIFIED,
	o.STATUS_VERIFY,o.DATE_VERIFY,o.STATUS_PAYMENT,o.DATE_PAYMENT,o.STATUS_CONFIRM,o.DATE_CONFIRM,
	o.STATUS_PROCESS,o.DATE_PROCESS,o.STATUS_APPROVE,o.DATE_APPROVE,o.STATUS_SUCCESS,o.DATE_SUCCESS,
	o.FLAG_ADDRESS,decode(o.SEND_CODE,null,'-',o.SEND_CODE) SEND_CODE,o.DOCUMENT_CODE,o.TOTAL,decode(o.SEND_DATE,null,'0000-00-00',o.SEND_DATE) SEND_DATE,
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
	decode(p.countpayment,null,0,p.countpayment) countpayment,
	NVL(rc.countreceipt, 0) AS countreceipt,
	NVL(rc.pending, 0) AS pending,
    NVL(rc.operate, 0) AS operate,
    NVL(rc.success, 0) AS success,
	NVL(rc.cancel, 0) AS cancel,
	NVL(note.note, '-') AS NOTE,
	NVL(rc.none, 0) AS none,
	s.NAME_THAI,
    case
        when O.FLAG_ADDRESS  = '1'  then decode(a.postcode,null,'-',a.postcode)
        when O.FLAG_ADDRESS  = '2'  then decode(FA.POSTCODE,null,'-',FA.POSTCODE)
        else decode(a.postcode,null,'-',a.postcode)
    end zipcode,st.sumtotal,NVL(A.MOBILE_TELEPHONE, '-') MOBILE_TELEPHONE, NVL(FA.MOBILE,'-') MOBILE,
	NVL(fr.counterA, 0) AS counterA, NVL(fr.counterB, 0) AS counterB
	FROM
    fees_order o
LEFT JOIN
    (
        SELECT
            r.order_id,
            SUM(r.amount * r.price) AS sumtotal
        FROM
            fees_receipt r
        WHERE
            R.STATUS_OPERATE <> 'CANCEL'
        GROUP BY
            r.order_id
    ) st ON st.order_id = o.order_id
LEFT JOIN (SELECT
            r.order_id,
            COUNT(r.order_id) AS countreceipt,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'PENDING' THEN 1 END) AS PENDING,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'OPERATE' THEN 1 END) AS OPERATE,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'SUCCESS' THEN 1 END) AS SUCCESS,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'CANCEL' THEN 1 END) AS CANCEL,
            COUNT(CASE WHEN r.STATUS_OPERATE = 'NONE' THEN 1 END) AS NONE
        FROM
            fees_receipt r
        WHERE R.STATUS_OPERATE IN ('PENDING', 'OPERATE','SUCCESS', 'CANCEL','NONE')
        and code != 40 

        GROUP BY
            r.order_id) rc ON rc.order_id = o.order_id
LEFT JOIN
    (
        SELECT
            std_code,
            document_code,
            COUNT(*) AS countpayment
        FROM
            regis000.VM_MNY_BANK_MACTH_RU
        WHERE
            MATCH_RECEIPT IS NOT NULL
        GROUP BY
            std_code,
            document_code
    ) p ON o.std_code = p.std_code AND o.document_code = p.document_code
LEFT JOIN
    (
        SELECT
            DISTINCT std_code,
            document_code,
            NVL(NOTE, '-') AS NOTE
        FROM
            regis000.VM_MNY_BANK_MACTH_RU
        WHERE
            MATCH_RECEIPT IS NOT NULL
    ) note ON o.std_code = note.std_code AND o.document_code = note.document_code
LEFT JOIN
    DBBACH00.VM_STUDENT_MOBILE s ON o.std_code = s.std_code
LEFT JOIN
    DBBACH00.UGB_STUDENT_ADDRESS a ON o.std_code = a.std_code
LEFT JOIN
    fees_address fa ON o.std_code = fa.std_code
LEFT JOIN
    (
        SELECT
            r.order_id,
            COUNT(CASE WHEN s.fee_role = 'B' THEN 1 END) AS counterB,
            COUNT(CASE WHEN s.fee_role = 'A' THEN 1 END) AS counterA
        FROM
            fees_receipt r
        LEFT JOIN
            fees_sheet s ON s.FEE_NO = r.CODE
        WHERE
            s.fee_role IN ('A', 'B')
        GROUP BY
            r.order_id
    ) fr ON o.order_id = fr.order_id
		where o.ORDER_CODE = :1 `

		row := h.oracle_db.QueryRow(sql, id)

		err := row.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
				&order.StatusVerify, &order.DateVerify, &order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm,
				&order.StatusProcess, &order.DateProcess, &order.StatusApprove, &order.DateApprove, &order.StatusSuccess, &order.DateSuccess,
				&order.FlagAddress, &order.SendCode, &order.OrderSlip, &order.Total, &order.SendDate, 
				&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
				&order.Status,&order.StatusPackage,
				&order.CountPayment, &order.CountReceipt, &order.Pending, &order.Operate, &order.Success,&order.Cancel, &order.Note,&order.CountNone,
				&order.NameThai, &order.ZipCode, &order.SumTotal,&order.MobileTelephone,&order.Mobile,&order.CounterA,&order.CounterB)
		if err != nil {
			param := fmt.Sprintf("scan %v",id)
			h.logAndNotifyError(err,param)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		sql = `select r.RECEIPT_ID,r.CODE,r.STD_CODE,r.AMOUNT,r.PRICE,r.STATUS,r.CREATED,r.MODIFIED,r.ORDER_CODE,r.ORDER_ID,
		STATUS_OPERATE,
		decode(r.USER_UPDATE, null, '-', r.USER_UPDATE) USER_UPDATE,
		decode(r.ADDITION_DOCUMENT, null, '-', r.ADDITION_DOCUMENT) ADDITION_DOCUMENT,
		decode(r.YEAR, null, '-', r.YEAR) YEAR,
		decode(r.SEMESTER, null, '-', r.SEMESTER) SEMESTER,
		decode(sh.FEE_AMOUNT, null, '-', sh.FEE_AMOUNT) FEE_AMOUNT,
		decode(sh.FEE_FORM, null, 'X', sh.FEE_FORM) FEE_FORM,
		decode(sh.FEE_ROLE, null, 'G', sh.FEE_ROLE) FEE_ROLE,
		decode(sh.FEE_DESCRIPTION, null, 'ไม่พบข้อมูล.', sh.FEE_DESCRIPTION) FEE_DESCRIPTION,
		decode(sh.FEE_SEND, null, '-', sh.FEE_SEND) FEE_SEND,
		decode(fee.MONEY_NOTE, null, decode(r.CODE,'40','ค่าจัดส่งเอกสาร',sh.fee_name) , FEE.MONEY_NOTE) description
		from (select RECEIPT_ID ,CODE ,STD_CODE ,AMOUNT ,PRICE ,STATUS , CREATED ,MODIFIED ,ORDER_CODE ,ORDER_ID ,STATUS_OPERATE ,USER_UPDATE ,ADDITION_DOCUMENT,YEAR,SEMESTER  from fees_receipt where ORDER_CODE = :1 ) r 
		left join dbeng000.vm_feesem_money_web fee on ( r.std_code = fee.std_code and r.code = fee.fee_no )
		left join fees_sheet sh on r.code = sh.fee_no 
		order by 1` 

		rows, err := h.oracle_db.Query(sql, id)
		if err != nil {
			param := fmt.Sprintf("scan %v",id)
			h.logAndNotifyError(err,param)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&receipt.ReceiptId, &receipt.Code, &receipt.StudentCode, &receipt.Amount,
				&receipt.Price, &receipt.Status, &receipt.Created, &receipt.Modified, &receipt.OrderCode,
				&receipt.OrderId, &receipt.StatusOperate, &receipt.UserUpdate, &receipt.AdditionDocument,
				&receipt.Year, &receipt.Semester,
				&receipt.FeeAmount, &receipt.FeeForm, &receipt.FeeRole,
				&receipt.FeeDescription, &receipt.FeeSend, &receipt.Description)
			order.Cart = append(order.Cart, receipt)
		}

		if err = rows.Err(); err != nil {
			param := fmt.Sprintf("scan %v",id)
			h.logAndNotifyError(err,param)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		cache.SetOrder(key, &order)

		return c.JSON(http.StatusOK, order)
	}

	fmt.Println("order redis :",key)

	return c.JSON(http.StatusOK, ordercache)
}

func (h *backendRepoDB) UpdateOrder(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")

	if c.Request().Method != "PUT" {
		s := fmt.Sprintf("เกิดข้อผิดพลาด %s .", "Error Status Method Not Allowed")
		return c.JSON(http.StatusBadRequest, map[string]string{"message": s})
	}

	Order := new(Order)

	if err := c.Bind(Order); err != nil {
		s := fmt.Sprintf("เกิดข้อผิดพลาด %s .", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"message": s})
	}

	if err := c.Validate(Order); err != nil {
		s := fmt.Sprintf("เกิดข้อผิดพลาด %s .", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"message": s})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		s := fmt.Sprintf("เกิดข้อผิดพลาด %s .", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"message": s})
	}

	// var sendcode string
	// var orderid int64
	// sqlStatement := `SELECT SEND_CODE,ORDER_ID FROM fees_order WHERE SEND_CODE = :1`
	// row := h.oracle_db.QueryRow(sqlStatement, Order.SendCode)
	// err = row.Scan(&sendcode, &orderid)

	// if sendcode == Order.SendCode && orderid != Order.OrderId {
	// 	s := fmt.Sprintf("เลขที่ใบส่งพัสดุ %s มีในระบบแล้ว. เป็นของคำสั่งซื้อเลขที่ %d .", sendcode, orderid)
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"message": s})
	// }

	switch Order.Status {
	case "QR":
		_, err = tx.ExecContext(ctx, "UPDATE fees_order SET DOCUMENT_CODE  = TO_CHAR(SYSDATE+3, 'YYMMDDHH24MISS'), STATUS_PAYMENT= :1 , DATE_PAYMENT = sysdate, MODIFIED = sysdate WHERE ORDER_CODE = :2 ", Order.Status, Order.OrderCode)
	case "CONFIRM":
		_, err = tx.ExecContext(ctx, "UPDATE fees_order SET STATUS_CONFIRM = :1 , DATE_CONFIRM = sysdate, MODIFIED = sysdate WHERE ORDER_CODE = :2 ", Order.Status, Order.OrderCode)
	case "APPROVE":
		_, err = tx.ExecContext(ctx, "UPDATE fees_order SET STATUS_APPROVE = :1 , DATE_APPROVE = sysdate, MODIFIED = sysdate WHERE ORDER_CODE = :2 ", Order.Status, Order.OrderCode)
	case "PROCESS":
		_, err = tx.ExecContext(ctx, "UPDATE fees_order SET STATUS_PROCESS = :1 , DATE_PROCESS = sysdate, MODIFIED = sysdate WHERE ORDER_CODE = :2 ", Order.Status, Order.OrderCode)
	case "SUCCESS":
		_, err = tx.ExecContext(ctx, "UPDATE fees_order SET STATUS_SUCCESS = :1 , SEND_CODE = :2 , SEND_DATE = :3, DATE_SUCCESS = sysdate, MODIFIED = sysdate, STATUS_PACKAGE = :4 WHERE ORDER_CODE = :5 ", Order.Status, Order.SendCode, Order.SendDate, Order.StatusPackage, Order.OrderCode)
	default:
		s := fmt.Sprintf("ไม่พบสถานะ %s ของ คำสั่งซื้อเลขที่ : %s ", Order.Status, Order.OrderCode)
		return c.JSON(http.StatusOK, map[string]string{"message": s})
	}

	if err != nil {
		param := fmt.Sprintf("update %s",Order.OrderCode)
		h.logAndNotifyError(err,param)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		param := fmt.Sprintf("commit %s",Order.OrderCode)
		h.logAndNotifyError(err,param)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	s := fmt.Sprintf("ปรับสถานะ %s คำสั่งซื้อเลขที่ ...: %s ", Order.Status, Order.OrderCode)

	return c.JSON(http.StatusOK, map[string]string{"message": s})

}

func (h *backendRepoDB) FindOrderDate(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")

	if c.Request().Method != "GET" {
		return c.JSON(http.StatusMethodNotAllowed, "Error Status Method Not Allowed")
	}

	student := new(StudentForm)

	if err := c.Bind(student); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(student); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		orders []Order
		order  Order
		cache  = NewRedisCache(h.redis_cache,time.Second*60)
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
               and o.status_approve is null
               and o.status_process is null
               and o.status_success is null
          then
             o.status_payment
          when     o.status_verify is not null
			   and o.status_payment is not null
               and o.status_confirm is not null
               and o.status_approve is null
               and o.status_process is null
               and o.status_success is null
          then
             o.status_confirm
          when     o.status_verify is not null
			   and o.status_payment is not null
               and o.status_confirm is not null               
               and o.status_approve is not null
               and o.status_success is null
          then
             o.status_approve
          when     o.status_verify is not null
			   and o.status_payment is not null
               and o.status_confirm is not null
               and o.status_process is not null
               and o.status_success is null
          then
             o.status_process
          when     o.status_verify is not null
			   and o.status_payment is not null
               and o.status_confirm is not null
               and o.status_approve is not null
               and o.status_process is not null
               and o.status_success is not null
          then
             o.status_success
          else
             'ERORR'
       end status,
	decode(p.countpayment,null,0,p.countpayment) countpayment,
	decode(ro.countreceipt,null,0,ro.countreceipt) countreceipt,
	decode(pe.pending,null,0,pe.pending) pending,
	decode(op.operate,null,0,op.operate) operate,
	decode(s.success,null,0,s.success) success,
	decode(c.cancel,null,0,c.cancel) cancel,
	st.sumtotal
	from fees_order o 
	left join (select r.order_id , sum(r.amount*r.price) sumtotal from fees_receipt r
	where R.STATUS_OPERATE <> 'CANCEL'
	group by r.order_id) st on st.order_id = o.order_id
	left join (select r.order_id ro_order_id,count(STATUS_OPERATE) countreceipt from fees_receipt r group by order_id) ro on ro.ro_order_id = o.order_id
	left join (select r.order_id p_order_id,count(STATUS_OPERATE) pending from fees_receipt r where STATUS_OPERATE = 'PENDING' group by order_id) pe on pe.p_order_id = o.order_id
	left join (select r.order_id op_order_id,count(STATUS_OPERATE) operate from fees_receipt r where STATUS_OPERATE = 'OPERATE' group by order_id) op on op.op_order_id = o.order_id
	left join (select r.order_id s_order_id,count(STATUS_OPERATE) success from fees_receipt r where STATUS_OPERATE = 'SUCCESS' group by order_id) s on s.s_order_id = o.order_id
	left join (select r.order_id c_order_id,count(STATUS_OPERATE) cancel from fees_receipt r where STATUS_OPERATE = 'CANCEL' group by order_id) c on c.c_order_id = o.order_id
	left join (select std_code,substr(qrid,-12) doccode,count(*) countpayment from qr_payment_confirm_tmb qr where qr.system_id = 161 group by std_code,substr(qrid,-12)) p 
	on o.std_code = p.std_code and o.document_code = p.doccode
	where o.ORDER_ID > 789 and 1=1 `

	switch student.Status {
	case "VERIFY":
		sql += ` and o.STATUS_PAYMENT = 'VERIFY' and o.date_verify between to_date(:1,'yyyy-mm-dd') and to_date(:2,'yyyy-mm-dd')`
	case "QR":
		sql += ` and o.STATUS_PAYMENT = 'QR' and o.date_payment between to_date(:1,'yyyy-mm-dd') and to_date(:2,'yyyy-mm-dd')`
	case "CONFIRM":
		sql += ` and o.STATUS_CONFIRM = 'CONFIRM' and o.date_confirm between to_date(:1,'yyyy-mm-dd') and to_date(:2,'yyyy-mm-dd')`
	/*case "APPROVE":
	sql += ` and o.STATUS_APPROVE = 'APPROVE' and o.STATUS_PROCESS is null and o.STATUS_SUCCESS is null `*/
	case "PROCESS":
		sql += ` and o.STATUS_PROCESS = 'PROCESS' and o.STATUS_SUCCESS is null `
	case "SUCCESS":
		sql += ` and o.STATUS_SUCCESS = 'SUCCESS' `
	}

	sql += ` order by 1 desc `

	key := fmt.Sprintf("%s-%s-%s", student.Status, student.StartDate, student.EndDate)

	fmt.Println(key)

	ordercache := cache.GetOrderAll(key)

	if ordercache == nil {

		fmt.Println("order backend oracle")

		rows, err := h.oracle_db.Query(sql, student.StartDate, student.EndDate)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&order.OrderId, &order.OrderCode, &order.StudentCode, &order.Created, &order.Modified,
				&order.StatusPayment, &order.DatePayment, &order.StatusConfirm, &order.DateConfirm,
				&order.StatusProcess, &order.DateProcess, &order.StatusSuccess, &order.DateSuccess,
				&order.FlagAddress, &order.SendCode, &order.OrderSlip, &order.Total, &order.StatusApprove, &order.DateApprove, 
				&order.FiscalYear, &order.CounterNo, &order.ReceiptNo, &order.CheckBill,
				&order.Status,
				&order.CountPayment, &order.CountReceipt, &order.Pending, &order.Operate, &order.Success, &order.Cancel, &order.SumTotal)
			orders = append(orders, order)
		}

		if err = rows.Err(); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if len(orders) < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูล."})
		}

		cache.SetOrderAll(key, &orders)

		return c.JSON(http.StatusOK, orders)
	}

	fmt.Println("order backend redis")

	return c.JSON(http.StatusOK, cache)

}

func DeleteOrder(c echo.Context) error {
	return c.String(http.StatusOK, "delete receipt.")
}
