package pkg

import (
	"encoding/json"
	"github.com/pranitbauva1997/splitwise-demo/pkg/store"
	"io"
	"net/http"
)

func addBillPost(w http.ResponseWriter, r *http.Request, app *Application) {
	defer r.Body.Close()
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var newBill store.AddBillRequest
	err = json.Unmarshal(buf, &newBill)
	if err != nil {
		app.Log.err.Printf("couldn't parse the request body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// TODO: validate Bill

	err = app.StorageClient.AddBill(newBill)
	if err != nil {
		app.Log.err.Printf("couldn't insert the bill in db: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type Status struct {
		Message string `json:"message"`
	}
	w.Header().Set(ContentType, ContentType_ApplicationJson)
	buf, _ = json.Marshal(Status{Message: "success"})
	_, _ = w.Write(buf)
}