# API
API at golang with gin framework

On this project I am realizing API for pseudo-store of plants.
Now you can test these HTTP methods: GET, PUT, POST and DELETE,
of course you also can use method GET like get item by ID.

The data is stored at local postgres in docker, the credentials is stored in code
for now.


**Launch project**

First, you need to download *the docker image* from the docker hub and to run it;
to do this, go the command in the terminal:

`make docker`

Second, you can set your own variables on the config(host, port), for this change them in:
    
`/API/config/local.yaml`

After, you can start the project with:

`make`

P.S. for now commands for the start works only for MacOs and linux systems


**API request for test**
(Examples shown on localhost and port 8080)

*GET*
1) Get all plants from table:

`curl http://localhost:8080/plants`

2) Get plant by ID:

`curl http://localhost:8080/plants/<ID>` 

*POST*

* here we add plant with ID = 6, name = Mint, Amount = 10 and Price = 9.99

`curl http://localhost:8080/plants \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"id": "6","Product": "Mint","Amount": "10","price": 9.99}'`

*PUT*

* here we change plant by ID = 5 to the name = Avocado

`curl http://localhost:8080/plants/5 \
--include \
--header "Content-Type: application/json" \
--request "PUT" \
--data '{"id": "5","Product": "Avocado"}'`

*DELETE*

* here we delete plant by ID = 6

`curl http://localhost:8080/plants/6 \
  --header "Content-Type: application/json" \
  --request "DELETE"
  `