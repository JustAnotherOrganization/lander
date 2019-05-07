package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"justanother.org/lander/data"
)

const (
	fromURL = "justanother.org"
	toURL   = "github.com/justanotherorganization"
)

// Network Bindings
var (
	flagPort int
	flagHost string

	//GitComHash and BuildStamp are used to display version and build time
	GitComHash = "undefined"
	//BuildStamp is used to display  build time
	BuildStamp = "empty"
)

func init() {
	// Check if we want to display the version information.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "-v":
			os.Stdout.WriteString("Git Commit Hash: " + GitComHash + "\n")
			os.Stdout.WriteString("Build Time: " + BuildStamp + "\n")
			os.Exit(0)
			return
		}
	}

	flag.IntVar(&flagPort, "port", 8080, "Specify the port to run the application.")
	flag.StringVar(&flagHost, "host", "", "Specify the host to bind the application to.")
	flag.Parse()
}

func constructPage(url string) (page string) {
	githubURL := fmt.Sprintf("https://%v/%v", toURL, strings.Replace(url, "/", "-", -1))
	return fmt.Sprintf(`<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
			<meta name="go-import" content="justanother.org/%v git %s">
			<meta http-equiv="refresh" content="0; url=%s">
		</head>
	</html>`, url, githubURL, githubURL)
}

func handler(rw http.ResponseWriter, r *http.Request) {
	url := r.URL.String()[1:len(r.URL.String())]

	if url == "" {
		if bytes, _ := data.Asset("index.html"); bytes != nil {
			rw.WriteHeader(http.StatusOK)
			rw.Write(bytes)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rw, "Page not found!")
		}
		return
	}

	if bytes, _ := data.Asset(url); bytes != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write(bytes)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(rw, constructPage(url))
}

func main() {
	fmt.Printf("Starting server %v:%v\n", flagHost, flagPort)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", flagHost, flagPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
