# adapted from https://docs.docker.com/language/golang/build-images/#multi-stage-builds
FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /authduck

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /authduck /authduck

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/authduck"]
