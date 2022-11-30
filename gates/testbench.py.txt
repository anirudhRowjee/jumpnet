# Testbench file for Producer and Consumer libraries
from communicator import JumpnetProducer, JumpnetConsumer
# NOTE assumes you have the mock server running at localhost:1234

producer = JumpnetProducer("test", "http://localhost:1234")
print(producer.send_message("lol this is a test"))


consumer = JumpnetConsumer("test", "http://localhost:1234")
print(consumer.receiver_message_all())

