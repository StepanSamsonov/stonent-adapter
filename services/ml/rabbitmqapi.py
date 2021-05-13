import rabbitmq

import config


def consume_events():
    connection = rabbitmq.Connection()

    connection.open(config.rabbit_host, config.rabbit_port)
    connection.login(f'amqp://{config.rabbit_login}:{config.rabbit_pass}@{config.rabbit_host}:{config.rabbit_port}/', 0)
    connection.declare_exchange(1, config.rabbit_queue, 'direct')
    connection.declare_queue(1, config.rabbit_queue)
    connection.bind_queue(1, config.rabbit_queue, config.rabbit_queue, config.rabbit_queue)

    while True:
        envelope = connection.consume_message()
        print(envelope.message.body)
        yield "contract_address", "nft_id", "source"

    return None
