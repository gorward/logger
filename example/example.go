package main

import (
	"fmt"
	"github.com/gorward/logger"
	"net/http"
	"time"
)

var log *logger.Logger

func main() {
	log = logger.New(logger.Config{
		Level:  logger.All,
		Err:    "error.log",
		Warn:   "error.log",
		Debug:  "rec.log",
		Info:   "rec.log",
		Access: "access.log",
	})

	/*log.Error("Error", logger.Data{
		"data": "Data details here",
	})

	log.Warn("Warning", logger.Data{
		"data": "Data details here",
		"code": 123,
	})

	log.Debug("Debug", logger.Data{
		"data": "Data details here",
		"code": 123,
	})

	log.Info("Info", logger.Data{
		"data": "Data details here",
		"code": 123,
	})*/

	http.ListenAndServe("10.10.13.25:3000", SampleLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "helllo")
	})))

}

func SampleLogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			log.Access(start, w, r)
		}()

		h.ServeHTTP(w, r)
	})

}
