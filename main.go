package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

const (
	DatabaseUrlKey = "/workmachine/database_url"
)

var (
	Db gorm.DB
)

func dbConnect(databaseUrl string) {
	log.Println("Connecting to database:", databaseUrl)
	var err error
	Db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Println(err)
	}
	Db.LogMode(true)
}

func init() {
	etcdHosts := os.Getenv("ETCD_HOSTS")
	if etcdHosts == "" {
		etcdHosts = "http://127.0.0.1:4001"
	}

	etcdClient := etcd.NewClient([]string{etcdHosts})

	resp, err := etcdClient.Get(DatabaseUrlKey, false, false)
	if err != nil {
		panic(err)
	}

	databaseUrl := resp.Node.Value
	dbConnect(databaseUrl)
}

func main() {
	go AssignmentExpirer()

	log.Println("WorkMachine Starting...")

	r := mux.NewRouter()
	r.HandleFunc("/v1/work", NewWorkflowsHandler).Methods("PUT")
	r.HandleFunc("/v1/work/{workflow}", ShowWorkflowsHandler).Methods("GET")

	//r.HandleFunc("/v1/workflow/{workflow}/tasks", TaskHandler).Methods("PUT")
	//r.HandleFunc("/v1/workflow/{workflow}/tests", TaskHandler).Methods("PUT")
	//r.HandleFunc("/v1/tasks", newTaskHandler).Methods("POST")
	//r.HandleFunc("/v1/assignments", newTaskHandler).Methods("POST")

	// r.HandleFunc("/v1/assignments", func(w http.ResponseWriter, req *http.Request) {
	// 	assign := AvailableAssignments.GetUnfinished()
	// 	if assign == nil {
	// 		renderJson(w, false)
	// 		return
	// 	}

	// 	if !assign.TryToAssign() {
	// 		renderJson(w, false)
	// 		return
	// 	}

	// 	renderJson(w, assign)
	// }).Methods("GET")

	// r.HandleFunc("/v1/assignments", func(w http.ResponseWriter, req *http.Request) {
	// 	log.Println("Posting", req.FormValue("id"))

	// 	assign := AvailableAssignments.Find(req.FormValue("id"))
	// 	if assign != nil {
	// 		assign.Finish(req.FormValue(assign.InputField.Value))
	// 	}

	// 	renderJson(w, true)
	// }).Methods("POST")

	r.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path[1:])
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.Handle("/", r)
	http.ListenAndServe(":3002", nil)
}
