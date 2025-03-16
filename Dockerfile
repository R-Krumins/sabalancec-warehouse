FROM golang:1.24

# Install SQLite dependencies
RUN apt-get update && apt-get install -y sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*

# Create data directory for SQLite database with proper permissions
RUN mkdir -p /data && chmod 755 /data

WORKDIR /usr/src/warehouse

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/warehouse ./src


ENV PORT=3000
ENV DB_PATH="/data/db.db"

EXPOSE 3000

CMD ["warehouse"]

