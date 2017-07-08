package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(os.Getenv("API_HOST"))

	if err == nil && res != nil {
		if res.StatusCode > 400 {
			err = fmt.Errorf("Error: %d", res.StatusCode)
		} else {

			body, err := ioutil.ReadAll(res.Body)

			if err == nil {
				fmt.Fprint(w, string(string(body)))
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprint(w, err)
}

func main() {
	goji.Get("/*", ping)
	goji.Serve()
}
