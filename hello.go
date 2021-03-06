package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	// disable timestamps in the log format -- that's handled by syslog
	log.SetFlags(0)
}

func withLogs(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf(":%s", port)
}

func main() {
	for _, x := range os.Environ() {
		log.Println("ENVIRON", x)
	}
	http.Handle("/", http.FileServer(http.Dir("assets")))
	addr := getAddr()
	log.Printf("Listening on %s", addr)
	http.HandleFunc("/cool", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(
			`<DOCTYPE html>
<html>
<head>
<title>You're cool.</title>
</head>
<body>
<h2>You seem cool and nice</h2>
</body>
</html>
`))
	})
	if err := http.ListenAndServe(addr, withLogs(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
