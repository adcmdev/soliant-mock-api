# Builder
FROM golang:1.22 AS builder

ARG GH_PACKAGES_TOKEN
ENV GOPRIVATE=github.com/adcmdev/*
ENV GONOPROXY=github.com/adcmdev/*
ENV GONOSUMDB=github.com/adcmdev/*

RUN git config --global url."https://${GH_PACKAGES_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOARCH=arm go build -a -o app ./cmd/main.go

# App
FROM scratch

COPY --from=builder /app/app /app

CMD ["/app"]
