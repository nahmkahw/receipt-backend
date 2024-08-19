package handlers

import (
	"receipt-backend/nahmkahw/services"
	"net/http"

	"github.com/labstack/echo"
)

type (
	ReportHandlers struct {
		reportServices services.ReportServiceInterface
	}
)

func NewReportHandlers(reportServices services.ReportServiceInterface) ReportHandlers {
	return ReportHandlers{reportServices: reportServices}
}

func (h *ReportHandlers) GetReportFees(c echo.Context) error {
	RequestBody := new(services.ReportFeeRequest)

	if err := c.Bind(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	if err := c.Validate(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	reportFeesResponse, err := h.reportServices.GetReportFees(RequestBody)

	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	return c.JSON(http.StatusOK, reportFeesResponse)

}


func (h *ReportHandlers) GetReport(c echo.Context) error {
	RequestBody := new(services.ReportRequest)

	if err := c.Bind(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	if err := c.Validate(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	reportResponse, err := h.reportServices.GetReport(RequestBody)

	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	return c.JSON(http.StatusOK, reportResponse)

}