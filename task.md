# YAPAY ZEKA İLE GO ÖĞRENİYORUM

---
# Go Hexagonal Microservices Commerce Platform

Bu proje; **Go**, **Hexagonal Architecture**, **Microservice Architecture** ve **gRPC** kullanılarak geliştirilecek kapsamlı bir backend projesidir.

Amaç; staj sürecinde öğrenilmesi istenen modern backend mimarilerini gerçek bir proje üzerinde uygulamaktır.

---

## Proje Amacı

Bu proje ile aşağıdaki konuların öğrenilmesi ve uygulanması hedeflenmektedir:

* Go ile backend servis geliştirme
* Hexagonal Architecture kullanımı
* Microservice mimarisi
* Servisler arası gRPC iletişimi
* PostgreSQL ile veri yönetimi
* Redis ile sepet/cache yönetimi
* RabbitMQ ile event-driven yapı
* JWT tabanlı authentication
* Docker Compose ile çok servisli çalışma ortamı
* Test, logging ve health check yapıları

---

## Kullanılacak Teknolojiler

* Go
* gRPC
* Protocol Buffers
* PostgreSQL
* Redis
* RabbitMQ
* Docker
* Docker Compose
* JWT
* GitHub Actions

---

## Servisler

Projede aşağıdaki servisler yer alacaktır:

| Servis               | Görevi                                                  |
| -------------------- | ------------------------------------------------------- |
| API Gateway          | Dış istekleri karşılar ve ilgili servislere yönlendirir |
| Auth Service         | Register, login ve JWT işlemlerini yönetir              |
| User Service         | Kullanıcı profil bilgilerini yönetir                    |
| Product Service      | Ürün ve kategori işlemlerini yönetir                    |
| Inventory Service    | Stok kontrolü ve stok rezervasyonu yapar                |
| Cart Service         | Sepet işlemlerini yönetir                               |
| Order Service        | Sipariş oluşturma ve sipariş durumlarını yönetir        |
| Payment Service      | Mock ödeme işlemlerini yönetir                          |
| Notification Service | Sipariş ve ödeme bildirimlerini yönetir                 |

---

## Genel Mimari

```text
Client
  |
  v
API Gateway
  |
  |-- gRPC --> Auth Service
  |-- gRPC --> User Service
  |-- gRPC --> Product Service
  |-- gRPC --> Cart Service
  |-- gRPC --> Order Service
```

Sipariş oluşturma sırasında:

```text
Order Service
  |
  |-- gRPC --> User Service
  |-- gRPC --> Product Service
  |-- gRPC --> Inventory Service
  |-- gRPC --> Payment Service
  |
  |-- Event --> Notification Service
```

---

## Hexagonal Architecture

Her servis kendi içinde Hexagonal Architecture yapısına göre geliştirilecektir.

Örnek servis yapısı:

```text
service-name/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── domain/
│   ├── application/
│   ├── ports/
│   ├── adapters/
│   ├── infrastructure/
│   └── config/
│
├── proto/
├── migrations/
├── Dockerfile
└── go.mod
```

Katmanlar:

* `domain`: İş modelleri ve temel kurallar
* `application`: Use-case ve iş akışları
* `ports`: Interface yapıları
* `adapters`: PostgreSQL, gRPC, Redis, RabbitMQ implementasyonları
* `infrastructure`: Config, logger, server ve database bağlantıları

---

## Sipariş Akışı

```text
1. Kullanıcı sisteme kayıt olur.
2. Kullanıcı login olur ve JWT token alır.
3. Kullanıcı ürünleri listeler.
4. Kullanıcı ürünü sepete ekler.
5. Kullanıcı sipariş oluşturur.
6. Order Service kullanıcıyı kontrol eder.
7. Order Service ürün bilgisini alır.
8. Inventory Service stok rezervasyonu yapar.
9. Payment Service mock ödeme işlemini yapar.
10. Ödeme başarılıysa sipariş onaylanır.
11. Notification Service bildirim eventi alır.
```

---

## Veri Yönetimi

Her servis kendi verisini yönetir.

| Servis               | Veri Kaynağı          |
| -------------------- | --------------------- |
| Auth Service         | PostgreSQL            |
| User Service         | PostgreSQL            |
| Product Service      | PostgreSQL            |
| Inventory Service    | PostgreSQL            |
| Cart Service         | Redis                 |
| Order Service        | PostgreSQL            |
| Payment Service      | PostgreSQL            |
| Notification Service | PostgreSQL / RabbitMQ |

Önemli kural:

```text
Bir servis başka bir servisin veritabanına doğrudan bağlanmaz.
Gerekli bilgiler gRPC veya event üzerinden alınır.
```

---

## Proje Klasör Yapısı

```text
go-hexagonal-microservices-commerce-platform/
│
├── api-gateway/
├── auth-service/
├── user-service/
├── product-service/
├── inventory-service/
├── cart-service/
├── order-service/
├── payment-service/
├── notification-service/
│
├── proto/
│   ├── auth.proto
│   ├── user.proto
│   ├── product.proto
│   ├── inventory.proto
│   ├── cart.proto
│   ├── order.proto
│   ├── payment.proto
│   └── notification.proto
│
├── deployments/
│   └── docker-compose.yml
│
├── docs/
├── scripts/
├── Makefile
├── README.md
└── .env.example
```

---

## Geliştirme Aşamaları

### 1. Proje Temeli

* Monorepo yapısı kurulacak
* Docker Compose hazırlanacak
* Ortak proto klasörü oluşturulacak
* PostgreSQL, Redis ve RabbitMQ servisleri eklenecek

### 2. Temel Servisler

* Auth Service
* User Service
* Product Service
* Inventory Service

### 3. Sipariş Sistemi

* Cart Service
* Order Service
* Payment Service
* Notification Service

### 4. API Gateway

* REST endpointleri yazılacak
* JWT middleware eklenecek
* Gateway üzerinden gRPC servislerine bağlanılacak

### 5. Test ve Kalite

* Unit testler yazılacak
* Integration testler yazılacak
* Logging eklenecek
* Health check eklenecek

---

## Hedeflenen Ana Akış

Proje sonunda aşağıdaki akış çalışır durumda olacaktır:

```text
Register
  ↓
Login
  ↓
Product List
  ↓
Add to Cart
  ↓
Create Order
  ↓
Reserve Stock
  ↓
Mock Payment
  ↓
Confirm Order
  ↓
Send Notification
```

---

## Proje Durumu

```text
[ ] Project Setup
[ ] Auth Service
[ ] User Service
[ ] Product Service
[ ] Inventory Service
[ ] Cart Service
[ ] Order Service
[ ] Payment Service
[ ] Notification Service
[ ] API Gateway
[ ] Docker Compose
[ ] Tests
[ ] Documentation
```

---

## Kısa Özet

Bu proje, Go ile geliştirilecek kapsamlı bir microservice backend sistemidir.

Projede her servis Hexagonal Architecture ile tasarlanacak, servisler arası iletişim gRPC ile sağlanacak ve sistem Docker Compose üzerinden çalıştırılacaktır.

Bu yapı; Go, microservice, gRPC ve temiz mimari konularını öğrenmek için güçlü bir portfolyo projesi olarak tasarlanmıştır.
