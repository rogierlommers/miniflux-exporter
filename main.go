package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/feeds"
	"miniflux.app/client"
)

const AppVersion = "2024-march-12"

var (
	targetOPMLFile     string
	targetBookmarkFile string
	username           string
	password           string
	hostname           string
	silent             bool
	verbose            bool
)

func main() {

	// parse flags
	flag.StringVar(&targetOPMLFile, "output-opml", "", "output filename, f.e. /tmp/opml.xml")
	flag.StringVar(&targetBookmarkFile, "output-bookmarks", "", "output filename, f.e. /tmp/bookmarks.txt")
	flag.StringVar(&username, "user", "", "miniflux username")
	flag.StringVar(&password, "pass", "", "miniflux password")
	flag.StringVar(&hostname, "host", "http://localhost:8080", "miniflux hostname, f.e. http://localhost:8080")
	flag.BoolVar(&silent, "s", false, "if flag -s is provided, the happy-flow won't display any output")
	flag.BoolVar(&verbose, "v", false, "if flag -v is provided, debugging info is printed")
	version := flag.Bool("version", false, "prints current version")
	flag.Parse()

	// version
	if *version {
		fmt.Println("miniflux-exporter, version " + AppVersion + ".")
		os.Exit(0)
	}

	// verbose
	if verbose {
		fmt.Println("running in verbose mode")
	}

	// get miniflux client
	c := client.New(hostname, username, password)

	// start export to opml
	if len(targetOPMLFile) > 0 {
		exportOPML(c)
	} else {
		message("skipping opml export (see -help for more info)")
	}

	// start export bookmarks/starred entries
	if len(targetBookmarkFile) > 0 {
		exportStarredEntries(c)
	} else {
		message("skipping export of bookmarks/starred entries (see -help for more info)")
	}

}

func exportOPML(c *client.Client) {
	opmlFile, err := c.Export()
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	verboseMessage("opmlFile fetched from miniflux")

	err = os.WriteFile(targetOPMLFile, opmlFile, 0644)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return
	}
	verboseMessage("opmlFile written to disk")

	message(fmt.Sprintf("export OPML done, %s written to file %s", humanize.Bytes(uint64(len(opmlFile))), targetOPMLFile))
}

func exportStarredEntries(c *client.Client) {
	var number int

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Miniflux starred entries",
		Description: "RSS feed from all starred entries in Miniflux",
		Link:        &feeds.Link{Href: hostname},
		Created:     now,
	}

	entries, err := c.Entries(&client.Filter{})
	if err != nil {
		fmt.Println("error fetching starred entries: " + err.Error())
		return
	}
	verboseMessage(fmt.Sprintf("%d starred/bookmarked entries fetched", len(entries.Entries)))

	for _, entry := range entries.Entries {
		if entry.Starred {

			newItem := feeds.Item{
				Title:       entry.Title,
				Link:        &feeds.Link{Href: entry.URL},
				Author:      &feeds.Author{Name: entry.Author},
				Description: entry.Content,
				Id:          strconv.Itoa(int(entry.ID)),
			}

			feed.Items = append(feed.Items, &newItem)
			number++
			verboseMessage(fmt.Sprintf("Entry with ID '%s' added", newItem.Id))
		}
	}

	rss, err := feed.ToRss()
	if err != nil {
		fmt.Println("error exporting starred items to RSS feed: " + err.Error())
		return
	}
	verboseMessage(fmt.Sprintf("rss file created in-memory, size: %d bytes", len(rss)))

	err = os.WriteFile(targetBookmarkFile, []byte(rss), 0644)
	if err != nil {
		fmt.Println("error writing target bookmark file: " + err.Error())
		return
	}

	message(fmt.Sprintf("export %d bookmarks done, %s written to file %s", number, humanize.Bytes(uint64(len(rss))), targetBookmarkFile))
}

func message(m string) {
	if !silent {
		fmt.Println(m)
	}
}

func verboseMessage(m string) {
	if verbose {
		fmt.Println(m)
	}
}
