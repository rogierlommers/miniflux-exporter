## miniflux-exporter

### installation

If you just want the (linux, 64 bit) binary: [miniflux-exporter](https://github.com/rogierlommers/miniflux-exporter/releases/download/9/miniflux-exporter)

### usage
```
  -host string
    	miniflux hostname, f.e. http://localhost:8080 (default "http://localhost:8080")
  -output-opml string
    	optional, output filename, f.e. /tmp/opml.xml
  -output-stars string
    	optional, output filename, f.e. /tmp/starred-articles.xml
  -output-unread string
    	optional, output filename, f.e. /tmp/unread-articles.xml
  -pass string
    	miniflux password
  -s	if flag -s is provided, the happy-flow won't display any output
  -user string
    	miniflux username
  -v	if flag -v is provided, debugging info is printed
  -version
    	prints current version
```

### in your crontab
Put miniflux-exporter in your crontab to frequently make a backup of all your feeds, f.e.:

```
@weekly        /usr/bin/miniflux-exporter -s -user YOUR_NAME -pass YOUR_PASS -host http://miniflux2-server -output-opml "/my-backups/feeds-opml.xml" -output-stars "/my-backups/miniflux-starred-articles.xml" -output-unread "/my-backups/miniflux-unread-articles.xml"
```

This will backup once a week the starred items, the unread items and exports all feeds to an OPML file and will only display error messages. Please note that the different outputs are all options. So if you only want to export the starred articles, then you should only provide the `-output`-stars flag.

## If you have a working go environment

Building the binary locally:

```
GOOS=linux GOARCH=amd64 go build -o miniflux-exporter-linux64  *.go
```
