package services

import (
	"receipt-backend/nahmkahw/errs"
	"encoding/json"
	"fmt"
	"log"
	"time"

)

func (g *reportServices) GetReport(reportRequest *ReportRequest) (*ReportResponse, error) {

	reportResponse := ReportResponse{
		StartDate : reportRequest.StartDate,
		EndDate   : reportRequest.EndDate,
		FeeRole  : reportRequest.FeeRole,
	}

	key := "report::" + reportRequest.FeeRole + reportRequest.StartDate + reportRequest.EndDate
	reportCache, err := g.redis_cache.Get(key).Result()
	if err == nil {
		_ = json.Unmarshal([]byte(reportCache), &reportResponse)
		fmt.Println("report-cache")
		return &reportResponse, nil
	}

	fmt.Println("report-database")

	reportsRepo , err := g.reportRepo.FindReport(reportRequest.StartDate,reportRequest.EndDate,reportRequest.FeeRole)

	if err != nil {
		return &reportResponse, err
	}

	var reports []map[string]interface{}

	for _, item := range reportsRepo {
		reports = append(reports, item)
	}

	reportResponse = ReportResponse{
		StartDate :reportRequest.StartDate,
		EndDate  : reportRequest.EndDate,
		FeeRole  : reportRequest.FeeRole,
		Reports:      reports,
		Count:       len(reports),
	}

	if len(reports) < 1 {
		errStr := fmt.Sprintf("ไม่พบข้อมูลรายงาน %s ถึง %s กลุ่มงาน : %s",reportRequest.StartDate,reportRequest.EndDate,reportRequest.FeeRole)
		return nil, errs.NewNotFoundError(errStr)
	}

	reportsJSON, _ := json.Marshal(&reportResponse)
	timeNow := time.Now()
	redisCachereport := time.Unix(timeNow.Add(time.Second * 5).Unix(), 0)
	_ = g.redis_cache.Set(key, reportsJSON, redisCachereport.Sub(timeNow)).Err()

	return &reportResponse, nil
}

func (g *reportServices) GetReportFees(reportFeeRequest *ReportFeeRequest) (*ReportFeeResponse, error) {

	reportFeeResponse := ReportFeeResponse{
		FeeRole: reportFeeRequest.FeeRole,
	}

	key := "reportfees::" + reportFeeRequest.FeeRole
	reportCache, err := g.redis_cache.Get(key).Result()
	if err == nil {
		log.Println(err)
		_ = json.Unmarshal([]byte(reportCache), &reportFeeResponse)
		fmt.Println("reportfee-cache")
		return &reportFeeResponse, nil
	}

	fmt.Println("reportfee-database")

	reportfeeRepo , err := g.reportRepo.FindReportFees(reportFeeRequest.FeeRole)

	if err != nil {
		return &reportFeeResponse, err
	}

	fees := []ReportFee{}

	for _, item := range reportfeeRepo {
		fees = append(fees, ReportFee{
			FeeNo:     item.FEE_NO,
			FeeName:      item.FEE_NAME,
			FeeRole:   item.FEE_ROLE,
		})
	}

	reportFeeResponse = ReportFeeResponse{
		FeeRole: reportFeeRequest.FeeRole,
		Fees:      fees,
		Count:       len(fees),
	}

	if len(fees) < 1 {
		return nil, errs.NewNotFoundError("ไม่พบข้อมูลค่าธรรมเนียม" + reportFeeRequest.FeeRole)
	}

	reportfeeJSON, _ := json.Marshal(&reportFeeResponse)
	timeNow := time.Now()
	redisCachereportfee := time.Unix(timeNow.Add(time.Second * 5).Unix(), 0)
	_ = g.redis_cache.Set(key, reportfeeJSON, redisCachereportfee.Sub(timeNow)).Err()

	return &reportFeeResponse, nil
}