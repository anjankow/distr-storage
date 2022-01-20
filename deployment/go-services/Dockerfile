FROM golang:1.17


ARG SERVICE_NAME=

WORKDIR ${GOPATH}/${SERVICE_NAME}/

COPY ${SERVICE_NAME}/ .
RUN go build -o ${SERVICE_NAME} cmd/main.go 
RUN chmod +x ${SERVICE_NAME} && cp ${SERVICE_NAME} /out

ENTRYPOINT [ "/out" ]
