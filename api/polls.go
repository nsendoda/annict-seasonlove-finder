package main

import (
	"errors"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type record struct {
	ID            int         `json:"id"`
	Rating        interface{} `json:"rating"`
	RatingState   string      `json:"rating_state"  bson:"rating_state"`
	LikesCount    int         `json:"likes_count"   bson:"likes_count"`
	CommentsCount int         `json:"comments_count"bson:"comments_count"`
	IsModified    bool        `json:"is_modified"   bson:"is_modified"`
	Comment       string      `json:"comment"`
	CreatedAt     string      `json:"created_at"    bson:"created_at"`
	Work          struct {
		ID              int    `json:"id"`
		Title           string `json:"title"`
		EpisodesCount   int    `json:"episodes_count"   bson:"episodes_count"`
		SeasonName      string `json:"season_name"      bson:"season_name"`
		SeasonNameText  string `json:"season_name_text" bson:"season_name_text"`
	} `json:"work"`
	User struct {
		ID                 int       `json:"id"`
		Username           string    `json:"username"`
		Name               string    `json:"name"`
	} `json:"user"`
	Episode struct {
		ID                  int    `json:"id"`
		Number              string `json:"number"`
		NumberText          string `json:"number_text" bson:"number_text"`
		SortNumber          int    `json:"sort_number" bson:"sort_number"`
		Title               string `json:"title"`
	} `json:"episode"`
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
	p := NewPath(r.URL.Path)
        var q *mgo.Query
	if p.HasParameter() {
	  q = c.Find(bson.M{"user.username": p.UserName,
                            "work.season_name": p.SeasonName})
          
	} else {
		log.Println("パラメータが足りません.")
	}
        result := []record{}
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
