#!/usr/bin/env bash

curl -X POST http://localhost:8080/v1/display-part -F 'file=@resources/singlepart.eml'