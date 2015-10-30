Simple Logger for GO
====

###How to import
```
import (
	"github.com/gorward/logger"
)
```

###Sample Usage
```
// Creating intance of logger
log = logger.New(logger.Config{
	Level:  logger.All,		// Values: All, Access, Info, Debug, Warn, Error, None
	Err:    "error.log",	// File name assignment where logs will be saved
	Warn:   "error.log",
	Debug:  "rec.log",
	Info:   "rec.log",
	Access: "access.log",
})

log.Error("Error", logger.Data{
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
})
```

###Sample Usage for Access log level
```
var log *logger.Logger

func main() {
	http.ListenAndServe(":3000", SampleLogMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
```