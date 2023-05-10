package student

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>sada to API by darpan</h1>"))
}
