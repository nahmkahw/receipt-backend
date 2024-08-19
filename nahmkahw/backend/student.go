package backend

import (
	"fmt"
	"net/http"
	"time"

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
}

func (h *backendRepoDB) FindStudentId(c echo.Context) error {

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Max-Age", "86400")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")

	id := c.Param("id")

	sql := `select STD_CODE ,NAME_THAI ,(FIRST_NAME_ENG||' '||MIDDLE_NAME_ENG||' '||LAST_NAME_ENG) NAME_ENG ,
				CITIZEN_NO ,BIRTH_DATE ,
				DECODE(FACULTY_NAME_THAI,null,'-',FACULTY_NAME_THAI) FACULTY_NAME_THAI,
				DECODE(CURR_NAME_THAI,null,'-',CURR_NAME_THAI) CURR_NAME_THAI,
				DECODE(MAJOR_NAME_THAI,null,'-',MAJOR_NAME_THAI) MAJOR_NAME_THAI, 
				DECODE(MAJOR_NAME_ENG,null,'-',MAJOR_NAME_ENG) MAJOR_NAME_ENG 
				from DBBACH00.VM_STUDENT_MOBILE where STD_CODE = :1`

	var (
		student Student
		cache  = NewRedisCache(h.redis_cache,time.Second*20)
	)

	key := fmt.Sprintf("BACKEND-STUDENT-%s", id)

	studentcache := cache.GetStudent(key)

	if studentcache == nil {
		fmt.Printf("student backend oracle: %s", key)

		row := h.oracle_db.QueryRow(sql, id)
		err := row.Scan(&student.StudentCode, &student.NameThai, &student.NameEng, &student.CitizenNo,
			&student.BirthDate, &student.FacultyName, &student.CurrName, &student.MajorNameThai, &student.MajorNameEng)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		cache.SetStudent(key, &student)

		return c.JSON(http.StatusOK, student)

	}

	fmt.Printf("student backend redis: %s", key)

	return c.JSON(http.StatusOK, studentcache)
}
