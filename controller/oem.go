package controller

import (
	"github.com/just1689/scrapy/model"
	"github.com/just1689/scrapy/remote"
	"github.com/just1689/scrapy/tools"
	"github.com/sirupsen/logrus"
	"strings"
)

var OEMURL = tools.Base64ToStr("aHR0cHM6Ly93d3cuY2Fycy5jby56YS9zZWFyY2gvYi9mYWNldF9jb3VudD9uZXdfb3JfdXNlZD0mbWFrZV9tb2RlbD0mdmZzX2FyZWE9JmFnZW50X2xvY2FsaXR5PSZwcmljZV9yYW5nZT0mb3M9JnBhcmFtPW1ha2VfbW9kZWw=")

func GetOEMs() {
	outOEM := make(chan remote.PutItem, 5)
	outModel := make(chan remote.PutItem, 5)
	logrus.Infoln("...remote.GetPoster()")
	remote.StartPutter(outOEM, true)
	remote.StartPutter(outModel, true)

	logrus.Infoln("...remote.GetOEMsByURL()")
	oems := remote.GetOEMsByURL(OEMURL)
	for row := range oems {
		outOEM <- remote.PutItem{
			Item: tools.StructToIOBody(model.IDItem{ID: row.Oem}),
			URL:  (model.DBUrlOEM + "?id=eq." + row.Oem),
		}

		id := row.Oem + "." + row.Model
		id = strings.ReplaceAll(id, " ", "")
		outModel <- remote.PutItem{
			Item: tools.StructToIOBody(model.ModelItem{
				ID:    id,
				Oem:   row.Oem,
				Model: row.Model,
			}),
			URL: (model.DBUrlModel + "?id=eq." + id),
		}
	}
	logrus.Infoln("...done")
	close(outOEM)
	close(outModel)
}
