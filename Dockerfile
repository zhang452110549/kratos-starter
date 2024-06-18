FROM golang:1.22-alpine AS build-env
ARG BRANCH_NAME

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /build
COPY . /build
RUN make build-$BRANCH_NAME

FROM alpine:3.17
ARG PROJECT_NAME
ARG BRANCH_NAME
ENV PROJECT_NAME=$PROJECT_NAME
ENV BRANCH_NAME=$BRANCH_NAME
WORKDIR /opt

COPY --from=build-env /build/$PROJECT_NAME /opt/
COPY --from=build-env /build/configs /opt/configs

EXPOSE 9000
CMD /opt/$PROJECT_NAME -conf /opt/configs/$BRANCH_NAME.yaml