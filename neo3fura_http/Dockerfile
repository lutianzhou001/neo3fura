FROM golang:1.21

ENV GO111MODULE="on"

ENV GOPROXY="https://goproxy.io"

ARG RT

RUN echo $RT

ENV RUNTIME=$RT

RUN mkdir application

COPY . ./application

WORKDIR "application"

RUN  wget https://packages.microsoft.com/config/ubuntu/21.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb

RUN  dpkg -i packages-microsoft-prod.deb

RUN  rm packages-microsoft-prod.deb

RUN  apt-get clean && \
     apt-get update

RUN  wget wget http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.1f-1ubuntu2_amd64.deb && \
     dpkg -i libssl1.1_1.1.1f-1ubuntu2_amd64.deb

RUN  apt-get install -y apt-transport-https

RUN  apt-get install -y dotnet-sdk-5.0 && \
     apt-get install -y aspnetcore-runtime-5.0

RUN  tar -xvf compiler.tar.gz

RUN  chmod 777 compiler/a/nccs && \
     chmod 777 compiler/c/nccs && \
     chmod 777 compiler/b/nccs

RUN  go build -o main app/neo3fura/src.go

EXPOSE 1926

CMD ["./main"]