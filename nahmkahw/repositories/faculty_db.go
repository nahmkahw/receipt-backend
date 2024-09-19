
package repositories

import (
	"receipt-backend/nahmkahw/util"
	"fmt"
    "github.com/sirupsen/logrus"
    "runtime"
)

func (r *facultyRepoDB) FindFaculty(std_code string) ([]Faculty ,error) {
	var (
		facultys []Faculty
		faculty  Faculty
	)

	sql := `select distinct STD_CODE,FAC_SEL FACULTY_NO,FAC_NAME_SEL FACULTY_NAME,CURR_SEL CURR_NO,CURR_NAME_SEL CURR_NAME,MAJOR_NO_1 MAJOR_NO,MAJOR_NAME_SEL MAJOR_NAME from DBBACH00.VM_FACULTY_TRANSFER where STD_CODE = :1 order by FAC_SEL,CURR_SEL,MAJOR_NO_1`

	rows, err := r.oracle_db.Query(sql, std_code)
	defer rows.Close()

	if err != nil {
		param := fmt.Sprintf("faculty:%s",std_code)
		r.logAndNotifyError(err,param)
        return nil, err
	}

	for rows.Next() {
		rows.Scan(&faculty.STD_CODE,&faculty.FACULTY_NO,&faculty.FACULTY_NAME,&faculty.CURR_NO,&faculty.CURR_NAME,&faculty.MAJOR_NO,&faculty.MAJOR_NAME)
        //fmt.Println(faculty)
		facultys = append(facultys,faculty)
	}

	return facultys, nil
}

func (r *facultyRepoDB) logAndNotifyError(err error,param string) {
    oraCode := extractORACode(err.Error())

    pc, file, line, ok := runtime.Caller(1)
    if !ok {
        r.logger.Error("Failed to retrieve caller information")
    }
    funcName := runtime.FuncForPC(pc).Name()

    r.logger.WithFields(logrus.Fields{
        "func_name": funcName,
        "file":      file,
        "line":      line,
        "error":     err.Error(),
        "ORA_CODE":  oraCode,
    }).Error("SQL Error")

    message := fmt.Sprintf("SQL Error in %s File: %s Line: %d ORA_CODE: %s Parameter: %s", funcName, file, line, oraCode, param)
    if err := util.SendToDiscord(message); err != nil {
        r.logger.Error("Failed to send message to Discord: ", err)
    }

	if err := util.SendToTeams(message); err != nil {
        r.logger.Error("Failed to send message to Discord: ", err)
    }

	
}