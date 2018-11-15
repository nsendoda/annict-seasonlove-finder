package main

import (
	"errors"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type record struct {
	ID          bson.ObjectId `bson:"_id"                 json:"id"`
	RecordID    string        `bson:"record_id"           json:"record_id"`
	RatingState string        `bson:"rating_state"        json:"rating_state"`
	CreatedAt   string        `bson:"created_at"          json:"created_at"`
	UserID      string        `bson:"user_id"             json:"user_id"`
	Name        string        `bson:"user_name"           json:"user_name"`
	WorkID      string        `bson:"work_id"             json:"work_id"`
	Title       string        `bson:"work_title"          json:"work_title"`
	SeasonName  string        `bson:"season_name"         json:"season_name"`
	EpisodeID   string        `bson:"episode_id"          json:"episode_id"`
	Number      string        `bson:"episode_number"      json:"episode_number"`
	SortNumber  string        `bson:"episode_sort_number" json:"episode_sort_number"`
}

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
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasParameter() {
		q = c.Find(bson.M{"user_name": p.UserName, "season_name": p.SeasonName})
	} else {
		log.Println("パラメータが足りません.")
	}
	var result []*record
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, &result)
}

func handlePollPost(w http.ResponseWriter, r *http.Request) {
	respondErr(w, r, http.StatusInternalServerError, errors.New("未実装です"))
}

func handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	respondErr(w, r, http.StatusInternalServerError, errors.New("未実装です"))
}
