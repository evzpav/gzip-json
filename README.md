# Gzip JSON

## Run
```
    go run main.go
```

## Endpoints:

- Returns data normally as json. Data size via network ~1.2 Mb:  
    http://localhost:8888/normal


- Returns data normally as gzip. Data size via network ~250 Kb:  
    http://localhost:8888/zip