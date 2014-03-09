package main

import (
	"flag"
	"fmt"
	"github.com/norcalli/megashares"
	"io"
	"log"
	"net/http"
	"os"
)

func init() {
}

func main() {
	usernamePt := flag.String("u", "", "Megashares username.")
	passwordPt := flag.String("p", "", "Megashares password.")
	queryPt := flag.String("q", "", "Query to search for.")
	nPt := flag.Int("n", -1, "Which result to download. n < 0 => Print results")

	flag.Parse()
	username, password, query, n := *usernamePt, *passwordPt, *queryPt, *nPt
	if username == "" || password == "" || query == "" {
		flag.Usage()
		os.Exit(0)
	}

	m := megashares.New()
	if err := m.Login(username, password); err != nil {
		log.Fatalf("Couldn't login! Reason: %s\n", err)
	}
	entries, _ := m.SearchEntries(query)
	if n < 0 {
		for i, entry := range entries {
			fmt.Printf("%d: %s\n", i, entry.Filename)
		}
	} else {
		n = n % len(entries) // I'm lazy and don't feel like checking for invalid.
		entry := entries[n]
		fmt.Printf("Downloading entry %d: %s...\n", n, entry.Filename)
		if err := DownloadEntry(entry); err != nil {
			log.Fatalf("Failed to download because %s\n", err)
		}
	}
}

func DownloadEntry(entry *megashares.Entry) error {
	if response, err := http.Get(entry.Url); err != nil {
		return err
	} else {
		defer response.Body.Close()
		if file, err := os.OpenFile(entry.Filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666); err != nil {
			return err
		} else {
			if _, err := io.Copy(file, response.Body); err != nil {
				return err
			}
			return nil
		}
	}
}
