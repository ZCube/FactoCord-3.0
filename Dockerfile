# syntax=docker/dockerfile:1.4
ARG FACTORIO_VERSION=latest
FROM golang:1.19-bullseye AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN version="$(git describe --tags) ($(git rev-parse --short HEAD))" \
    && echo "Version: $version" \
    && CGO_ENABLED=0 go build -ldflags "-X 'github.com/maxsupermanhd/FactoCord-3.0/v3/support.FactoCordVersion=$version'" -o factorio-cord .

FROM factoriotools/factorio:${FACTORIO_VERSION}

COPY --from=builder /app/factorio-cord /factorio-cord

COPY config-docker.json /config.json

RUN cat /docker-entrypoint.sh | grep -v -E "^exec" > /docker-entrypoint-discord.sh \
    && echo 'if [[ ! -f "/factorio/config.json" ]]; then cp /config.json /factorio/config.json; fi' >> /docker-entrypoint-discord.sh \
    && echo 'set +x' >> /docker-entrypoint-discord.sh \
    && echo 'cd factorio' >> /docker-entrypoint-discord.sh \
    && echo 'DISCORD_TOKEN=${DISCORD_TOKEN:-""}' >> /docker-entrypoint-discord.sh \
    && echo 'FACTORIO_CHANNEL_ID=${FACTORIO_CHANNEL_ID:-""}' >> /docker-entrypoint-discord.sh \
    && echo 'USERNAME=${USERNAME:-""}' >> /docker-entrypoint-discord.sh \
    && echo 'TOKEN=${TOKEN:-""}' >> /docker-entrypoint-discord.sh \
    && echo 'DISCORD_TOKEN_FILE=${DISCORD_TOKEN_FILE:-"/discord/token"}' >> /docker-entrypoint-discord.sh \
    && echo 'FACTORIO_CHANNEL_ID_FILE=${FACTORIO_CHANNEL_ID_FILE:-"/discord/factorio_channel_id"}' >> /docker-entrypoint-discord.sh \
    && echo 'USERNAME_FILE=${USERNAME_FILE:-"/account/username"}' >> /docker-entrypoint-discord.sh \
    && echo 'TOKEN_FILE=${TOKEN_FILE:-"/account/token"}' >> /docker-entrypoint-discord.sh \
    && echo 'exec $EXEC /factorio-cord \' >> /docker-entrypoint-discord.sh \
    && echo '    --discord-token=${DISCORD_TOKEN} --factorio-channel-id=${FACTORIO_CHANNEL_ID} --username=${USERNAME} --mod-portal-token=${TOKEN} \' >> /docker-entrypoint-discord.sh \
    && echo '    --discord-token-file=${DISCORD_TOKEN_FILE} --factorio-channel-id-file=${FACTORIO_CHANNEL_ID_FILE} --username-file=${USERNAME_FILE} --mod-portal-token-file=${TOKEN_FILE} \' >> /docker-entrypoint-discord.sh \
    && echo '    -- /opt/factorio/bin/x64/factorio "${FLAGS[@]}" "$@"' >> /docker-entrypoint-discord.sh \
    && chmod +x /docker-entrypoint-discord.sh

ENTRYPOINT ["/docker-entrypoint-discord.sh"]
