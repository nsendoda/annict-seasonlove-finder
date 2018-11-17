package annict

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/json-iterator/go"
)

const (
	SORT_PARAMETER     = "sort_id=desc"
	PER_PAGE_PARAMETER = "per_page=50"
)

type Annict struct {
	APIKey string
}

func (a *Annict) httpGetAnnictAPIRaw(parameters []string) (*http.Response, error) {
	return http.Get("https://api.annict.com/v1/records?" + strings.Join(parameters, "&"))
}

func (a *Annict) httpGetAnnictAPI(page int) (*http.Response, error) {
	ss := []string{SORT_PARAMETER, PER_PAGE_PARAMETER}
	ss = append(ss, "page="+fmt.Sprint(page))
	ss = append(ss, "access_token="+a.APIKey)
	response, err := a.httpGetAnnictAPIRaw(ss)
	return response, err
}

// response -> []string
// 1. json形式 -> bytes形式に変換
// 2. bytes形式をUnmarshalでinterface{}型(map)へ
// 3. map["records"][0]をMarshalでjson形式へ
// 4. json -> string にして[]stringへappend
func (a *Annict) decode(response *http.Response) ([]string, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	s := buf.Bytes()
	var map_records interface{}
	err := jsoniter.Unmarshal(s, &map_records)
	var res []string
	for _, map_rec := range map_records.(map[string]interface{})["records"].([]interface{}) {
		str_rec, _ := jsoniter.Marshal(map_rec)
		res = append(res, string(str_rec))
	}
	return res, err
}

func (a *Annict) Records(page int) ([]string, error) {
	response, err := a.httpGetAnnictAPI(page)
	defer response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("annict: %dページ目の記録取得に失敗しました: %v", page, err)
	}
	// fmt.Print(response)
	return a.decode(response)
}
