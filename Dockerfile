FROM dlckdgk4858/go-face:1.0

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY main.go main.go
COPY go.mod go.mod

#ENTRYPOINT ["go","run", "."]
