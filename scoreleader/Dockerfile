FROM registry01.idc-sginfra.net/stove-platform/stove-tomcat:2.0.3

MAINTAINER Cloudtechteam <seungjhan@smilegate.com>

USER root

RUN sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-Linux-*
RUN sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-Linux-*

RUN yum -y install openldap-clients && yum clean all  && rm -rf /var/cache

ADD ./bin /stove/apps/template/

RUN chown -R stove:stove /stove && \
    chmod -R 700 /stove

CMD ["/stove/apps/template/template", "-c=/stove/apps/template/config.yml"]

