package main

import (
	"log"
	"net/http"

	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/course"
	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/faculty"
	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/helper"
	"github.com/dyte-submissions/vit-hiring-2023-phase-1-SkullCreek/internals/api/slot"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/admin/slot", slot.SlotHandler).Methods("POST")
	r.HandleFunc("/admin/faculty", faculty.FacultyHandler).Methods("POST")
	r.HandleFunc("/admin/course", course.CourseHandler).Methods("POST")
	r.HandleFunc("/faculty", helper.Faculties).Methods("POST")

	handler := HandleCORS(r)

	log.Fatal(http.ListenAndServe(":3000", handler))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to FFCS API</h1><p>Made by Darpan</p>"))
}

func HandleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
