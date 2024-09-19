package repositories

import (
	"receipt-backend/nahmkahw/errs"
	"receipt-backend/nahmkahw/util"
	"fmt"
	"strings"
	"net/http"
    "github.com/sirupsen/logrus"
    "runtime"
)

func (r *reportRepoDB) FindReportFees(feerole string) ([]ReportFee ,error) {
	var (
		fees []ReportFee
		fee  ReportFee
	)

	sqlfee := `select fee_no,fee_name,fee_role from fees_sheet where fee_role = :1`

	rows, err := r.oracle_db.Query(sqlfee, feerole)
	defer rows.Close()

	if err != nil {
		param := fmt.Sprintf("%s",feerole)
		r.logAndNotifyError(err,param)
        return nil, err
	}

	for rows.Next() {
		rows.Scan(&fee.FEE_NO,&fee.FEE_NAME,&fee.FEE_ROLE)

		fees = append(fees,fee)
	}

	return fees, nil
}

func (r *reportRepoDB) GetDateString(str string) string {
		length := len(str)
		switch length {
			case 4 :
				return `'YYYY'`
			case 7 :
				return `'YYYY-MM'`
			case 10 :
				return `'YYYY-MM-DD'`
			default:
				return `'YYYY-MM-DD'`
	}
}

func (r *reportRepoDB) FindReport(startdate,enddate,feerole string) ([]map[string]interface{},[]ReportFee,error) {

	fees, err := r.FindReportFees(feerole)

	if len(fees) < 1 {
		errStr := fmt.Sprintf("ไม่พบข้อมูลค่าธรรมเนียม กลุ่มงาน : %s",feerole)
		return nil,nil, errs.NewMessageAndStatusCode(http.StatusNotFound,errStr)
	}

	dateString := r.GetDateString(startdate)

	if err != nil {
		return nil , nil ,err
	}

	 sql := `SELECT ROW_NUMBER() OVER (ORDER BY date_report) AS id, date_report, `
		var sqlParts []string
		for _, fee := range fees {
			// สร้าง SQL code สำหรับแต่ละฟิลด์
			sqlcode := fmt.Sprintf("SUM(CASE WHEN code_report = '%s' THEN count_report ELSE 0 END) AS code_%s", fee.FEE_NO, fee.FEE_NO)
			sqlParts = append(sqlParts, sqlcode)
		}
	sql += strings.Join(sqlParts, ", ")
	

	sql +=  ` FROM (select TO_CHAR(o.DATE_SUCCESS,` + dateString + `) DATE_REPORT,r.code CODE_REPORT,count(r.RECEIPT_ID) COUNT_REPORT 
				from (select f.receipt_id, f.code,f.order_id from fees_receipt f 
				where f.STATUS_OPERATE = 'SUCCESS' and std_code != '6299999991') r inner join fees_order o on r.order_id = o.order_id and o.STATUS_SUCCESS = 'SUCCESS' and O.DATE_SUCCESS is not null
		where TO_CHAR(o.DATE_SUCCESS,` + dateString + `) between :1 and :2 group by r.code,TO_CHAR(o.DATE_SUCCESS,` + dateString + `))
		GROUP BY 
			date_report
		ORDER BY
    date_report`

		rows, err := r.oracle_db.Query(sql,startdate,enddate)

		defer rows.Close()

		if err != nil {
			param := fmt.Sprintf("%s,%s,%s",startdate,enddate,feerole)
			r.logAndNotifyError(err,param)
			return nil,nil, err
		}

		// ดึงชื่อคอลัมน์
		columns, err := rows.Columns()
		if err != nil {
			return nil, nil , err
		}



		var reports []map[string]interface{}

		for rows.Next() {
			// สร้าง slice ของ interface{} ขนาดเท่ากับจำนวนคอลัมน์
			values := make([]interface{}, len(columns))
			// สร้าง slice ของ pointers สำหรับ scan
			scanArgs := make([]interface{}, len(columns))
			for i := range values {
				scanArgs[i] = &values[i]
			}

			// สแกนข้อมูล
			err := rows.Scan(scanArgs...)
			if err != nil {
				return nil,nil,err
			}

			// สร้าง map สำหรับเก็บข้อมูล
			report := make(map[string]interface{})
			for i, col := range columns {
				if values[i] != nil {
					// เก็บข้อมูลลงใน map โดยใช้ชื่อคอลัมน์เป็นคีย์
					report[col] = values[i]
					// for _, fee := range fees {
					// 	if(strings.Contains(col,fee.FEE_NO)){
					// 		report[col+"_name"] = fee.FEE_NAME
					// 	}
					// }
				} else {
					report[col] = nil
				}
			}

			reports = append(reports, report)
		}

		if err = rows.Err(); err != nil {
			return nil,nil,err
		}

	return reports,fees, nil
}

func (r *reportRepoDB) logAndNotifyError(err error,param string) {
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
