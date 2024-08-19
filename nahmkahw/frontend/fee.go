package frontend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	_ "github.com/godror/godror"
)

type Fee struct {
	FeeNo         string  `json:"code"`
	MoneyNote     string  `json:"description"`
	FeeAmount     int64   `json:"amount"`
	FeePrice      float64 `json:"price"`
	FeeStatus     string  `json:"status"`
	FeeRole       string  `json:"role"`
	FeeSend       string  `json:"sentstatus"`
	FeeComent     string  `json:"comment"`
	StatusCurrent string  `json:statuscurrent`
}

type YearSemester struct {
	Year     string `json:"year"`
	Semester string `json:"semester"`
}

func (h *frontendRepoDB) FindFeesId(c echo.Context) error {

	id := c.Param("id")
	fmt.Println(id)

	sql := `select fee.fee_no,NVL(sh.fee_form,'-') status,NVL(sh.FEE_AMOUNT,0) amount,fee.fee_amount price,money_note,
	NVL(sh.fee_role,'-') fee_role ,NVL(sh.fee_description ,'-') fee_comment,NVL(sh.fee_send,'-') fee_send,s.std_status_current
		from dbeng000.vm_feesem_money_web fee
		left join fees_sheet sh on FEE.FEE_NO = sh.fee_no
		left join DBBACH00.VM_STUDENT_MOBILE s on fee.std_code = s.std_code
		where fee.std_code = :1`

	var (
		fees []Fee
		fee  Fee
	)

	rows, err := h.oracle_db.Query(sql, id)
	defer rows.Close()

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	for rows.Next() {
		err = rows.Scan(&fee.FeeNo, &fee.FeeStatus, &fee.FeeAmount, &fee.FeePrice, &fee.MoneyNote, &fee.FeeRole, &fee.FeeComent, &fee.FeeSend, &fee.StatusCurrent)
		if err != nil {
			c.Logger().Error(err.Error())
			fmt.Println(err.Error())
		}

		// if fee.StatusCurrent == "E" {
		// 	if fee.FeeNo == "35" || fee.FeeNo == "36" {
		// 		fees = append(fees, fee)
		// 	}
		// } else {
		// 	fees = append(fees, fee)
		// }

		fees = append(fees, fee)

	}

	if len(fees) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูลค่าธรรมเนียมของ : " + id})
	}

	c.Logger().Info("frontend-fee-id")
	return c.JSON(http.StatusOK, fees)
}

func (h *frontendRepoDB) FindYearSemester(c echo.Context) error {

	sql := `select year,semester from 
		(select /*+ INDEX (e INDEX_YEAR_SEMESTER )*/ distinct year,semester from tr000.ria_regis_ru24_his e where semester is not null and year is not null and semester != 0) 
		group by year,semester order by 1 desc,2`

	var (
		yearsemesters []YearSemester
		yearsemester  YearSemester
		cache = NewRedisCacheFrontEnd(h.redis_cache,time.Second*60)
	)

	currentTime := time.Now()

	key := fmt.Sprintf("yearsemester-%s", currentTime.Format("01-02-2006"))

	yearsemestercache := cache.GetYearSemesterAll(key)

	if yearsemestercache == nil {

		fmt.Println("year frontend oracle")

		rows, err := h.oracle_db.Query(sql)
		defer rows.Close()

		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		for rows.Next() {
			rows.Scan(&yearsemester.Year, &yearsemester.Semester)
			yearsemesters = append(yearsemesters, yearsemester)
		}

		if err = rows.Err(); err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		cache.SetYearSemesterAll(key, &yearsemesters)

		c.Logger().Info("frontend-year-semester-database")

		return c.JSON(http.StatusOK, yearsemesters)
	}

	fmt.Println("year frontend redis:" + key)

	c.Logger().Info("frontend-year-semester-redis")

	return c.JSON(http.StatusOK, yearsemestercache)

}
