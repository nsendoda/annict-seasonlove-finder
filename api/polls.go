package main

import (
	"errors"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/naninunenosi/annict-seasonlove-finder/recfilter"
)

func handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePollsGet(w, r)
		return
	case "POST":
		handlePollPost(w, r)
		return
	case "DELETE":
		handlePollsDelete(w, r)
		return
	}
	// 未対応のHTTPメソッド
	respondHTTPErr(w, r, http.StatusNotFound)
}

func handlePollsGet(w http.ResponseWriter, r *http.Request) {
	db := GetVar(r, "db").(*mgo.Database)
	c := db.C("rating")
	p := NewPath(r.URL.Path)
	var q *mgo.Query
	if p.HasParameter() {
		q = c.Find(bson.M{"user.username": p.UserName,
			"work.season_name": p.SeasonName})

	} else {
		log.Println("パラメータが足りません.")
	}
	result := []recfilter.Record{}
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	log.Println(result[0])
	log.Println("internal output end")
	respond(w, r, http.StatusOK, &result)
}

func handlePollPost(w http.ResponseWriter, r *http.Request) {
	respondErr(w, r, http.StatusInternalServerError, errors.New("未実装です"))
}

func handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	respondErr(w, r, http.StatusInternalServerError, errors.New("未実装です"))
}
