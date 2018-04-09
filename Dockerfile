#docker 17.05+版本支持
#sudo docker build -t  gopub .
#sudo docker run --name gopub -p 8192:8192  --restart always  -d   gopub:latest 
FROM golang:1.10.1-alpine3.7 as golang
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add  bash  && \ 
    rm -rf /var/cache/apk/*   /tmp/*     
ADD src/ /data/gopub/src/
ADD control /data/gopub/control
WORKDIR /data/gopub/
RUN ./control build

FROM node:9.11.1-alpine as node 
ADD ./ /data/gopub/
WORKDIR /data/gopub/vue-gopub
RUN npm install && npm run build

FROM alpine:3.7
MAINTAINER Linc "13579443@qq.com"
ENV TZ='Asia/Shanghai' 
RUN TERM=linux && export TERM
USER root 
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates bash tzdata sudo curl wget openssh git && \ 
    echo "Asia/Shanghai" > /etc/timezone && \
    cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /var/cache/apk/*   /tmp/*  && \ 
    mkdir -p /data/htdocs && \
    mkdir -p /data/logs && \
    ssh-keygen -q -N "" -f /root/.ssh/id_rsa && \
    #输出的key需要加入发布目标机的 ~/.ssh/authorized_keys
    cat ~/.ssh/id_rsa.pub  
WORKDIR /data/gopub
ADD control /data/gopub/control
COPY --from=golang /data/gopub/src/gopub /data/gopub/src/gopub
COPY --from=golang /data/gopub/src/conf /data/gopub/src/conf
COPY --from=golang /data/gopub/src/logs /data/gopub/src/logs
COPY --from=golang /data/gopub/src/agent /data/gopub/src/agent
COPY --from=node /data/gopub/src/views /data/gopub/src/views
COPY --from=node /data/gopub/src/static /data/gopub/src/static
CMD ["./control","rundocker"]
