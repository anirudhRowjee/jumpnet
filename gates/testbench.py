# Testbench file for Producer and Consumer libraries
from communicator import JumpnetProducer, JumpnetConsumer
# NOTE assumes you have the mock server running at localhost:1234

producer = JumpnetProducer("test", "http://127.0.0.1:49249")
producer2 = JumpnetProducer("test", "http://127.0.0.1:49243")
producer3 = JumpnetProducer("test", "http://127.0.0.1:49245")

print(producer.send_message("lol this is a test"))
print(producer2.send_message("lol this is a test - 2"))
print(producer3.send_message("lol this is a test - 3"))


consumer = JumpnetConsumer("test", "http://localhost:49247")
print(consumer.receiver_message_all())
