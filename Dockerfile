FROM ubuntu
LABEL description="Resume from Rogier Lommers"
LABEL maintainer="Rogier Lommers <rogier@lommers.org>"

# add binary and assets
COPY --chown=1000:1000 ./bin/resume /resume/
COPY --chown=1000:1000 ./src/assets /assets

# binary will serve on 8080
EXPOSE 8080

# make binary executable
RUN chmod +x /resume/resume

# run binary
CMD ["/resume/resume"]
