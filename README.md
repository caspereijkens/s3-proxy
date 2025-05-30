# s3-proxy
This is a proxy for s3 uploads. The point is to be able to migrate s3 object stores without the uploading party having to migrate directly. 

## Docker demo 
```
# Create a docker network
docker network create my-bridge-network

# Run the minio backend
docker run --name minio -p 9001:9001 --network my-bridge-network \
  quay.io/minio/minio server /data --console-address ":9001"

# Login to MinIO Console (localhost:9001) and create bucket (or use mc)
mc alias set local-minio localhost:9001 minioadmin minioadmin
mc mb proxy-bucket

# Run the proxy
docker run -p 8080:80 -e MINIO_ACCESS_KEY_ID=minioadmin -e MINIO_SECRET_ACCESS_KEY=minioadmin -e MINIO_ENDPOINT=minio:9000 -e MINIO_BUCKET_NAME=proxy-bucket --network my-bridge-network s3-proxy:v0.0.1
```

