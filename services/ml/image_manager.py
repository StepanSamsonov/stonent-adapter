from PIL import Image
from io import BytesIO
import rabbitmqapi


class ImageManager:
    def __init__(self, image_checker):
        self.image_checker = image_checker

    def register_new_image(self, contract_address, nft_id, bytes_source):
        try:
            pil_image = Image.open(BytesIO(bytes_source))
            description = f'{str(contract_address)}-{str(nft_id)}'

            self.image_checker.add_image_to_storage(pil_image, description)
            print(f'Consumed by NN: {contract_address} {nft_id}')
        except Exception as e:
            print("error in registering new image", e)

    def register_new_images(self, mutex):
        for contract_address, nft_id, bytes_source in rabbitmqapi.consume_events():
            with mutex:
                self.register_new_image(contract_address, nft_id, bytes_source)
