FROM postgres:13.3-alpine

ENV POSTGRES_USER=guest
ENV POSTGRES_PASSWORD=guest

COPY init.sql /docker-entrypoint-initdb.d/
