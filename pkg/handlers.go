package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

func homeGet(w http.ResponseWriter, _ *http.Request, _ *Application) {
	type Status struct {
		Message string `json:"message"`
	}
	buf, err := json.Marshal(Status{Message: "Server is up"})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set(ContentType, ContentType_ApplicationJson)
	_, _= w.Write(buf)
}

func home(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Log.info.Println("calling the home route")
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		} else {
			homeGet(w, r, app)
		}
	}
}

func dashboardGet(w http.ResponseWriter, r *http.Request, app *Application) {
	// TODO: Get user_id and show the summary
}

func dashboard(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err := w.Write([]byte("Method not allowed"))
			if err != nil {
				log.Println("error while writing the message to response body:", err)
			}
			return
		} else {
			dashboardGet(w, r, app)
		}
	}
}

func signUp(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			SignUpPost(w, r, app)
		} else {
			w.Header().Set("Allow", "POST")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

func allUsers(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			allUsersGet(w, r, app)
		} else {
			w.Header().Set("Allow", "GET")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

func addBill(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addBillPost(w, r, app)
		} else if r.Method == http.MethodOptions {
			w.Header().Set("Allow", "POST")
			http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
			return
		} else {
			w.Header().Set("Allow", "GET")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}