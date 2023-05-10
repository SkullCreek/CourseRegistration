package slot

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/utilities"
)

type Slots struct {
	Slot_id  string  `json:"id"`
	Timmings []Times `json:"timings"`
}

type Times struct {
	Day   string `json:"day"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// Checks and stores slots in database
func SlotHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	auth := r.Header.Get("Authorization")

	if utilities.IsTokenValid(auth) {

		if r.Body == nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var slot Slots
		_ = json.NewDecoder(r.Body).Decode(&slot)

		if slot.IsEmpty() {
			http.Error(w, "Some data missing", http.StatusBadRequest)
			return
		}
		db := utilities.InitDB()
		response_data := make(map[string]interface{})
		if AddSlot(db, slot) == "success" {

			response_data["success"] = true
			response_data["data"] = slot
			db.Close()
			json.NewEncoder(w).Encode(response_data)
		} else {
			response_data["success"] = false
			response_data["data"] = "Slot Already Exists"
			db.Close()
			json.NewEncoder(w).Encode(response_data)
		}

	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

}

// Store Slot in Database
func AddSlot(db *sql.DB, slot Slots) string {
	if IsSlotValid(db, slot) {
		prepare_query, err := db.Prepare("INSERT INTO slots VALUES($1, $2)")
		utilities.HandleError(err)
		jsonString, err := json.Marshal(slot.Timmings)
		utilities.HandleError(err)
		response, err := prepare_query.Exec(slot.Slot_id, jsonString)
		utilities.HandleError(err)
		numRows, err := response.RowsAffected()
		utilities.HandleError(err)
		if numRows > 0 {
			return "success"
		} else {
			return "failure"
		}
	} else {
		return "Slot Already Exists"
	}

}

// Check if ID is valid
func IsSlotValid(db *sql.DB, slot Slots) bool {
	prepare_query, err := db.Prepare("SELECT * FROM slots WHERE slot_id = $1")
	utilities.HandleError(err)
	response, err := prepare_query.Exec(slot.Slot_id)
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

func GetSlots(s []string) []Slots {
	query := "SELECT * FROM slots WHERE "
	for i := 0; i < len(s); i++ {
		query += "slot_id = $" + strconv.Itoa(i+1) + " OR "
	}
	db := utilities.InitDB()
	prepare_query, err := db.Prepare(strings.TrimSuffix(query, " OR "))
	utilities.HandleError(err)
	var args []interface{}
	for _, value := range s {
		args = append(args, value)
	}
	rows, err := prepare_query.Query(args...)
	utilities.HandleError(err)
	dataList := make([]Slots, 0)
	for rows.Next() {
		var slot Slots
		var jsonString string
		if err := rows.Scan(&slot.Slot_id, &jsonString); err != nil {
			panic(err)
		}
		err := json.Unmarshal([]byte(jsonString), &slot.Timmings)
		utilities.HandleError(err)
		dataList = append(dataList, slot)
	}
	return dataList
}

func (s *Slots) IsEmpty() bool {
	return s.Slot_id == ""
}
