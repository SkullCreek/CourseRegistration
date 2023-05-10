package faculty

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/utilities"
)

type Faculties struct {
	Faculty_Id   string `json:"id"`
	Faculty_Name string `json:"name"`
}

func FacultyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	auth := r.Header.Get("Authorization")
	if utilities.IsTokenValid(auth) {
		if r.Body == nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		var faculty Faculties
		_ = json.NewDecoder(r.Body).Decode(&faculty)
		if faculty.IsEmpty() {
			http.Error(w, "Some data missing", http.StatusBadRequest)
			return
		}
		db := utilities.InitDB()
		response_data := make(map[string]interface{})
		if AddFaculty(db, faculty) == "success" {
			response_data["success"] = true
			response_data["data"] = faculty
			db.Close()
			json.NewEncoder(w).Encode(response_data)
		} else {
			response_data["success"] = false
			response_data["data"] = "Faculty Already Exists"
			db.Close()
			json.NewEncoder(w).Encode(response_data)
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func AddFaculty(db *sql.DB, faculty Faculties) string {
	if IsValidFaculty(db, faculty) {
		prepare_query, err := db.Prepare("INSERT INTO faculties VALUES($1, $2)")
		utilities.HandleError(err)
		response, err := prepare_query.Exec(faculty.Faculty_Id, faculty.Faculty_Name)
		utilities.HandleError(err)
		numRows, err := response.RowsAffected()
		utilities.HandleError(err)
		if numRows > 0 {
			return "success"
		} else {
			return "failure"
		}
	} else {
		return "error"
	}
}

func IsValidFaculty(db *sql.DB, faculty Faculties) bool {
	prepare_query, err := db.Prepare("SELECT * FROM faculties WHERE faculty_id = $1")
	utilities.HandleError(err)
	response, err := prepare_query.Exec(faculty.Faculty_Id)
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
func (f Faculties) IsEmpty() bool {
	return f.Faculty_Id == ""
}
