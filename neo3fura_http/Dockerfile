FROM golang:1.17

ENV GO111MODULE="on"

ENV GOPROXY="https://goproxy.cn"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  wget https://packages.microsoft.com/config/ubuntu/21.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb

RUN  dpkg -i packages-microsoft-prod.deb

RUN  rm packages-microsoft-prod.deb

RUN  apt-get update

RUN  apt-get install -y apt-transport-https

RUN  apt-get update

RUN  apt-get install -y dotnet-sdk-5.0

RUN  apt-get install -y aspnetcore-runtime-5.0

RUN  tar -xvf compiler.tar.gz

RUN  chmod 777 compiler/a/nccs

RUN  chmod 777 compiler/c/nccs

RUN  chmod 777 compiler/b/nccs

RUN  go build -o main app/neo3fura/src.go

EXPOSE 1926

CMD ["./main"]