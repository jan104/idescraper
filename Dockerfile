FROM golang AS builder

WORKDIR $GOPATH/src/github.com/jan104/idescraper
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-w -extldflags -static -X main.GitCommit=$(git rev-parse HEAD) -X main.BuildTime=$(date -u --iso-8601=seconds)" -o /go/bin/idescraper

FROM alpine:3.15
COPY --from=builder /go/bin/idescraper /go/bin/idescraper
RUN chmod +x /go/bin/idescraper
USER guest
ENTRYPOINT ["/go/bin/idescraper"]


