# Antenna
Antenna is a telegram bot implemented in Golang, broadcasting message to subscribers. 

# Work with Docker
## Build Docker Image
It is easy to build docker image by simply run
```
docker build -t antenna .
```

## Run Docker Container
After building the docker image, you can run the container by
```
docker run -it antenna
```
Configs can be overridden by environment variables. For instance,
```
docker run -it --env ANTENNA_TELEGRAM_TOKEN="<telegram_bot_token>" antenna
```
