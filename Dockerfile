# syntax=docker/dockerfile:1.4
ARG FACTORIO_VERSION=latest
FROM golang:1.19-bullseye AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN export version=$(git describe --tags $(git rev-parse --short HEAD)) \
 && echo "Version: $version" \
 && CGO_ENABLED=0 go build -ldflags "-X 'github.com/maxsupermanhd/FactoCord-3.0/support.FactoCordVersion=$version'" -o factorio-cord .

FROM factoriotools/factorio:${FACTORIO_VERSION}

COPY --from=builder /app/factorio-cord /factorio-cord

COPY config-docker.json /config.json

RUN cat /docker-entrypoint.sh | grep -v -E "^exec" > /docker-entrypoint-discord.sh \
 && echo 'if [[ ! -f "/factorio/config.json" ]]; then cp /config.json /factorio/config.json; fi' >> /docker-entrypoint-discord.sh \
 && echo 'cd factorio' >> /docker-entrypoint-discord.sh \
 && echo 'exec $SU_EXEC /factorio-cord -- /opt/factorio/bin/x64/factorio "${FLAGS[@]}" "$@"' >> /docker-entrypoint-discord.sh \
 && chmod +x /docker-entrypoint-discord.sh

ENTRYPOINT ["/docker-entrypoint-discord.sh"]
