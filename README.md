# My AI sandbox

```
make up
make get-models
make go-build
make chat [--config filename.yaml] [--chat chatname]
```
Browser: http://localhost:xxxx (depends on chat name in config file)


# Fill database
```
# all above inside docker container (make shell)
# create Cocktail class id DB
./bin/go-client DB init

# load CSV data to DB
./bin/go-client db learn
```


# Chat procedure
```
# all above inside docker container (make shell)

# 1th terminal
./bin/go-client http --chat waiter

# 2nd terminal
./bin/go-client http --chat bartender

# 3th terminal
./bin/go-client chat

```
browser: http://localhost:3000


# VSC + docker
- Install (Ctrl + Shift + X): Dev - Containers (Microsoft)
- On left bottom corner click >< icon and select Attach to running container... and select container $(APP_NAME)
- Install (Ctrl + Shift + X): Go (Go Team at Google)
- Run command (Ctrl + Shift + P) Go: Install/Update tools, select all and click OK

# VSC + Weaviate 
- Instal Weaviate Studio