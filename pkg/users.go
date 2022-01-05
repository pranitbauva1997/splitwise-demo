package pkg

import (
	"log"
	"net/http"
)

func SignUpPost(w http.ResponseWriter, r *http.Request, app *Application) {
	firstName := r.FormValue("fname")
	lastName := r.FormValue("lname")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordRepeat := r.FormValue("password_repeat")

	logRequest := func() {
		log.Print("firstName:", firstName, ";", "lastName:", lastName, ";", "username:", username, ";",
			"email:", email, ";", "username:", username, ";", "password:", password, ";", "passwordRepeat:", passwordRepeat)
	}

	badRequestErrorHandle := func(message string) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(message))
		if err != nil {
			log.Println("error while writing message to response body:", err)
		}
	}

	internalServerErrorHandle := func() {
		logRequest()
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		if err != nil {
			log.Println("error while writing message to response body:", err)
		}
	}

	// Validate if we want to create a user
	if !isSignUpInputValid(firstName, lastName, username, email, password, passwordRepeat) {
		logRequest()
		badRequestErrorHandle("Bad Request")
		return
	}

	isUsernameAvailable, err := app.StorageClient.IsUsernameAvailable(username)
	if err != nil {
		logRequest()
		log.Println("couldn't query db to check if username is available:", err)
		internalServerErrorHandle()
		return
	}

	if !isUsernameAvailable {
		logRequest()
		log.Println(username, "is already taken")
		badRequestErrorHandle("Username already taken")
		return
	}

	isEmailAvailable, err := app.StorageClient.IsEmailAvailable(email)
	if err != nil {
		logRequest()
		log.Println("couldn't query db to check if email is available:", err)
		internalServerErrorHandle()
		return
	}

	if !isEmailAvailable {
		logRequest()
		log.Println(email, "is already taken")
		badRequestErrorHandle("Email already taken")
		return
	}

	err = app.StorageClient.InsertUser(firstName, lastName, username, email)
	if err != nil {
		logRequest()
		log.Println("couldn't insert user in db:", err)
		internalServerErrorHandle()
		return
	}

	// TODO: Show a response
}
