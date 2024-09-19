
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

)

func (g *reportServices) GetReportSummary() (*ReportSummaryResponse, error)  {
	reportSummaryResponse := ReportSummaryResponse{
		SUMMARY_CODE : "", 
		SUMMARY_NAME : "",
		SUMMARY_COUNT   : "",
		MAX_CODE : "", 
		MAX_NAME : "",
		MAX_COUNT   : "",
		MIN_CODE : "", 
		MIN_NAME : "",
		MIN_COUNT   : "",
		SUCCESS_TODAY : "", 
		SUCCESS_THIS_MONTH : "",
		SUCCESS_THIS_YEAR   : "",
		CANCEL_TODAY : "", 
		CANCEL_THIS_MONTH : "",
		CANCEL_THIS_YEAR : "",
	}

	key := "reportsummary"
	reportCache, err := g.redis_cache.Get(key).Result()
	if err == nil {
		log.Println(err)
		_ = json.Unmarshal([]byte(reportCache), &reportSummaryResponse)
		fmt.Println("reportsummary-cache")
		return &reportSummaryResponse, nil
	}

	fmt.Println("reportsummary-database")

	reportRepo , err := g.reportRepo.FindReportSummary()

	if err != nil {
		return &reportSummaryResponse, err
	}

	reportSuccessRepo , err := g.reportRepo.FindReportSuccess()

	if err != nil {
		return &reportSummaryResponse, err
	}

	reportCancelRepo , err := g.reportRepo.FindReportCancel()

	if err != nil {
		return &reportSummaryResponse, err
	}

	reportSummaryResponse = ReportSummaryResponse{
		SUMMARY_CODE : reportRepo.SUMMARY_CODE, 
		SUMMARY_NAME : reportRepo.SUMMARY_NAME,
		SUMMARY_COUNT   : reportRepo.SUMMARY_COUNT,
		MAX_CODE : reportRepo.MAX_CODE, 
		MAX_NAME : reportRepo.MAX_NAME,
		MAX_COUNT   : reportRepo.MAX_COUNT,
		MIN_CODE : reportRepo.MIN_CODE, 
		MIN_NAME : reportRepo.MIN_NAME,
		MIN_COUNT   : reportRepo.MIN_COUNT,
		SUCCESS_TODAY : reportSuccessRepo.SUCCESS_TODAY, 
		SUCCESS_THIS_MONTH : reportSuccessRepo.SUCCESS_THIS_MONTH,
		SUCCESS_THIS_YEAR   : reportSuccessRepo.SUCCESS_THIS_YEAR,
		CANCEL_TODAY : reportCancelRepo.CANCEL_TODAY, 
		CANCEL_THIS_MONTH : reportCancelRepo.CANCEL_THIS_MONTH,
		CANCEL_THIS_YEAR   : reportCancelRepo.CANCEL_THIS_YEAR,
	}

	reportfeeJSON, _ := json.Marshal(&reportSummaryResponse)
	timeNow := time.Now()
	redisCachereportfee := time.Unix(timeNow.Add(time.Minute * 15).Unix(), 0)
	_ = g.redis_cache.Set(key, reportfeeJSON, redisCachereportfee.Sub(timeNow)).Err()

	return &reportSummaryResponse, nil
}