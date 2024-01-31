
#
# docker build --no-cache --progress=plain -t blockchain .
# docker run --rm -it --privileged -p 7011:7011 -p 8801:8801  blockchain
#

FROM docker:24-dind

WORKDIR /blockchain
COPY . /blockchain

COPY <<-"EOT" /etc/supervisor/conf.d/supervisord.conf
[supervisord]
logfile=/dev/stdout 
logfile_maxbytes=0  
loglevel=info
pidfile=/tmp/supervisord.pid
nodaemon=true
user=root

[unix_http_server]
file=/tmp/supervisor.sock

[program:dockerd]
command=/usr/local/bin/dockerd-entrypoint.sh
autorestart=false
startretries=0

[program:fabric]
command=/blockchain/setup-fablo.sh entrypoint
autorestart=false
startretries=0
EOT

RUN \
  apk add --no-cache bash curl jq supervisor git go && \
  rm -rf /var/cache/apk/*

# Rest
EXPOSE 7011
# Explorer
EXPOSE 8801
# Peer 0
EXPOSE 8041
EXPOSE 7041
# Peer 2
EXPOSE 8042
EXPOSE 7042

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
