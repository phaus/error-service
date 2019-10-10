#!/bin/bash

echo "open http://127.0.0.1:9000"

docker run -p 9000:9000 phaus/error-service

