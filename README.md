# ccvc-go-app
Simple WebApp for Demonstrating Kubernetes Deployment

#### Running the application

Export environment variables

```
> export MYSQL_USER=cvcc-user \
    MYSQL_PASS=cvcc-pass \
    MYSQL_HOST=127.0.0.1 \
    MYSQL_PORT=3307 \
    MYSQL_DB=sample
```

Start a MySQL Docker container

```
> docker \
  run \
  --detach \
  --env MYSQL_ROOT_PASSWORD=${MYSQL_PASS} \
  --env MYSQL_USER=${MYSQL_USER} \
  --env MYSQL_PASSWORD=${MYSQL_PASS} \
  --env MYSQL_DATABASE=${MYSQL_DB} \
  --name cvcc-mysql \
  --publish 3307:3306 \
  mysql:5.7
```

Start the Application

```
> go run main.go
```

Profit

```
> open http://localhost:8080
```