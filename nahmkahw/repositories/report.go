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

	ReportSummary struct {
		SUMMARY_CODE string `db:"SUMMARY_CODE"` 
		SUMMARY_NAME string `db:"SUMMARY_NAME"`
		SUMMARY_COUNT   string `db:"SUMMARY_COUNT"`
		MAX_CODE string `db:"MAX_CODE"` 
		MAX_NAME string `db:"MAX_NAME"`
		MAX_COUNT   string `db:"MAX_COUNT"`
		MIN_CODE string `db:"MIN_CODE"` 
		MIN_NAME string `db:"MIN_NAME"`
		MIN_COUNT   string `db:"MIN_COUNT"`
	}

	ReportSuccess struct {
		SUCCESS_TODAY string `db:"SUCCESS_TODAY"` 
		SUCCESS_THIS_MONTH string `db:"SUCCESS_THIS_MONTH"`
		SUCCESS_THIS_YEAR   string `db:"SUCCESS_THIS_YEAR"`
	}

	ReportCancel struct {
		CANCEL_TODAY string `db:"CANCEL_TODAY"` 
		CANCEL_THIS_MONTH string `db:"CANCEL_THIS_MONTH"`
		CANCEL_THIS_YEAR   string `db:"CANCEL_THIS_YEAR"`
	}


	ReportRepoInterface interface {
		FindReportFees(feerole string) ([]ReportFee ,error)
		FindReport(startdate,enddate,feerole string) ([]map[string]interface{},[]ReportFee,error)

		FindReportSummary() (*ReportSummary ,error)
		FindReportSuccess() (*ReportSuccess ,error)
		FindReportCancel() (*ReportCancel ,error)
	}
)

func NewReportRepo(oracle_db *sqlx.DB, logger *logrus.Logger) ReportRepoInterface {
	return &reportRepoDB{oracle_db: oracle_db, logger : logger}
}