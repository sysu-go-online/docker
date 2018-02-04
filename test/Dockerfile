# 运行go
FROM golang:latest

# 创建服务端运行环境
WORKDIR app

# 添加文件
ADD . app

RUN ["ls", "-l"]

# 将8080端口暴露出来
#EXPOSE 8080

# 添加GOPATH环境后，容器运行
CMD ["go", "run", "app/helloworld.go"]