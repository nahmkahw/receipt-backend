package backend

import (
	"context"
	//"database/sql"
	"net/http"
	//"strconv"
	//"fmt"
	"time"
	//"encoding/json"

	"github.com/labstack/echo"
)

type (
	Log struct {
		Id       string `json:"Id"`
		Code     int64  `json:"Code" validate:"required"`
		Module   string `json:"Module" validate:"required"`
		Task     string `json:"Task" validate:"required"`
		Username string `json:"Username" validate:"required"`
		Created  string `json:"Created"`
		Modified string `json:"Modified"`
	}

	LogForm struct {
		Code   int64  `json:"Code" validate:"required"`
		Module string `json:"Module" validate:"required"`
	}
)

func (h *backendRepoDB) CreateLogs(c echo.Context) error {

	form := new(Log)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		logs []Log
		log  Log
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	tx, err := h.oracle_db.BeginTx(ctx, nil)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = tx.ExecContext(ctx, "insert into fees_logs (ID,CODE,MODULE,TASK,USERNAME,CREATED,MODIFIED) values (FEES_LOGS_SEQ.NEXTVAL,:1,:2,:3,:4,sysdate,sysdate)", form.Code, form.Module, form.Task, form.Username)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = tx.Commit()

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sql := `select * from fees_logs where code = :1 and module = :2 order by id desc `

	rows, err := h.oracle_db.Query(sql, form.Code, form.Module)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&log.Id, &log.Code, &log.Module, &log.Task, &log.Username,
			&log.Created, &log.Modified)

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(logs) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูล."})
	}

	return c.JSON(http.StatusOK, logs)

}

func (h *backendRepoDB) FindLogs(c echo.Context) error {

	form := new(LogForm)

	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		logs []Log
		log  Log
	)

	sql := `select * from fees_logs where code = :1 and module = :2 order by id desc `

	rows, err := h.oracle_db.Query(sql, form.Code, form.Module)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&log.Id, &log.Code, &log.Module, &log.Task, &log.Username,
			&log.Created, &log.Modified)

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(logs) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ไม่พบข้อมูล."})
	}

	return c.JSON(http.StatusOK, logs)

}
