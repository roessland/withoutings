oapi-codegen \
  -generate types,client \
  -package openapi \
  withings-swagger-v3.0.3-fixed.json \
   > ../generated.go
