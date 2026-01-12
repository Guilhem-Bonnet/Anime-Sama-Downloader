FROM node:20-alpine AS web-build

WORKDIR /webapp
COPY webapp/package.json webapp/package-lock.json ./
RUN npm ci
COPY webapp/ ./
RUN npm run build


FROM python:3.12-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    ASD_WEB_HOST=0.0.0.0 \
    ASD_WEB_PORT=8000 \
    ASD_DOWNLOAD_ROOT=/data/videos

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

COPY . /app

# Embed the built SPA for production usage.
COPY --from=web-build /webapp/dist /app/webapp/dist

EXPOSE 8000

CMD ["python", "main.py", "--ui", "web"]
