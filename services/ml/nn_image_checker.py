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

        self._index = nmslib.init(method='hnsw', space='cosinesimil')
        self._feature_dict = {}
        self._pages_counter = 1
        self._page_size_counter = 0
        self._page_saving_file = self._get_saving_indexes_filename(self._pages_counter)

        if not os.path.exists(config.nn_features_dict_dir):
            os.mkdir(config.nn_features_dict_dir)
        if not os.path.exists(config.nn_index_dir):
            os.mkdir(config.nn_index_dir)

    @staticmethod
    def _get_saving_indexes_filename(page_number):
        return f'{config.nn_index_file_prefix}{page_number}{config.nn_index_file_postfix}'

    @staticmethod
    def _get_saving_features_filename(page_number):
        return f'{config.nn_features_dict_file_prefix}{page_number}{config.nn_features_dict_file_postfix}'

    def _save_features_dict(self):
        with open(f'{self._get_saving_features_filename(self._pages_counter)}', 'w') as f:
            json.dump(self._feature_dict, f)

    def _load_features_dict(self, page_number):
        with open(f'{self._get_saving_features_filename(page_number)}') as f:
            return json.load(f)

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

    def add_image_to_storage(self, pil_image, description):
        """
        :param pil_image: image loaded py PIL library.
        :param description: description of the image. Will be returned if image will be chosen as neighbour
        :return: None
        """
        features = self._get_features(pil_image)
        index = len(self._feature_dict)

        self._feature_dict[index] = description
        self._index.addDataPoint(data=features, id=index)

        self._page_size_counter += 1
        self._index.createIndex({'post': 2})

        saving_filename = self._get_saving_indexes_filename(self._pages_counter)

        self._index.saveIndex(saving_filename, True)
        self._save_features_dict()

        if self._page_size_counter == config.nn_page_size:
            self._index = nmslib.init(method='hnsw', space='cosinesimil')
            self._feature_dict = {}
            self._page_size_counter = 0
            self._pages_counter += 1

    def find_most_similar_images(self, pil_image, num=5):
        """
        :param pil_image:  image loaded py PIL library.
        :param num: number of neighbours to return
        :return: scores, nearest_descriptions.
                 scores - Scores of simularity between pil_image and neighbours
                 nearest_descriptions — descriptions of neighbours
        """
        features = self._get_features(pil_image)
        candidates_descriptions = []

        for page in range(self._pages_counter + 1):
            if page == self._pages_counter and self._page_size_counter == 0:
                continue

            index = nmslib.init(method='hnsw', space='cosinesimil')
            feature_dict = self._load_features_dict(page)

            index.loadIndex(self._get_saving_indexes_filename(page))

            indexes, scores = index.knnQuery(features, k=num)
            nearest_descriptions = [feature_dict[str(i)] for i in indexes]

            candidates_descriptions.extend(nearest_descriptions)

        index = nmslib.init(method='hnsw', space='cosinesimil')
        feature_dict = {}

        for i, description in enumerate(candidates_descriptions):
            feature_dict[i] = description
            [address, nft_id] = description.split('-')
            candidate_source, candidate_error = loader.get_image_source(address, nft_id)

            if not candidate_source or candidate_error:
                raise candidate_error

            pil_candidate = Image.open(BytesIO(candidate_source))
            candidate_features = self._get_features(pil_candidate)

            index.addDataPoint(data=candidate_features, id=i)

        index.createIndex({'post': 2})

        indexes, scores = index.knnQuery(features, k=num)
        nearest_descriptions = [feature_dict[i] for i in indexes]
        scores = self._transform_scores(scores)

        return scores, nearest_descriptions
