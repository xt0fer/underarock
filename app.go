package underarock

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *Driver
}

type User struct {
	UserId string `json:"userid"`
	Name   string `json:"name"`
	Github string `json:"github"`
}

type Message struct {
	Sequence  string    `json:"sequence"`
	Timestamp time.Time `json:"timestamp"`
	FromId    string    `json:"fromid"`
	ToId      string    `json:"toid"`
	Message   string    `json:"message"`
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(accountEmail string) {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.DB, _ = New(pwd, nil)
	a.Router = mux.NewRouter()
	a.setRouters()

	a.LoadAllMessages()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// #### /ids
	// * `GET` : Get all github ids registered
	// * `POST` : add your github id / name to be registered
	// * `PUT` : change the name linked to your github id

	// json payload for /ids/
	// ```json
	// {
	//     "userid": "-", // gets filled w id
	//     "name": "Kris",
	//     "githubid": "xt0fer"
	// }
	// ```

	a.Get("/ids", a.GetAllUsers)
	a.Post("/ids", a.PostNewUser)
	a.Put("/ids", a.PutUser)

	a.Get("/ids/{myid}", a.GetUserId)

	// #### /messages/
	// * `GET` : Get last 20 msgs - returns an JSON array of message objects
	a.Get("/messages", a.GetAllMessages)

	// #### /ids/:myid/messages/
	// * `GET` : Get last 20 msgs for myid  - returns an JSON array of message objects
	// * `POST` : Create a new message in timeline - need to POST a new message object, and will get back one with a message sequence number
	a.Get("/ids/{myid}/messages", a.GetMyMessages)
	a.Post("/ids/{myid}/messages", a.PostNewMessage)

	// #### /ids/:myid/messages/:sequence
	// * `GET` : Get msg with a sequence  - returns a JSON message object for a sequence number
	a.Get("/ids/{myid}/messages/{sequence}", a.GetAMessage)

	// #### /ids/:myid/from/:friendid
	// * `GET` : Get last 20 msgs for myid from friendid
	// * `POST` : Create a new message in timeline
	a.Get("/ids/{myid}/from/{friendid}", a.GetFriendMessages)
	//a.Post("/ids/{myid}/messages/{friendid}", a.PostFriendMessage)

}

func makeUserHash(foo *User) string {
	s := foo.Name + foo.Github
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func makeMessageHash(foo *Message) string {
	s := foo.FromId + foo.ToId + foo.Message + foo.Timestamp.String()
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func (db *Driver) FetchUser(id string) (*User, error) {
	onefish := User{}
	if err := db.Read("user", id, &onefish); err != nil {
		fmt.Println("Error", err)
	}
	return &onefish, nil
}

// GetUserByGithub returns the first User found with the given GithubId
func (db *Driver) GetUserByGithub(id string) (string, error) {
	records, err := db.ReadAll("user")
	if err != nil {
		return "", err
	}

	for _, f := range records {
		fishFound := User{}
		if err := json.Unmarshal([]byte(f), &fishFound); err != nil {
			return "", err
		}

		if fishFound.Github == id {
			return fishFound.UserId, nil
		}
	}
	return "", errors.New("GithubID not found")
}

func (db *Driver) AllUsers() (*[]User, error) {
	records, err := db.ReadAll("user")
	if err != nil {
		fmt.Println("Error", err)
	}

	fishies := []User{}
	for _, f := range records {
		fishFound := User{}
		if err := json.Unmarshal([]byte(f), &fishFound); err != nil {
			fmt.Println("Error", err)
		}
		fishies = append(fishies, fishFound)
	}

	return &fishies, nil
}
func (db *Driver) NewUser(onefish *User) (*User, error) {
	onefish.UserId = makeUserHash(onefish)
	if err := db.Write("user", onefish.UserId, onefish); err != nil {
		fmt.Println("Error", err)
	}
	return onefish, nil
}

func (db *Driver) NewMessage(onefish *Message) (*Message, error) {
	onefish.Timestamp = time.Now()
	onefish.Sequence = makeMessageHash(onefish)
	if err := db.Write("message", onefish.Sequence, onefish); err != nil {
		fmt.Println("Error", err)
	}
	return onefish, nil
}

func (db *Driver) FetchMessage(id string) (*Message, error) {
	onefish := Message{}
	if err := db.Read("message", id, &onefish); err != nil {
		fmt.Println("Error", err)
	}
	return &onefish, nil
}
func (db *Driver) AllMessages() ([]Message, error) {
	records, err := db.ReadAll("message")
	if err != nil {
		fmt.Println("Error", err)
	}

	fishies := []Message{}
	for _, f := range records {
		fishFound := Message{}
		if err := json.Unmarshal([]byte(f), &fishFound); err != nil {
			fmt.Println("Error", err)
		}
		fishies = append(fishies, fishFound)
	}

	return fishies, nil
}

func (a *App) LoadAllMessages() {
	mm, _ := a.DB.AllMessages()
	for _, m := range mm {
		AddMessage(&m)
	}
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

/*
** Project Handlers
** Map requests to handlers.
 */
// patterned on https://github.com/mingrammer/go-todo-rest-api-example

func (a *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	GetAllUsers(a.DB, w, r)
}

func (a *App) PostNewUser(w http.ResponseWriter, r *http.Request) {
	PostNewUser(a.DB, w, r)
}

func (a *App) PutUser(w http.ResponseWriter, r *http.Request) {
	PutUser(a.DB, w, r)
}

// GetUserByGithub
func (a *App) GetUserId(w http.ResponseWriter, r *http.Request) {
	GetUserId(a.DB, w, r)
}

// GetAllMessages
func (a *App) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	GetAllMessages(a.DB, w, r)
}

// GetMyMessages
func (a *App) GetMyMessages(w http.ResponseWriter, r *http.Request) {
	GetMyMessages(a.DB, w, r)
}

// PostNewMessage
func (a *App) PostNewMessage(w http.ResponseWriter, r *http.Request) {
	PostNewMessage(a.DB, w, r)
}

// GetAMessage
func (a *App) GetAMessage(w http.ResponseWriter, r *http.Request) {
	GetAMessage(a.DB, w, r)
}

// GetFriendMessages
func (a *App) GetFriendMessages(w http.ResponseWriter, r *http.Request) {
	GetFriendMessages(a.DB, w, r)
}

// Run the app on its router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
