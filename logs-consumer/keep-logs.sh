#!/bin/bash

mkdir -p logs
cd logs
touch requests.log
echo "Starting logging to file requests.log"
/logs-client error ok &>> requests.log