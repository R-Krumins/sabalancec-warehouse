# 📦🥦 Sabalancec Warehouse

> [!IMPORTANT]
> Watch this [video](https://www.youtube.com/watch?v=8lGpZkjnkt4) to learn how to make pull requests to this repo.

# 🛠️ How to run

## With Docker

Define `PORT` in .env

Run:

```bash
docker compose up
```

## Without Docker

Define `PORT` and `DB_PATH` in .env

Run:

```bash
go build -o warehouse ./src && ./warehouse
```
