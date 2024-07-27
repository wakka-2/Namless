## Data service

This µService offers a (key, value) store. Exposes a REST API with CRUD operations over a DB.

## Functionality
- a user can issue REST calls to create/updated/retrieve/dele a (key, value) entry
- data is stored locally, in a postgres DB

## How to run locally
- produce the appropriate configs, similar to /configs/data.json (i.e.: /etc/data/data.json)
- type _go run cmd/main/main.go -config=/etc/data/data.json_
