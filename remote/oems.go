package remote

import (
	"encoding/json"
	"github.com/just1689/scrapy/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetOEMsByURL(url string) (c chan model.ModelItem) {
	c = make(chan model.ModelItem)
	go func() {
		defer close(c)

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorln(err)
		}

		r := make(map[string]json.RawMessage)
		err = json.Unmarshal(b, &r)
		if err != nil {
			logrus.Errorln(err)
		}

		list := r["make-model"]

		m := make(map[string]json.RawMessage)
		err = json.Unmarshal(list, &m)
		if err != nil {
			logrus.Errorln(err)
		}

		for x, y := range m {
			bundle := make(map[string]json.RawMessage)
			json.Unmarshal(y, &bundle)

			models := (bundle["children"])
			modelI := make(map[string]Item)
			json.Unmarshal(models, &modelI)

			for _, x2 := range modelI {
				c <- model.ModelItem{
					Oem:   strings.ToUpper(x),
					Model: x2.Name,
				}
			}
		}
	}()
	return c

}

type Item struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
