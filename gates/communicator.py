import requests
import logging 
import json

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
        req = requests.put(self.server_URL + "/publish", json= {
            "Key": self.topic_name,
            "Value": str(message)
        })
        return req.status_code
