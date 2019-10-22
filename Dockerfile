FROM golang:1.13 as one

WORKDIR /github.com/fasmide/god

# Download dependencies and hopefully cache them
ADD go.* /github.com/fasmide/god/
RUN go mod download

ADD . /github.com/fasmide/god
RUN CGO_ENABLED=0 go build -o /god

FROM scratch
COPY --from=one /god /god