package handlers

import (
	"receipt-backend/nahmkahw/services"
	"receipt-backend/nahmkahw/errs"
	"net/http"

	"github.com/labstack/echo"
)

type (
	FacultyHandlers struct {
		facultyServices services.FacultyServiceInterface
	}
)

func NewFacultyHandlers(facultyServices services.FacultyServiceInterface) FacultyHandlers {
	return FacultyHandlers{facultyServices: facultyServices}
}

func (h *FacultyHandlers) GetFacultys(c echo.Context) error {

	id := c.Param("id")

	if id == "" {
		c.Logger().Error("ไม่พบข้อมูลรหัสนักศึกษา.")
		return c.JSON(http.StatusInternalServerError, errs.NewMessageAndStatusCode(http.StatusInternalServerError,"ไม่พบข้อมูลรหัสนักศึกษา."))
	}

	facultyFeesResponse, err := h.facultyServices.GetFaculty(id)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, errs.NewMessageAndStatusCode(http.StatusInternalServerError,err.Error()))
	}

	return c.JSON(http.StatusOK, facultyFeesResponse)

}