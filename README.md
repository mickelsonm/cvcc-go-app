# ccvc-go-app
Simple WebApp for Demonstrating Kubernetes Deployment

#### Running the application

Start a MySQL Kubernetes Deployment and Service

```
> kubectl create -f mysql.deployment.json
> kubectl create -f mysql.service.json
```

Build the Docker Image

```
> docker build -t cvcc-go:production .
REPOSITORY      TAG                 IMAGE ID            CREATED             SIZE
cvcc-go         production          cb82fc78c390        15 minutes ago      13.9MB
cvcc-go         latest              0fecbc3178e9        17 minutes ago      796MB
```

Start the applications Kubernetes Deployment and Service

```
> kubectl create -f app.deployment.json
> kubectl create -f app.service.json
```

Profit

```
> open http://localhost:8080
```