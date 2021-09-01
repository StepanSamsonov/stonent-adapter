server_port = 9090
statistics_server_port = 9191

mode = 'PROD'

loader_host = 'localhost' if mode == 'DEV' else 'loader'
loader_url = f'http://{loader_host}:8080'

rabbit_host = 'localhost' if mode == 'DEV' else 'rabbitmq'
rabbit_login = 'guest'
rabbit_password = 'guest'
rabbit_queue = 'indexing'
rabbit_port = 5672

postgres_host = 'localhost' if mode == 'DEV' else 'postgres'
postgres_db_name = 'postgres'
postgres_schema = 'schema'
postgres_user = 'guest'
postgres_password = 'guest'
postgres_port = 5432

nn_page_size = 500_000
nn_existed_blocks_file = 'nn_blocks.txt'
nn_existed_blocks_dir = './'
nn_index_dir = 'nn_indexes'
nn_index_file_prefix = f'./{nn_index_dir}/indexes-'
nn_index_file_postfix = '.hnsw'
nn_descriptions_dir = 'nn_descriptions'
nn_descriptions_file_prefix = f'./{nn_descriptions_dir}/descriptions-'
nn_descriptions_file_postfix = '.json'
