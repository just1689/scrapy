package remote

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type PutItem struct {
	Item io.Reader
	URL  string
}

func StartPutter(c chan PutItem, panicOnErr bool) {

	go func() {
		for i := range c {

			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

			req, err := http.NewRequest(http.MethodPut, i.URL, i.Item)
			if err != nil {
				logrus.Panic(err)
			}
			req.Header.Set("Content-Type", "application/json; charset=utf-8")

			resp, err := http.DefaultClient.Do(req.WithContext(ctx))
			if err != nil {
				if panicOnErr {
					logrus.Panic(err)
				} else {
					logrus.Errorln(err)
				}
			}
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				logrus.Errorln("tried sending to: ", i.URL)
				b, _ := ioutil.ReadAll(resp.Body)
				logrus.Errorln(string(b))
				logrus.Errorln(resp.StatusCode)
			}
		}
	}()
}
