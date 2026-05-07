FROM node:24-alpine AS web-build

WORKDIR /app

RUN npm install -g pnpm@10

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY web/package.json web/package.json
COPY packages/brand/package.json packages/brand/package.json
COPY packages/ui/package.json packages/ui/package.json

RUN pnpm install --frozen-lockfile --filter @netstamp/web...

COPY web web
COPY packages/brand packages/brand
COPY packages/ui packages/ui

RUN pnpm --filter @netstamp/web build

FROM nginx:1.27-alpine

COPY deployments/docker/nginx.conf /etc/nginx/nginx.conf
COPY --from=web-build /app/web/dist /usr/share/nginx/html

EXPOSE 80 9090
