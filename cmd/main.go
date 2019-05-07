package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"google.golang.org/appengine"
	"justanother.org/lander/data"
)

const (
	fromURL = "justanother.org"
	toURL   = "github.com/justanotherorganization"
)

// Network Bindings
var (
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

}

func constructPage(url string) (page string) {
	// We don't have multi later repositories, so only get the first item if
	// seperated
	urlParts := strings.Split(url, "/")
	url = strings.ToLower(urlParts[0])

	githubURL := fmt.Sprintf("https://%v/%v", toURL, url)

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

	cleanURL(r.URL)

	if r.URL.Path == "" {
		if bytes, _ := data.Asset("index.html"); bytes != nil {
			rw.WriteHeader(http.StatusOK)
			rw.Write(bytes)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Fprint(rw, "Page not found!")
		}
		return
	}

	if bytes, _ := data.Asset(r.URL.Path); bytes != nil {
		rw.WriteHeader(http.StatusOK)
		rw.Write(bytes)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(rw, constructPage(r.URL.Path))
}

func main() {
	http.HandleFunc("/", handler)
	appengine.Main() // Start the server
}

func cleanURL(url *url.URL) {
	if url == nil {
		return
	}

	// Sanitize
	switch {
	case strings.HasPrefix(url.Path, "https://"):
		url.Path = url.Path[8:]
	case strings.HasPrefix(url.Path, "http://"):
		url.Path = url.Path[7:]
	case strings.HasPrefix(url.Path, "/"):
		url.Path = url.Path[1:]
	}

	return
}
