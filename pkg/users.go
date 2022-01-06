package pkg

import (
	"encoding/json"
	"net/http"
)

func SignUpPost(w http.ResponseWriter, r *http.Request, app *Application) {
	firstName := r.FormValue("fname")
	lastName := r.FormValue("lname")
	username := r.FormValue("username")
	email := r.FormValue("email")

	badRequestErrorHandle := func() {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	internalServerErrorHandle := func() {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Validate if we want to create a user
	if !isSignUpInputValid(firstName, lastName, username, email) {
		badRequestErrorHandle()
		return
	}

	isUsernameAvailable, err := app.StorageClient.IsUsernameAvailable(username)
	if err != nil {
		app.Log.err.Println("couldn't query db to check if username is available:", err)
		internalServerErrorHandle()
		return
	}

	if !isUsernameAvailable {
		app.Log.err.Println(username, "is already taken")
		badRequestErrorHandle()
		return
	}

	isEmailAvailable, err := app.StorageClient.IsEmailAvailable(email)
	if err != nil {
		app.Log.err.Println("couldn't query db to check if email is available:", err)
		internalServerErrorHandle()
		return
	}

	if !isEmailAvailable {
		app.Log.err.Println(email, "is already taken")
		badRequestErrorHandle()
		return
	}

	err = app.StorageClient.InsertUser(firstName, lastName, username, email)
	if err != nil {
		app.Log.err.Println("couldn't insert user in db:", err)
		internalServerErrorHandle()
		return
	}

	type Status struct {
		Message string `json:"message"`
	}
	buf, _:= json.Marshal(Status{Message: "success"})
	w.Header().Set(ContentType, ContentType_ApplicationJson)
	_, _ = w.Write(buf)
}

func allUsersGet(w http.ResponseWriter, _ *http.Request, app *Application) {
	allUsers, err := app.StorageClient.GetAllUsers()
	if err != nil {
		app.Log.err.Printf("could't get all users from db: %s", err)
	}
	buf, _ := json.Marshal(allUsers)
	w.Header().Set(ContentType, ContentType_ApplicationJson)
	_, _ = w.Write(buf)
}
