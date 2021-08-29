FROM smartcontract/chainlink:0.10.3

RUN mkdir /chainlink

ENV DATABASE_URL="postgresql://guest:guest@postgres:5432/chainlink?sslmode=disable"
ENV ROOT="/chainlink"
ENV LOG_LEVEL="debug"
ENV MIN_OUTGOING_CONFIRMATIONS="2"
ENV CHAINLINK_TLS_PORT="0"
ENV SECURE_COOKIES="false"
ENV GAS_UPDATER_ENABLED="false"
ENV ALLOW_ORIGINS="*"
