package frontend

import (
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/godror/godror"
)

type Student struct {
	StudentCode   string `json:"StudentCode" validate:"required"`
	NameThai      string `json:"NameThai" validate:"required`
	NameEng       string `json:"NameEng" validate:"required"`
	CitizenNo     string `json:"CitizenNo" validate:"required"`
	BirthDate     string `json:"BirthDate" validate:"required"`
	FacultyName   string `json:"FacultyName" validate:"required`
	CurrName      string `json:"CurrName" validate:"required"`
	MajorNameThai string `json:"MajorNameThai" validate:"required"`
	MajorNameEng  string `json:"MajorNameEng" validate:"required"`
	StatusCurrent string `json:"StatusCurrent" validate:"required"`
}

func (h *frontendRepoDB) FindStudentId(c echo.Context) error {

	id := c.Param("id")

	sql := `select STD_CODE ,NAME_THAI ,(FIRST_NAME_ENG||' '||MIDDLE_NAME_ENG||' '||LAST_NAME_ENG) NAME_ENG ,
				CITIZEN_NO ,BIRTH_DATE ,FACULTY_NAME_THAI ,NVL(CURR_NAME_THAI,'-') CURR_NAME_THAI,NVL(MAJOR_NAME_THAI,'-') MAJOR_NAME_THAI,NVL(MAJOR_NAME_ENG,'-') MAJOR_NAME_ENG,STD_STATUS_CURRENT 
				from DBBACH00.VM_STUDENT_MOBILE where STD_CODE = :1`

	var (
		student Student
	)

	row := h.oracle_db.QueryRow(sql, id)
	err := row.Scan(&student.StudentCode, &student.NameThai, &student.NameEng, &student.CitizenNo,
		&student.BirthDate, &student.FacultyName, &student.CurrName, &student.MajorNameThai, &student.MajorNameEng, &student.StatusCurrent)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.Logger().Info("frontend-order")
	return c.JSON(http.StatusOK, student)
}
