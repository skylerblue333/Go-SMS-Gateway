FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /sky-sms-envelope .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /sky-sms-envelope /sky-sms-envelope
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/sky-sms-envelope"]
