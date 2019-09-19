package main

import (
	"database/sql"
	"io/ioutil"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"

	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
)

// HITOKOTOAMOUNT is Number of databases
var HITOKOTOAMOUNT int64
var db *sqlx.DB
var err error
var conn redis.Conn

// FormatMap xxxxx
var FormatMap map[string]HTTPFormat
var config *UserConfig

type UserConfig struct {
	Postgres struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	ListenPort string `yaml:"listenport"`
}

// MakeReturnMap xxxxxxx
func MakeReturnMap() {
	FormatMap = make(map[string]HTTPFormat)
	FormatMap["js"] = HTTPFormat{Charset: "text/javascript; charset=", Text: "var hitokoto=\"%s——「%s」\";var dom=document.querySelector('.hitokoto');Array.isArray(dom)?dom[0].innerText=hitokoto:dom.innerText=hitokoto;"}
	FormatMap["json"] = HTTPFormat{Charset: "application/json; charset=", Text: "{\"hitokoto\": \"%s\", \"source\": \"%s\"}"}
	FormatMap["text"] = HTTPFormat{Charset: "text/plain; charset=", Text: "%s——「%s」"}
}

func initConfig(filename string) {
	config = new(UserConfig)
	yamlFile, _ := ioutil.ReadFile("./config.yaml")
	_ = yaml.Unmarshal(yamlFile, config)
}

func initRedis() {
	conn, _ = redis.Dial("tcp", config.Redis.Address, redis.DialDatabase(config.Redis.DB))
	if config.Redis.Password != "" {
		conn.Do("auth", config.Redis.Password)
	}
	_, err := conn.Do("ping")
	checkErr(err)
}

func initHitokotoDB() {
	// db, err := sql.Open("postgres", "postgres://p")
	url := fmt.Sprintf("postgres://%s:%s@localhost/api?sslmode=disable", config.Postgres.User, config.Postgres.Password)
	db, err = sqlx.Connect("postgres", url)
	checkErr(err)
	// err1 := db.QueryRow("SELECT COUNT(id) FROM main;").Scan(&HITOKOTOAMOUNT)
	// checkErr(err1)
}

func initLogFile() {
	file, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	defer file.Close()
	syscall.Dup2(int(file.Fd()), 1)
	syscall.Dup2(int(file.Fd()), 2)
}

//DisallowMethod is allow current method
func DisallowMethod(w http.ResponseWriter, allow string, method string) bool {
	if allow != method && allow != "HEAD" {
		w.WriteHeader(405)
		fmt.Fprintln(w, "<h1>405 Not Allowed</h1>")
		return true
	}
	return false
}

func IsLimited(r *http.Request) []interface{} {
	header := r.Header
	xforwared := header.Get("X-Forwarded-For")
	if xforwared == "" {
		xforwared = "NoForwaredIP"
	}
	ret, err := redis.Values(conn.Do("CL.THROTTLE", xforwared, "35", "36", "360"))
	if err != nil {
		log.Printf("cl error level 1: %s", err)
		initRedis()
		ret, err = redis.Values(conn.Do("CL.THROTTLE", xforwared, "35", "36", "360"))
		log.Printf("cl error level 2: %s", err)
	}
	return ret
}

func checkErr(err error) {
	switch {
	case err == sql.ErrNoRows:
		// handleError("queryError")
		cnt.Hito = "哦~"
		cnt.Source = "袴田日向"
		log.Println("None Query")
	case err != nil:
		// handleError("connectError")
		cnt.Hito = "哦~"
		cnt.Source = "袴田日向"
		log.Println(err)
	}
}

func setCheating(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "%s", "{'cheating': true}")
}
