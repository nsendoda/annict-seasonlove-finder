package main

import (
	"context"
	"log"
	"my_project/seasonlove-finder/annict"
	"os"
	"time"
  "strings"
)

const (
	CURL_PAGE_NUMBER       = 2000
  READY_CURL_PAGE_NUMBER = 1000
	TEST_CURL_PAGE_NUMBER  = 3
)

func getAndDecodeAnnictRatings(page int) ([][]string, error) {
	apiKey := os.Getenv("ANNICT_APIKEY")
	annict := &annict.Annict{APIKey: apiKey}
  log.Printf("page:%dの取得を開始します.\n", page)
	recs, err := annict.Records(page)
	if err != nil {
		log.Fatalf("page:%dの記録取得に失敗しました: %v\n", page, err)
		return recs, err
	}
	if len(recs) == 0 {
		log.Fatalf("page:%dにデータはありませんでした\n", page)
		return recs, nil
	}
	/*	for _, rec := range recs {
			for _, s := range rec {
				log.Println(s)
			}
		}
	*/
	return recs, nil
}

func readFromAnnict(ctx context.Context, votes chan<- string,
	page int) {

	done := make(chan struct{})
	defer func() { <-done }()

	// defer resp.Body.Close() // respが加工されて帰ってくるかどうかが問題

	go func() {
		defer close(done)
		//log.Println("resp:", resp.StatusCode)
		//  if resp.StatusCode != 200 {
		ratings, err := getAndDecodeAnnictRatings(page)
		if err != nil {
			log.Println("取得リクエストに失敗しました:", err)
			return
		}
		for _, rating := range ratings {
      joined_rating := strings.Join(rating, ",")
      log.Println("joined_rating is :", joined_rating)
			votes <- joined_rating
		}
	}()
	select {
	case <-ctx.Done(): // 終了要求がきた
	case <-done: // goroutineが終了した
	}
}

// タイムアウトの実装
// @how 1分後にキャンセルされ、ctx.Done()が呼ばれるように
func readFromAnnictWithTimeout(ctx context.Context, timeout time.Duration,
	votes chan<- string,
	page int) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	readFromAnnict(ctx, votes, page)
}

// 10秒ごとに、readFromAnnictWithTimeoutを呼び出す
// ctx.Done()が呼ばれると終了
// @note votesのクローズはここで行う
func annictStream(ctx context.Context, votes chan<- string) {
	defer close(votes)
	for page_i := 1; page_i <= READY_CURL_PAGE_NUMBER; page_i++ {
		log.Println("Annictに問い合わせます...")
		readFromAnnictWithTimeout(ctx, 1*time.Minute, votes, page_i)
		log.Println("　待機中")
		select {
		case <-ctx.Done():
			log.Println("Annictへの問い合わせを終了します...")
			return
		case <-time.After(10 * time.Second):
		}
	}
	log.Println("全てのデータ取得が正常に終了しました.")
}
