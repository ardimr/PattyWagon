export GOOSE_DRIVER=postgres
export GOOSE_MIGRATION_DIR=db/sql/migrations
export GOOSE_DBSTRING="user=postgres password=postgres dbname=patty-wagon-dev host=localhost port=5432 sslmode=disable"

export PORT=8080
export APP_ENV=local
export DB_HOST=localhost
export DB_PORT=5432
export DB_DATABASE=patty-wagon-dev
export DB_USERNAME=postgres
export DB_PASSWORD=postgres
export DB_SCHEMA=public

export DB_MAX_OPEN_CONNS=20
export DB_MAX_IDLE_CONNS=5
export DB_CONN_MAX_IDLE_TIME_IN_SECONDS=60
export DB_CONN_MAX_LIFE_TIME_IN_SECONDS=300

export JWT_SIGNATURE_KEY=solidteam

export S3_ACCESS_KEY_ID=team-solid
export S3_SECRET_ACCESS_KEY=@team-solid
export S3_ENDPOINT=localhost:9000
export S3_BUCKET=images
export S3_MAX_CONCURRENT_UPLOAD=5

# Image Compression
export MAX_CONCURRENT_COMPRESS=10

export OTLP_ENDPOINT=localhost:4317
