FROM dlckdgk4858/go-face:1.0

RUN mkdir /go/src/app
WORKDIR /go/src/app

RUN apt-get install vim -y

COPY main.go main.go
COPY common.go common.go
COPY signup-face.go signup-face.go
COPY signin-face.go signin-face.go
COPY models models
COPY go.mod go.mod

#ENTRYPOINT ["go","run", "."]
