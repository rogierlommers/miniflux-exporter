FROM alpine
LABEL description="Miniflux exporter"
LABEL maintainer="Rogier Lommers <rogier@lommers.org>"

# needed to do GETs through https
RUN apk update
RUN apk --no-cache add tzdata zip ca-certificates && update-ca-certificates

# add binary and assets
COPY --chown=1000:1000 ./bin/miniflux-exporter /app

# make binary executable
RUN chmod +x /app/miniflux-exporter

# run binary
CMD ["/app/miniflux-exporter"]
