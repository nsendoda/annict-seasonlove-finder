package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type record struct {
	ID            int         `json:"id"`
	Rating        interface{} `json:"rating"`
	RatingState   string      `json:"rating_state"`
	LikesCount    int         `json:"likes_count"`
	CommentsCount int         `json:"comments_count"`
	IsModified    bool        `json:"is_modified"`
	Comment       string      `json:"comment"`
	CreatedAt     time.Time   `json:"created_at"`
	Work          struct {
		ID              int    `json:"id"`
		ReleasedOnAbout string `json:"released_on_about"`
		ReviewsCount    int    `json:"reviews_count"`
		NoEpisodes      bool   `json:"no_episodes"`
		TitleKana       string `json:"title_kana"`
		Title           string `json:"title"`
		ReleasedOn      string `json:"released_on"`
		OfficialSiteURL string `json:"official_site_url"`
		WikipediaURL    string `json:"wikipedia_url"`
		EpisodesCount   int    `json:"episodes_count"`
		Media           string `json:"media"`
		WatchersCount   int    `json:"watchers_count"`
		SeasonName      string `json:"season_name"`
		MediaText       string `json:"media_text"`
		MalAnimeID      string `json:"mal_anime_id"`
		Images          struct {
			RecommendedURL string `json:"recommended_url"`
			Facebook       struct {
				OgImageURL string `json:"og_image_url"`
			} `json:"facebook"`
			Twitter struct {
				BiggerAvatarURL   string `json:"bigger_avatar_url"`
				OriginalAvatarURL string `json:"original_avatar_url"`
				ImageURL          string `json:"image_url"`
				MiniAvatarURL     string `json:"mini_avatar_url"`
				NormalAvatarURL   string `json:"normal_avatar_url"`
			} `json:"twitter"`
		} `json:"images"`
		SeasonNameText  string `json:"season_name_text"`
		TwitterUsername string `json:"twitter_username"`
		TwitterHashtag  string `json:"twitter_hashtag"`
	} `json:"work"`
	User struct {
		ID                 int       `json:"id"`
		WatchingCount      int       `json:"watching_count"`
		WatchedCount       int       `json:"watched_count"`
		WannaWatchCount    int       `json:"wanna_watch_count"`
		BackgroundImageURL string    `json:"background_image_url"`
		Username           string    `json:"username"`
		Name               string    `json:"name"`
		FollowersCount     int       `json:"followers_count"`
		StopWatchingCount  int       `json:"stop_watching_count"`
		OnHoldCount        int       `json:"on_hold_count"`
		RecordsCount       int       `json:"records_count"`
		CreatedAt          time.Time `json:"created_at"`
		Description        string    `json:"description"`
		AvatarURL          string    `json:"avatar_url"`
		FollowingsCount    int       `json:"followings_count"`
		URL                string    `json:"url"`
	} `json:"user"`
	Episode struct {
		ID                  int    `json:"id"`
		Number              string `json:"number"`
		NumberText          string `json:"number_text"`
		SortNumber          int    `json:"sort_number"`
		Title               string `json:"title"`
		RecordsCount        int    `json:"records_count"`
		RecordCommentsCount int    `json:"record_comments_count"`
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
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasParameter() {
		q = c.Find(bson.M{"user.username": p.UserName, "work.season_name": p.SeasonName})
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
