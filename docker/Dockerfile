#***********************************************************************
# Stage 1 - Build the application in a Linux container. Cross compiling
# on a Mac is a pain since sqlite3 is a CGO enabled library. It
# requires that we have a C compiler set up for cross compilation.
#***********************************************************************
FROM golang as builder

WORKDIR /go/src

# Copy over the project files. This is better than cloning the project, because
# it allows us to break the docker cache when developing.
COPY *.go cmd github.com/ameske/nfl-pickem/
COPY cmd github.com/ameske/nfl-pickem/cmd/
COPY http github.com/ameske/nfl-pickem/http/
COPY parser github.com/ameske/nfl-pickem/parser/
COPY sqlite3 github.com/ameske/nfl-pickem/sqlite3
COPY vendor github.com/ameske/nfl-pickem/vendor/

RUN go build github.com/ameske/nfl-pickem/cmd/nfld 


#***********************************************************************
# Stage 2 - Construct the image with the artifacts from stage 1
#***********************************************************************
FROM debian

LABEL maintainer="Kyle Ames"
LABEL version="1.0"
LABEL description="nfl-pickem daemon served over caddy"

RUN apt-get update && apt-get install -y curl libcap2-bin supervisor && rm -rf /var/lib/apt/lists/*

# Install the caddy webserver & setup the web asset directory
RUN mkdir /caddy_install && \
    curl https://caddyserver.com/download/linux/amd64 > /caddy_install/caddy_server.tar.gz && \
    tar xzf /caddy_install/caddy_server.tar.gz -C /caddy_install && \
    cp /caddy_install/caddy /usr/local/bin && \
    rm -r /caddy_install && \
    chown root:root /usr/local/bin/caddy && \
    chmod 755 /usr/local/bin/caddy && \ 
    setcap 'cap_net_bind_service=+ep' /usr/local/bin/caddy && \
    mkdir /etc/caddy && \
    chown -R root:www-data /etc/caddy && \
    mkdir /etc/ssl/caddy && \
    chown -R www-data:root /etc/ssl/caddy && \
    chmod 0770 /etc/ssl/caddy && \
    mkdir /var/www && \
    chown www-data:www-data /var/www && \
    chmod 555 /var/www && \
    mkdir /var/www/nfl

# Setup the Caddyfile
COPY docker/Caddyfile /etc/caddy/
RUN chown www-data:www-data /etc/caddy/Caddyfile && \
    chmod 444 /etc/caddy/Caddyfile

# Copy the supervisord artifacts
RUN mkdir -p /var/log/supervisor
COPY docker/supervisord.conf /etc/supervisor/supervisord.conf
COPY docker/nfld.conf /etc/supervisor/conf.d/nfld.conf

# Install the nfl-pickem daemon
COPY --from=builder /go/src/nfld /usr/local/bin
COPY debug.config /nfl-pickem/nfld.conf
COPY faketest.db /nfl-pickem/nfl.db

# Set-up the website assets
COPY www /var/www/nfl/
COPY logos /var/www/nfl/
RUN chown -R www-data:www-data /var/www/nfl && \
    chmod -R 555 /var/www/nfl

EXPOSE 80

CMD ["/usr/bin/supervisord"]
