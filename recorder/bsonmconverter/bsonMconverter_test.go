package bsonmconverter

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	b, err := ioutil.ReadFile("./input")
	if err != nil {
		t.Error("input error\n", err)
	}
	query := string(b)
	bsonmconverter := NewBsonMConverter()
	log.Println(bsonmconverter.BsonMRating(query))
	log.Println(bsonmconverter.BsonMAfterDateRating(query))
	log.Println(bsonmconverter.BsonMUntilDateRating(query))

}
