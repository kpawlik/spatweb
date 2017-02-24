# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang



RUN apt-get update
RUN apt-get install --yes libsqlite3-dev sqlite3
RUN apt-get install --yes libspatialite-dev spatialite-bin

ENV PATH /usr/lib/x86_64-linux-gnu:$PATH



# Copy the local package files to the container's workspace.
ADD . /go/src/fileserver
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/gorilla/mux

RUN go install fileserver

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/fileserver

# Document that the service listens on port 8080.
EXPOSE 8080