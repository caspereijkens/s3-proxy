#!/bin/bash

# Usage: ./minio-upload my-bucket my-file.zip 8080

bucket=$1
file=$2
port=$3

host=localhost:$port
s3_key=minioadmin
s3_secret=minioadmin
filename=$(basename "$file")
resource="/${bucket}/${filename}"
content_type="application/octet-stream"
date=`date -R`
_signature="PUT\n\n${content_type}\n${date}\n${resource}"
signature=`echo -en ${_signature} | openssl sha1 -hmac ${s3_secret} -binary | base64`

curl -X PUT -T "${file}" \
          -H "Host: ${host}" \
          -H "Date: ${date}" \
          -H "Content-Type: ${content_type}" \
          -H "Authorization: AWS ${s3_key}:${signature}" \
          http://${host}${resource} 
          
