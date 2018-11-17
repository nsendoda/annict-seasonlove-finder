package bsonmconverter

import (
	"log"

	"gopkg.in/mgo.v2/bson"
)

const CREATEDAT = "created_at"

type bsonMConverter struct {
	episode_identify_key map[string][]string
}

func NewBsonMConverter() *bsonMConverter {
	b := bsonMConverter{}
	b.episode_identify_key = map[string][]string{"user": {"id"}, "work": {"id"}, "episode": {"id"}}
	return &b
}

func (b bsonMConverter) BsonMRating(record string) (m bson.M) {
	var bdoc bson.M
	err := bson.UnmarshalJSON([]byte(record), &bdoc)
	if err != nil {
		log.Println("bsonMconverter.go:BsonMRating: BSON.Mへの変換が失敗しました.", err)
	}
	return bdoc
}

func Map(v interface{}) bson.M {
	js, err := bson.MarshalJSON(v)
	if err != nil {
		log.Println("bsonMconverter.go:Map: JSONへの変換が失敗しました.", err)
	}
	var bdoc bson.M
	bdoc = bson.M{}
	err = bson.UnmarshalJSON(js, &bdoc)
	if err != nil {
		log.Println("bsonMconverter.go:Map: BSON.Mへの変換が失敗しました.", err)
	}
	return bdoc
}

func (b bsonMConverter) BsonMEpisodeIdentify(record string, s bson.M) (m bson.M) {
	m = bson.M{}
	for first_key, second_identify_key := range b.episode_identify_key {
		m[first_key] = bson.M{}
		for _, second_key := range second_identify_key {
			s_f := Map(s[first_key])
			m_f := Map(m[first_key])
			m_f[second_key] = s_f[second_key]
			m[first_key] = m_f
		}
	}
	return m
}

func (b bsonMConverter) BsonMAfterDateRating(record string) (m bson.M) {
	s := b.BsonMRating(record)
	m = b.BsonMEpisodeIdentify(record, s)
	m[CREATEDAT] = bson.M{"$gt": s[CREATEDAT]}
	return m
}

func (b bsonMConverter) BsonMUntilDateRating(record string) (m bson.M) {
	s := b.BsonMRating(record)
	m = b.BsonMEpisodeIdentify(record, s)
	m[CREATEDAT] = bson.M{"$lte": s[CREATEDAT]}
	return m
}
