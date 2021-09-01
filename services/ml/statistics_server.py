from urllib.parse import urlparse, parse_qs
from postgres import RegisteredImages, RejectedImagesByIPFS, RejectedImagesByNN
import http.server
import socketserver
import json
import config
import globals
import loader


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

        if parsed_url.path == '/check_registered_image':
            contract_address = extract_query(parsed_query, 'contract_address')
            nft_id = extract_query(parsed_query, 'nft_id')

            status_code, response_body = check_registered_image(contract_address, nft_id)
        elif parsed_url.path == '/registered_images':
            status_code, response_body = get_registered_images()
        elif parsed_url.path == '/rejected_images_by_nn':
            status_code, response_body = get_rejected_images_by_nn()
        elif parsed_url.path == '/rejected_images_by_ipfs':
            status_code, response_body = get_rejected_images_by_ipfs()
        elif parsed_url.path == '/statistics':
            status_code, response_body = get_statistics()

        response = json.dumps(response_body)
        response = bytes(response, 'utf8')

        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(response)
        return


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


def get_statistics():
    with globals.mutex:
        try:
            registered_images = RegisteredImages.filter()
            registered_images_count = registered_images.total

            rejected_by_ipfs_images = RejectedImagesByIPFS.filter()
            rejected_by_ipfs_images_count = rejected_by_ipfs_images.total

            rejected_by_nn_images = RejectedImagesByNN.filter()
            rejected_by_nn_images_count = rejected_by_nn_images.total

            loader_response, loader_error = loader.get_statistics()

            if not loader_response or loader_error:
                raise loader_error
            else:
                found_images_count = loader_response.get('CountOfFound')
                processed_images_count = loader_response.get('CountOfDownloaded')

            return 200, {
                'statistics': {
                    'found_images_count': found_images_count,
                    'processed_images_count': processed_images_count,
                    'registered_images_count': registered_images_count,
                    'rejected_by_ipfs_images_count': rejected_by_ipfs_images_count,
                    'rejected_by_nn_images_count': rejected_by_nn_images_count,
                    'is_completed': globals.all_images_has_been_downloaded,
                }
            }
        except Exception as e:
            return 500, {'error': e}


def run_statistics_server():
    with socketserver.TCPServer(('', config.statistics_server_port), RequestHandler) as httpd:
        print(f'Statistics server started at port: {str(config.statistics_server_port)}', flush=True)
        httpd.serve_forever()
