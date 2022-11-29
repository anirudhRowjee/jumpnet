#!/usr/bin/env python3
#
# A mock server implementation for testing the Gates library

from flask import Flask, request, jsonify
import json

data = {
    "test": ["v3", "v2", "v1"],
    "test2": ["v3-2", "v2", "v1"],
    "test4": ["v3-4", "v2", "v1"],
}

app = Flask(__name__)


# TODO Store a message when published (PUT Req on "/")
@app.route("/publish", methods=["PUT"])
def index_get_topic():
    # Parse the JSON Body
    # req_body = request.data
    print("This is a test", str(request.data))
    mydata = json.loads(request.data.decode())
    print(mydata)
    # Update data
    if mydata["Key"] not in data.keys():
        data[mydata["Key"]] = [mydata["Value"],]
    else:
        data[mydata["Key"]] = [mydata["Value"], *data[mydata["Key"]]]

    # Function to say hello
    return "hello, world"

@app.route("/<topic>")
def get_latest_mesg_topic(topic):
    if topic not in data.keys():
        return "no messages"
    else:
        return data[topic][0]

@app.route("/history/<topic>")
def get_all_mesg_topic(topic):
    if topic not in data.keys():
        return "no messages"
    else:
        return data[topic]

# TODO Implement :/topic and :/history/topic methods

app.run(host='0.0.0.0', port=1234)
