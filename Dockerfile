#sudo docker build -t centos7-gopub .
#sudo docker run --name centos7-gopub   -d   centos7-gopub:latest
From 192.168.176.2:81/test/centos7-gopub
MAINTAINER Linc "13579443@qq.com"
ADD src/ /data/gopub/src/
ADD control /data/gopub/control
RUN  cd /data/gopub  &&./control build
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
