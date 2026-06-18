# build
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /backend .

# run
FROM gcr.io/distroless/static-debian12
COPY --from=build /backend /backend
ENV BACKEND_ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["/backend"]
