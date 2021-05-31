from threading import Thread
from globals import image_manager, mutex
import server


if __name__ == '__main__':
    print('START ML', flush=True)
    registerer_thread = Thread(target=image_manager.register_new_images, args=[mutex])
    server_thread = Thread(target=server.run_server)

    registerer_thread.start()
    server_thread.start()
