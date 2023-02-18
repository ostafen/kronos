FROM golang:1.19-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -buildvcs=false -o /kronos cmd/main.go

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /kronos /kronos
USER nonroot:nonroot
ENTRYPOINT ["/kronos"]