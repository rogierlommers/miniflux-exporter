FROM ubuntu
LABEL description="Miniflux exporter"
LABEL maintainer="Rogier Lommers <rogier@lommers.org>"

# add binary and assets
COPY --chown=1000:1000 ./bin/miniflux-exporter /app

# make binary executable
RUN chmod +x /app/miniflux-exporter

# run binary
CMD ["/app/miniflux-exporter"]
