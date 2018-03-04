## miniflux-exporter

### installation

If you have a running Go environment:

```
go get github.com/rogierlommers/miniflux-exporter
```

If you just want the (linux, 64 bit) binary: [miniflux-exporter](https://github.com/rogierlommers/miniflux-exporter/releases/download/4/miniflux-exporter)

### usage
```
-host string
  	miniflux hostname, f.e. http://localhost:8080 (default"http://localhost:8080")
-output-bookmarks string
  	output filename, f.e. /tmp/bookmarks.xml
-output-opml string
  	output filename, f.e. /tmp/opml.xml
-pass string
  	miniflux password
-user string
  	miniflux username
-s	if flag -s is provided, the happy-flow won't display any output
```

### in your crontab
Put miniflux-exporter in your crontab to frequently make a backup of all your feeds, f.e.:

```
@weekly        /usr/bin/miniflux-exporter -s -user YOUR_NAME -pass YOUR_PASS -host http://miniflux2-server -output-bookmarks "/my-backups/miniflux-opml.xml" -output-opml "/my-backups/miniflux-bookmarks.xml"
```

This will backup once a week both starred items and exports to an OPML file and will only display error messages
