package underarock

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAllUsers
func GetAllUsers(db *Driver, w http.ResponseWriter, r *http.Request) {
	users, _ := db.AllUsers()
	respondJSON(w, http.StatusOK, users)
}

// PostNewUser
func PostNewUser(db *Driver, w http.ResponseWriter, r *http.Request) {
	newuser := &User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newuser); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	newuser, _ = db.NewUser(newuser)
	respondJSON(w, http.StatusOK, newuser)
}

// GetUserId returns the first entry found with the supplied githubid.
func GetUserId(db *Driver, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	myid := vars["myid"]

	userId, _ := db.GetUserByGithub(myid)
	respondJSON(w, http.StatusOK, userId)
}

// PutUser
func PutUser(db *Driver, w http.ResponseWriter, r *http.Request) {
	user := &User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Write("user", user.UserId, user); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// GetAllMessages
func GetAllMessages(db *Driver, w http.ResponseWriter, r *http.Request) {
	mm := Top20()
	respondJSON(w, http.StatusOK, mm)
}

// GetMyMessages
func GetMyMessages(db *Driver, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	myid := vars["myid"]
	mm := Top20For(myid)
	respondJSON(w, http.StatusOK, mm)
}

// PostNewMessage
func PostNewMessage(db *Driver, w http.ResponseWriter, r *http.Request) {
	newm := &Message{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newm); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	db.NewMessage(newm)
	AddMessage(newm)
	log.Println("adding", newm)
	respondJSON(w, http.StatusOK, newm)
}

// GetAMessage
func GetAMessage(db *Driver, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	msgid := vars["sequence"]
	m, err := db.FetchMessage(msgid)
	if err != nil {
		respondError(w, 501, "unable to find message")
	}
	respondJSON(w, http.StatusOK, m)
}

// GetFriendMessages
func GetFriendMessages(db *Driver, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	myid := vars["myid"]
	fid := vars["friendid"]
	mm := Top20From(myid, fid)
	respondJSON(w, http.StatusOK, mm)
}
