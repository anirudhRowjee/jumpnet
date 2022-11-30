import requests
import logging
import json
import time

# Base Class for the producer


class JumpnetProducer:
    # [X] connection to the broker
    # [X] Accepting the message in the same manner
    # [] something you send (?)

    def __init__(self, topic_name: str, server_URL: str):
        # Initialize the producer
        # NOTE The server URL will consist of a port as well as an IP Adddress.
        # WARNING there will
        #   Example: http://localhost:9091
        logging.debug("Starting jumpnet producer...")
        self.topic_name = topic_name
        self.server_URL = server_URL

    def send_message(self, message):
        # Make an HTTP PUT request to the Broker
        # Check the response status, and reply appropriately
        req = requests.put(self.server_URL + "/publish", json={
            "Key": self.topic_name,
            "Value": str(message)
        })
        return req.status_code


class JumpnetConsumer:
    def __init__(self, topic_name: str, server_URL: str):
        # Initialize the producer
        logging.debug("Starting jumpnet consumer...")
        self.topic_name = topic_name
        self.server_URL = server_URL

    def receiver_message_latest(self):
        # Make an HTTP PUT request to the Broker
        # Check the response status, and reply appropriately
        req = requests.get(self.server_URL + "/" + self.topic_name)
        return req.content

    def receiver_message_all(self):
        # parse the json ie get the array to normal text that can be used in python
        # Make an HTTP PUT request to the Broker
        # Check the response status, and reply appropriately
        req = requests.get(self.server_URL + "/history/" + self.topic_name)
        # return req.content
        all_msgs = json.loads(req.content)
        return all_msgs

    def poll(self):
        count = 60  # 60 for 60 seconds
        old_value = None
        current_value = None
        while(count != 0):
            req = requests.get(self.server_URL + "/" + self.topic_name)
            time.sleep(1)
            if(req == old_value):
                pass
            else:
                old_value = current_value
                current_value = req
                print("new value detected")
            count -= 1
