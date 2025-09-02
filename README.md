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


# TODO
Prompt structure example:
```
[ROLE]
You are an AI SQL assistant. Your job is to translate user requests into SQL queries.

[CONTEXT]
Available tables:
- users(id, name, email)
- orders(id, user_id, total, date)
- products(id, name, stock, price)

[EXAMPLES]
Example 1:
[USER]: "Find all users with no orders."
[ASSISTANT]: "SELECT * FROM users WHERE id NOT IN (SELECT user_id FROM orders);"

Example 2:
[USER]: "List products with stock below 10."
[ASSISTANT]: "SELECT * FROM products WHERE stock < 10;"

[TASK]
User request: "Show me the top 5 customers by spending."

[INSTRUCTIONS]
- Return only valid SQL query.
- Do not explain your reasoning.
- Use only tables provided in context.
```