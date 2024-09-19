package services

import (
	"receipt-backend/nahmkahw/errs"
	"encoding/json"
	"net/http"
	"fmt"
	"time"

)


func (g *facultyServices) GetFaculty(std_code string) (*[]FacultyResponse, error) {

	facultyResponse := []FacultyResponse{}

	key := "faculty::" + std_code
	facultyCache, err := g.redis_cache.Get(key).Result()
	if err == nil {
		_ = json.Unmarshal([]byte(facultyCache), &facultyResponse)
		fmt.Println("cache:" + key)
		return &facultyResponse, nil
	}

	fmt.Println("database:" + key)

	facultysRepo , err := g.facultyRepo.FindFaculty(std_code)

	if err != nil {
		return &facultyResponse, err
	}

	// Map เพื่อเก็บข้อมูล Faculty โดยจัดกลุ่ม
	studentMap := make(map[string]map[string]*Faculty)

	for _, item := range facultysRepo {
		STD_CODE := item.STD_CODE
		FACULTY_NO := item.FACULTY_NO

		// ตรวจสอบว่า STD_CODE มีอยู่ใน map หรือไม่
		if studentMap[STD_CODE] == nil {
			studentMap[STD_CODE] = make(map[string]*Faculty)
		}

		// ตรวจสอบว่า faculty มีอยู่แล้วหรือไม่ ถ้าไม่มีก็สร้างใหม่
		if _, ok := studentMap[STD_CODE][FACULTY_NO]; !ok {
			studentMap[STD_CODE][FACULTY_NO] = &Faculty{
				FACULTY_NO:   item.FACULTY_NO,
				FACULTY_NAME: item.FACULTY_NAME,
				Majors:      []Major{},
			}
		}

		// เพิ่ม major เข้าไปใน faculty ที่ตรงกับ FACULTY_NO
		studentMap[STD_CODE][FACULTY_NO].Majors = append(studentMap[STD_CODE][FACULTY_NO].Majors, Major{
			CURR_NO : item.CURR_NO,
			CURR_NAME: item.CURR_NAME,
			MAJOR_NO:   item.MAJOR_NO,
			MAJOR_NAME: item.MAJOR_NAME,
		})
	}

	
		// สร้าง array สำหรับ JSON ผลลัพธ์

		for STD_CODE, faculties := range studentMap {
			var facs []Faculty
			for _, fac := range faculties {
				facs = append(facs, *fac)
			}
			facultyResponse = append(facultyResponse, FacultyResponse{
				STD_CODE:   STD_CODE,
				Faculties: facs,
			})
		}

	if len(facultyResponse) < 1 {
		errStr := fmt.Sprintf("ไม่พบข้อมูลคณะและสาขาวิชา ของ %s",std_code)
		return &facultyResponse, errs.NewMessageAndStatusCode(http.StatusNotFound,errStr)
	}

	facultysJSON, _ := json.Marshal(&facultyResponse)
	timeNow := time.Now()
	redisCachefaculty := time.Unix(timeNow.Add(time.Second * 5).Unix(), 0)
	_ = g.redis_cache.Set(key, facultysJSON, redisCachefaculty.Sub(timeNow)).Err()

	return &facultyResponse, nil
}