package main

import (
	"flag"
	"io/ioutil"

	humanize "github.com/dustin/go-humanize"
	"github.com/miniflux/miniflux-go"
	"github.com/sirupsen/logrus"
)

var (
	targetFile string
	username   string
	password   string
	hostname   string
	silent     bool
)

func main() {
	flag.StringVar(&targetFile, "o", "opml.xml", "output filename, f.e. /tmp/opml.xml")
	flag.StringVar(&username, "user", "", "miniflux username")
	flag.StringVar(&password, "pass", "", "miniflux password")
	flag.StringVar(&hostname, "host", "http://localhost:8080", "miniflux hostname, f.e. http://localhost:8080")
	flag.BoolVar(&silent, "s", false, "if flag -s is provided, the happy-flow won't display any output")
	flag.Parse()

	// start export
	client := miniflux.NewClient(hostname, username, password)

	opmlFile, err := client.Export()
	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile(targetFile, opmlFile, 0644)
	if err != nil {
		logrus.Fatal(err)
	}

	if !silent {
		logrus.Infof("export done, %s written to file %s", humanize.Bytes(uint64(len(opmlFile))), targetFile)
	}

}
