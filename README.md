# Antenna
Antenna is a telegram bot implemented in Golang, broadcasting message to subscribers.

# Work with Go
## Build
It is easy to build
```bash
go build ./cmd/antenna
```

## Run Serve
After building the code, you can start the serve.
```bash
./antenna -c configs/config.yml serve
```
You can customize the config.yml

## Run App
You can run this code to manage your app to create or update or delete
```bash
./antenna -c configs/config.yml app
```

# Work with Docker
## Build Docker Image
It is easy to build docker image by simply run
```bash
docker build -t antenna .
```

## Run Docker Container
After building the docker image, you can run the container by
```bash
docker run -it antenna
```
Configs can be overridden by environment variables. For instance,
```bash
docker run -it --env ANTENNA_TELEGRAM_TOKEN="<telegram_bot_token>" antenna
```

If you want, you can change the port,just like
```bash
--env ANTENNA_HTTP_ADDR=":2333"
```

If you run to mysql or azure, you need add conn config
```bash
--env ANTENNA_STORAGE_OPTIONS='{"conn": "root:root@tcp(localhost)/antrnna?charset=utf8"}'
```

## Admin model
You can open admin model
```bash
--env ANTENNA_ADM_ENABLED=true
```
If app private is true, Add this app need admin user

And add admin usernames
```bash
--env ANTENNA_ADMIN_USERNAMES="user1 user2"
```

## Private model
You can app update to private
```bash
app update -i {AppID} -p
```
Or return to public
```bash
app update -i {AppID} -p=false
```

# Push messages
```
POST http://localhost:8080/antenna/{AppID}
```
Add header authorization
```
Authorization Basic {Base64(AppID:Secret)}
```
