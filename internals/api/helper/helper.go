package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/slot"
	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/utilities"
)

type Faculty struct {
	Faculty_id   string `json:"id"`
	Faculty_name string `json:"name"`
}

type Slots struct {
	Slot_id  string  `json:"id"`
	Timmings []Times `json:"timings"`
}
type Times struct {
	Day   string `json:"day"`
	Start string `json:"start"`
	End   string `json:"end"`
}
type Courses struct {
	Course_id   string   `json:"id"`
	Course_name string   `json:"name"`
	Slot_ids    []string `json:"slot_ids"`
	Faculty_ids []string `json:"faculty_ids"`
	Course_type string   `json:"course_type"`
}

func GetFacuties(f []string) []Faculty {
	query := "SELECT * FROM faculties WHERE "
	for i := 0; i < len(f); i++ {
		query += "faculty_id = $" + strconv.Itoa(i+1) + " OR "
	}
	db := utilities.InitDB()
	prepare_query, err := db.Prepare(strings.TrimSuffix(query, " OR "))
	utilities.HandleError(err)
	var args []interface{}
	for _, value := range f {
		args = append(args, value)
	}
	rows, err := prepare_query.Query(args...)
	result := make([]Faculty, 0)
	utilities.HandleError(err)

	for rows.Next() {
		var temp Faculty
		if err := rows.Scan(&temp.Faculty_id, &temp.Faculty_name); err != nil {
			fmt.Println(err)
		}
		result = append(result, temp)
	}
	return result
}

func Faculties(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	var course Courses
	_ = json.NewDecoder(r.Body).Decode(&course)
	response_data := make(map[string]interface{})
	response_data["darpan"] = slot.GetSlots(course.Slot_ids)
	json.NewEncoder(w).Encode(response_data)
}
