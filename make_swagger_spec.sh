#!/bin/bash

# use goas by mikunalpha
goas --module-path . --main-file-path ./api/server.go --output docs/swagger.json
cp docs/swagger.json docs/swaggerui/swagger.json