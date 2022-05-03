oapi-codegen \
  -generate types,client \
  -package openapi2 \
  withings-swagger-v3.0.3-fixed.json \
   > openapi2/generated.go
