from PIL import Image
from urllib.parse import urlparse, parse_qs
from postgres import RegisteredImages, RejectedImagesByIPFS, RejectedImagesByNN
import http.server
import socketserver
import json
import loader
import numpy
import cv2
import config
import globals


class RequestHandler(http.server.BaseHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

    def do_GET(self):
        def extract_query(query, key):
            res = None

            if query.get(key) and len(query.get(key)):
                res = query.get(key)[0]

            return res

        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)

        status_code = 404
        response_body = {'error': 'Not found'}

        if parsed_url.path == '/register_new_image':
            contract_address = extract_query(parsed_query, 'contract_address')
            nft_id = extract_query(parsed_query, 'nft_id')

            status_code, response_body = register_new_image(contract_address, nft_id)
        elif parsed_url.path == '/check_registered_image':
            contract_address = extract_query(parsed_query, 'contract_address')
            nft_id = extract_query(parsed_query, 'nft_id')

            status_code, response_body = check_registered_image(contract_address, nft_id)
        elif parsed_url.path == '/registered_images':
            status_code, response_body = get_registered_images()
        elif parsed_url.path == '/rejected_images_by_nn':
            status_code, response_body = get_rejected_images_by_nn()
        elif parsed_url.path == '/rejected_images_by_ipfs':
            status_code, response_body = get_rejected_images_by_ipfs()

        response = json.dumps(response_body)
        response = bytes(response, 'utf8')

        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(response)
        return

    def do_POST(self):
        parsed_url = urlparse(self.path)

        status_code = 404
        response_body = {'error': 'Not found'}

        if parsed_url.path == '/check':
            content_len = int(self.headers.get('Content-Length'))
            body = self.rfile.read(content_len)
            body = json.loads(body)

            job_id = body.get('id')
            data = body.get('data') or dict()
            contract_address = data.get('contract_address')
            nft_id = data.get('nft_id')

            response_body = get_adapter_result(job_id, contract_address, nft_id)
            status_code = 200

        response = json.dumps(response_body)
        response = bytes(response, 'utf8')

        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(response)


def get_adapter_result(job_id, contract_address, nft_id):
    def get_result(code, data=None, error=None):
        return {
            'jobRunID': job_id,
            'statusCode': code,
            'data': data,
            'error': str(error),
        }

    with globals.mutex:
        try:
            if not contract_address:
                return get_result(400, None, 'Invalid contract address')
            if not nft_id or not nft_id.isdigit():
                return get_result(400, None, 'Invalid nft id')

            image_source, block_number, image_source_error = loader.get_image_source(contract_address, nft_id)

            if not image_source or not block_number or image_source_error:
                return get_result(500, None, image_source_error)

            np_image_source = numpy.frombuffer(image_source, numpy.uint8)
            image = cv2.imdecode(np_image_source, cv2.IMREAD_COLOR)
            image = Image.fromarray(image)

            scores, descriptions = globals.image_checker.find_most_similar_images(image, block_number)

            if not len(scores) or not len(descriptions):
                scores = [1]
                descriptions = ['']

            if descriptions[0] == f'{contract_address}-{nft_id}':
                score = scores[1]
            else:
                score = scores[0]

            result = {
                'score': int(score * 100),
                'detailed_information': {
                    'scores': [*scores],
                    'descriptions': [*descriptions]
                }
            }

            return get_result(200, result)
        except Exception as e:
            return get_result(200, None, e)


def register_new_image(contract_address, nft_id):
    with globals.mutex:
        try:
            if not contract_address:
                return 400, {'error': 'Invalid contract address'}
            if not nft_id or not nft_id.isdigit():
                return 400, {'error': 'Invalid nft id'}

            already_registered_filter = RegisteredImages.filter(
                filter={
                    'contract_address': contract_address,
                    'nft_id': nft_id,
                }
            )
            already_registered = already_registered_filter.total

            if already_registered:
                return 400, {'error': 'Already registered'}

            image_source, block_number, image_source_error = loader.get_image_source(contract_address, nft_id)

            if not image_source or image_source_error:
                return 500, {'error': image_source_error or 'Image source is empty'}

            globals.image_manager.register_new_image(contract_address, nft_id, block_number, image_source)

            return 200, {'error': None}
        except Exception as e:
            return 500, {'error': e}


def check_registered_image(contract_address, nft_id):
    with globals.mutex:
        try:
            if not contract_address:
                return 400, {'error': 'Invalid contract address'}
            if not nft_id or not nft_id.isdigit():
                return 400, {'error': 'Invalid nft id'}

            already_registered_filter = RegisteredImages.filter(
                filter={
                    'contract_address': contract_address,
                    'nft_id': nft_id,
                }
            )
            already_registered = already_registered_filter.total

            return 200, {'is_registered': not not already_registered}
        except Exception as e:
            return 500, {'error': e}


def get_registered_images():
    with globals.mutex:
        try:
            registered_images = RegisteredImages.all()
            registered_images = map(lambda x: x.to_dict(), registered_images)
            registered_images = list(registered_images)

            return 200, {'registered_images': registered_images}
        except Exception as e:
            return 500, {'error': e}


def get_rejected_images_by_ipfs():
    with globals.mutex:
        try:
            rejected_images_by_ipfs = RejectedImagesByIPFS.all()
            rejected_images_by_ipfs = map(lambda x: x.to_dict(), rejected_images_by_ipfs)
            rejected_images_by_ipfs = list(rejected_images_by_ipfs)

            return 200, {'rejected_images_by_ipfs': rejected_images_by_ipfs}
        except Exception as e:
            return 500, {'error': e}


def get_rejected_images_by_nn():
    with globals.mutex:
        try:
            rejected_images_by_nn = RejectedImagesByNN.all()
            rejected_images_by_nn = map(lambda x: x.to_dict(), rejected_images_by_nn)
            rejected_images_by_nn = list(rejected_images_by_nn)

            return 200, {'rejected_images_by_nn': rejected_images_by_nn}
        except Exception as e:
            return 500, {'error': e}


def run_server():
    with socketserver.TCPServer(('', config.server_port), RequestHandler) as httpd:
        print(f'Server started at port: {str(config.server_port)}', flush=True)
        httpd.serve_forever()
