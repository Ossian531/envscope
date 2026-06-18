# build
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /frontend .

# run
FROM gcr.io/distroless/static-debian12
COPY --from=build /frontend /frontend
ENV FRONTEND_ADDR=:3000
ENV BACKEND_URL=http://backend:8080/
EXPOSE 3000
ENTRYPOINT ["/frontend"]
