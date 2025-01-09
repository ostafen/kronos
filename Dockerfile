FROM node:20 AS ui-builder

WORKDIR /app
COPY ui/package.json ui/package-lock.json ./

RUN npm install

COPY ui/ ./
RUN npm run build

FROM golang:1.23.4-bullseye AS build

ARG VERSION
ARG COMMIT
ARG BUILD_TIME

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
COPY --from=ui-builder /app/web ./webbuild/web
RUN go build -a -installsuffix cgo -ldflags "-w -s -X main.version=$VERSION -X main.commit=$COMMIT -X 'main.buildTime=$BUILD_TIME'" -buildvcs=false -o /kronos cmd/main.go

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /kronos /kronos
ENTRYPOINT ["/kronos"]
