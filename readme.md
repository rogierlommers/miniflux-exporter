## miniflux-exporter

### installation

If you have a running Go environment:

```
go get github.com/rogierlommers/miniflux-exporter
```

If you just want the (linux, 64 bit) binary: [miniflux-exporter](https://github.com/rogierlommers/miniflux-exporter/releases/download/2/miniflux-exporter)

### usage
```
-host string
  	miniflux hostname, f.e. http://localhost:8080 (default"http://localhost:8080")
-o string
  	output filename, f.e. /tmp/opml.xml (default"opml.xml")
-pass string
  	miniflux password
-user string
  	miniflux username
```

### in your crontab
Put miniflux-exporter in your crontab to frequently make a backup of all your feeds, f.e.:

```
@weekly        /usr/bin/miniflux-exporter -user YOUR_NAME -pass YOUR_PASS -host http://miniflux2-server -o "/my-backups/miniflux-opml.xml"
```

This will backup once a week.
