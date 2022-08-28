# URL Shortner using Redis written in Go

## Running Locally
- Service URL: `http://localhost:5564`
- Redis URL: `http://localhost:6379`

## Endpoints
```
POST http://localhost:5564/v1/surl
{
    "original_url": "https://google.com"
}
```