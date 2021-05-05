import cv2
import numpy as np
import json
from adapter import Adapter
from flask import Flask, request, jsonify
from nn_image_checker import NNModelChecker
from PIL import Image
from image_manager import ImageManager


app = Flask(__name__)
image_checker = NNModelChecker()
image_manager = ImageManager(image_checker)


@app.before_request
def log_request_info():
    app.logger.debug('Headers: %s', request.headers)
    app.logger.debug('Body: %s', request.get_data())


def load_image(data):
    nparr = np.fromstring(data, np.uint8)
    # decode image
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    img = Image.fromarray(img)
    return img


@app.route('/register_image', methods=['POST'])
def register_image():
    image_manager.register_new_images()

    img = load_image(request.data)
    image_checker.add_image_to_storage(img, 'None')

    # build a response dict to send back to client
    response = {'message': 'image received.'}
    # encode response using jsonpickle
    response = jsonify(response)
    return response


@app.route('/image_score', methods=['POST'])
def image_score():
    image_manager.register_new_images()

    img = load_image(request.data)
    scores, descriptions = image_checker.find_most_simular_images(img)

    # build a response dict to send back to client
    response = {'scores': str(scores), 'descriptions': str(descriptions)}
    # encode response using jsonpickle
    response = jsonify(response)
    return response


@app.route('/check', methods=['POST'])
def call_adapter():
    body = json.loads(request.data)
    adapter = Adapter(body, image_checker, image_manager)
    return jsonify(adapter.result)


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port='9090', threaded=True)
