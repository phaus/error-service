#!/bin/bash

docker build --pull -t phaus/error-service .

docker images | grep "phaus/error-service"