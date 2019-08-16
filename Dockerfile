FROM golang:latest as builder
WORKDIR /go/src/github.com/9thGear/fc2-mock-bank
COPY * /go/src/github.com/9thGear/fc2-mock-bank/
RUN make deps
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/9thGear/fc2-mock-bank/fc2-mock-bank.macos .
EXPOSE 8080
CMD ["./fc2-mock-bank.macos"]

# To start the mockbank-api in a Docker container run:
#   docker run --name mockbank-api --rm --tty --interactive --publish 5000:5000 --detach \
#        --env DB_HOST="localhost" \
#        --env DB_PORT="5432"  \
#        --env DB_USERNAME="bankwriter" \
#        --env DB_PASSWORD="change_me" \
#        --env DB_NAME="mockbank" \
#        --env API_PORT=5000 \
#        mockbank-api
#
# On MacOS use:
#      --env DB_HOST="host.docker.internal"