FROM golang:1.19-buster AS build

ARG VERSION
ARG COMMIT
ARG BUILD_TIME

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix cgo -ldflags "-w -s -X main.version=$VERSION -X main.commit=$COMMIT -X 'main.buildTime=$BUILD_TIME'" -buildvcs=false -o /kronos cmd/main.go

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /kronos /kronos
USER nonroot:nonroot
ENTRYPOINT ["/kronos"]