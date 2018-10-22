package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("I am a sample CloudKinds provider that does absolutely nothing useful! ☁️🌈")
	fmt.Println("Listening on *:8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handled %s\n", r.URL.Path)
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "%s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("\t%v\n", string(b))
		fmt.Fprintf(w, "%v", string(b))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
