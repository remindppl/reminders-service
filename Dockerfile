FROM golang:1.19 as sheets-to-db-builder
COPY ./sheets-to-db /app/src/
WORKDIR /app/src
RUN go mod tidy
RUN go fmt .
RUN go vet .
RUN go build -o sheets-to-db .

FROM gcr.io/distroless/base-debian11
WORKDIR /app/
COPY --from=sheets-to-db-builder /app/src/sheets-to-db /app/sheets-to-db
ENTRYPOINT ["/app/sheets-to-db"]
