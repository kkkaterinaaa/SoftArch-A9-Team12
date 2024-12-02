import json
import time

STOP_WORDS = {"bird-watching", "ailurophobia", "mango"}

def filter_message(msg):
    return not any(stop_word in msg.lower() for stop_word in STOP_WORDS)

def filter_messages(input_queue, output_queue):
    while True:
        message = input_queue.get()
        if message is None: 
            break

        content = message['content']
        if filter_message(content):
            print(f"[Filter Service] Message passed: {message['content']}")
            output_queue.put(message)
        else:
            print(f"[Filter Service] Message filtered: {message['content']}")
