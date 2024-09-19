package services

import (
	"receipt-backend/nahmkahw/repositories"
	"github.com/go-redis/redis/v7"
)

type (
	reportServices struct {
		reportRepo repositories.ReportRepoInterface
		redis_cache  *redis.Client
	}

	ReportRequest struct {
		StartDate string `json:"start_date" validate:"required"`
		EndDate   string `json:"end_date" validate:"required"`
		FeeRole   string `json:"fee_role" validate:"required"`
	}

	ReportResponse struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		FeeRole   string `json:"fee_role"`
		Count      int  `json:"count"`
		Reports []map[string]interface{} `json:"reports"`
	}

	ReportFeeRequest struct {
		FeeRole   string `json:"fee_role" validate:"required"`
	}

	ReportFeeResponse struct {
		FeeRole string `json:"fee_role" db:"fee_role"` 
		Count int `json:"count"` 
		Fees []ReportFee `json:"fees"` 
	}

	ReportSummaryResponse struct {
		SUMMARY_CODE string `json:"summary_code"` 
		SUMMARY_NAME string `json:"summary_name"`
		SUMMARY_COUNT   string `json:"summary_count"`
		MAX_CODE string `json:"max_code"` 
		MAX_NAME string `json:"max_name"`
		MAX_COUNT   string `json:"max_count"`
		MIN_CODE string `json:"min_code"` 
		MIN_NAME string `json:"min_name"`
		MIN_COUNT   string `json:"min_count"`
		SUCCESS_TODAY string `json:"success_today"` 
		SUCCESS_THIS_MONTH string `json:"success_this_month"`
		SUCCESS_THIS_YEAR   string `json:"success_this_year"`
		CANCEL_TODAY string `json:"cancel_today"` 
		CANCEL_THIS_MONTH string `json:"cancel_this_month"`
		CANCEL_THIS_YEAR   string `json:"cancel_this_year"`
	}

	ReportSuccessResponse struct {
		SUCCESS_TODAY string `json:"success_today"` 
		SUCCESS_THIS_MONTH string `json:"success_this_month"`
		SUCCESS_THIS_YEAR   string `json:"success_this_year"`
	}

	ReportCancelResponse struct {
		CANCEL_TODAY string `json:"cancel_today"` 
		CANCEL_THIS_MONTH string `json:"cancel_this_month"`
		CANCEL_THIS_YEAR   string `json:"cancel_this_year"`
	}

	ReportFee struct {
		FeeNo string `json:"fee_no" db:"fee_no"` 
		FeeName string `json:"fee_name" db:"fee_name"`
		FeeRole   string `json:"fee_role" db:"fee_role"`
	}

	ReportServiceInterface interface {
		GetReportFees(reportFeeRequest *ReportFeeRequest) (*ReportFeeResponse, error)
		GetReport(reportRequest *ReportRequest) (*ReportResponse, error)

		GetReportSummary() (*ReportSummaryResponse, error)
	}
)

func NewReportServices(reportRepo repositories.ReportRepoInterface, redis_cache *redis.Client) ReportServiceInterface {
	return &reportServices{
		reportRepo: reportRepo,
		redis_cache:  redis_cache,
	}
}