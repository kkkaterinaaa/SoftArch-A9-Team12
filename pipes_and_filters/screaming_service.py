import json
import time

def scream_message(input_queue, output_queue):
    while True:
        message = input_queue.get()
        if message is None:
            break
        
        print(f"[Screaming Service] Before screaming: {message['content']}")
        message['content'] = message['content'].upper()
        print(f"[Screaming Service] After screaming: {message['content']}")
        output_queue.put(message)
