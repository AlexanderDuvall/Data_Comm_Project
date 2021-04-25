package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var db *sql.DB
var connErr error

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func openConnection() {
	db, connErr = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/rigby")
	if connErr != nil {
		fmt.Println(connErr.Error())
		panic(connErr.Error())
	}
}
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	mac := r.Form.Get("Mac")
	rand.Seed(time.Now().UnixNano())
	auth := RandStringRunes(20)
	rand.Seed(time.Now().UnixNano())
	physicalHash := RandStringRunes(25)
	addUserRoute(email, password, mac, auth, physicalHash)
	fmt.Fprintf(w,"Your new Authentication key is: %v \n Make note of this. If lost you will not be able to access your files.",auth)
}
func addUserRoute(email, password, mac, auth, physicalHash string) {
	insertUser, err := db.Prepare("INSERT INTO users(email,pw,mac,authKey,physicalHash)(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertUser.Exec(email, password, mac, auth, physicalHash)
	if err != nil {
		fmt.Println(err)
	}
}
func addFile(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	name := r.Form.Get("name")
	file := r.Form.Get("file")
	id := r.Form.Get("user_id")
	insert, err := db.Prepare("INSERT INTO folders(fileName,files,user_id) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	insert.Exec(name, file, id)
}
func getFile(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	file_id := r.Form.Get("file_id")
	user_id := r.Form.Get("user_id")
	//getFile//
	addAccess(user_id, file_id)
}
func addAccess(user_id, file_id string) {
	insert, err := db.Prepare("INSERT INTO accesses(user_id,file_id) VALUES(?,?)")
	if err != nil {
		fmt.Println(err)
	} else {
		insert.Exec(user_id, file_id)
	}
}
func GetId(table, identifier, value string) string {
	openConnection()
	defer db.Close()
	command := fmt.Sprintf("select id from %s where %s = \"%s\"", table, identifier, value)
	fmt.Println(command)
	rows, err := db.Query(command)
	var id = "-1"
	if err == nil {
		for rows.Next() {
			var temp string
			if err := rows.Scan(&temp); err != nil {
				log.Fatal(err)
			}
			fmt.Println("....why the issue: " + temp)
			id = temp
		}
	} else {
		log.Fatal(err)
	}
	return id
}

func main() {
http.HandleFunc("/addUser",addUser)
http.HandleFunc("/addfile",addFile)
http.HandleFunc("/getfile",getFile)
}
