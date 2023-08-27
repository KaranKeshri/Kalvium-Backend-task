package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Post struct {
	ID    string `json:"PersonID"`
	LName string `json:"LastName"`
}
type Ans struct {
	Question string `json:"question"`
	Answer   int    `json:"answer"`
	DateTime string `json:"date,omitempty"`
}

func getAnswer(str string) Ans {
	s := strings.TrimSpace(str)
	res := strings.Split(s, "/")
	res1 := res[1:len(res)]
	//5,plus,3,plus,6
	res = res1
	ouptutString := ""
	ans := 0
	fmt.Println(res)
	for i, j := range res {
		if i == 0 {
			fmt.Println(j)
			ouptutString += j
			x, _ := strconv.Atoi(j)
			ans = x
		}
		if j == "by" {
			x, _ := strconv.Atoi(res[i+1])
			ans = ans / x
			ouptutString += "/"
			ouptutString += res[i+1]

		} else if j == "into" {
			x, _ := strconv.Atoi(res[i+1])
			ans = ans * x
			ouptutString += "*"
			ouptutString += res[i+1]
		} else if j == "plus" {
			x, _ := strconv.Atoi(res[i+1])
			ans = ans + x
			ouptutString += "+"
			ouptutString += res[i+1]
		} else if j == "minus" {
			x, _ := strconv.Atoi(res[i+1])
			ans = ans - x
			ouptutString += "-"
			ouptutString += res[i+1]
		}
	}
	var answer Ans
	answer.Question = ouptutString
	answer.Answer = ans
	return answer
}

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s := r.URL.RequestURI()
	data := getAnswer(s)
	question := data.Question
	ansVal := data.Answer
	stmt, err := db.Prepare("INSERT INTO arithmatic(question,answer) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(question, ansVal)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(data)
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, err := db.Query("SELECT * FROM arithmatic ORDER BY date desc LIMIT 20")
	if err != nil {
		panic(err.Error())
	}
	var queries []Ans
	for result.Next() {
		var query Ans
		err := result.Scan(&query.Question, &query.Answer, &query.DateTime)
		if err != nil {
			panic(err.Error())
		}
		queries = append(queries, query)
	}
	json.NewEncoder(w).Encode(queries)
}

func main() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/kalvium")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/history", getHistory).Methods("GET")
	router.HandleFunc("/5/plus/3", getData).Methods("GET")
	router.HandleFunc("/3/minus/5", getData).Methods("GET")
	router.HandleFunc("/3/minus/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/3/into/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/8/plus/3", getData).Methods("GET")
	router.HandleFunc("/7/minus/5", getData).Methods("GET")
	router.HandleFunc("/9/plus/3", getData).Methods("GET")
	router.HandleFunc("/5/minus/5", getData).Methods("GET")
	router.HandleFunc("/7/minus/5/plus/9", getData).Methods("GET")
	router.HandleFunc("/2/into/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/8/minus/5/plus/5", getData).Methods("GET")
	router.HandleFunc("/6/into/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/10/into/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/20/minus/5/plus/5", getData).Methods("GET")
	router.HandleFunc("/8/into/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/7/into/5/plus/8", getData).Methods("GET")
	router.HandleFunc("/9/minus/5/plus/5", getData).Methods("GET")
	router.HandleFunc("/3/into/5/plus/8/into/6", getData).Methods("GET")

	http.ListenAndServe(":3000", router)
}
