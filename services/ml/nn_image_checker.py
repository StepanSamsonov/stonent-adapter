import torch
import nmslib
import numpy as np
import json
import os

import config
import loader

from torch import nn
from PIL import Image
from torchvision import transforms
from scipy.stats import logistic
from io import BytesIO


class NNImageChecker:
    def __init__(self):
        """
        We will use renset50 trained on ImageNet as feature extractor.
        To get features we remove last classification layer of the nn.
        To find nearest features we use nmslib index.
        """
        model = torch.hub.load('pytorch/vision:v0.9.0', 'resnet50', pretrained=True)

        self.feature_extractor = nn.Sequential(*list(model._modules.values())[:-1])

        self.preprocess = transforms.Compose([
                          transforms.Resize(256),
                          transforms.CenterCrop(224),
                          transforms.ToTensor(),
                          transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
                        ])

        if not os.path.exists(config.nn_descriptions_dir):
            os.mkdir(config.nn_descriptions_dir)
        if not os.path.exists(config.nn_index_dir):
            os.mkdir(config.nn_index_dir)
        if not os.path.isfile(f'{config.nn_existed_blocks_dir}{config.nn_existed_blocks_file}'):
            open(NNImageChecker._get_existed_blocks_file(), 'w').close()

    @staticmethod
    def _get_indexes_filename(block_number):
        return f'{config.nn_index_file_prefix}{block_number}{config.nn_index_file_postfix}'

    @staticmethod
    def _get_descriptions_filename(block_number):
        return f'{config.nn_descriptions_file_prefix}{block_number}{config.nn_descriptions_file_postfix}'

    @staticmethod
    def _save_descriptions(descriptions, block_number):
        with open(f'{NNImageChecker._get_descriptions_filename(block_number)}', 'w') as f:
            json.dump(descriptions, f)

    @staticmethod
    def _load_descriptions(block_number):
        with open(f'{NNImageChecker._get_descriptions_filename(block_number)}') as f:
            return json.load(f)

    @staticmethod
    def _get_existed_blocks_file():
        return f'{config.nn_existed_blocks_dir}{config.nn_existed_blocks_file}'

    @staticmethod
    def _get_existed_blocks():
        with open(NNImageChecker._get_existed_blocks_file()) as f:
            return map(lambda x: int(x), f.read().rstrip().split('\n'))

    @staticmethod
    def _add_existed_block(block_number):
        with open(NNImageChecker._get_existed_blocks_file(), 'a') as f:
            f.write(f'{block_number}\n')

    def _get_image_features(self, pil_image):
        """
        :param pil_image: image loaded py PIL library.
        :return: array of features for the image
        """
        with torch.no_grad():
            image = np.array(pil_image)

            if image.ndim == 2:
                image = image[..., None]
                image = np.concatenate([image, image, image], -1)
            if image.shape[-1] == 4:
                image = image[..., :3]

            input_image = self.preprocess(Image.fromarray(image))

            return self.feature_extractor(input_image[None, :])[0].reshape(-1)

    def _get_features(self, pil_image):
        """
        :param pil_image: image loaded py PIL library.
        :return: array of features for the image
        """
        with torch.no_grad():
            image = np.array(pil_image)
            if image.ndim == 2:
                image = image[..., None]
                image = np.concatenate([image, image, image], -1)
            if image.shape[-1] == 4:
                image = image[..., :3]
            input_image = self.preprocess(Image.fromarray(image))
            return self.feature_extractor(input_image[None, :])[0].reshape(-1)

    def _transform_scores(self, scores):
        """
        :param scores: raw cosine distance scores
        :return: scores scaled to [0, 1]

        mean was calculated on classical art dataset
        temp was choose to make sigmiod output close to 0 or 1
        """
        mean = 0.0037
        temp = 10000

        return logistic.cdf((scores - mean) * temp)

    @staticmethod
    def cosine_distance(input1, input2):
        """
        :param input1: first feature vector
        :param input2: second feature vector
        :return: cosine distance between vectors.
        """
        return np.dot(input1, input2.T) / np.sqrt(np.dot(input1, input1.T) * np.dot(input2, input2.T))

    def add_image_to_storage(self, pil_image, description, block_number):
        """
        :param pil_image: image loaded py PIL library.
        :param description: description of the image. Will be returned if image will be chosen as neighbour
        :param block_number: Ethereum block number
        :return: None
        """
        descriptions_filename = self._get_descriptions_filename(block_number)
        descriptions = dict()

        if os.path.isfile(descriptions_filename):
            descriptions = self._load_descriptions(block_number)
        else:
            self._add_existed_block(block_number)

        descriptions_count = len(descriptions)
        features = self._get_features(pil_image)
        index = nmslib.init(method='hnsw', space='cosinesimil')
        indexes_filename = self._get_indexes_filename(block_number)

        if os.path.isfile(indexes_filename):
            index.loadIndex(indexes_filename)

        descriptions[descriptions_count] = description
        index.addDataPoint(data=features, id=descriptions_count)

        index.createIndex({'post': 2})

        index.saveIndex(indexes_filename, True)
        self._save_descriptions(descriptions, block_number)

    def find_most_similar_images(self, pil_image, block_number, num=5):
        """
        :param pil_image:  image loaded py PIL library.
        :param block_number: Ethereum block number
        :param num: number of neighbours to return
        :return: scores, nearest_descriptions.
                 scores - Scores of simularity between pil_image and neighbours
                 nearest_descriptions — descriptions of neighbours
        """
        features = self._get_features(pil_image)
        candidates_descriptions = []

        for current_block_number in self._get_existed_blocks():
            if current_block_number >= block_number:
                continue

            index = nmslib.init(method='hnsw', space='cosinesimil')
            descriptions = self._load_descriptions(current_block_number)

            index.loadIndex(self._get_indexes_filename(current_block_number))

            indexes, scores = index.knnQuery(features, k=num)
            nearest_descriptions = [descriptions[str(i)] for i in indexes]

            if not len(nearest_descriptions):
                continue

            candidates_descriptions.append(nearest_descriptions[0])
            print(f'Nearest description for block {current_block_number}: {nearest_descriptions[0]}')

        index = nmslib.init(method='hnsw', space='cosinesimil')
        descriptions = {}

        for i, description in enumerate(candidates_descriptions):
            descriptions[i] = description
            [address, nft_id] = description.split('-')
            candidate_source, _, candidate_error = loader.get_image_source(address, nft_id)

            if not candidate_source or candidate_error:
                raise candidate_error

            pil_candidate = Image.open(BytesIO(candidate_source))
            candidate_features = self._get_features(pil_candidate)

            index.addDataPoint(data=candidate_features, id=i)

        index.createIndex({'post': 2})

        indexes, scores = index.knnQuery(features, k=num)
        nearest_descriptions = [descriptions[i] for i in indexes]
        scores = self._transform_scores(scores)

        return scores, nearest_descriptions
