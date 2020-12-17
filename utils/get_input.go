package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	dayFlag       = flag.Int("day", 0, "day num (or current if unset)")
	yearFlag      = flag.Int("year", 0, "year (or current if unset)")
	sessionCookie = flag.String("session", "", "session cookie")
	sessionFile   = flag.String("session_file", "",
		"file containing session cookie")
)

func fetchURL(url, sessionCookie string) ([]byte, error) {
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionCookie,
	}

	req, err := http.NewRequestWithContext(
		context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.AddCookie(cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed request: %v", resp.Status)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	return contents, nil
}

func main() {
	flag.Parse()

	var cookieValue string
	if *sessionCookie != "" {
		cookieValue = *sessionCookie
	} else if *sessionFile != "" {
		cArr, err := ioutil.ReadFile(*sessionFile)
		if err != nil {
			log.Fatalf("failed to read cookie from %v: %v",
				*sessionFile, err)
		}
		cookieValue = strings.TrimSpace(string(cArr))
	} else {
		log.Fatalf("--session(_file) is required with --url")
	}

	day := *dayFlag
	if day == 0 {
		day = time.Now().Day()
	}

	year := *yearFlag
	if year == 0 {
		year = time.Now().Year()
	}

	fmt.Fprintf(os.Stderr, "Using year %v day %v\n", year, day)

	url := fmt.Sprintf(`https://adventofcode.com/%v/day/%v/input`,
		year, day)

	contents, err := fetchURL(url, cookieValue)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(contents)
}
