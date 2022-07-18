package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func dbRequest(message string) []byte {
	println("START CLIENT")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	host := os.Getenv("dgraph.Host")
	port := os.Getenv("dgraph.Port")
	path := "http://" + host + ":" + port

	r := strings.NewReader(message)

	resp, err := client.Post(path+"/query?timeout=5s", "application/dql", r)
	if err != nil {
		log.Print(err)
		return []byte("ERROR_CONNECT DB")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return []byte("ERROR_READ FROM DB")
	}

	defer resp.Body.Close()
	return body
}

func ManuHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	const q = `{
        results(func: uid(	%s))
        @recurse(depth: 5, loop: true)
        {
          uid
          content.group
          content.group.fragment
          content.group.path
          content.group.children
        }
      }`
	request := fmt.Sprintf(q, id)
	response := dbRequest(request)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	const q = `{ results(func: uid(%s))
      {uid   content.group.fragment  {uid fragment  versions  (orderdesc: version_date,first: 1)
      {uid   version_date blocks @facets(orderasc: ord)  {uid type   value before }}}}
      }`

	request := fmt.Sprintf(q, id)
	response := dbRequest(request)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/menu/{id}", ManuHandler)
	r.HandleFunc("/api/articles/{id}", ArticlesHandler)
	http.Handle("/", r)

	host := os.Getenv("server.Host")
	port := os.Getenv("server.Port")
	addr := host + ":" + port

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
