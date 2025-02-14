# Distributed URL Shortener with Rate Limiting & Analytics

This project implements a distributed URL shortener service built in GoLang with the following features:

- **URL Shortening API:**
  - `POST /shorten` – Accepts a long URL and returns a shortened URL.
  - `GET /{shortID}` – Redirects to the original URL and logs analytics.
- **Storage & Expiry:**
  - Uses Redis to store shortID-to-URL mappings with a 30-day expiration.
- **Analytics Logging:**
  - Logs every redirect (shortID, timestamp, and user IP) asynchronously into MongoDB.
- **Rate Limiting:**
  - Limits clients to 10 requests per minute (implemented via Redis).
- **Worker Goroutines:**
  - Offloads analytics logging from API requests to maintain responsiveness.
- **Deployment:**
  - Containerized with Docker.
  - Kubernetes manifests provided for deploying:
    - 2 replicas for the API,
    - 1 replica for Redis,
    - 1 replica for MongoDB,
    - and Nginx as a reverse proxy.
- **Testing:**
  - Unit tests are written using [testify](https://github.com/stretchr/testify).

## File structure

- main.go – The entry point of the application.
- analytics.go – Contains the analytics logging logic.
- rate-limit.go – Contains the rate limiting logic.
- handlers.go – Contains the API handlers.
- main_test.go – Contains the unit tests.
- Dockerfile – Builds the Go binary.
- k8s folder – Contains the Kubernetes manifests.
- README.md – Contains the project description.

## Usage

1. clone the repository
2. if you will run the main_test.go, you will need to run local docker redis and mongodb container through these commands :  
    - docker run -d -p 6379:6379 --name redis redis
    - docker run -d -p 27017:27017 --name mongo mongo
    - then you can run the tests through the following command : go test -v
3. to run via k8s, you can use the k8s folder and run the following commands:
    - kubectl apply -f redis-deployment.yaml
    - kubectl apply -f mongo-deployment.yaml
    - kubectl apply -f api-deployment.yaml
    - kubectl apply -f nginx-deployment.yaml
4. then expose the nginx service to be able to access the api through this command : kubectl port-forward deployment/nginx 8080:80
5. then you can access the api through the following: 
    - POST http://localhost:8080/shorten with a body for example like this : {"url":"https://www.google.com"}
    - GET http://localhost:8080/{shortID} where shortID is the short id that you got from the previous step
    - you can also try to access the analytics through the following : GET http://localhost:8080/analytics

## technologies used
- GoLang
- Redis
- MongoDB
- Docker
- Kubernetes
- Nginx
- Testify


