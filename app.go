package main

import (
	"github.com/just1689/scrapy/controller"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {

	logrus.Infoln("Starting!")
	http.HandleFunc("/oems", handleOEMs)
	http.HandleFunc("/details", handleDetails)
	panic(http.ListenAndServe(":8080", nil))

}

func handleDetails(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("handleDetails()")
	controller.ModelSearch()
	w.Write([]byte("done"))

}

func handleOEMs(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("handleOEMs()")
	controller.GetOEMs()
	w.Write([]byte("done"))

}
