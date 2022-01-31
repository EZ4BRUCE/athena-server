# athena

docker run -p 27018:27017 --name mongo-docker -v mongo-data:/data -d mongo

 protoc --go_out=plugins=grpc:./proto ./proto/*.proto