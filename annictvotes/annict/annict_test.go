package annict

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	apiKey := os.Getenv("ANNICT_APIKEY")
	annict := Annict{APIKey: apiKey}
	_, err := annict.Records(1)
	if err != nil {
		t.Error("エラーが発生しました.", err)
	}
	ss := []string{"access_token=" + apiKey, "filter_ids=1825247"}
	response, err := annict.httpGetAnnictAPIRaw(ss)
	defer response.Body.Close()
	if err != nil {
		t.Error("curlに失敗しました.\n", err)
	}
	_, err = annict.decode(response)
}
