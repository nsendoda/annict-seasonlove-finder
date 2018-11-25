package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/naninunenosi/annict-seasonlove-finder/recorder/bsonmconverter"

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
	var records []string

	log.Println("NSQに接続します")
	q, err := nsq.NewConsumer("votes", "recorder", nsq.NewConfig())
	if err != nil {
		return err
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		recordsLock.Lock()
		defer recordsLock.Unlock()
		if records == nil {
			records = make([]string, 0)
		}
		records = append(records, string(m.Body))
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
		log.Printf("recordsinfo records length: %d", len(records))
		ok := true
		for _, record := range records {
			// when record is after the latest, insert
			// when record is equal to the latest and modified, upsert

			m, _ := bconverter.BsonMRating(record)
			episode_identify_m := bconverter.BsonMEpisodeIdentify(m)
			has_data := ratingCollection.Find(episode_identify_m)
			has_data_ct, err := has_data.Count()
			if err != nil {
				log.Println("first findの計数に失敗しました。")
				continue
			}
			// 該当するエピソードの記録がまだないならばInsert
			if has_data_ct == 0 {
				log.Println("new doc. Ready Insert.")
				err := ratingCollection.Insert(m)
				if err != nil {
					log.Println("挿入に失敗しました.", err)
					ok = false
				} else {
					log.Println("inserted:", m["id"])
				}

			// 該当するエピソードの記録がある時は、
			// 受け取ったjsonの方が新しいならUpdate
			// 全く同じ時間（つまり記録IDが一致）しているなら、is_modifiedがtrueであればUpdate
			// 受け取ったjsonの方が古いならSkip
			} else {
				query := ratingCollection.Find(bconverter.BsonMAfterDateRating(record))
				ct, err := query.Count()
				query_e := ratingCollection.Find(bconverter.BsonMAfterOrEqualDateRating(record))
				ct_e, err2 := query_e.Count()
				if err != nil || err2 != nil {
					log.Println("queryの計数に失敗しました.", err, err2)
					ok = false
					continue
				}
				// Update
				if ct_e == 0 || (ct == 0 && bconverter.IsModified(record)) {
					log.Println("detect modified. Ready Upsert.")
					err := ratingCollection.Update(episode_identify_m, m)
					if err != nil {
						log.Println("更新に失敗しました.", err)
						ok = false
					} else {
						log.Println("updated:", m["id"])
					}
				// Skip
				} else {
					log.Println("更新対象ではないため、skipをします.\n")
				}
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
