FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o proxy-broadcast

# ---

FROM ubuntu:24.04

# 비root 사용자 생성
ARG USERNAME=appuser
ARG USER_UID=1001
ARG USER_GID=1001

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    # 그룹 및 사용자 생성, 로그인 불가
    && groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID --no-create-home --shell /usr/sbin/nologin $USERNAME

WORKDIR /app

COPY --from=builder /app/proxy-broadcast .

# HEALTHCHECK 추가_
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --spider -q http://localhost:8083/health || exit 1

EXPOSE 8083

# 비root 사용자 권한으로 실행
USER $USERNAME

CMD ["./proxy-broadcast"]