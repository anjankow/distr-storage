FROM golang:1.17


ARG SERVICE_NAME=
ENV SERVICE_NAME=${SERVICE_NAME}

WORKDIR ${GOPATH}/${SERVICE_NAME}/

# RUN go build -o ${SERVICE_NAME} cmd/main.go 
# RUN chmod +x ${SERVICE_NAME} && cp ${SERVICE_NAME} /out

RUN go install github.com/cosmtrek/air@latest
CMD ${GOPATH}/bin/air -c .air.toml
