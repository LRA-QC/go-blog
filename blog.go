package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"

	"net/http"
	"os"

	"github.com/LRA-QC/blog/woxcache"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type PageVariables struct {
	Date string
	Time string
}

const version = "1.0"

//TICKER EVER 60 SECONDS
func startCacheTicker(f func()) chan bool {
	done := make(chan bool, 1)
	go func() {
		ticker := time.NewTicker(time.Second * 60)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-done:
				fmt.Println("done")
				return
			}
		}
	}()
	return done
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var p string

	now := time.Now()              // find the time right now
	HomePageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	switch r.URL.Path {
	case "/page/about":
		p = "www/about.html"
	default:
		p = "www/index.html"
	}
	var tpl = template.Must(template.ParseFiles(p))
	tpl.Execute(w, HomePageVars)
	//	w.Write([]byte("<h1>Hello World!</h1>"))
}

func main() {
	db, errsql := sql.Open("sqlite3", "./foo.db")
	fmt.Print(errsql)
	defer db.Close()
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	url := fmt.Sprintf("http://%s:%s/", address, port)
	fmt.Println("==> Starting WoxCMS server 🌍 " + version + " 🚀:" + url + "]")
	fs := http.FileServer(http.Dir("static/"))

	//mux := http.NewServeMux()
	mux := mux.NewRouter()

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", indexHandler)

	done := startCacheTicker(func() {
		woxcache.CachePurge()
	})

	woxcache.CacheDump()
	woxcache.CacheSet("this_is_a_date", "02-01-2006")
	res := woxcache.CacheGet("this_is_a_date")
	fmt.Println("Cache entry for [this_is_a_date] result: " + res)
	woxcache.CacheDump()
	//	time.Sleep(30 * time.Second)
	//	woxcache.CacheDump()

	http.ListenAndServe(address+":"+port, mux)
	//close cache ticker
	close(done)
}
