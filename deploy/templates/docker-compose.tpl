version: '3.8'

services:
  # PostgreSQL Database
  postgres-{PROJECT_NAME}:
    image: postgres:18-alpine
    container_name: postgres-{PROJECT_NAME}
    environment:
      POSTGRES_USER: {DB_USER}
      POSTGRES_PASSWORD: {DB_PASSWORD}
      POSTGRES_DB: {DB_NAME}
      POSTGRES_INITDB_ARGS: "--encoding=UTF8"
    ports:
      - "{DB_PORT}:5432"
    volumes:
      - postgres_{PROJECT_NAME}_data:/var/lib/postgresql/data
    networks:
      - {DOCKER_NETWORK}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U {DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    profiles:
      - db

  # Redis Cache
  redis-{PROJECT_NAME}:
    image: redis:7-alpine
    container_name: redis-{PROJECT_NAME}
    command: redis-server --appendonly yes --requirepass {REDIS_PASSWORD}
    ports:
      - "{REDIS_PORT}:6379"
    volumes:
      - redis_{PROJECT_NAME}_data:/data
    networks:
      - {DOCKER_NETWORK}
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    profiles:
      - cache

  # API Service
  api-{PROJECT_NAME}:
    image: {API_IMAGE_REPO}:{API_IMAGE_TAG}
    container_name: api-{PROJECT_NAME}
    environment:
      # 数据库配置
      DB_HOST: postgres-{PROJECT_NAME}
      DB_PORT: "5432"
      DB_USER: {DB_USER}
      DB_PASSWORD: {DB_PASSWORD}
      DB_NAME: {DB_NAME}
      DB_SSL_MODE: {DB_SSL_MODE}
      DB_TYPE: {DB_TYPE}
      
      # Redis 配置
      REDIS_HOST: redis-{PROJECT_NAME}:6379
      REDIS_PASSWORD: {REDIS_PASSWORD}
      REDIS_DB: {REDIS_DB}
      REDIS_POOL_SIZE: {REDIS_POOL_SIZE}
      
      # 认证配置
      AUTH_ACCESS_SECRET: {AUTH_ACCESS_SECRET}
      AUTH_ACCESS_EXPIRE: {AUTH_ACCESS_EXPIRE}
      OAUTH_STATE_SECRET: {OAUTH_STATE_SECRET}
      
      # 验证码配置
      CAPTCHA_TYPE: {CAPTCHA_TYPE}
      CAPTCHA_STORE_TYPE: {CAPTCHA_STORE_TYPE}

      # i18n 配置 (国际化)
      # 默认值: /app/etc/i18n/locale (内置 i18n 文件)
      # 设置为 /app/etc/i18n/custom 表示使用自定义 i18n 文件
      I18N_DIR: /app/etc/i18n/locale

      # RPC 服务
      RPC_TARGET: rpc-{PROJECT_NAME}:8080

      # 部署环境
      DEPLOY_ENV: {DEPLOY_ENV}
    ports:
      - "{API_PORT}:8889"
    volumes:
      # 可选: 挂载自定义 i18n 文件
      # - ./i18n:/app/etc/i18n/custom
    networks:
      - {DOCKER_NETWORK}
    depends_on:
      postgres-{PROJECT_NAME}:
        condition: service_healthy
      redis-{PROJECT_NAME}:
        condition: service_healthy
      rpc-{PROJECT_NAME}:
        condition: service_started
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8889/health"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # RPC Service
  rpc-{PROJECT_NAME}:
    image: {RPC_IMAGE_REPO}:{RPC_IMAGE_TAG}
    container_name: rpc-{PROJECT_NAME}
    environment:
      # 数据库配置
      DB_HOST: postgres-{PROJECT_NAME}
      DB_PORT: "5432"
      DB_USER: {DB_USER}
      DB_PASSWORD: {DB_PASSWORD}
      DB_NAME: {DB_NAME}
      DB_SSL_MODE: {DB_SSL_MODE}
      DB_TYPE: {DB_TYPE}
      
      # Redis 配置
      REDIS_HOST: redis-{PROJECT_NAME}:6379
      REDIS_PASSWORD: {REDIS_PASSWORD}
      REDIS_DB: {REDIS_DB}
      REDIS_POOL_SIZE: {REDIS_POOL_SIZE}
      
      # 认证配置
      OAUTH_STATE_SECRET: {OAUTH_STATE_SECRET}
      
      # 服务端口
      SERVER_PORT: "8080"
      
      # 部署环境
      DEPLOY_ENV: {DEPLOY_ENV}
    ports:
      - "{RPC_PORT}:8080"
    networks:
      - {DOCKER_NETWORK}
    depends_on:
      postgres-{PROJECT_NAME}:
        condition: service_healthy
      redis-{PROJECT_NAME}:
        condition: service_healthy
    restart: unless-stopped

networks:
  {DOCKER_NETWORK}:
    driver: bridge

volumes:
  postgres_{PROJECT_NAME}_data:
    driver: local
  redis_{PROJECT_NAME}_data:
    driver: local

