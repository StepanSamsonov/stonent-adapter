from PIL import Image
from io import BytesIO
from postgres import RegisteredImages, RejectedImagesByNN
import rabbitmq


class ImageManager:
    def __init__(self, image_checker):
        self.image_checker = image_checker

    def register_new_image(self, contract_address, nft_id, block_number, bytes_source):
        try:
            if not bytes_source:
                return

            duplicated_images = RegisteredImages.filter(
                filter={
                    'contract_address': contract_address,
                    'nft_id': nft_id,
                }
            )

            if duplicated_images.total:
                return

            pil_image = Image.open(BytesIO(bytes_source))
            description = f'{str(contract_address)}-{str(nft_id)}'

            self.image_checker.add_image_to_storage(pil_image, description, block_number)

            RegisteredImages.add({
                'contract_address': contract_address,
                'nft_id': nft_id,
                'format': pil_image.format,
            })
            print(f'Consumed by NN: {contract_address} {nft_id}', flush=True)
        except Exception as e:
            RejectedImagesByNN.add({
                'contract_address': contract_address,
                'nft_id': nft_id,
                'description': e,
            })
            print("Error in registering new image", e, flush=True)

    def register_new_images(self, mutex):
        for contract_address, nft_id, block_number, bytes_source in rabbitmq.consume_events():
            with mutex:
                self.register_new_image(contract_address, nft_id, block_number, bytes_source)
