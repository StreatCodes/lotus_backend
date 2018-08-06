FROM alpine
WORKDIR /root
COPY target/debug/lotus_backend .
COPY tables.sql .
CMD [ "./lotus_backend" ]