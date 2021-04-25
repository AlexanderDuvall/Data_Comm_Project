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
	fmt.Println(email)
	password := r.Form.Get("password")
	mac := r.Form.Get("Mac")
	rand.Seed(time.Now().UnixNano())
	auth := RandStringRunes(20)
	rand.Seed(time.Now().UnixNano())
	physicalHash := RandStringRunes(25)
	addUserRoute(email, password, mac, auth, physicalHash)
	fmt.Fprintf(w, "Auth:%v\n", auth)
	fmt.Fprintf(w, "PhysicalHash:%v\n", physicalHash)
}
func addUserRoute(email, password, mac, auth, physicalHash string) {
	openConnection()
	defer db.Close()
	s := fmt.Sprintf("INSERT INTO users(email,password,mac,authKey,physicalHash) VALUES (\"%v\",md5(\"%v\"),md5(\"%v\"),md5(\"%v\"),md5(\"%v\"))", email, password, mac, auth, physicalHash)
	fmt.Println(s)
	insertUser, err := db.Prepare(s)
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertUser.Exec()
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
	openConnection()
	defer db.Close()
	fileName := r.Form.Get("fileName")
	var isAuth bool = false
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	mac := r.Form.Get("Mac")
	authKey := r.Form.Get("auth")
	physicalHash := r.Form.Get("physHash")
	err := db.QueryRow("Select IF(COUNT(*),'true','false') from users where email = ? and password = md5(?) and mac = md5(?) and authkey = md5(?) and physicalHash = md5(?) ", email, password, mac, authKey, physicalHash).Scan(&isAuth)
	if err != nil {
		panic(err.Error())
	}
	if isAuth {
		fmt.Fprintf(w, "1")
		user_id := GetId("users", "email", email)
		retrieveFile(user_id, fileName)
		//addAccess(user_id, GetId("folder", "fileName", fileName))
	} else {
		fmt.Fprintf(w, "0")
	}
}

func addAccess(user_id, file_id string) {
	insert, err := db.Prepare("INSERT INTO accesses(user_id,file_id) VALUES(?,?)")
	if err != nil {
		fmt.Println(err)
	}
	insert.Exec(user_id, file_id)
}
func GetId(table, identifier, value string) string {
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
func retrieveFile(user_id, value string) bool {
	fmt.Println("auth passed...")
	command := fmt.Sprintf("Select IF(COUNT(*),'true','false') from folder where fileName = \"%s\" and user_id = \"%s\" ", value, user_id)
	var isAuth bool = false
	err := db.QueryRow(command).Scan(&isAuth)
	if err != nil {
		panic(err)
	}
	addAccess(user_id, GetId("folder", "fileName", value))
	fmt.Println(command)
	return isAuth
}

func main() {
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/addfile", addFile)
	http.HandleFunc("/getfile", getFile)
	http.ListenAndServe("127.0.0.1:3000", nil)
}
