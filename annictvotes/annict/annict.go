package annict

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	NONERATING_STRING = "none"

	SORT_PARAMETER     = "sort_id=desc"
	PER_PAGE_PARAMETER = "per_page=50"
)

type Annict struct {
	APIKey string
}

type records struct {
	Record []record `json:"records"`
}

type record struct {
	ID          int     `json:"id"`
	RatingState string  `json:"rating_state"`
	CreatedAt   string  `json:"created_at"`
	User        user    `json:"user"`
	Work        work    `json:"work"`
	Episode     episode `json:"episode"`
}
type user struct {
	ID   int    `json:"id"`
	Name string `json:"username"`
}
type work struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	SeasonName string `json:"season_name"`
}
type episode struct {
	ID         int    `json:"id"`
	Number     string `json:"number"`
	SortNumber int    `json:"sort_number"`
}

func noneRatingFilter(rating_state string) string {
	if rating_state == "" {
		return NONERATING_STRING
	}
	return rating_state
}

func (r record) PrintFormat() {
	log.Println("ID                :", r.ID)
	log.Println("RatingState       :", r.RatingState)
	log.Println("CreatedAt         :", r.CreatedAt)
	log.Println("User.ID           :", r.User.ID)
	log.Println("User.Name         :", r.User.Name)
	log.Println("Work.ID           :", r.Work.ID)
	log.Println("Work.Title        :", r.Work.Title)
	log.Println("Work.SeasonName   :", r.Work.SeasonName)
	log.Println("Episode.ID        :", r.Episode.ID)
	log.Println("Episode.Number    :", r.Episode.Number)
	log.Println("Episode.SortNumber:", r.Episode.SortNumber)
  log.Println()
}

func (a *Annict) httpGetAnnictAPIRaw(parameters []string) (*http.Response, error) {
	return http.Get("https://api.annict.com/v1/records?" + strings.Join(parameters, "&"))
}

func (a *Annict) httpGetAnnictAPI(page int) (*http.Response, error) {
	ss := []string{SORT_PARAMETER, PER_PAGE_PARAMETER}
	ss = append(ss, "page="+fmt.Sprint(page))
	ss = append(ss, "access_token="+a.APIKey)
	response, err := a.httpGetAnnictAPIRaw(ss)
	return response, err
}

func (a *Annict) decode(response *http.Response) ([][]string, error) {
	var data records
	var recs [][]string
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Print("annict: json decodeに失敗しました\n", err)
		return recs, err
	}
	for _, r := range data.Record {
		r.PrintFormat()
		var rec []string
		rec = append(rec, fmt.Sprintf("%v", r.ID))
		rec = append(rec, noneRatingFilter(fmt.Sprintf("%v", r.RatingState)))
		rec = append(rec, fmt.Sprint(r.CreatedAt))
		rec = append(rec, fmt.Sprint(r.User.ID))
		rec = append(rec, fmt.Sprint(r.User.Name))
		rec = append(rec, fmt.Sprint(r.Work.ID))
		rec = append(rec, fmt.Sprint(r.Work.Title))
		rec = append(rec, fmt.Sprint(r.Work.SeasonName))
		rec = append(rec, fmt.Sprint(r.Episode.ID))
		rec = append(rec, fmt.Sprint(r.Episode.Number))
		rec = append(rec, fmt.Sprint(r.Episode.SortNumber))
		recs = append(recs, rec)
	}
	return recs, nil
}

func (a *Annict) Records(page int) ([][]string, error) {
	response, err := a.httpGetAnnictAPI(page)
	defer response.Body.Close()
	if err != nil {
		return make([][]string, 0), fmt.Errorf("annict: %dページ目の記録取得に失敗しました: %v", page, err)
	}
	// fmt.Print(response)
	return a.decode(response)
}
