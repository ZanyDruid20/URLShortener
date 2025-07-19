# URLShortener

A simple, scalable URL shortener service built with Go, Gin, Redis, Docker, and deployed on AWS EC2.

## Features

- Shorten long URLs and generate unique short links
- Redirect short URLs to original URLs
- Fast key-value storage with Redis
- RESTful API endpoints
- Containerized with Docker and orchestrated with Docker Compose
- Deployed on AWS EC2 (Amazon Linux 2023)

## Tech Stack

- Go (Golang)
- Gin Web Framework
- Redis
- Docker & Docker Compose
- AWS EC2

## API Endpoints

- `GET /` — Health check
- `POST /create-short-url` — Create a short URL  
  **Body:** `{ "url": "https://example.com" }`
- `GET /:shortUrl` — Redirect to original URL

## License

This project is licensed under the MIT License.

---
## Resources:
https://dev.to/aws-builders/how-to-deploy-a-multi-container-docker-compose-application-on-amazon-ec2-59n2
https://www.eddywm.com/lets-build-a-url-shortener-in-go/

**Author:**
[Furnom Dam](https://github.com/ZanyDruid20)