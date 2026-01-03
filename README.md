# MIME API

MIME 이메일에서 Display Part(본문)를 추출하는 REST API 서버입니다.

## Requirements

- Go 1.25+
- Docker (optional)

## Build

```bash
make build
```

## Run

```bash
make run
```

서버가 `http://localhost:8080`에서 실행됩니다.

## API

### POST /v1/display-part

이메일 파일(.eml)을 업로드하면 본문 콘텐츠를 추출하여 반환합니다.

```bash
curl -X POST http://localhost:8080/v1/display-part \
  -F "file=@email.eml"
```

### GET /health

헬스 체크 엔드포인트입니다.

```bash
curl http://localhost:8080/health
```

## Docker

```bash
# 빌드
make docker-build

# 실행
make docker-run
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `build` | 바이너리 빌드 (dist/) |
| `clean` | 빌드 아티팩트 삭제 |
| `test` | 테스트 실행 |
| `lint` | golangci-lint 실행 |
| `fmt` | 코드 포맷팅 |
| `run` | 서버 실행 |
| `docker-build` | Docker 이미지 빌드 |
| `docker-run` | Docker 컨테이너 실행 |
