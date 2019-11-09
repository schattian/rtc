FROM golang:1.13

LABEL maintainer="Sebastián Chamena <sebachamena@gmail.com>"

ARG env=development
 
ENV ENV $env
ENV WORKDIR /git-crud

WORKDIR ${WORKDIR} 

COPY go.* ./

RUN go mod download

COPY . .

RUN go build .

CMD if [ ${ENV} = development ]; \
	then \
    go get github.com/pilu/fresh && \
	fresh; \
	fi

EXPOSE 8888 