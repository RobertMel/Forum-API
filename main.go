package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db *sql.DB

type Test struct {
	ID       json.Number `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

type Session struct {
	User_id             json.Number `json:"id"`
	Name                string      `json:"pseudo"`
	Token               string      `jspn:"token"`
	ExpirationTimeStamp time.Time   `json:"expirationtimestamp"`
}

type Topic struct {
	User_id       json.Number `json:"id"`
	Name          string      `json:"pseudo"`
	Sujet         string      `json:"sujet"`
	Content       string      `json:"content"`
	Commentaires  string      `json:"commentaires"`
	DateTimeStamp time.Time   `json:"datetimestamp"`
}

type User struct {
	User_id         int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmation"`
}

type Response struct {
	User_id       json.Number `json:"id"`
	Name          string      `json:"name"`
	TopicID       string      `json:"topicid"`
	Content       string      `json:"content"`
	DateTimeStamp time.Time   `json:"datetimestamp"`
}
type TokenResponse struct {
	Token    string `json:"token"`
	UserName string `json:"username"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var user User
	dec := json.NewDecoder(r.Body)
	dec.Decode(&user)

	query := fmt.Sprintf("SELECT id, user_pseudo FROM user WHERE Email = '%s' AND WHERE user_mdp = '%s' ", user.Email, user.Password)

	rows := db.QueryRow(query)

	rows.Scan(&user.User_id, &user.Name)

	if user.User_id != 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	token := uuid.New().String()

	querysession := fmt.Sprintf("INSERT INTO session (user_id, token) VALUES ('%d','%s'  ", user.User_id, token)

	db.QueryRow(querysession)

	response, _ := json.Marshal(TokenResponse{Token: token, UserName: user.Name})
	w.Write(response)

}

func Registerhandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	var reg User

	if reg.Password == reg.ConfirmPassword {
		dec := json.NewDecoder(r.Body)
		dec.Decode(&reg)
		query := fmt.Sprintf("INSERT INTO user (user_pseudo, Email, user_mdp) VALUES ('%s','%s','%s' ", reg.Name, reg.Email, reg.Password)
		db.Exec(query)
		w.WriteHeader(http.StatusCreated)
	}
	w.WriteHeader(http.StatusBadRequest)
}

/*
func topicshandler(w http.ResponseWriter, r *http.Request) { //  Sujet

	fmt.Fprintf(w, "2e page = %q\n", r.URL.Path)
}
func topichandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "4e page = %q\n", r.URL.Path)
}
*/

/*func selectTestHandler(w http.ResponseWriter, r *http.Request) {
query := "SELECT id, name, email, password FROM test_user"
rows, _ := db.Query(query)

var tests []Test
for rows.Next() {
	var test Test
	rows.Scan(&test.ID, &test.Name, &test.Email, &test.Password)
	tests = append(tests, test)
}

response, _ := json.Marshal(tests)
w.Write(response)
*/

func main() {
	db, _ = sql.Open("mysql", "root:root@tcp(localhost:8889)/DataBaseFRM")

	// http.handler

	http.HandleFunc("/api/login", loginHandler) // each request calls handler
	//http.HandleFunc("/api/topics", topicshandler)
	//http.HandleFunc("/api/topic/", topichandler)
	http.HandleFunc("/api/register", Registerhandler)
	//http.HandleFunc("/api/test/select", selectTestHandler)

	http.ListenAndServe(":8084", nil)
}
