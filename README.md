# TestTask-Mobile-API
The REST API for mobile app

<h2>Configuration</h2>

```
sudo docker run --name=CV -e POSTGRES_PASSWORD=11111 -p 5320:5432 -d postgres:16.2
```

<h3>Export env variables</h3>

```
export DB="postgres:11111@localhost:5320/postgres?sslmode=disable" 
export KEY="imagine your own secret key"
export addr_port="your free port"
```

<h2>How to run</h2>

``` make mig-up --> make run ```

