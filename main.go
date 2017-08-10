package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	bindPort   string
	backendUrl string
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func FrontendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("FrontendHandler\n")
	logRequestInfo(r)

	body := getRequestBody(r)
	fmt.Printf("BODY: %s\n", body)

	backendStatusCode, backendResponse := sendRequestToBackend(body)

	response := fmt.Sprintf("%s | %d | %s\n", getRequestInfo(r), backendStatusCode, backendResponse)
	fmt.Fprintf(w, response)
}

func BackendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("BackendHandler\n")
	logRequestInfo(r)

	body := getRequestBody(r)
	fmt.Printf("BODY: %s\n", body)

	response := fmt.Sprintf("%s %s", getRequestInfo(r), body)
	fmt.Fprintf(w, response)
}

func getRequestBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return []byte("ERR_READING_BODY")
	}
	return body
}

func logRequestInfo(r *http.Request) {
	fmt.Printf("REQUEST: %s %s %s %s\n",
		r.Host,
		r.URL,
		r.Method,
		r.Proto,
	)
}

func sendRequestToBackend(body []byte) (int, string) {
	fmt.Printf("%s\n", backendUrl)
	rs, err := http.Post(backendUrl+"/backend", "text/html", bytes.NewBuffer(body))
	defer rs.Body.Close()
	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}
	return rs.StatusCode, string(bodyBytes)
}

func getRequestInfo(r *http.Request) string {
	t := time.Now()
	return fmt.Sprintf("%s \"%s%s\"", t.Format("2006-01-02T15:04:05"), r.Host, r.URL.Path)
}

func main() {
	bindPort = os.Getenv("PORT")
	backendUrl = os.Getenv("BACKEND_URL")

	fmt.Printf("bindPort %s\n", bindPort)
	fmt.Printf("backendUrl %s\n", backendUrl)

	http.HandleFunc("/hello", HelloHandler)
	http.HandleFunc("/frontend", FrontendHandler)

	// All other requests will simply reply with whatever data was sent
	http.HandleFunc("/", BackendHandler)

	fmt.Printf("Listening\n")
	http.ListenAndServe(":"+bindPort, nil)
}
