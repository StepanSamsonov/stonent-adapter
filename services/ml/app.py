from threading import Thread
from globals import image_manager, mutex
from PIL import Image
import time
import server
import postgres


if __name__ == '__main__':
    print('Start ML', flush=True)

    # Disable image size limitation
    Image.MAX_IMAGE_PIXELS = None

    postgres.connect()

    registerer_thread = Thread(target=image_manager.register_new_images, args=[mutex])
    registerer_thread.daemon = True
    server_thread = Thread(target=server.run_server)
    server_thread.daemon = True

    registerer_thread.start()
    server_thread.start()

    while True:
        time.sleep(1)
