# AI ile Öğrenilen Kapsamlı Yolculuk

## Tag Bazlı Gelişim Özeti

Bu depo, AI ile adım adım ilerleyerek geliştirilen bir öğrenme projesidir. Aşağıdaki etiketler, yolculuğun ana dönüm noktalarını gösterir:

| Tag | Öğrenilen / Tamamlanan Konu |
|---|---|
| `phase-01.00-baseline` | İlk proje temeli ve başlangıç yapısı |
| `phase-01.01-monorepo-structure` | Monorepo klasör düzeni |
| `phase-01.02-docker-compose-foundation` | PostgreSQL, Redis ve RabbitMQ ile compose altyapısı |
| `phase-01.03-proto-foundation` | Proto dosyaları ve gRPC temel düzeni |
| `phase-01.04-documentation-foundation` | Proje dokümantasyonunun ilk güçlü hali |
| `phase-01.05-ide-cleanup` | Çalışma alanı ve IDE düzeni |
| `phase-02.01-auth-service-skeleton` | Auth servis iskeleti |
| `phase-02.02-auth-domain-model` | Domain model tasarımı |
| `phase-02.03-auth-domain-validation` | Domain doğrulama kuralları |
| `phase-02.04-auth-ports` | Hexagonal ports tasarımı |
| `phase-02.05-auth-register-use-case` | Register use case |
| `phase-02.06-auth-register-use-case-tests` | Register testleri |
| `phase-02.07-auth-login-use-case` | Login use case |
| `phase-02.08-auth-login-use-case-tests` | Login testleri |
| `phase-02.09-bcrypt-password-hasher` | Şifre hashleme adaptörü |
| `phase-02.10-jwt-token-generator` | JWT üretimi |
| `phase-02.11-grpc-auth-handler` | gRPC handler entegrasyonu |
| `phase-02.12-auth-service-runnable` | Auth servisinin çalışır hale gelmesi |
| `phase-03.01-user-service-runnable` | User servisinin çalışır hale gelmesi |
| `phase-03.02-product-service-runnable` | Product servisinin çalışır hale gelmesi |
| `phase-03.03-auth-user-grpc-communication-runnable` | Auth ve User arasında gRPC iletişimi |
| `phase-03.04-inventory-service-runnable` | Inventory servisinin çalışır hale gelmesi |
| `phase-03.05-cart-service-runnable` | Cart servisinin Redis ile çalışması |
| `phase-03.06-order-payment-notification-runnable` | Order, Payment ve Notification servislerinin birlikte çalışması |
| `phase-04.01-api-gateway-complete` | API Gateway tamamlanması |

## Öğrenilen Ana Temalar

- Monorepo ve servis sınırları
- Hexagonal Architecture
- Domain, application, ports ve adapters ayrımı
- gRPC ile servisler arası iletişim
- PostgreSQL ve Redis ile veri sahipliği
- RabbitMQ ile olay tabanlı iletişim
- Docker Compose ile çoklu servis orkestrasyonu
- Tek ortak Dockerfile ile build standardizasyonu
- Container runtime içinde dosya ve migration yönetimi
- `docker compose config` ile doğrulama alışkanlığı

## En Önemli Ders

- Build tamamlanmış olsa bile runtime image içinde ihtiyaç duyulan dosyalar yoksa servis yine çöker.
- Bu yüzden uygulama kodu, Dockerfile, compose ve veri/migration path’leri birlikte düşünülmelidir.

# GoCommerceX

GoCommerceX, Go ile geliştirilmiş; Hexagonal Architecture, Microservice Architecture ve gRPC temelli öğrenme odaklı bir commerce backend platformudur.

## Amaç

Bu projenin amacı, gerçekçi bir commerce platformu üzerinden modern backend geliştirme pratiklerini öğrenmek ve uygulamaktır.

Ana öğrenme hedefleri:

- Go backend geliştirme
- Hexagonal Architecture
- Microservice Architecture
- gRPC iletişimi
- Protocol Buffers
- PostgreSQL
- Redis
- RabbitMQ
- JWT authentication
- Docker Compose
- Test, logging ve health check

## Servisler

- API Gateway
- Auth Service
- User Service
- Product Service
- Inventory Service
- Cart Service
- Order Service
- Payment Service
- Notification Service

## Mevcut Faz

Faz 01: Proje Temeli

Tamamlanan adımlar:

- 01.00 Baseline
- 01.01 Monorepo Structure
- 01.02 Docker Compose Foundation
- 01.03 Proto Foundation

Güncel adım:

- 01.04 Project Documentation Foundation

## Proje Yapısı

```text
GoCommerceX/
├── api-gateway/
├── auth-service/
├── user-service/
├── product-service/
├── inventory-service/
├── cart-service/
├── order-service/
├── payment-service/
├── notification-service/
├── proto/
├── deployments/
├── docs/
├── scripts/
├── .env.example
├── go.mod
└── task.md
```

## Altyapı

Yerel altyapı şu dosyada tanımlıdır:

```text
deployments/docker-compose.yml
```

İçerdiği servisler:

- PostgreSQL
- Redis
- RabbitMQ ve yönetim arayüzü

Docker Compose dosyasını doğrulamak için:

```bash
docker compose --env-file .env.example -f deployments/docker-compose.yml config
```

## Git İlerleme Etiketleri

- `phase-01.00-baseline`
- `phase-01.01-monorepo-structure`
- `phase-01.02-docker-compose-foundation`
- `phase-01.03-proto-foundation`
- `phase-01.04-documentation-foundation`
- `phase-01.05-ide-cleanup`
- `phase-02.01-auth-service-skeleton`
- `phase-02.02-auth-domain-model`
- `phase-02.03-auth-domain-validation`
- `phase-02.04-auth-ports`
- `phase-02.05-auth-register-use-case`
- `phase-02.06-auth-register-use-case-tests`
- `phase-02.07-auth-login-use-case`
- `phase-02.08-auth-login-use-case-tests`
- `phase-02.09-bcrypt-password-hasher`
- `phase-02.10-jwt-token-generator`
- `phase-02.11-grpc-auth-handler`
- `phase-02.12-auth-service-runnable`
- `phase-03.01-user-service-runnable`
- `phase-03.02-product-service-runnable`
- `phase-03.03-auth-user-grpc-communication-runnable`
- `phase-03.04-inventory-service-runnable`
- `phase-03.05-cart-service-runnable`
- `phase-03.06-order-payment-notification-runnable`
- `phase-04.01-api-gateway-complete`

## Geliştirme Kuralı

Her servis kendi verisine sahiptir.

Bir servis başka bir servisin veritabanına doğrudan bağlanmamalıdır. Servisler arası iletişim gRPC veya event mekanizması ile yapılmalıdır.
