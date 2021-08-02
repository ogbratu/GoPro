package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	var u User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respond(w, http.StatusBadRequest, "Invalid login data")
		return
	}
	defer r.Body.Close()

	// Create read-only transaction
	txn := Db.Txn(false)
	defer txn.Abort()

	// Lookup by name
	var name interface{} = u.Name
	raw, err := txn.First("user", "name", name)
	if err != nil {
		panic(err)
	}

	if raw.(*User).Name == u.Name && raw.(*User).Password == u.Password {
		log.Printf("Logged in as user %s\n", u.Name)
		respond(w, http.StatusOK, "Login successful")
	} else {
		log.Printf("Invalid username / password for user %s\n", u.Name)
		respond(w, http.StatusUnauthorized, "Invalid login details")
	}
	u.Id = raw.(*User).Id
	u.Email = raw.(*User).Email
	Users = append(Users, u)
}

func logout(w http.ResponseWriter, r *http.Request)  {
	loggedInUser := "no one"
	if len(Users) > 0 {
		Users = Users[:len(Users) - 1 ]
		if len(Users) > 0 {
			loggedInUser = Users[len(Users) - 1].Name
		}
	} else {
		Users = nil
	}
	log.Printf("Logout successful, logged in user now is " + loggedInUser)
	respond(w, http.StatusOK, "Logout successful, logged in user now is " + loggedInUser)
}

func getLoggedInUser() int {
	return Users[len(Users) - 1].Id
}