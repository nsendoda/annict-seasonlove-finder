package bsonmconverter

import (
	"bufio"
	"log"
	"os"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestNew(t *testing.T) {
  // input
	log.Println("testデータベースに接続します...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		t.Errorf("mgodial failed")
	}
	defer func() {
		log.Println("データベース接続を閉じます..")
		db.Close()
	}()
	C := db.DB("annict").C("test")
	C.RemoveAll(bson.M{"user.username":"K"})
	C.RemoveAll(bson.M{"user.username":"naninunenosi"})
  if n, err := C.Count(); n != 0 || err != nil {
      t.Errorf("databese test collection count error.")
  }
	fp, err := os.Open("./input")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)

	var s []string
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}

	bsonmconverter := NewBsonMConverter()
	for i := range s {
		log.Println("insert: ", i)
		m, err := bsonmconverter.BsonMRating(s[i])
		if err != nil {
			t.Errorf("bsonmrating error")
		}
		err = C.Insert(m)
		if err != nil {
			t.Errorf("insert error")
		}
	}

  // test
  // time: s[0] -> s[2] -> s[1]
  // user and episode: s[1] == s[2]
  if c, err := C.Count(); c != 3 || err != nil {
    if c != 3 {
      t.Errorf("column count: %d, but expected 3\n", c)
    } else {
      t.Error("column count error")
    }
  } else {
    log.Println("Column count is valid. PASSED")
  }
  q1 := C.Find(bsonmconverter.BsonMAfterDateRating(s[2]))
  q2 := C.Find(bsonmconverter.BsonMAfterOrEqualDateRating(s[2]))
  if ct1, err := q1.Count(); ct1 != 1 || err != nil {
    if ct1 != 1 {
      t.Errorf("count : %d, but expected : 1", ct1)
    } else {
      t.Error("count error", err)
    }
  } else {
    log.Println("AfterData Test OK")
  }

  if ct2, err := q2.Count(); ct2 != 2 || err != nil {
    if ct2 != 2 {
      t.Errorf("count : %d, but expected : 2", ct2)
    } else {
      t.Error("count error", err)
    }
  } else {
    log.Println("AfterOrEqualData Test OK")
  }
  q3 := C.Find(bsonmconverter.BsonMUntilDateRating(s[1]))
  if ct3, err := q3.Count(); ct3 != 2 || err != nil {
    if ct3 != 2 {
      t.Errorf("count : %d, but expected : 2", ct3)
    } else {
      t.Error("count error", err)
    }
  } else {
    log.Println("UntilData Test OK")
  }
  if bsonmconverter.IsModified(s[1]) != true {
    t.Error("IsModified error. expected: true")
  } else {
    log.Println("IsModified check s[1] ok.")
  }
  if bsonmconverter.IsModified(s[2]) != false {
    t.Error("IsModified error. expected: false")
  } else {
    log.Println("IsModified check s[2] ok.")
  }
}
