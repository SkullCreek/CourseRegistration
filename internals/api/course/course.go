package course

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/helper"
	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/slot"
	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/utilities"
)

type Courses struct {
	Course_id   string   `json:"id"`
	Course_name string   `json:"name"`
	Slot_ids    []string `json:"slot_ids"`
	Faculty_ids []string `json:"faculty_ids"`
	Course_type string   `json:"course_type"`
}

func CourseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	auth := r.Header.Get("Authorization")

	if utilities.IsTokenValid(auth) {
		if r.Body == nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		var course Courses
		_ = json.NewDecoder(r.Body).Decode(&course)
		if course.IsEmpty() {
			http.Error(w, "Some data missing", http.StatusBadRequest)
			return
		}
		db := utilities.InitDB()
		response_data := make(map[string]interface{})
		if AddCourse(db, course) == "success" {

			temp := make(map[string]interface{})
			db.Close()
			temp["id"] = course.Course_id
			temp["name"] = course.Course_name
			temp["course_type"] = course.Course_type
			temp["faculties"] = helper.GetFacuties(course.Faculty_ids)
			temp["allowed_slots"] = slot.GetSlots(course.Slot_ids)
			response_data["success"] = true
			response_data["data"] = temp
			json.NewEncoder(w).Encode(response_data)
		} else {
			db.Close()
			response_data["success"] = false
			response_data["data"] = "Slot Already Exists"
			json.NewEncoder(w).Encode(response_data)
		}

	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func AddCourse(db *sql.DB, course Courses) string {
	if IsSlotValid(db, course) {
		prepare_query, err := db.Prepare("INSERT INTO courses VALUES($1, $2, $3)")
		utilities.HandleError(err)
		response, err := prepare_query.Exec(course.Course_id, course.Course_name, course.Course_type)
		utilities.HandleError(err)
		numRows, err := response.RowsAffected()
		utilities.HandleError(err)
		if numRows > 0 {
			if AddSlotCourse(db, course) == "success" {
				if AddFacultyCourse(db, course) == "success" {
					return "success"
				} else {
					return "failure"
				}
			} else {
				return "failure"
			}
		} else {
			return "failure"
		}
	} else {
		return "Slot Already Exists"
	}

}

func AddSlotCourse(db *sql.DB, course Courses) string {
	query := "INSERT INTO slots_course VALUES"
	for i := 0; i < len(course.Slot_ids); i++ {
		query += "($" + strconv.Itoa((i + 2)) + ",$1),"
	}
	prepare_query, err := db.Prepare(strings.TrimSuffix(query, ","))
	utilities.HandleError(err)
	var args []interface{}
	args = append(args, course.Course_id)
	for _, value := range course.Slot_ids {
		args = append(args, value)
	}
	response, err := prepare_query.Exec(args...)
	utilities.HandleError(err)
	numRows, err := response.RowsAffected()
	utilities.HandleError(err)
	if numRows > 0 {
		return "success"
	} else {
		return "failure"
	}
}

func AddFacultyCourse(db *sql.DB, course Courses) string {
	query := "INSERT INTO faculty_course VALUES"
	for i := 0; i < len(course.Faculty_ids); i++ {
		query += "($" + strconv.Itoa((i + 2)) + ",$1),"
	}
	prepare_query, err := db.Prepare(strings.TrimSuffix(query, ","))
	utilities.HandleError(err)
	var args []interface{}
	args = append(args, course.Course_id)
	for _, value := range course.Faculty_ids {
		args = append(args, value)
	}
	response, err := prepare_query.Exec(args...)
	utilities.HandleError(err)
	numRows, err := response.RowsAffected()
	utilities.HandleError(err)
	if numRows > 0 {
		return "success"
	} else {
		return "failure"
	}
}

func IsSlotValid(db *sql.DB, course Courses) bool {
	prepare_query, err := db.Prepare("SELECT * FROM courses WHERE course_id = $1")
	utilities.HandleError(err)
	response, err := prepare_query.Exec(course.Course_id)
	utilities.HandleError(err)
	numRows, err := response.RowsAffected()
	utilities.HandleError(err)
	if numRows > 0 {
		return false
	} else if numRows == 0 {
		return true
	} else {
		return false
	}
}

func (c *Courses) IsEmpty() bool {
	return c.Course_id == "" && c.Course_name == "" && c.Course_type == ""
}
