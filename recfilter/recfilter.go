package recfilter

import (
	"log"

	"github.com/json-iterator/go"
)

type Record struct {
	ID            int         `json:"id"`
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
		OfficialSiteURL string `json:"official_site_url"bson:"official_site_url"`
	} `json:"work"`
	User struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	} `json:"user"`
	Episode struct {
		ID         int    `json:"id"`
		Number     string `json:"number"`
		NumberText string `json:"number_text" bson:"number_text"`
		SortNumber int    `json:"sort_number" bson:"sort_number"`
		Title      string `json:"title"`
	} `json:"episode"`
}

func Squeeze(rec string) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	unmarshaled_r := Record{}
	err := json.Unmarshal([]byte(rec), &unmarshaled_r)
	if err != nil {
		log.Println("Unmarshal failed in Squeezing.", err)
		return nil, err
	}
	b, err := json.Marshal(unmarshaled_r)
	if err != nil {
		log.Println("Marshal failed in Squeezing.", err)
		return nil, err
	}
	return b, nil
}
