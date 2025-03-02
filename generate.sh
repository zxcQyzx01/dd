#!/bin/bash

# Создаем директории для сгенерированного кода
mkdir -p pkg/{auth,user,geo}

# Генерируем код из proto файлов
protoc --go_out=. \
       --go_opt=module=dd \
       --go-grpc_out=. \
       --go-grpc_opt=module=dd \
       proto/auth.proto

protoc --go_out=. \
       --go_opt=module=dd \
       --go-grpc_out=. \
       --go-grpc_opt=module=dd \
       proto/user.proto

protoc --go_out=. \
       --go_opt=module=dd \
       --go-grpc_out=. \
       --go-grpc_opt=module=dd \
       proto/geo.proto 