package handlers

import (
	"receipt-backend/nahmkahw/services"
	"receipt-backend/nahmkahw/errs"
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
		return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
	}

	if err := c.Validate(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
	}

	reportFeesResponse, err := h.reportServices.GetReportFees(RequestBody)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, errs.NewMessageAndStatusCode(http.StatusInternalServerError,err.Error()))
	}

	return c.JSON(http.StatusOK, reportFeesResponse)

}


func (h *ReportHandlers) GetReport(c echo.Context) error {
	RequestBody := new(services.ReportRequest)

	if err := c.Bind(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
	}

	if err := c.Validate(RequestBody); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, errs.NewBadRequestError())
	}

	reportResponse, err := h.reportServices.GetReport(RequestBody)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, errs.NewMessageAndStatusCode(http.StatusInternalServerError,err.Error()))
	}

	return c.JSON(http.StatusOK, reportResponse)

}

func (h *ReportHandlers) GetReportSummary(c echo.Context) error {

	reportResponse, err := h.reportServices.GetReportSummary()

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, errs.NewMessageAndStatusCode(http.StatusInternalServerError,err.Error()))
	}

	return c.JSON(http.StatusOK, reportResponse)

}