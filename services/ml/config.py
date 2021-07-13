server_port = 9090

loader_url = 'http://localhost:8080'

rabbit_login = 'guest'
rabbit_password = 'guest'
rabbit_host = 'localhost'
rabbit_queue = 'indexing'
rabbit_port = 5672

postgres_db_name = 'postgres'
postgres_schema = 'schema'
postgres_user = 'guest'
postgres_password = 'guest'
postgres_host = 'localhost'
postgres_port = 5432

nn_page_size = 500_000
nn_index_dir = 'nn_indexes'
nn_index_file_prefix = f'./{nn_index_dir}/indexes-'
nn_index_file_postfix = '.hnsw'
nn_features_dict_dir = 'nn_features_dicts'
nn_features_dict_file_prefix = f'./{nn_features_dict_dir}/features-dict-'
nn_features_dict_file_postfix = '.json'
