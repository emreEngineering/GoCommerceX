# GoCommerceX – Proje Referans Rehberi

> Son güncelleme: 2026-06-29
> Bu dosya, projenin amacını, mimarisini, mevcut durumunu ve sonraki adımları tek kaynakta toplar.
> Yeni bir yapay zeka asistanına devrederken bu dosyayı okutman yeterlidir.

---

## 1. Proje Amacı

GoCommerceX, Go ile modern backend geliştirme öğrenmek için yapılan bir **öğrenme platformudur**.

Kapsanan konular:
- Go, Hexagonal Architecture, Microservice, gRPC, Protobuf
- PostgreSQL, Redis, RabbitMQ, Docker, JWT

---

## 2. Mimari Kararlar

- Her servis **Hexagonal Architecture** (domain, application, ports, adapters, transport) ile yazılır.
- Her servis **kendi veritabanına sahiptir**; başka servisin DB'sine doğrudan bağlanmaz.
- Servisler arası iletişim: **gRPC** (senkron) veya **RabbitMQ** (asenkron).
- Testler **fake implementasyonlarla** izole yazılır.
- Git akışı: `main` branch + faz branch'leri + tag'ler.

---

## 3. Mevcut Durum

### Tamamlanan Servisler

| Servis | Durum | gRPC Port | Tag |
|--------|-------|-----------|-----|
| **Auth Service** | ✅ Çalışıyor | 50051 | `phase-02.12-auth-service-runnable` |
| **User Service** | ✅ Çalışıyor | 50052 | `phase-03.01-user-service-runnable` |
| **Product Service** | ✅ Geliştirildi | 50053 | - |

### Auth Service Detay
- **Metodlar:** Register, Login
- **Domain:** `User (ID, Email, PasswordHash, CreatedAt, UpdatedAt)`
- **Ports:** `UserRepository`, `PasswordHasher`, `TokenGenerator`
- **Adapters:** PostgreSQL (`pgx`), bcrypt, JWT
- **Tablo:** `auth-service` → `users` (id, email, password_hash, created_at, updated_at)

### User Service Detay
- **Metodlar:** CreateUser, GetUser, GetUserByEmail, UpdateUser, DeleteUser
- **Domain:** `User (ID, Email, FirstName, LastName, Phone, CreatedAt, UpdatedAt)`
- **Ports:** `UserRepository`
- **Adapters:** PostgreSQL (`pgx`)
- **Tablo:** `user-service` → `user_profiles` (id, email, first_name, last_name, phone, created_at, updated_at)

### Altyapı
- **Docker Compose:** PostgreSQL (5432), Redis (6379), RabbitMQ (5672, 15672)
- **Veritabanı:** PostgreSQL, kullanıcı: `gocommerce`, şifre: `gocommerce_password`, db: `gocommerce`
- **Proto tanımları:** 8 `.proto` dosyası var, sadece `auth.proto` ve `user.proto` için Go kodu üretildi.

### Bekleyen Servisler
- Inventory, Cart, Order, Payment, Notification, API Gateway (sadece boş klasörler var)

---

## 4. Proje Klasör Yapısı (Önemli Dosyalar)
GoCommerceX/
├── task.md ← bu dosya
├── go.mod
├── deployments/docker-compose.yml
├── proto/
│ ├── auth.proto
│ ├── user.proto
│ └── ... (diğer proto'lar)
├── auth-service/
│ ├── cmd/auth-service/main.go
│ └── internal/
│ ├── domain/user.go
│ ├── ports/{user_repository,password_hasher,token_generator,errors}.go
│ ├── application/{register_user,login_user}.go
│ ├── adapters/{bcrypt,jwt,postgres}.go
│ ├── transport/grpc/auth_handler.go
│ ├── config/config.go
│ └── infrastructure/migrate.go
├── user-service/
│ ├── cmd/user-service/main.go
│ └── internal/
│ ├── domain/user.go
│ ├── ports/{user_repository,errors}.go
│ ├── application/{create_user,get_user,update_delete_user}.go
│ ├── adapters/postgres_user_repository.go
│ ├── transport/grpc/user_handler.go
│ ├── config/config.go
│ └── infrastructure/migrate.go
└── (diğer boş servis klasörleri)



---

## 5. Çalışma Kuralları

1. Asistan açıklar, kullanıcı kodu yazar.
2. Asistan repo işlerini (branch, commit, tag, push) üstlenir.
3. Her faz tamamlandığında `phase-XX.YY-açıklama` formatında tag basılır.
4. Tag'li commit `main` branch'e merge edilir.

---

## 6. Sıradaki Adımlar (Önerilen)

1. Servisler arası iletişim (Auth → User gRPC çağrısı)
2. Inventory Service (stok)
3. Cart Service (Redis ile sepet)
4. Order + Payment + Notification (sipariş akışı)
5. API Gateway

---

## 7. Yeni Konuşma Başlatma Komutu

> GoCommerceX projesine devam et. `task.md` dosyasını oku, `main` branch'teyiz.
> Auth Service (50051), User Service (50052) ve Product Service (50053) hazır. Sıradaki adım servisler arası iletişim.
