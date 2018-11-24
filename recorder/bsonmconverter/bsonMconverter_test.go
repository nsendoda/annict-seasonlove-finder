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
	m, err := bsonmconverter.BsonMRating(query)
	log.Println(m)
	if err != nil {
		t.Error(err)
	}
	a := bsonmconverter.BsonMAfterDateRating(query)
	log.Println(a)
	c := bsonmconverter.BsonMUntilDateRating(query)
	log.Println(c)

}
