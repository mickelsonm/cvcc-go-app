# ccvc-go-app
Simple WebApp for Demonstrating Kubernetes Deployment

#### Running the application

Export environment variables

```
> export MYSQL_USER=cvcc-user \
    MYSQL_PASS=cvcc-pass \
    MYSQL_DB=sample \
    MYSQL_HOST=localhost \
    MYSQL_PORT=3307
```

Start a MySQL Docker container

```
> mkdir -p ~/tmp/mysql
> docker \
  run \
  -d \
  -e MYSQL_ROOT_PASSWORD=${MYSQL_PASS} \
  -e MYSQL_USER=${MYSQL_USER} \
  -e MYSQL_PASSWORD=${MYSQL_PASS} \
  -e MYSQL_DATABASE=${MYSQL_DB} \
  -v ~/tmp/mysql:/var/lib/mysql \
  -p 3307:3306 \
  --name cvcc-mysql \
  mysql:5.7
```

Start the Application

```
> go run main.go
```

Build the Docker Image

```
> docker build -t cvcc-go:production .
REPOSITORY      TAG                 IMAGE ID            CREATED             SIZE
cvcc-go         production          cb82fc78c390        15 minutes ago      13.9MB
cvcc-go         latest              0fecbc3178e9        17 minutes ago      796MB
```

Notice the difference in image size

```
> docker images | grep cvcc
```

Run the Docker Image

```
> docker \
    run \
    -p 8080:8080 \
    -e MYSQL_USER=${MYSQL_USER} \
    -e MYSQL_PASS=${MYSQL_PASS} \
    -e MYSQL_HOST=mysql \
    -e MYSQL_PORT=3306 \
    -e MYSQL_DB=${MYSQL_DB} \
    --name cvcc-go \
    --link cvcc-mysql:mysql \
    cvcc-go:production
```

Profit

```
> open http://localhost:8080
```