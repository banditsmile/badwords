package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"

	"./wordfilter"
)

var wT = wordfilter.Init("badwords.txt")

func run(w http.ResponseWriter, r *http.Request, replace bool) {
	t1, _ := strconv.ParseFloat(time.Now().Format("05.9999"), 32)
	result := make(map[string]interface{})
	if r.Method == "POST" {
		r.ParseForm()
		for k, v := range r.PostForm {
			txtSlice := strings.Split(v[0], "")
			badWords := wordfilter.Search(wT, &txtSlice, "*")
			if replace {
				result[k] = strings.Join(txtSlice, "")
			} else {
				result[k] = badWords
			}
		}
		enc := json.NewEncoder(w)
		enc.Encode(result)
	}
	t2, _ := strconv.ParseFloat(time.Now().Format("05.9999"), 32)
	if replace {
		fmt.Printf("replace: %1.4fms\n", t2-t1)
	} else {
		fmt.Printf("match: %1.4fms\n", t2-t1)
	}
}

func reload(w http.ResponseWriter, r *http.Request){
	wT=wordfilter.Init("badwords.txt");
	if r.Method == "POST"{
		r.ParseForm()
		for k, v := range r.PostForm{
			fmt.Println(k,v)
		}
		enc :=json.NewEncoder(w)
		enc.Encode(r.PostForm)
	}
	w.Write([]byte("hello"))
}

func replace(w http.ResponseWriter, r *http.Request) {
	run(w, r, true)
}

func match(w http.ResponseWriter, r *http.Request) {
	run(w, r, false)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	
	http.HandleFunc("/match", match)
	http.HandleFunc("/replace", replace)
	http.HandleFunc("/reload", reload)
	server := &http.Server{
		Addr:	strings.Join([]string{os.Args[1], os.Args[2]}, ":"),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	server.SetKeepAlivesEnabled(true)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe err:", err)
	}
}
