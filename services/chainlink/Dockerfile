FROM smartcontract/chainlink:0.10.3

# Create the chainlink node root path
RUN mkdir /chainlink

# Arg for api user email, with default value
#ARG API_USER_EMAIL="user@example.com"
#ENV API_USER_EMAIL=$API_USER_EMAIL

# Arg for api user password, with default value
#ARG API_USER_PASSWORD="PA@SSword1234!567"
#ENV API_USER_PASSWORD=$API_USER_PASSWORD

# Arg for node wallet password, with default value
#ARG WALLET_PASSWORD="PA@SSword1234!567"
#ENV WALLET_PASSWORD=$WALLET_PASSWORD

# Create chainlink node required values to initialize with
RUN echo $API_USER_EMAIL > /chainlink/.api
RUN echo $API_USER_PASSWORD >> /chainlink/.api
RUN echo $WALLET_PASSWORD > /chainlink/.password

ENV DATABASE_URL="postgresql://guest:guest@postgres:5432/chainlink?sslmode=disable"
ENV ROOT="/chainlink"
ENV LOG_LEVEL="debug"
ENV ETH_CHAIN_ID="4"
ENV MIN_OUTGOING_CONFIRMATIONS="2"
ENV CHAINLINK_TLS_PORT="0"
ENV SECURE_COOKIES="false"
ENV GAS_UPDATER_ENABLED="false"
ENV ALLOW_ORIGINS="*"
