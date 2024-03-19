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

const (
	AppVersion        = "2024-march-12"
	ExportTypeStarred = "starred-articles"
	ExportTypeUnread  = "unread-articles"
)

var (
	targetOPMLFile    string
	targetStarredFile string
	targetUnreadsFile string
	username          string
	password          string
	hostname          string
	silent            bool
	verbose           bool
)

func main() {

	// parse filename flags
	flag.StringVar(&targetOPMLFile, "output-opml", "", "output filename, f.e. /tmp/opml.xml")
	flag.StringVar(&targetStarredFile, "output-stars", "", "output filename, f.e. /tmp/starred-articles.xml")
	flag.StringVar(&targetUnreadsFile, "output-unread", "", "output filename, f.e. /tmp/unread-articles.xml")

	// parse usenanme/pass/host
	flag.StringVar(&username, "user", "", "miniflux username")
	flag.StringVar(&password, "pass", "", "miniflux password")
	flag.StringVar(&hostname, "host", "http://localhost:8080", "miniflux hostname, f.e. http://localhost:8080")

	// parse options
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
	if len(targetStarredFile) > 0 {
		exportEntries(ExportTypeStarred, c)
	} else {
		message("skipping export of bookmarks/starred entries (see -help for more info)")
	}

	// start export unread entries
	if len(targetUnreadsFile) > 0 {
		exportEntries(ExportTypeUnread, c)
	} else {
		message("skipping export of unread entries (see -help for more info)")
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

func exportEntries(exportType string, c *client.Client) {

	// configuration of use case specific stuff
	var (
		feedTitle       string
		feedDescription string
		targetFilename  string
		minifluxFilter  *client.Filter
	)

	switch exportType {
	case "starred-articles":
		verboseMessage("applying settings for starred-articles")
		feedTitle = "Miniflux starred entries"
		feedDescription = "RSS feed from all starred entries in Miniflux"
		targetFilename = targetStarredFile
		minifluxFilter = &client.Filter{
			Starred: client.FilterOnlyStarred,
		}

	case "unread-articles":
		verboseMessage("applying settings for unread-articles")
		feedTitle = "Miniflux starred entries"
		feedDescription = "RSS feed from all starred entries in Miniflux"
		targetFilename = targetUnreadsFile
		minifluxFilter = &client.Filter{
			Status: client.EntryStatusUnread,
		}
	}

	// start actual export
	var number int

	now := time.Now()
	feed := &feeds.Feed{
		Title:       feedTitle,
		Description: feedDescription,
		Link:        &feeds.Link{Href: hostname},
		Created:     now,
	}

	entries, err := c.Entries(minifluxFilter)
	if err != nil {
		fmt.Println("error entries: " + err.Error())
		return
	}
	verboseMessage(fmt.Sprintf("%d entries fetched", len(entries.Entries)))

	for _, entry := range entries.Entries {

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

	rss, err := feed.ToRss()
	if err != nil {
		fmt.Println("error items to RSS feed: " + err.Error())
		return
	}
	verboseMessage(fmt.Sprintf("rss file created in-memory, size: %d bytes", len(rss)))

	err = os.WriteFile(targetFilename, []byte(rss), 0644)
	if err != nil {
		fmt.Println("error writing target file: " + err.Error())
		return
	}

	message(fmt.Sprintf("exported %d articles, %s written to file %s", number, humanize.Bytes(uint64(len(rss))), targetFilename))
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
