package handler

import (
	"fmt"
	"net/http"
)

func OutWriter(out chan string, w http.ResponseWriter, exit chan int) {
	f, flushable := w.(http.Flusher)

	for s := range out {
		fmt.Fprintf(w, "%s", s)
		if flushable {
			f.Flush()
		}
	}
	exit <- 0
}

func IgnoreAll(out chan string) {
	for _ = range out {
	}
}
