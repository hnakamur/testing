package main

import (
	"flag"
	"net/http"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	local := flag.String("l", ":8000", "listening address")
	repeat := flag.Int("repeat", 10, "count to repeat output")
	interval := flag.Int("interval", 1, "interval to loop output")

	flag.Parse()

	http.ListenAndServe(*local, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Cache-Control", "max-age=30")
		w.WriteHeader(http.StatusOK)
		for i := 0; i < *repeat; i++ {
			w.Write([]byte("hello\n"))
			if w, ok := w.(http.Flusher); ok {
				w.Flush()
			}
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}))
}
