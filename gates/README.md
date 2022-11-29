# Jumpnet - Gates

Gates is the Python library to interact with the jumpnet mux.


## Request Mocks

Mock HTTP requests for the test server
```
PUT http://127.0.0.1:1234/publish
Content-Type: application/json
{
    "Key": "Topic4",
    "Value": "2"
}

GET http://127.0.0.1:1234/history/Topic4
```
