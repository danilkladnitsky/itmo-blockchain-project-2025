~/go/bin/swag fmt  -d cmd/analyze-service,\
internal/handlers/check_alive,\
internal/handlers/analyze,\
internal/storage/moralis,\
internal/lib/api/error,\
internal/lib/api/response

~/go/bin/swag init -d cmd/analyze-service,\
internal/handlers/check_alive,\
internal/handlers/analyze,\
internal/storage/moralis,\
internal/lib/api/error,\
internal/lib/api/response