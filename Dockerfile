FROM centos:centos7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY docker/tpl-go-http /var/docker/tpl-go-http
RUN mkdir -p /var/docker/tpl-go-http/log

EXPOSE 7060

WORKDIR /var/docker/tpl-go-http
CMD [ "bin/echo", "-c", "configs/echo.json" ]
