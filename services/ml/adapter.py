from bridge import Bridge
from PIL import Image
import cv2
import numpy as np

import loader


class Adapter:
    id_params = ['id', 'nft', '_id']

    def __init__(self, input, image_checker, image_manager):
        self.image_cheker = image_checker
        self.image_manager = image_manager
        self.id = input.get('id', '1')
        self.result = ""   #kostil
        self.request_data = input.get('data')
        if self.validate_request_data():
            self.bridge = Bridge()
            self.set_params()
            self.create_request()
        else:
            self.result_error('No data provided')

    def validate_request_data(self):
        if self.request_data is None:
            return False
        if self.request_data == {}:
            return False
        return True

    def set_params(self):
        for param in self.id_params:
            self.id_params = self.request_data.get(param)
            if self.id_params is not None:
                break

    def create_request(self):
        try:
            #params = {
            #    'fsym': self.from_param,
            #    'tsyms': self.to_param,
            #}
            #response = self.bridge.request(self.base_url, params)

            # TODO: get contract address

            nft_id = int(self.id_params)
            contract_address = '0xd07dc4262bcdbf85190c01c996b4c06a461d2430'

            image_source, error = loader.get_image_source(contract_address, nft_id)

            if not image_source or error:
                raise error

            nparr = np.fromstring(image_source, np.uint8)
            img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
            img = Image.fromarray(img)

            self.image_manager.register_new_images()
            scores, descriptions = self.image_cheker.find_most_simular_images(img)

            if int(descriptions[0]) == nft_id:
                score = scores[1]
            else:
                score = scores[0]

            score = int(score * 100)


            #data = response.json()
            #self.result = data[self.to_param]
            data = {'score': score, 'detailed information': {
                'scores': [*scores],
                'descriptions': [*descriptions]
            }}
            self.result_success(data)
        except Exception as e:
            self.result_error(e)
        finally:
            self.bridge.close()

    def result_success(self, data):
        self.result = {
            'jobRunID': self.id,
            'data': data,
            'result': self.result,
            'statusCode': 200,
        }

    def result_error(self, error):
        self.result = {
            'jobRunID': self.id,
            'status': 'errored',
            'error': f'There was an error: {error}',
            'statusCode': 500,
        }
