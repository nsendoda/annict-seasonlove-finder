package annict

import (
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	apiKey := os.Getenv("ANNICT_APIKEY")
	annict := Annict{APIKey: apiKey}
	// test all
	res, err := annict.Records(1)
	if err != nil {
		t.Error("エラーが発生しました.", err)
	}
  for i, str := range res {
	  log.Printf("Recordsから帰ってきた文字列結果 index:%d, res:%s\n", i, str)
  }
	// test httpGetAnnictAPIRaw
	ss := []string{"access_token=" + apiKey, "filter_ids=1825247"}
	response, err := annict.httpGetAnnictAPIRaw(ss)
	defer response.Body.Close()
	if err != nil {
		t.Error("curlに失敗しました.\n", err)
	}
	// test decode
	_, err = annict.decode(response)
}
