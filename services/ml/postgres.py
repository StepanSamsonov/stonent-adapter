import sql
import config


class RegisteredImage:
    def __init__(self):
        self.id = None
        self.contract_address = None
        self.nft_id = None
        self.format = None

    def to_dict(self):
        return {
            'contract_address': self.contract_address,
            'nft_id': self.nft_id,
            'format': self.format,
        }


class RegisteredImages(sql.Table):
    name = 'registered_images'
    type = RegisteredImage
    fields = {
        'id': {'type': 'int'},
        'contract_address': {'type': 'string'},
        'nft_id': {'type': 'string'},
        'format': {'type': 'string'},
    }


class RejectedImageByIPFS:
    def __init__(self):
        self.id = None
        self.contract_address = None
        self.nft_id = None
        self.ipfs_path = None
        self.description = None

    def to_dict(self):
        return {
            'contract_address': self.contract_address,
            'nft_id': self.nft_id,
            'ipfs_path': self.ipfs_path,
            'description': self.description,
        }


class RejectedImagesByIPFS(sql.Table):
    name = 'rejected_images_by_ipfs'
    type = RejectedImageByIPFS
    fields = {
        'id': {'type': 'int'},
        'contract_address': {'type': 'string'},
        'nft_id': {'type': 'string'},
        'ipfs_path': {'type': 'string'},
        'description': {'type': 'string'},
    }


class RejectedImageByNN:
    def __init__(self):
        self.id = None
        self.contract_address = None
        self.nft_id = None
        self.description = None

    def to_dict(self):
        return {
            'contract_address': self.contract_address,
            'nft_id': self.nft_id,
            'description': self.description,
        }


class RejectedImagesByNN(sql.Table):
    name = 'rejected_images_by_nn'
    type = RejectedImageByNN
    fields = {
        'id': {'type': 'int'},
        'contract_address': {'type': 'string'},
        'nft_id': {'type': 'string'},
        'description': {'type': 'string'},
    }


def connect():
    connection_string = ' '.join([
        f'dbname={config.postgres_db_name}',
        f'user={config.postgres_user}',
        f'password={config.postgres_password}',
        f'host={config.postgres_host}',
        f'port={config.postgres_port}',
    ])
    sql.db = sql.Db(connection_string)
    sql.Table.schema = config.postgres_schema

    sql.db.init()

    sql.query(f'DROP SCHEMA IF EXISTS {config.postgres_schema} CASCADE')
    sql.query(f'CREATE SCHEMA IF NOT EXISTS {config.postgres_schema}')
    sql.query(f'''
    CREATE TABLE IF NOT EXISTS {config.postgres_schema}.{RegisteredImages.name} (
        id BIGSERIAL PRIMARY KEY NOT NULL,
        contract_address VARCHAR(128) NOT NULL,
        nft_id VARCHAR(128) NOT NULL,
        format VARCHAR(128) NOT NULL
    )''')
    sql.query(f'''
    CREATE TABLE IF NOT EXISTS {config.postgres_schema}.{RejectedImagesByIPFS.name} (
        id BIGSERIAL PRIMARY KEY NOT NULL,
        contract_address VARCHAR(128) NOT NULL,
        nft_id VARCHAR(128) NOT NULL,
        ipfs_path VARCHAR(128) NOT NULL,
        description VARCHAR(128) NOT NULL
    )''')
    sql.query(f'''
    CREATE TABLE IF NOT EXISTS {config.postgres_schema}.{RejectedImagesByNN.name} (
        id BIGSERIAL PRIMARY KEY NOT NULL,
        contract_address VARCHAR(128) NOT NULL,
        nft_id VARCHAR(128) NOT NULL,
        description VARCHAR(128) NOT NULL
    )''')

    print('Successfully connected to the Postgres')
