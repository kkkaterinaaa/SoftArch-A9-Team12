from flask import Flask, request, jsonify
from multiprocessing import Queue
import json
import os


app = Flask(__name__)

input_queue = Queue()
output_queue = Queue()

@app.route('/send', methods=['POST'])
def handle_post_message():
    data = request.get_json()
    message = {
        'alias': data['alias'],
        'content': data['content']
    }

    input_queue.put(message)
    return jsonify({"message": "Message sent successfully"}), 200


if __name__ == '__main__':
    from multiprocessing import Process
    from filter_service import filter_messages
    from screaming_service import scream_message
    from publish_service import publish_message
    
    processes = [
        Process(target=filter_messages, args=(input_queue, output_queue)),
        Process(target=scream_message, args=(output_queue, output_queue)),
        Process(target=publish_message, args=(output_queue,))
    ]
    
    for p in processes:
        p.start()

    app.run(debug=True, port=8080)
