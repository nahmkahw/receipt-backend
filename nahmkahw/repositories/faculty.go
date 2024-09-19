package repositories

import (
	"github.com/jmoiron/sqlx"
    "github.com/sirupsen/logrus"
)

type (
	facultyRepoDB struct {
		oracle_db *sqlx.DB
		logger *logrus.Logger
		discordURL string
	}

	Faculty struct {
		STD_CODE string `db:"STD_CODE"`
		FACULTY_NO string `db:"FACULTY_NO"`
		FACULTY_NAME   string `db:"FACULTY_NAME"`
		CURR_NO string `db:"CURR_NO"` 
		CURR_NAME string `db:"CURR_NAME"`
		MAJOR_NO string `db:"MAJOR_NO"` 
		MAJOR_NAME string `db:"MAJOR_NAME"`
	}

	FacultyRepoInterface interface {
		FindFaculty(std_code string) ([]Faculty ,error)
	}
)

func NewFacultyRepo(oracle_db *sqlx.DB, logger *logrus.Logger) FacultyRepoInterface {
	return &facultyRepoDB{oracle_db: oracle_db, logger : logger}
}