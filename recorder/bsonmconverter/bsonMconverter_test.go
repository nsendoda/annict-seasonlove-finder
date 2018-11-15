package bsonmconverter

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

const (
	RECORD_LEN                     = 11
	RECORD_IDENTI_LEN              = 4
	RECORD_IDENTI_ADDCREATEDAT_LEN = RECORD_IDENTI_LEN + 1
)

func EqualExpectedMapLength(m map[string]int, expected int, t *testing.T) {
	if len(m) != expected {
		t.Errorf("想定される長さが違います. expected: %d, ans: %d\n", expected, len(m))
	}
}

func EqualExpectedBsonMLength(m bson.M, expected int, t *testing.T) {
	if len(m) != expected {
		t.Errorf("想定される長さが違います. expected: %d, ans: %d\n", expected, len(m))
	}
}

func TestNew(t *testing.T) {
	var query []string
	for i := 0; i < RECORD_LEN; i++ {
		query = append(query, string(i))
	}
	bsonmconverter := NewBsonMConverter()
	m1 := bsonmconverter.BsonMRating(query)
	m2 := bsonmconverter.BsonMAfterDateRating(query)
	m3 := bsonmconverter.BsonMUntilDateRating(query)

	EqualExpectedBsonMLength(m1, RECORD_LEN, t)
	EqualExpectedBsonMLength(m2, RECORD_IDENTI_ADDCREATEDAT_LEN, t)
	EqualExpectedBsonMLength(m3, RECORD_IDENTI_ADDCREATEDAT_LEN, t)
	EqualExpectedMapLength(bsonmconverter.record_identify_map, RECORD_IDENTI_LEN, t)

}
