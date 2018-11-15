package main

import (
	"log"
	"my_project/seasonlove-finder/recorder/bsonmconverter"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	mgo "gopkg.in/mgo.v2"
)

const updateDuration = 1 * time.Second

func recorderMain() error {
	log.Println("データベースに接続します...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		return err
	}
	defer func() {
		log.Println("データベース接続を閉じます..")
		db.Close()
	}()
	ratingCollection := db.DB("annict").C("rating")

	var recordsLock sync.Mutex
	var records [][]string

	log.Println("NSQに接続します")
	q, err := nsq.NewConsumer("votes", "recorder", nsq.NewConfig())
	if err != nil {
		return err
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		recordsLock.Lock()
		defer recordsLock.Unlock()
		if records == nil {
			records = make([][]string, 0)
		}
		parsed_string := strings.Split(string(m.Body), ",")
		records = append(records, parsed_string)
		return nil
	}))

	if err := q.ConnectToNSQLookupd("localhost:4161"); err != nil {
		return err
	}

	ticker := time.NewTicker(updateDuration)
	defer ticker.Stop()

	bconverter := bsonmconverter.NewBsonMConverter()

	update := func() {
		recordsLock.Lock()
		defer recordsLock.Unlock()
		if len(records) == 0 {
			log.Println("新しい記録はありません。データベースの更新をスキップします。")
			return
		}
		log.Println("データベースを更新します...")
		log.Printf("recordsinfo records length: %d, record size: %d", len(records), len(records[0]))
		log.Println(records[0])
		ok := true
		for _, record := range records {
			query := ratingCollection.Find(bconverter.BsonMAfterDateRating(record))
			ct, err := query.Count()
			if err != nil {
				log.Println("queryの計数に失敗しました.", err)
				ok = false
				continue
			}
			if ct == 0 {
				changeinfo, err := ratingCollection.Upsert(bconverter.BsonMUntilDateRating(record),
					bconverter.BsonMRating(record))
				if err != nil {
					log.Println("更新に失敗しました.", err)
					ok = false
				}
        log.Println("Upserted. ", record)
				log.Printf("is_updated: %d, number of doc remove: %d, number of doc matched: %d",
					changeinfo.Updated, changeinfo.Removed, changeinfo.Matched)
			} else {
				log.Println("更新対象ではないため、skipをします.\n", record)
			}
		}
		if ok {
			log.Println("データベースの更新が完了しました")
			records = nil
		}
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case <-ticker.C:
			update()
		case <-termChan:
			q.Stop()
		case <-q.StopChan:
			// 完了しました
			return nil
		}
	}
}

func main() {
	err := recorderMain()
	if err != nil {
		log.Fatal(err)
	}
}
