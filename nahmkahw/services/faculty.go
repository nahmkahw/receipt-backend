package services

import (
	"receipt-backend/nahmkahw/repositories"
	"github.com/go-redis/redis/v7"
)

type (
	facultyServices struct {
		facultyRepo repositories.FacultyRepoInterface
		redis_cache  *redis.Client
	}

	Major struct {
		CURR_NO string `json:"curr_no"` 
		CURR_NAME string `json:"curr_name"`
		MAJOR_NO string `json:"major_no"` 
		MAJOR_NAME string `json:"major_name"`
	}
	
	Faculty struct {
		FACULTY_NO string `json:"faculty_no"`
		FACULTY_NAME   string `json:"faculty_name"`
		Majors      []Major `json:"majors"`
	}
	
	FacultyResponse struct {
		STD_CODE string `json:"STD_CODE"`
		Faculties []Faculty `json:"faculties"`
	}

	FacultyServiceInterface interface {
		GetFaculty(std_code string) (*[]FacultyResponse, error)
	}
)

func NewFacultyServices(facultyRepo repositories.FacultyRepoInterface, redis_cache *redis.Client) FacultyServiceInterface {
	return &facultyServices{
		facultyRepo: facultyRepo,
		redis_cache:  redis_cache,
	}
}