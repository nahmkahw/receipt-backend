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

	ReportFee struct {
		FeeNo string `json:"fee_no" db:"fee_no"` 
		FeeName string `json:"fee_name" db:"fee_name"`
		FeeRole   string `json:"fee_role" db:"fee_role"`
	}

	ReportServiceInterface interface {
		GetReportFees(reportFeeRequest *ReportFeeRequest) (*ReportFeeResponse, error)
		GetReport(reportRequest *ReportRequest) (*ReportResponse, error)
	}
)

func NewReportServices(reportRepo repositories.ReportRepoInterface, redis_cache *redis.Client) ReportServiceInterface {
	return &reportServices{
		reportRepo: reportRepo,
		redis_cache:  redis_cache,
	}
}