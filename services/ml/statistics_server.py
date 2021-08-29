from urllib.parse import urlparse
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
        parsed_url = urlparse(self.path)

        status_code = 404
        response_body = {'error': 'Not found'}

        if parsed_url.path == '/statistics':
            status_code, response_body = get_statistics()

        response = json.dumps(response_body)
        response = bytes(response, 'utf8')

        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(response)
        return


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
                precessed_images_count = loader_response.get('CountOfDownloaded')

            return 200, {
                'statistics': {
                    'found_images_count': found_images_count,
                    'precessed_images_count': precessed_images_count,
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
