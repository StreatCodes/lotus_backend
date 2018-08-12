FROM alpine
WORKDIR /root
COPY target/x86_64-unknown-linux-musl/debug/lotus_backend .
COPY tables.sql .
CMD [ "./lotus_backend" ]