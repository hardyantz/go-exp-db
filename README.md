

Simple Application / Experiment golang with rethinkDB & trailDB .

# How to 
1. brew install rethinkdb (https://rethinkdb.com/)
2. brew install traildb (http://traildb.io/)
3. go run traildb.go
4. go run rethinkdb.go

# Endpoint traildb
1. `POST /get` Create new trail
2. `GET /create` Get all trail
3. `GET /get-wiki` Get all data from sample data (500.000 record) 

# Endpoint rethinkdb
1. `POST /create` Create new data
2. `GET /get` Get all data
3. `GET /get/:id` Get data by ID
