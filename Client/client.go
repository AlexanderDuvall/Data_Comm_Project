package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var path = "F:\\Data Comm\\auth.txt"

func AddUser(email, password, mac string) {
	postBody := url.Values{
		"email":    {email},
		"password": {password},
		"Mac":      {mac},
	}
	resp, err := http.PostForm("http://127.0.0.1:3000/addUser", postBody)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("%s\n", string(body))
		}
		array := strings.Split(string(body), "\n")
		auth := strings.Split(string(array[0]), ":")[1]
		phys := strings.Split(string(array[1]), ":")[1]
		fmt.Println(auth)
		fmt.Println(phys)
		buildAuthenticator(phys)
	}
}

func GetFile(email, password, auth, phys, mac, fileName string) {
	postBody := url.Values{
		"email":    {email},
		"password": {password},
		"auth":     {auth},
		"physHash": {phys},
		"Mac":      {mac},
		"fileName": {fileName},
	}
	resp, err := http.PostForm("http://127.0.0.1:3000/getfile", postBody)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.Body)
	}
}

func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func buildAuthenticator(auth string) {

	var file, err = os.Create(path)
	defer file.Close()
	if err == nil {
		file.WriteString(auth + "\n")
	}
}
func getPhys() string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
		return ""
	}
	fmt.Println("Contents of file:", string(data))
	return strings.Split(string(data), "\n")[0]
}

func main() {
	as, _ := getMacAddr()
	email := "newmail.com"
	password := "password"
	mac := as[0]
	//AddUser(email, password, mac)
	GetFile(email, password, "LpDdhkKWLpwMRzPAoLAZ", getPhys(), mac, "testfile")
	//LpDdhkKWLpwMRzPAoLAZ
}
