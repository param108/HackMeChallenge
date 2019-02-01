package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var mutex *sync.Mutex

type LoginMessage struct {
	Password string
	Username string
}

type LoginResponse struct {
	Error   string
	Message string
}

var count = 1

var logs = []string{}

var USERNAME = "bhulakkad"
var PASSWORD = "192"

func main() {
	mutex = &sync.Mutex{}
	http.HandleFunc("/login/", loginUser)
	http.HandleFunc("/logs/", showLogs)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.ListenAndServe(":3500", nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func showLogs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ShowLogs")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	msg := LoginMessage{}
	// Unmarshal
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	message := ""
	if strings.Compare(msg.Username, USERNAME) == 0 &&
		strings.Compare(msg.Password, PASSWORD) == 0 {
		message = strings.Join(logs, ",\n")
	} else {
		message = "Invalid username and password"
	}

	enableCors(&w)

	success := LoginResponse{}
	success.Error = "false"
	success.Message = message
	output, err := json.Marshal(success)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginUser")
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	msg := LoginMessage{}
	// Unmarshal
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	message := "Invalid Username"
	switch msg.Username {
	case USERNAME:
		if strings.Compare(msg.Password, PASSWORD) == 0 {
			message = "Success. You Wish. This is just the next level of control." +
				"You need to login as root to find the real truth." +
				"Try posting the same username and password to /logs/ to get some clues"
		} else {
			/*logs = append(logs,
				fmt.Sprintf("Login to  %v failed. Expected %v got %v", USERNAME, PASSWORD, msg.Password))
			if len(logs) > 30 {
				logs = logs[len(logs)-30:]
			}*/
			message = "Fail, Invalid Password"
		}
	case "root":
		expected := strconv.FormatUint(PRIMES[(count%300)+2000], 10)
		if strings.Compare(msg.Password, expected) == 0 {
			message = "# You seek the answer to this question: Which book is famous for proclaiming that the answer to life, the universe and everything is 42 ? and who was the author ?"
		} else {
			message = "Fail, Invalid Password"
			logs = append(logs,
				fmt.Sprintf("Login to root failed. Expected %v got %v", expected, msg.Password))
			if len(logs) > 30 {
				logs = logs[len(logs)-30:]
			}
			count++
		}
	}
	enableCors(&w)

	success := LoginResponse{}
	success.Error = "false"
	success.Message = message
	output, err := json.Marshal(success)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
