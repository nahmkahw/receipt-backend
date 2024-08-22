package repositories

import (
	"github.com/jmoiron/sqlx"
    "github.com/sirupsen/logrus"
)

type (
	reportRepoDB struct {
		oracle_db *sqlx.DB
		logger *logrus.Logger
		discordURL string
	}

	ReportDay struct {
		StartDate string `json:"startdate"`
		EndDate   string `json:"enddate"`
		FeeRole   string `json:"FeeRole"`
		Reportdays []map[string]interface{} `json:"reportdays"`
		Count      int                       `json:"count"`
	}

	ReportFee struct {
		FEE_NO string `json:"FEE_NO"` 
		FEE_NAME string `json:"FEE_NAME"`
		FEE_ROLE   string `json:"FEE_ROLE"`
	}


	ReportRepoInterface interface {
		FindReportFees(feerole string) ([]ReportFee ,error)
		FindReport(startdate,enddate,feerole string) ([]map[string]interface{},[]ReportFee,error)
	}
)

func NewReportRepo(oracle_db *sqlx.DB, logger *logrus.Logger) ReportRepoInterface {
	return &reportRepoDB{oracle_db: oracle_db, logger : logger}
}