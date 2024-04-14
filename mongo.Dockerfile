FROM alpine:3.19
ENV MONGOMS_SYSTEM_BINARY=/usr/bin/mongod

# Have to revert to Alpine 3.9 repo caches to get MongoDB
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.9/main" >> /etc/apk/repositories
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.9/community" >> /etc/apk/repositories
RUN cat /etc/apk/repositories

RUN apk add --update mongodb yaml-cpp=0.6.2-r2
RUN apk add bash netcat-openbsd
COPY ./scripts/mongodb_entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
