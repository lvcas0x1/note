---
date: 2026-09-22
---
### installation

npm install @openapitools/openapi-generator-cli -g

### openapi-generator

npx openapi-generator-cli generate -i openapi/schema/schema.yaml -g typescript-fetch -o openapi/output/

### redoc

npm i -g @redocly/cli@latest

### run redoc

redocly build-docs openapi/schema/schema.yaml -o openapi/redoc/index.html
