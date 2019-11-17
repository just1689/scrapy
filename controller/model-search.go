package controller

import (
	"encoding/json"
	"fmt"
	"github.com/just1689/scrapy/model"
	"github.com/just1689/scrapy/remote"
	"github.com/just1689/scrapy/tools"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var search = `<div class="resultsnum pagination__page-number pagination__page-number_right">`
var search2 = `of`
var search3 = `</div>`
var diff = len(search)

func ModelSearch() {
	resp, err := http.Get(model.DBUrlModel)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	arr := make([]model.ModelItem, 0)
	err = json.Unmarshal(b, &arr)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	c := make(chan model.ListingItem)

	//Writes model.ListingItem s until the channel is closed
	wgWriter := sync.WaitGroup{}
	go func() {
		wgWriter.Add(1)
		putterChan := make(chan remote.PutItem)
		remote.StartPutter(putterChan, true)
		for row := range c {
			fmt.Println("writing!")
			putterChan <- remote.PutItem{
				Item: tools.StructToIOBody(row),
				URL:  model.DBUrlListing + "?id=eq." + row.ID,
			}
		}
		wgWriter.Done()
	}()

	//Finds every model.ListingItem for a ...&P=n for 1..count
	for _, row := range arr {
		count := ModelSearchSpecific(row.Oem, row.Model)
		GetPages(row.Oem, row.Model, count, c)
	}
	close(c)

	//Block until the writer is done
	wgWriter.Wait()
}

func ModelSearchSpecific(oem, model string) (pages int) {

	start := `aHR0cHM6Ly93d3cuY2Fycy5jby56YS9zZWFyY2hWZWhpY2xlLnBocD9uZXdfb3JfdXNlZD0mbWFrZV9tb2RlbD0=`
	afterMake := `JTVC`
	afterModel := `JTVEJnZmc19hcmVhPSZhZ2VudF9sb2NhbGl0eT0mcHJpY2VfcmFuZ2U9Jm9zPSZsb2NhbGl0eT0mYm9keV90eXBlX2V4YWN0PSZ0cmFuc21pc3Npb249JmZ1ZWxfdHlwZT0mbG9naW5fdHlwZT0mbWFwcGVkX2NvbG91cj0mdmZzX3llYXI9JnZmc19taWxlYWdlPSZ2ZWhpY2xlX2F4bGVfY29uZmlnPSZrZXl3b3JkPSZzb3J0PXZmc19wcmljZQ==`
	url := tools.Base64ToStr(start) + oem + tools.Base64ToStr(afterMake) + model + tools.Base64ToStr(afterModel)

	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorln(err)
		return 0
	}

	return pageToCount(resp)

}

func pageToCount(resp *http.Response) int {
	b, _ := ioutil.ReadAll(resp.Body)
	s := string(b)

	//Find right section
	i := strings.Index(s, search)
	if i == -1 {
		//When the search string is not found, there is only one page
		return 1
	}
	s = s[i+diff : i+diff+30]

	//Remove whitespace
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, " ", "")

	//Cut off everything until of
	i = strings.Index(s, search2) + 2
	s = s[i : len(s)-i]
	//fmt.Println(s)

	//Cut off everything from </div>
	i = strings.Index(s, search3)
	s = s[:i]

	i, err := strconv.Atoi(s)
	if err != nil {
		logrus.Errorln(err)
	}
	return i
}

func GetPages(oem, mdl string, count int, resultOut chan model.ListingItem) {

	var pages int
	pages = count/20 + 1

	for page := 1; page <= pages; page++ {
		start := `aHR0cHM6Ly93d3cuY2Fycy5jby56YS9zZWFyY2hWZWhpY2xlLnBocD9uZXdfb3JfdXNlZD0mbWFrZV9tb2RlbD0=`
		afterMake := `JTVC`
		afterModel := `JTVEJnZmc19hcmVhPSZhZ2VudF9sb2NhbGl0eT0mcHJpY2VfcmFuZ2U9Jm9zPSZsb2NhbGl0eT0mYm9keV90eXBlX2V4YWN0PSZ0cmFuc21pc3Npb249JmZ1ZWxfdHlwZT0mbG9naW5fdHlwZT0mbWFwcGVkX2NvbG91cj0mdmZzX3llYXI9JnZmc19taWxlYWdlPSZ2ZWhpY2xlX2F4bGVfY29uZmlnPSZrZXl3b3JkPSZzb3J0PXZmc19wcmljZQ==`
		url := tools.Base64ToStr(start) + url.QueryEscape(oem) + tools.Base64ToStr(afterMake) + url.QueryEscape(mdl) + tools.Base64ToStr(afterModel) + "?P=" + strconv.Itoa(page)

		resp, err := http.Get(url)
		if err != nil {
			logrus.Errorln(err)
			continue
		}

		doc, _ := ioutil.ReadAll(resp.Body)
		var i = 0
		s := string(doc)
		s = strings.ReplaceAll(s, "\r", "")
		s = strings.ReplaceAll(s, "\n", "")
		search := `vehicle-list__center-block"><a href="/`

		for n := 1; n <= 20; n++ {
			i = strings.Index(s, search)
			if i == -1 {
				break
			}
			s = s
			s = s[i+len(search) : len(s)-len(search)]

			i = strings.Index(s, `"`)
			snip := s[:i]

			resultOut <- model.ListingItem{
				ID:    tools.HashString(url),
				Oem:   oem,
				Model: mdl,
				Url:   snip,
			}

			s = s[i:]

		}

	}

}
