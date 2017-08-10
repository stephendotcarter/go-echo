#!/bin/bash

while true
do
curl -s ${FRONTEND_URL}/frontend -d "ECHO" 2>&1 | tee client.log
sleep 1
done
