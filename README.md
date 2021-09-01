# Stonent

## Start

You need to create some files in the current directory:
1. [.api](https://docs.chain.link/docs/miscellaneous/#use-password-and-api-files-on-startup) - credentials for Chainlink node with email on the first line and password on the second
2. [.password](https://docs.chain.link/docs/miscellaneous/#use-password-and-api-files-on-startup) - file with your wallet password
3. chainlink.env - specify ENV variables for Chainlink service:
    1. [ETH_CHAIN_ID](https://docs.chain.link/docs/configuration-variables/#eth_chain_id)
    2. [LINK_CONTRACT_ADDRESS](https://docs.chain.link/docs/configuration-variables/#link_contract_address)
    3. [ETH_URL](https://docs.chain.link/docs/configuration-variables/#eth_url)
4. loader.env - specify ENV variables for Stonent service:
    1. STONENT_CONTRACT_ADDRESS - our contract:
       1. Rinkeby: 0xFa9aF655Ef79445ECBb73389914e2ab16A31F62D
    2. COMMON_PROVIDER_URL - your Infura ws URL or something else
    3. COLLECTIONS_PROVIDER_URL - your Infura ws URL or something else
    4. DOWNLOAD_IMAGES_BUFFER_SIZE - count of images which can be downloading at the same time
    5. DOWNLOAD_IMAGES_MAX_COUNT - optional, specify max count of downloaded images from IPFS

You can see example configuration files [here](examples).

Then just run
```
docker compose up
```

To access Chainlink node:
```
http://localhost:6688
```

To check images indexing status:
```
GET http://localhost:9191/statistics
```

## What about stuffing?

There are five services which interact with each other:
1. [Loader](./services/loader/README.md) - downloading nft-images from IPFS
2. [ML](./services/ml/README.md) - neural network checking plagiarism
3. RabbitMQ - database-mediator between Loader and ML for images transfer
4. Postgres - database to store statistics and Chainlink stuff
5. Chainlink - Chainlink node

The whole system has such scheme:

![Scheme](./media/scheme.png)

## Using statistics server

ML service provides API to interact with statistics data.
It is available on port 9191.

You can retrieve the following information:
1. Check if image has been indexed
2. Get array of all indexed images
3. Get array of failed images while downloading by IPFS
4. Get array of failed images while indexing by NN
5. Get common statistics data

You can check [Swagger](./services/ml/statistics_server.yaml) for more details.

## Some comments

### About chainlink.env

1. Why do we need to specify different COMMON_PROVIDER_URL and COLLECTIONS_PROVIDER_URL?

Loader service interacts with ethereum provider and makes a lot of stuff.
We can divide this stuff into two categories:
- Reading and listening events from nft-collections
- Interacting with Stonent contract and performing some other minor stuff

Developing Loader we need to test functional.
We need to make sure that every category works correctly not only  with test network but also with mainnet.
Making your own configuration in production, COMMON_PROVIDER_URL and COLLECTIONS_PROVIDER_URL should be the same.

Not bug, just architecture feature.

2. What about DOWNLOAD_IMAGES_BUFFER_SIZE?

Computing performance of your machine can be different from other.
There are a lot of nft-images which should be indexed.
Of course, you can download all images at the same time, but IPFS node can die, or your machine can turn into pumpkin.
You can specify the optimal value of downloading images at the same time to prevent it.
For example, if you set variable to 1, all images will be downloaded sequentially.
If variable is 10, there will be 10 independent goroutines and every will download images in parallel.

3. What about DOWNLOAD_IMAGES_MAX_COUNT?

This variable is created for testing purposes only.
A lot of images start to download, when you start Stonent service.
If you want to test some configuration without downloading all images and spending traffic, you can specify the max count of images which will be downloaded after start.
Remove this variable or set to -1 in production.
