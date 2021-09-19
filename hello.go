package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var parsedTemplates *template.Template
var host = os.Getenv("HOST")
var onion = os.Getenv("ONION")

type SucessData struct {
	URL   string
	Onion string
}

func parseTemplates() {
	var allFiles []string
	files, err := ioutil.ReadDir("./template")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".tmpl") || strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./template/"+filename)
		}
	}

	parsedTemplates, err = template.ParseFiles(allFiles...) //#parses all .tmpl files in the 'templates' folder
	checkErr(err)
}

func createShortie(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		// logic part of log in
		var url = r.Form["url"][0] // todo validate - exists
		// if   {
		// 	http.Error(w, "url required", 400)
		// 	return
		// }
		var maxHits, parseErr = strconv.ParseInt(r.Form["maxHits"][0], 10, 32)
		if parseErr != nil || maxHits < 1 || maxHits > 10 {
			http.Error(w, "maxHits must be between 1 and 10", 400)
			return
		}
		var expiry, expiryErr = strconv.ParseInt(r.Form["expiry"][0], 10, 32)
		if expiryErr != nil || expiry < 1 || expiry > (24*30) {
			http.Error(w, "expiry must be valid and within 30 days of now", 400)
			return
		}
		var expiryDate = time.Now().Add(time.Hour * time.Duration(expiry))
		var err error
		b := make([]byte, 4)
		_, err = rand.Read(b)
		checkErr(err)

		id := fmt.Sprintf("%x", b[0:4])

		_, err = db.Exec("INSERT INTO shorties(id, url, hits, maxHits, created, expiry) values(?,?,0,?,date('now'),?)", id, url, maxHits, expiryDate)
		checkErr(err)
		t := parsedTemplates.Lookup("success.html")

		d := SucessData{
			URL:   fmt.Sprintf("http://%s/%s", host, id),
			Onion: fmt.Sprintf("http://%s/%s", onion, id),
		}
		t.Execute(w, d)
	} else {
		http.Error(w, "method not supported", 405)
		return
	}
}

var shortieReg = regexp.MustCompile(`[[:alnum:]]{8}`)

func handleMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path == "/" {
			t := parsedTemplates.Lookup("index.html")
			t.Execute(w, nil)
			// todo: set cache header on this
		} else if shortieReg.MatchString(r.URL.Path) {
			// look up path
			id := strings.Split(r.URL.Path, "/")[1]

			rows, err := db.Query("SELECT url, hits FROM shorties where id=? and expiry > ? and hits < maxHits", id, time.Now())
			checkErr(err)
			var url string
			var hits int
			if rows.Next() {
				err = rows.Scan(&url, &hits)
				checkErr(err)
				rows.Close()

				// update count
				_, err = db.Exec("UPDATE shorties set hits=? where id=?", hits+1, id)
				checkErr(err)

				http.Redirect(w, r, url, http.StatusFound) // or have intermediate page with JS redirect after 5s ?
				return
			} else {
				// 404
				t := parsedTemplates.Lookup("404.html")
				t.Execute(w, nil)
				return
			}
		} else {
			http.Error(w, "invalid url", 404)
			return
		}

	} else {
		http.Error(w, "method not supported", 405)
		return
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./mini-madeuce.db")

	checkErr(err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS shorties (
		id int primary key not null, 
		url char(256),
		hits int, 
		maxHits int,
		created datetime,
		expiry datetime
	);`)
	checkErr(err)

	err = mime.AddExtensionType(".css", "text/css")
	checkErr(err)

	http.HandleFunc("/", handleMain)

	http.HandleFunc("/createShortie", createShortie)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	parseTemplates()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	if host == "" {
		host = "localhost:" + port
	}
	fmt.Println("server up listing at :" + port)

	err = http.ListenAndServe(":"+port, nil)

	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
