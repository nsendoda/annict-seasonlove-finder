package bsonmconverter

import (
	"log"

	"github.com/naninunenosi/annict-seasonlove-finder/recfilter"
	"gopkg.in/mgo.v2/bson"
)

const CREATEDAT = "created_at"
const ISMODIFIED = "is_modified"

type bsonMConverter struct {
	episode_identify_key map[string][]string
}

func NewBsonMConverter() *bsonMConverter {
	b := bsonMConverter{}
	b.episode_identify_key = map[string][]string{"user": {"id"}, "episode": {"id"}}
	return &b
}

// convert json string to squeezed bson.M
func (b bsonMConverter) BsonMRating(record string) (bson.M, error) {
	var bdoc bson.M
	encoded_bytes, err := recfilter.Squeeze(record)
	err = bson.UnmarshalJSON(encoded_bytes, &bdoc)
	if err != nil {
		log.Println("bsonMconverter.go:BsonMRating: BSON.Mへの変換が失敗しました.", err)
		return bson.M{}, err
	}
	return bdoc, nil
}

// convert interface{} to bson.M
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

// bson.M型からユーザーのエピソードに対する記録を特定する要素のみを抽出したbson.M型を返す
// 現在はepisode.idとuser.username
func (b bsonMConverter) BsonMEpisodeIdentify(s bson.M) (m bson.M) {
  var join = func(a string, b string) string {
    return a + "." + b
  }
	m = bson.M{}
	for first_key, second_identify_key := range b.episode_identify_key {
		for _, second_key := range second_identify_key {
			s_f := Map(s[first_key])
			m[join(first_key, second_key)] = s_f[second_key]
		}
	}
	return m
}

func (b bsonMConverter) IsModified(record string) bool {
	m, _ := b.BsonMRating(record)
	if m[ISMODIFIED] == true {
		return true
	}
	return false
}

func (b bsonMConverter) BsonMAfterDateRating(record string) (m bson.M) {
	s, _ := b.BsonMRating(record)
	m = b.BsonMEpisodeIdentify(s)
	m[CREATEDAT] = bson.M{"$gt": s[CREATEDAT]}
	return m
}

func (b bsonMConverter) BsonMAfterOrEqualDateRating(record string) (m bson.M) {
	s, _ := b.BsonMRating(record)
	m = b.BsonMEpisodeIdentify(s)
	m[CREATEDAT] = bson.M{"$gte": s[CREATEDAT]}
	return m
}

func (b bsonMConverter) BsonMUntilDateRating(record string) (m bson.M) {
	s, _ := b.BsonMRating(record)
	m = b.BsonMEpisodeIdentify(s)
	m[CREATEDAT] = bson.M{"$lte": s[CREATEDAT]}
	return m
}
