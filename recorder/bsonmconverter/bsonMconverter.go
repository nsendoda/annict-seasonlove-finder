package bsonmconverter

import (
	"gopkg.in/mgo.v2/bson"
)

type bsonMConverter struct {
	record_order        []string
	record_identify_map map[string]int
	created_at_index    int
}

func NewBsonMConverter() *bsonMConverter {
	b := bsonMConverter{}
	b.record_order = []string{"record_id", "rating_state", "created_at", "user_id", "user_name",
		"work_id", "work_title", "season_name", "episode_id", "episode_number",
		"episode_sort_number"}
	record_identify_order := []string{"user_id", "work_id", "episode_id", "episode_number"}
	b.record_identify_map = make(map[string]int)
	for i := range b.record_order {
		for _, s := range record_identify_order {
			if b.record_order[i] == s {
				b.record_identify_map[s] = i
				break
			}
		}
		if b.record_order[i] == "created_at" {
			b.created_at_index = i
		}
	}
	return &b
}

func (b bsonMConverter) BsonMRating(record []string) (m bson.M) {
	m = bson.M{}
	for i := range b.record_order {
		m[b.record_order[i]] = record[i]
	}
	return m
}

func (b bsonMConverter) BsonMAfterDateRating(record []string) (m bson.M) {
	m = bson.M{}
	for key, value := range b.record_identify_map {
		m[key] = record[value]
	}
	m["created_at"] = bson.M{"$gt": record[b.created_at_index]}
	return m
}

func (b bsonMConverter) BsonMUntilDateRating(record []string) (m bson.M) {
	m = bson.M{}
	for key, value := range b.record_identify_map {
		m[key] = record[value]
	}
	m["created_at"] = bson.M{"$lte": record[b.created_at_index]}
	return m
}
