git pull

docker build -t proxy-broadcast .

docker run -p 8081:8081 proxy-broadcast
