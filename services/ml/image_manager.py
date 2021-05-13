from PIL import Image
import io
import rabbitmqapi


class ImageManager:
    def __init__(self, image_cheker):
        self.image_cheker = image_cheker

    def _register_new_image(self, contract_address, nft_id, source):
        try:
            pil_image = Image.open(io.BytesIO(source))
            description = f'{str(contract_address)}-{str(nft_id)}'

            self.image_cheker.add_image_to_storage(pil_image, description)
        except Exception as e:
            print("error in registring new image", e)

    def register_new_images(self):
        for contract_address, nft_id, source in rabbitmqapi.consume_events():
            self._register_new_image(contract_address, nft_id, source)
