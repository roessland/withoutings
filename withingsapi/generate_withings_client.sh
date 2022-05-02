docker run --rm \
    -v "${PWD}:/local" \
    --env GO_POST_PROCESS_FILE="/usr/local/bin/gofmt -w" \
    openapitools/openapi-generator-cli generate \
    -i /local/withings-swagger-v3.0.3-fixed.json \
    -g go \
    -o /local/internal