cargo build --target=x86_64-unknown-linux-musl;
docker-compose up --build\
	--force-recreate\
	--renew-anon-volumes;