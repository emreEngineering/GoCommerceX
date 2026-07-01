

# 1.  Mikroservis Mimarisi – Proje Örnekleri

---

###  Servis Bağımsızlığı

**Ne demek?** Her servis kendi başına çalışır, kendi koduna ve veritabanına sahiptir.

**Projedeki örnek:**
```
auth-service/              user-service/
├── cmd/                   ├── cmd/
├── internal/              ├── internal/
│   ├── domain/            │   ├── domain/
│   ├── application/       │   ├── application/
│   ├── ports/             │   ├── ports/
│   └── adapters/          │   └── adapters/
└── migrations/            └── migrations/
```
Auth Service ve User Service tamamen ayrı klasörlerdedir. Auth Service'i değiştirirken User Service'e dokunmayız. Hatta Auth Service çökse bile User Service çalışmaya devam eder.

---

### Tek Sorumluluk (Single Responsibility)

**Ne demek?** Her servis sadece bir işe odaklanır.

**Projedeki örnek:**

| Servis | Tek Sorumluluğu |
|--------|-----------------|
| Auth Service | Register, Login, JWT üretimi |
| User Service | Kullanıcı profili CRUD |
| Product Service | Ürün ekleme, listeleme, güncelleme |
| Cart Service | Sepete ürün ekleme/çıkarma |
| Order Service | Sipariş oluşturma ve takibi |

Auth Service "kullanıcı profili güncelleme" yapmaz, User Service "şifre doğrulama" yapmaz. Herkes kendi işine bakar.

---

###  Servisler Arası İletişim

**Ne demek?** Servisler birbiriyle iki şekilde konuşur: senkron (gRPC) veya asenkron (RabbitMQ).

**Projedeki örnek – Senkron (gRPC):**
```
Auth Service ──gRPC──> User Service
```
Bir kullanıcı kayıt olduğunda, Auth Service **gRPC client** ile User Service'e bağlanır ve `CreateUser` çağrısı yapar. User Service profili oluşturur ve cevap döner. Auth Service cevabı bekler (senkron).

**Projedeki örnek – Asenkron (RabbitMQ):**
```
Order Service ──event──> RabbitMQ ──event──> Notification Service
```
Sipariş oluşunca Order Service bir `OrderCreated` mesajı gönderir. Notification Service bu mesajı dinler ve bildirim oluşturur. Order Service, Notification Service'in işini bitirmesini beklemez (asenkron).

---

###  Veri Sahipliği (Data Ownership)

**Ne demek?** Her servis kendi veritabanına sahiptir. Başka servisin veritabanına doğrudan bağlanamaz.

**Projedeki örnek:**

```
Auth Service ──────> PostgreSQL (users tablosu)
User Service ──────> PostgreSQL (user_profiles tablosu)
Cart Service ──────> Redis (sepet verileri)
```

Auth Service `users` tablosuna yazar, User Service `user_profiles` tablosuna yazar. **Aynı veritabanını kullansalar bile farklı tablolar kullanırlar.** Auth Service asla `user_profiles` tablosuna doğrudan yazmaz, gerekirse gRPC ile User Service'i çağırır.

---

### Bağımsız Dağıtım (Independent Deployment)

**Ne demek?** Bir servisi güncellerken diğerlerini durdurmak gerekmez.

**Projedeki örnek:**
```bash
# Sadece Auth Service'i yeniden başlat
go run ./auth-service/cmd/auth-service/

# User Service hiç etkilenmez, çalışmaya devam eder
```

Gerçek hayatta her servis ayrı bir Docker konteyneri olarak çalışır. Auth Service'i güncelleyip yeniden başlattığında, User Service ve diğerleri çalışmaya devam eder.

---

###  Hata İzolasyonu (Fault Isolation)

**Ne demek?** Bir servis çökerse diğerleri çalışmaya devam eder.

**Projedeki örnek:**

```
Cart Service (Redis bağlantısı koptu, çöktü)
    ❌ Sepet işlemleri yapılamaz
    ✅ Auth Service çalışmaya devam eder (giriş yapılabilir)
    ✅ Product Service çalışmaya devam eder (ürünler listelenebilir)
```

Eğer tüm sistemi tek bir büyük uygulama (monolith) olarak yapsaydık, Cart Service'in çökmesi tüm siteyi çökertirdi. Mikroservis sayesinde sadece sepet özelliği devre dışı kalır.

---

###  Polyglot Persistence (Farklı Veritabanları)

**Ne demek?** Her servis ihtiyacına göre farklı veritabanı teknolojisi kullanabilir.

**Projedeki örnek:**

| Servis | Veritabanı | Neden? |
|--------|-----------|--------|
| Auth Service | PostgreSQL | Kalıcı, güvenilir |
| User Service | PostgreSQL | Kalıcı, ilişkisel |
| Cart Service | Redis | Hızlı, geçici veri |
| Order Service | PostgreSQL | Kalıcı, transactional |

Sepet verisi geçicidir, kullanıcı oturumu kapatınca silinebilir. O yüzden hızlı ve bellek içi Redis kullanılır. Sipariş verisi ise kalıcı olmalıdır, o yüzden PostgreSQL kullanılır.

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Servis Bağımsızlığı | `auth-service/` ve `user-service/` ayrı klasörler, ayrı kodlar |
| Tek Sorumluluk | Auth = şifre, User = profil, Cart = sepet |
| İletişim (Senkron) | Auth → User gRPC çağrısı (`CreateUser`) |
| İletişim (Asenkron) | Order → RabbitMQ → Notification |
| Veri Sahipliği | Auth → `users`, User → `user_profiles` |
| Bağımsız Dağıtım | `go run ./auth-service/...` sadece onu etkiler |
| Hata İzolasyonu | Cart çökse Auth çalışır |
| Polyglot Persistence | PostgreSQL (çoğu servis) + Redis (Cart) |

---
# 2.  Hexagonal Architecture – Proje Örnekleri

---

##  Önce Büyük Resim

Aşağıda **Auth Service**'in hexagonal yapısını görüyorsun. Her katmanın ne iş yaptığını, hangi dosyada olduğunu ve neden orada olduğunu anlatacağım.

```
auth-service/
├── cmd/auth-service/main.go          ← Transport + Bağımlılıkları bağlama
└── internal/
    ├── domain/user.go                 ← Domain
    ├── application/
    │   ├── register_user.go          ← Application (Use Case)
    │   └── login_user.go             ← Application (Use Case)
    ├── ports/
    │   ├── user_repository.go        ← Port (Interface)
    │   ├── password_hasher.go        ← Port (Interface)
    │   └── token_generator.go        ← Port (Interface)
    ├── adapters/
    │   ├── bcrypt_password_hasher.go ← Adapter
    │   ├── jwt_token_generator.go    ← Adapter
    │   └── postgres_user_repository.go ← Adapter
    ├── transport/grpc/
    │   └── auth_handler.go           ← Transport (gRPC)
    ├── config/config.go              ← Config
    └── infrastructure/migrate.go     ← Infrastructure
```

---

### Domain Katmanı

**Ne demek?** İş nesnelerinin ve temel kuralların bulunduğu, **hiçbir dış bağımlılığı olmayan** katman.

**Projedeki örnek:** `auth-service/internal/domain/user.go`

```go
package domain

type User struct {
    ID           string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Domain kuralı: email boş olamaz
func (u User) Validate() error {
    if u.Email == "" {
        return ErrUserEmailRequired
    }
    return nil
}
```

**Dikkat et:** Bu dosyada hiç `sql`, `grpc`, `http` import'u yok!  
Domain, veritabanına nasıl kaydedildiğini **bilmez**. Sadece "bir kullanıcının email'i olmalı" kuralını bilir.

---

### Application Katmanı

**Ne demek?** Use case'leri (iş akışlarını) içeren katman. Domain'i ve Port'ları kullanır, ama adaptörlerin gerçek implementasyonunu bilmez.

**Projedeki örnek:** `auth-service/internal/application/register_user.go`

```go
type RegisterUserUseCase struct {
    userRepo       ports.UserRepository    // ← Interface, gerçek implementasyon değil!
    passwordHasher ports.PasswordHasher    // ← Interface
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
    // 1. Email/password kontrolü
    // 2. Email zaten var mı?
    existingUser, _ := uc.userRepo.FindByEmail(ctx, email)
    // 3. Şifreyi hash'le
    hashedPassword, _ := uc.passwordHasher.Hash(ctx, password)
    // 4. User oluştur ve validate et
    user := domain.NewUser(email, hashedPassword)
    // 5. Kaydet
    uc.userRepo.Save(ctx, user)
    // ...
}
```

**Dikkat et:** Use case, `PostgresUserRepository` veya `BcryptPasswordHasher` **bilmez**.  
Sadece `ports.UserRepository` ve `ports.PasswordHasher` interface'lerini bilir.  
Yarın PostgreSQL yerine MongoDB kullansak, use case'e **dokunmayız bile**.

---

###  Ports (Arayüzler)

**Ne demek?** "Benim neye ihtiyacım var?" sorusunun cevabı olan interface'ler.

**Projedeki örnek:** `auth-service/internal/ports/user_repository.go`

```go
package ports

type UserRepository interface {
    Save(ctx context.Context, user domain.User) error
    FindByEmail(ctx context.Context, email string) (domain.User, error)
}
```

**Bu ne diyor?**  
"Bana bir şey lazım ki, User'ı kaydedebileyim ve email ile bulabileyim. Nasıl yaptığı umurumda değil."

Başka bir örnek: `auth-service/internal/ports/password_hasher.go`

```go
type PasswordHasher interface {
    Hash(ctx context.Context, plainPassword string) (string, error)
    Compare(ctx context.Context, plainPassword string, passwordHash string) error
}
```

---

###  Adapters (Uyarlayıcılar)

**Ne demek?** Port'ların **gerçek implementasyonları**. Dış dünyayla (veritabanı, şifreleme kütüphanesi) iletişimi sağlar.

**Projedeki örnek – PostgreSQL:** `auth-service/internal/adapters/postgres_user_repository.go`

```go
type PostgresUserRepository struct {
    pool *pgxpool.Pool    // ← Gerçek veritabanı bağlantısı
}

// Port'taki Save interface'ini implemente eder
func (r *PostgresUserRepository) Save(ctx context.Context, user domain.User) error {
    query := `INSERT INTO "users" (id, email, password_hash, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5)`
    _, err := r.pool.Exec(ctx, query, ...)
    return err
}
```

**Projedeki örnek – bcrypt:** `auth-service/internal/adapters/bcrypt_password_hasher.go`

```go
type BcryptPasswordHasher struct{}

func (h *BcryptPasswordHasher) Hash(ctx context.Context, plainPassword string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    return string(bytes), err
}
```

**Projedeki örnek – JWT:** `auth-service/internal/adapters/jwt_token_generator.go`

```go
type JWTTokenGenerator struct {
    secretKey []byte
}

func (g *JWTTokenGenerator) Generate(ctx context.Context, userID string, email string) (string, error) {
    // jwt oluşturma kodu...
}
```

---

###  Transport Katmanı

**Ne demek?** Dış dünyadan gelen istekleri karşılayan katman. Bizim projede **gRPC handler**.

**Projedeki örnek:** `auth-service/internal/transport/grpc/auth_handler.go`

```go
type AuthHandler struct {
    authv1.UnimplementedAuthServiceServer
    registerUseCase *application.RegisterUserUseCase
    loginUseCase    *application.LoginUserUseCase
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
    // 1. gRPC isteğini al
    input := application.RegisterUserInput{
        Email:    req.GetEmail(),
        Password: req.GetPassword(),
    }
    // 2. Use case'i çağır
    output, err := h.registerUseCase.Execute(ctx, input)
    // 3. gRPC cevabına dönüştür
    return &authv1.RegisterResponse{
        UserId: output.UserID,
        Email:  output.Email,
    }, nil
}
```

**Transport'un tek işi:** İsteği al, use case'e ilet, cevabı dön.  
Kendi başına iş mantığı içermez.

---

###  Infrastructure

**Ne demek?** Veritabanı bağlantısı, migration, logger gibi altyapısal kodlar.

**Projedeki örnek:** `auth-service/internal/infrastructure/migrate.go`

```go
func RunMigrations(pool *pgxpool.Pool) error {
    // migration dosyasını oku ve çalıştır
    sql, _ := os.ReadFile("migrations/001_create_users_table.sql")
    _, err := pool.Exec(context.Background(), string(sql))
    return err
}
```

---

### Config

**Ne demek?** Ortam değişkenlerinden ayarları okuyan katman.

**Projedeki örnek:** `auth-service/internal/config/config.go`

```go
func Load() *Config {
    return &Config{
        GRPCPort:   getEnv("GRPC_PORT", "50051"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBUser:     getEnv("DB_USER", "gocommerce"),
        JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
    }
}
```

**Neden ayrı bir config katmanı?**  
Ayarlar tek bir yerden yönetilir. Port değişince tüm kodda aramak gerekmez.

---

###  Dependency Inversion (Bağımlılık Tersine Çevirme)

**Ne demek?** Üst katmanlar alt katmanlara değil, interface'lere bağımlı olmalıdır.

**Projedeki örnek:**

```go
// ❌ YANLIŞ: Application doğrudan adapter'a bağımlı
type RegisterUserUseCase struct {
    userRepo *adapters.PostgresUserRepository  // kötü!
}

// ✅ DOĞRU: Application, port interface'ine bağımlı
type RegisterUserUseCase struct {
    userRepo ports.UserRepository  // iyi!
}
```

`main.go` dosyasında bağımlılıkları birleştiriyoruz:

```go
userRepo := adapters.NewPostgresUserRepository(pool)       // Adapter (gerçek)
registerUseCase := application.NewRegisterUserUseCase(userRepo)  // Application (interface olarak alır)
```

---

###  Test Edilebilirlik

**Ne demek?** Use case'leri gerçek veritabanı olmadan, fake adapter'larla test edebilme.

**Projedeki örnek:** `auth-service/internal/application/register_user_test.go`

```go
// Gerçek veritabanı yerine fake repository
type fakeUserRepository struct {
    users map[string]domain.User
}

func (f *fakeUserRepository) Save(ctx context.Context, user domain.User) error {
    f.users[user.Email] = user
    return nil
}

func TestRegisterUser_Success(t *testing.T) {
    fakeRepo := &fakeUserRepository{users: make(map[string]domain.User)}
    fakeHasher := &fakePasswordHasher{}
    useCase := application.NewRegisterUserUseCase(fakeRepo, fakeHasher)
    
    output, err := useCase.Execute(ctx, input)
    // test assertions...
}
```

Gerçek PostgreSQL çalışıyor olmasa bile testler çalışır!

---

## Özet Tablo

| Katman | Dosya Örneği | Bağımlı Olduğu Şey |
|--------|-------------|-------------------|
| **Domain** | `domain/user.go` | Hiçbir şey |
| **Application** | `application/register_user.go` | Domain + Ports |
| **Ports** | `ports/user_repository.go` | Domain |
| **Adapters** | `adapters/postgres_user_repository.go` | Ports + Dış kütüphaneler |
| **Transport** | `transport/grpc/auth_handler.go` | Application + Proto |
| **Infrastructure** | `infrastructure/migrate.go` | Dış kütüphaneler |
| **Config** | `config/config.go` | İşletim sistemi (`os.Getenv`) |

---
# 3.  gRPC – Proje Örnekleri

---

##  Önce Büyük Resim

**gRPC**, servislerin birbiriyle konuşmasını sağlayan hızlı bir iletişim protokolüdür.  
Projemizde her servis bir gRPC sunucusudur ve birbirlerini gRPC istemcisi olarak çağırır.

```text
 ┌──────────────┐       gRPC (50051)       ┌──────────────┐
 │ Auth Service │ ◄────────────────────── ► │ User Service │
 └──────────────┘                           └──────────────┘
       │ gRPC (50051)                              │ gRPC (50052)
       │                                           │
 ┌─────▼───────────────────────────────────────────▼──────┐
 │                     API Gateway (:8080)                 │
 └────────────────────────────────────────────────────────┘
```

---

###  Protobuf (Proto) Dosyası – Sözleşme

**Ne demek?** Servisin hangi metodları olduğunu, hangi verileri alıp verdiğini tanımlayan `.proto` dosyasıdır. Bu, servisler arasındaki **ortak dil**dir.

**Projedeki örnek:** `proto/auth.proto`

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
}

message RegisterRequest {
  string email = 1;
  string password = 2;
}

message RegisterResponse {
  string user_id = 1;
  string email = 2;
}
```

**Bu ne demek?**  
"Auth Service'in iki metodu var: Register ve Login. Register'a `RegisterRequest` verirsen `RegisterResponse` alırsın."  
Tüm servisler bu sözleşmeye uymak zorundadır. Değişiklik olursa herkesin haberi olur.

---

### Proto'dan Kod Üretme

**Ne demek?** `.proto` dosyasını `protoc` ile Go koduna çevirmek. Böylece elle request/response struct'ı yazmak zorunda kalmayız.

**Projedeki örnek:**

```bash
# auth.proto'dan Go kodu üret
protoc --go_out=. --go-grpc_out=. proto/auth.proto
```

Bu komut iki dosya oluşturur:
- `proto/auth/v1/auth.pb.go` → Mesaj struct'ları (`RegisterRequest`, `RegisterResponse`)
- `proto/auth/v1/auth_grpc.pb.go` → Servis interface'i ve client/server kodları

**Oluşan kod (auth.pb.go):**
```go
type RegisterRequest struct {
    Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
    Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}
```

Artık `RegisterRequest` struct'ını elle yazmıyoruz, protobuf'tan otomatik geliyor.

---

### gRPC Sunucu (Server)

**Ne demek?** Servisi implemente eden, istekleri dinleyip cevap veren taraftır.

**Projedeki örnek:** `auth-service/cmd/auth-service/main.go`

```go
// 1. TCP listener oluştur
lis, err := net.Listen("tcp", ":50051")

// 2. gRPC sunucusu oluştur
grpcServer := grpc.NewServer()

// 3. Auth handler'ı sunucuya kaydet
authv1.RegisterAuthServiceServer(grpcServer, authHandler)

// 4. Reflection ekle (grpcurl ile keşif için)
reflection.Register(grpcServer)

// 5. Dinlemeye başla
grpcServer.Serve(lis)
```

**Adım adım ne oldu?**
1. 50051 portunu dinlemeye başladık
2. gRPC sunucusu oluşturduk
3. "Bu sunucuda AuthService çalışacak, handler'ı bu" dedik
4. `grpcurl` ile keşfedilebilmesi için reflection ekledik
5. Sunucuyu başlattık

---

###  gRPC Handler – İstek Karşılama

**Ne demek?** Proto'da tanımlanan metodları implemente eden yapıdır. İstekleri alır, use case'e iletir, cevabı döndürür.

**Projedeki örnek:** `auth-service/internal/transport/grpc/auth_handler.go`

```go
type AuthHandler struct {
    authv1.UnimplementedAuthServiceServer  // Proto'nun ürettiği base
    registerUseCase *application.RegisterUserUseCase
    loginUseCase    *application.LoginUserUseCase
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
    // gRPC isteğini → application input'una çevir
    input := application.RegisterUserInput{
        Email:    req.GetEmail(),
        Password: req.GetPassword(),
    }
    
    output, err := h.registerUseCase.Execute(ctx, input)
    if err != nil {
        // Hataları gRPC durum kodlarına çevir
        return nil, status.Error(codes.AlreadyExists, err.Error())
    }
    
    // application output'unu → gRPC cevabına çevir
    return &authv1.RegisterResponse{
        UserId: output.UserID,
        Email:  output.Email,
    }, nil
}
```

---

###  gRPC İstemci (Client)

**Ne demek?** Başka bir servisi çağıran taraftır.

**Projedeki örnek – Auth Service, User Service'i çağırıyor:**

```go
// Auth Service içinde User Service client'ı oluşturma
conn, _ := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
userClient := userv1.NewUserServiceClient(conn)

// User Service'in CreateUser metodunu çağırma
resp, _ := userClient.CreateUser(ctx, &userv1.CreateUserRequest{
    Email:     "john@example.com",
    FirstName: "John",
    LastName:  "Doe",
})
// resp.User → oluşturulan kullanıcı bilgisi
```

**Bu ne yaptı?** Auth Service, User Service'in 50052 portuna bağlandı, `CreateUser` metodunu çağırdı ve cevabı aldı. Tıpkı yerel bir fonksiyon çağırır gibi!

---

###  gRPC Durum Kodları

**Ne demek?** Her cevap bir durum kodu ile döner. Başarılı ise `OK` (0), hata varsa özel kod kullanılır.

**Projedeki örnek:**

| Go Hatası | gRPC Durum Kodu | Anlamı |
|-----------|----------------|--------|
| `ErrRegisterEmailRequired` | `InvalidArgument` | Eksik parametre |
| `ErrUserAlreadyExists` | `AlreadyExists` | Kayıt zaten var |
| `ErrInvalidCredentials` | `Unauthenticated` | Yanlış şifre/email |
| `ErrUserNotFound` | `NotFound` | Kullanıcı bulunamadı |
| `bilinmeyen hata` | `Internal` | Sunucu hatası |

**Kodda nasıl kullanılıyor?**

```go
case errors.Is(err, application.ErrUserAlreadyExists):
    return nil, status.Error(codes.AlreadyExists, err.Error())
```

---

###  Reflection (Servis Keşfi)

**Ne demek?** Sunucunun hangi metodları sunduğunu dışarıya tanıtmasıdır. `grpcurl` gibi araçlar bu sayede servisi keşfeder.

**Projedeki örnek:**

```go
reflection.Register(grpcServer)  // main.go'da
```

**Test ederken:**
```bash
grpcurl -plaintext localhost:50051 list
# Çıktı: auth.v1.AuthService

grpcurl -plaintext localhost:50051 describe auth.v1.AuthService
# Çıktı: Register, Login metodları
```

Reflection olmasaydı, metodları tahmin etmek veya proto dosyasına bakmak zorunda kalırdık.

---

### grpcurl ile Test Etme

**Ne demek?** `curl`'ün gRPC versiyonu. Komut satırından gRPC isteği atmayı sağlar.

**Projedeki örnek:**

```bash
# Register testi
grpcurl -plaintext -d '{"email": "test@example.com", "password": "secret123"}' \
  localhost:50051 auth.v1.AuthService/Register

# Login testi
grpcurl -plaintext -d '{"email": "test@example.com", "password": "secret123"}' \
  localhost:50051 auth.v1.AuthService/Login
```

---

###  Servisler Arası İletişim (Tam Örnek)

**Senaryo:** Kullanıcı kayıt olduğunda, Auth Service otomatik olarak User Service'te profil oluştursun.

**Adım 1:** Auth Service'e gRPC isteği gelir (Register)  
**Adım 2:** Auth Service kullanıcıyı `users` tablosuna kaydeder  
**Adım 3:** Auth Service, User Service'e gRPC client ile `CreateUser` çağrısı yapar  
**Adım 4:** User Service profili `user_profiles` tablosuna kaydeder  
**Adım 5:** Cevap döner

```go
// Auth Service register use case içinde:
func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
    // ... kullanıcıyı auth veritabanına kaydet ...
    
    // User Service'e profil oluşturma isteği gönder
    _, err = uc.userClient.CreateUser(ctx, &userv1.CreateUserRequest{
        Email:     user.Email,
        FirstName: "", // başlangıçta boş
        LastName:  "",
    })
    // ...
}
```

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Proto Dosyası | `proto/auth.proto` – servis ve mesaj tanımı |
| Kod Üretme | `protoc --go_out=. --go-grpc_out=. proto/auth.proto` |
| Sunucu | `grpc.NewServer()` + `RegisterAuthServiceServer()` |
| Handler | `auth_handler.go` – isteği alır, use case'i çağırır |
| İstemci | `NewUserServiceClient(conn)` – başka servisi çağırır |
| Durum Kodları | `InvalidArgument`, `AlreadyExists`, `Unauthenticated` |
| Reflection | `reflection.Register(grpcServer)` – grpcurl keşfi |
| Test Aracı | `grpcurl` – komut satırından test |

---
# 4. Protocol Buffers (Protobuf) – Proje Örnekleri

---

##  Önce Büyük Resim

**Protobuf**, servisler arasında veri taşımak için kullanılan, JSON'dan çok daha hızlı ve küçük bir veri formatıdır. `.proto` dosyaları ile tanımlanır, sonra `protoc` ile Go (veya başka dil) koduna dönüştürülür.

```text
auth.proto (senin yazdığın)
    │
    ▼ protoc + protoc-gen-go + protoc-gen-go-grpc
    │
    ├── auth.pb.go       (mesaj yapıları)
    └── auth_grpc.pb.go  (servis interface + client)
```

---

###  Message (Mesaj) Tanımı

**Ne demek?** Taşınacak verinin yapısını tanımlar. Tıpkı Go'daki `struct` gibi.

**Projedeki örnek:** `proto/auth.proto`

```protobuf
message RegisterRequest {
  string email = 1;
  string password = 2;
}

message RegisterResponse {
  string user_id = 1;
  string email = 2;
}
```

**Go karşılığı (auth.pb.go):**
```go
type RegisterRequest struct {
    Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
    Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

type RegisterResponse struct {
    UserId string `protobuf:"bytes,1,opt,name=user_id,proto3" json:"user_id,omitempty"`
    Email  string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
}
```

**Dikkat et:** Proto'daki `user_id`, Go'da `UserId` olur. `GetUserId()` metodu da otomatik oluşur.

---

###  Service (Servis) Tanımı

**Ne demek?** Hangi RPC metodlarının olduğunu tanımlar.

**Projedeki örnek:** `proto/auth.proto`

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
}
```

**Go karşılığı (auth_grpc.pb.go):**
```go
// Sunucu tarafının implemente etmesi gereken interface
type AuthServiceServer interface {
    Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
    Login(context.Context, *LoginRequest) (*LoginResponse, error)
}

// İstemci tarafının kullanacağı client
type AuthServiceClient interface {
    Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
    Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}
```

---

###  Alan Numaraları (Field Numbers)

**Ne demek?** Her alanın benzersiz bir numarası vardır. Protobuf veriyi bu numaralarla tanır, alan adlarıyla değil.

**Projedeki örnek:** `proto/user.proto`

```protobuf
message User {
  string id = 1;         // numara: 1
  string email = 2;      // numara: 2
  string first_name = 3; // numara: 3
  string last_name = 4;  // numara: 4
  string phone = 5;      // numara: 5
}
```

**Neden numara kullanılır?**
- Alan adını değiştirsen bile numara aynı kalırsa eski istemciler çalışmaya devam eder.
- İkili formatta sadece numaralar taşınır, bu da veriyi küçültür.

**Kural:** 1-15 arası numaralar 1 byte yer kaplar, sık kullanılan alanlara verilmelidir.

---

###  Veri Tipleri

**Ne demek?** Protobuf'un desteklediği temel veri tipleridir.

**Projedeki örnekler:**

| Proto Tipi | Go Tipi | Örnek Kullanım |
|-----------|---------|---------------|
| `string` | `string` | `string email = 1;` |
| `int32` | `int32` | `int32 quantity = 1;` |
| `int64` | `int64` | `int64 price = 2;` |
| `bool` | `bool` | `bool success = 1;` |
| `bytes` | `[]byte` | `bytes photo = 5;` |
| `repeated` | `[]T` (slice) | `repeated string tags = 5;` |

**Projeden gerçek örnek – Order mesajı:**

```protobuf
message OrderItem {
  string product_id = 1;
  string product_name = 2;
  int32 quantity = 3;
  int64 price = 4;
}

message Order {
  string id = 1;
  string user_id = 2;
  repeated OrderItem items = 3;  // ← slice!
  string status = 4;
  int64 total_amount = 5;
}
```

---

###  Serileştirme (Serialization)

**Ne demek?** Veriyi ikili (binary) formata çevirme işlemidir. JSON gibi metin tabanlı değil, ikili formattır.

**JSON vs Protobuf boyut karşılaştırması:**

Aynı veri için:

**JSON (≈150 byte):**
```json
{"id":"abc123","email":"john@example.com","firstName":"John","lastName":"Doe","phone":"555-0100"}
```

**Protobuf (≈50 byte):**
```
\x0a\x06abc123\x12\x10john@example.com\x1a\x04John\x22\x03Doe\x2a\x08555-0100
```

Protobuf yaklaşık **3 kat daha küçük**tür. 1000 istekte bu ciddi fark yaratır.

---

### Geriye Dönük Uyumluluk

**Ne demek?** Proto'ya yeni alan eklediğinde eski istemciler bozulmaz.

**Projedeki örnek:**

```protobuf
// v1 - İlk sürüm
message User {
  string id = 1;
  string email = 2;
}

// v2 - Yeni alan eklendi, eski istemciler hâlâ çalışır
message User {
  string id = 1;
  string email = 2;
  string phone = 3;  // ← yeni alan, eski istemci görmezden gelir
}
```

**Kurallar:**
- Alan numarasını asla değiştirme
- Var olan alanın tipini değiştirme
- Silinen alanın numarasını tekrar kullanma (`reserved` ile işaretle)

---

###  go_package Seçeneği

**Ne demek?** Üretilen Go kodunun hangi pakette olacağını ve import yolunu belirtir.

**Projedeki örnek:** `proto/auth.proto`

```protobuf
option go_package = "GoCommerceX/proto/auth/v1;authv1";
```

**Bu ne demek?**
- `GoCommerceX/proto/auth/v1` → import yolu
- `authv1` → paket adı

**Kullanımı:**
```go
import authv1 "GoCommerceX/proto/auth/v1"
```

---

###  Proto Dosyalarının Projedeki Yeri

**Projedeki yapı:**

```
proto/
├── auth.proto
├── user.proto
├── product.proto
├── inventory.proto
├── cart.proto
├── order.proto
├── payment.proto
├── notification.proto
└── auth/v1/
    ├── auth.pb.go
    └── auth_grpc.pb.go
└── user/v1/
    ├── user.pb.go
    └── user_grpc.pb.go
└── ... (diğer servisler için de aynı)
```

Her servisin `.proto` dosyası ve üretilmiş `.pb.go` dosyaları ayrı ayrıdır.

---

###  Üretim Sürecini Tekrar Etme

**Komut:**
```bash
# auth servisi için
protoc --go_out=. --go-grpc_out=. proto/auth.proto

# user servisi için
protoc --go_out=. --go-grpc_out=. proto/user.proto
```

**Parametreler:**
- `--go_out=.` → Mesaj struct'larını bulunduğun dizine üret
- `--go-grpc_out=.` → Servis interface/client kodlarını bulunduğun dizine üret
- `proto/auth.proto` → Kaynak proto dosyası

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Message Tanımı | `RegisterRequest`, `User` gibi mesajlar |
| Service Tanımı | `AuthService`, `UserService` servisleri |
| Alan Numaraları | `string email = 2;` → 2 numarası |
| Veri Tipleri | `string`, `int32`, `bool`, `repeated` |
| Serileştirme | JSON'dan 3 kat küçük ikili format |
| Geriye Dönük Uyumluluk | Yeni alan eklenince eski istemciler bozulmaz |
| go_package | `option go_package = "GoCommerceX/proto/auth/v1;authv1";` |
| Üretim | `protoc --go_out=. --go-grpc_out=. proto/auth.proto` |

---

# 5.  JWT – Proje Örnekleri

---

##  Önce Büyük Resim

**JWT (JSON Web Token)**, kullanıcı giriş yaptıktan sonra kimliğini kanıtlamak için kullandığı dijital bir anahtardır. Sunucu oturum tutmaz (stateless), her istekte bu token kontrol edilir.

---

###  JWT'nin Yapısı

**Ne demek?** JWT üç parçadan oluşur: Header, Payload, Signature. Nokta ile ayrılır.

**Projedeki örnek (login testinden):**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3ODI0MjM0NjUsImlhdCI6MTc4MjMzNzA2NSwic3ViIjoiZjE5MTAxMWMtYzk3NS00OGVhLWI3NTYtYWVlODQzMzJhYWYzIn0.gdGPpMb_gqqw0pLSONARnGMkv7I_cIs15CdRRfNZ4pc
```

| Parça | İçerik |
|-------|--------|
| **Header** (kırmızı) | `{"alg": "HS256", "typ": "JWT"}` → Hangi algoritma ile imzalandığı |
| **Payload** (mor) | `{"sub": "user-id", "email": "...", "exp": ..., "iat": ...}` → Taşınan bilgiler |
| **Signature** (mavi) | Header + Payload + Secret'ın HMAC-SHA256 ile imzalanmış hali |

---

###  Payload (Claims) – Taşınan Bilgiler

**Ne demek?** Token'ın içinde taşınan kullanıcıya ait bilgilerdir.

**Projedeki örnek:** `auth-service/internal/adapters/jwt_token_generator.go`

```go
claims := jwt.MapClaims{
    "sub":   userID,    // subject = kullanıcı ID'si
    "email": email,     // kullanıcı email'i
    "iat":   time.Now().Unix(),                    // issued at = oluşturulma zamanı
    "exp":   time.Now().Add(24 * time.Hour).Unix(), // expiration = bitiş zamanı (24 saat)
}
```

**Standart Claim'ler:**

| Claim | Açıklama | Projedeki Değeri |
|-------|----------|-----------------|
| `sub` | Subject (kimlik) | Kullanıcı UUID'si |
| `email` | Özel claim | Kullanıcı email'i |
| `iat` | Issued At | Şimdiki zaman |
| `exp` | Expiration | 24 saat sonra |

---

###  Token Oluşturma (Generate)

**Ne demek?** Kullanıcı giriş yaptığında JWT token üretme işlemi.

**Projedeki örnek:** `auth-service/internal/adapters/jwt_token_generator.go`

```go
func (g *JWTTokenGenerator) Generate(ctx context.Context, userID string, email string) (string, error) {
    claims := jwt.MapClaims{
        "sub":   userID,
        "email": email,
        "iat":   time.Now().Unix(),
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(g.secret))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }
    return signedToken, nil
}
```

**Adım adım:**
1. Claim'leri oluştur (kimlik, email, süre)
2. HS256 algoritması ile yeni token oluştur
3. Secret key ile imzala
4. String olarak döndür

---

### Secret Key (Gizli Anahtar)

**Ne demek?** Token'ı imzalamak ve doğrulamak için kullanılan gizli anahtardır. Sadece sunucu bilir.

**Projedeki örnek:** `auth-service/internal/config/config.go`

```go
JWTSecret: getEnv("JWT_SECRET", "change-me-in-production")
```

**Uyarı:** `change-me-in-production` varsayılan değerdir. Gerçek ortamda güçlü bir secret kullanılmalıdır.

---

###  Login Akışında JWT Kullanımı

**Ne demek?** Kullanıcı email/şifre ile giriş yapar, doğruysa token alır.

**Projedeki örnek:** `auth-service/internal/application/login_user.go`

```go
func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
    // 1. Kullanıcıyı bul
    user, err := uc.userRepo.FindByEmail(ctx, email)
    
    // 2. Şifreyi kontrol et
    err = uc.passwordHasher.Compare(ctx, password, user.PasswordHash)
    
    // 3. Token üret
    token, err := uc.tokenGenerator.Generate(ctx, user.ID, user.Email)
    
    return LoginUserOutput{AccessToken: token}, nil
}
```

**Test çıktısı (grpcurl):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

### Stateless Auth (Oturumsuz Kimlik Doğrulama)

**Ne demek?** Sunucu, kimin giriş yaptığını bir yerde saklamaz. Her istekte token gelir, sunucu token'ı doğrular.

**Avantajı:**
- Veritabanında oturum tablosu tutmaya gerek yok
- Sunucu yeniden başlasa da token'lar geçerli kalır
- Yatay ölçekleme kolaydır (her sunucu aynı secret ile doğrulama yapabilir)

---

###  API Gateway'de JWT Doğrulama (Planlanan)

**Ne demek?** API Gateway, gelen HTTP isteklerinde JWT token'ı kontrol edecek.

**Planlanan yapı:**
```go
// API Gateway middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        // "Bearer <token>" formatından token'ı ayır
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Token Yapısı | Header.Payload.Signature |
| Claims | `sub`, `email`, `iat`, `exp` |
| Algoritma | HS256 (HMAC-SHA256) |
| Secret Key | `JWT_SECRET` config değişkeni |
| Token Süresi | 24 saat |
| Login Akışı | Email/şifre kontrolü → token üretimi |
| Stateless | Sunucuda oturum saklanmaz |

---

# 6.  Docker & Docker Compose – Proje Örnekleri

---

##  Önce Büyük Resim

Docker, uygulamaları paketleyip her yerde aynı şekilde çalıştırmaya yarar. Bizim projede PostgreSQL, Redis ve RabbitMQ'yu Docker ile çalıştırıyoruz.

---

###  İmaj (Image)

**Ne demek?** Hazır paket. İçinde uygulama ve tüm bağımlılıkları vardır.

**Projedeki örnek:** `deployments/docker-compose.yml`

```yaml
services:
  postgres:
    image: postgres:16-alpine   # ← Docker Hub'dan indirilen hazır imaj
  redis:
    image: redis:7-alpine        # ← Redis 7, Alpine Linux üzerinde
  rabbitmq:
    image: rabbitmq:3-management-alpine  # ← Yönetim arayüzü dahil
```

---

### Konteyner (Container)

**Ne demek?** İmajdan oluşturulmuş, çalışan örnek.

**Projedeki örnek (docker ps çıktısı):**
```
CONTAINER ID   IMAGE                          NAMES
abc123def456   postgres:16-alpine             gocommerce-postgres
def456ghi789   redis:7-alpine                 gocommerce-redis
ghi789jkl012   rabbitmq:3-management-alpine   gocommerce-rabbitmq
```

---

###  docker-compose.yml – Servis Tanımları

**Projedeki örnek:** `deployments/docker-compose.yml`

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: gocommerce-postgres    # ← konteyner adı
    environment:
      POSTGRES_USER: gocommerce            # ← ortam değişkenleri
      POSTGRES_PASSWORD: gocommerce_password
      POSTGRES_DB: gocommerce
    ports:
      - "5432:5432"                        # ← host:konteyner port eşlemesi
    volumes:
      - postgres_data:/var/lib/postgresql/data  # ← kalıcı veri
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gocommerce -d gocommerce"]
      interval: 10s
      timeout: 5s
      retries: 5
```

---

###  Port Mapping (Port Eşleme)

**Ne demek?** Host bilgisayarın portu ile konteynerin portunu eşleştirir.

**Projedeki örnek:**
```yaml
ports:
  - "5432:5432"   # Host'un 5432'si → Konteynerin 5432'si
  - "6379:6379"   # Host'un 6379'si → Redis
  - "15672:15672" # RabbitMQ yönetim paneli
```

---

###  Volume (Kalıcı Veri)

**Ne demek?** Konteyner silinse bile verinin kaybolmamasını sağlar.

**Projedeki örnek:**
```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data  # PostgreSQL verileri burada saklanır
  - redis_data:/data
  - rabbitmq_data:/var/lib/rabbitmq
```

`docker-compose down -v` yaparsak volume'ler de silinir, veritabanı sıfırlanır.

---

###  Healthcheck (Sağlık Kontrolü)

**Ne demek?** Docker, konteynerin sağlıklı olup olmadığını düzenli olarak kontrol eder.

**Projedeki örnek:**
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U gocommerce -d gocommerce"]
  interval: 10s   # Her 10 saniyede bir kontrol et
  timeout: 5s     # 5 saniye içinde cevap gelmezse başarısız
  retries: 5      # 5 kez başarısız olursa "unhealthy" işaretle
```

---

###  Sık Kullanılan Komutlar

| Komut | Açıklama |
|-------|----------|
| `docker-compose up -d` | Tüm servisleri arka planda başlat |
| `docker-compose down` | Tüm servisleri durdur |
| `docker-compose down -v` | Durdur + volume'leri sil (veritabanı sıfırlanır) |
| `docker ps` | Çalışan konteynerleri listele |
| `docker logs gocommerce-postgres` | PostgreSQL loglarını gör |
| `docker exec -it gocommerce-postgres psql -U gocommerce` | PostgreSQL'e bağlan |

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| İmaj | `postgres:16-alpine`, `redis:7-alpine` |
| Konteyner | `gocommerce-postgres`, `gocommerce-redis` |
| docker-compose.yml | `deployments/docker-compose.yml` |
| Port Mapping | `5432:5432`, `6379:6379` |
| Volume | `postgres_data`, `redis_data` |
| Healthcheck | `pg_isready` ile 10 saniyede bir kontrol |

---

# 7. ️ PostgreSQL – Proje Örnekleri

---

## Önce Büyük Resim

PostgreSQL, projemizdeki çoğu servisin (Auth, User, Product, Order, Payment) veritabanıdır. İlişkisel verileri saklar.

---

###  Bağlantı Havuzu (Connection Pool)

**Ne demek?** Veritabanı bağlantılarını önceden oluşturup havuzda tutar, tekrar tekrar bağlantı açıp kapamaktan kurtarır.

**Projedeki örnek:** `auth-service/cmd/auth-service/main.go`

```go
import "github.com/jackc/pgx/v5/pgxpool"

dbURL := "postgres://gocommerce:gocommerce_password@localhost:5432/gocommerce?sslmode=disable"
pool, err := pgxpool.New(ctx, dbURL)
defer pool.Close()
```

---

###  Migration (Şema Yönetimi)

**Ne demek?** Veritabanı tablolarını SQL dosyalarıyla yönetme.

**Projedeki örnek – Auth Service:** `auth-service/migrations/001_create_users_table.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Projedeki örnek – User Service:** `user-service/migrations/001_create_users_table.sql`

```sql
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

###  Parametreli Sorgu (SQL Injection Koruması)

**Ne demek?** Kullanıcı girdilerini doğrudan SQL'e yapıştırmak yerine `$1, $2` ile parametre olarak verme.

**Projedeki örnek:** `auth-service/internal/adapters/postgres_user_repository.go`

```go
// ✅ DOĞRU
query := `INSERT INTO "users" (id, email, password_hash, created_at, updated_at)
          VALUES ($1, $2, $3, $4, $5)`
pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, ...)

// ❌ YANLIŞ (asla yapma!)
query := fmt.Sprintf("INSERT INTO users VALUES ('%s', '%s')", user.ID, user.Email)
```

---

###  ErrNoRows – Sonuç Bulunamadı

**Ne demek?** Sorgu hiç satır döndürmezse `pgx.ErrNoRows` hatası döner. Bunu `errors.Is` ile kontrol ederiz.

**Projedeki örnek:**
```go
err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, ...)
if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return domain.User{}, ports.ErrUserNotFound  // özel hatamız
    }
    return domain.User{}, err  // gerçek veritabanı hatası
}
```

---

###  RowsAffected – Etkilenen Satır Sayısı

**Ne demek?** UPDATE veya DELETE sorgularında kaç satırın etkilendiğini kontrol etme.

**Projedeki örnek:** `user-service/internal/adapters/postgres_user_repository.go`

```go
result, err := r.pool.Exec(ctx, query, ...)
if result.RowsAffected() == 0 {
    return ports.ErrUserNotFound  // hiç satır etkilenmedi = kullanıcı yok
}
```

---

###  UUID ve Tablo Adlandırma

**Projedeki tablolar:**

| Servis | Tablo | ID Tipi |
|--------|-------|---------|
| Auth Service | `users` | UUID |
| User Service | `user_profiles` | UUID |
| Product Service | `products` | UUID |
| Order Service | `orders` | UUID |
| Payment Service | `payments` | UUID |

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Bağlantı Havuzu | `pgxpool.New(ctx, dbURL)` |
| Migration | `001_create_users_table.sql` |
| Parametreli Sorgu | `$1, $2, $3` |
| ErrNoRows | `errors.Is(err, pgx.ErrNoRows)` |
| RowsAffected | `result.RowsAffected() == 0` |

---

# 8.  Redis – Proje Örnekleri

---

##  Önce Büyük Resim

Redis, anahtar-değer tabanlı, bellek içi (in-memory) veri deposudur. Projemizde sadece **Cart Service** Redis kullanır çünkü sepet verisi geçicidir.

---

###  Key-Value Store

**Ne demek?** Veriler `anahtar → değer` şeklinde saklanır.

**Projedeki örnek:**
```
Key: "cart:user-123"  →  Value: {"items": [{"product_id": "abc", "quantity": 2}]}
```

---

###  In-Memory (Bellek İçi)

**Ne demek?** Veriler RAM'de tutulur, bu yüzden çok hızlıdır. Ama sunucu kapanırsa veriler kaybolabilir.

**Biz neden sepet için Redis kullandık?**
- Sepet geçici veridir, kalıcı olması gerekmez
- Çok hızlı okuma/yazma gerekir (her sayfada sepet kontrolü)
- Kullanıcı oturumu kapatınca sepet silinebilir

---

### Cart Service'te Redis Kullanımı

**Projedeki örnek – Sepet yapısı:**

```go
// Redis'te sepet anahtarı: "cart:<user_id>"
key := fmt.Sprintf("cart:%s", userID)

// Sepete ürün ekleme
redisClient.HSet(ctx, key, productID, quantity)

// Sepeti getirme
items := redisClient.HGetAll(ctx, key)

// Sepetten ürün çıkarma
redisClient.HDel(ctx, key, productID)

// Sepeti silme
redisClient.Del(ctx, key)
```

---

### TTL (Time To Live) – Süreli Veri

**Ne demek?** Veriye "ne kadar süre saklanacağı" bilgisi ekleme.

```go
// Sepeti 24 saat sonra otomatik sil
redisClient.Expire(ctx, key, 24*time.Hour)
```

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Key-Value | `cart:<user_id>` → sepet verisi |
| In-Memory | RAM'de tutulur, çok hızlı |
| HSET/HGET/HDEL | Sepete ekleme/getirme/çıkarma |
| TTL | 24 saat sonra otomatik silme |

---

# 9.  RabbitMQ – Proje Örnekleri

---

##  Önce Büyük Resim

RabbitMQ, servisler arası **asenkron** mesajlaşmayı sağlayan bir mesaj kuyruğudur. Bir servis mesaj gönderir, başka bir servis onu işler.

---

###  Producer (Üretici) & Consumer (Tüketici)

**Projedeki örnek:**

| Rol | Servis | Ne yapar? |
|-----|--------|-----------|
| **Producer** | Order Service | Sipariş oluşunca `OrderCreated` event'i gönderir |
| **Consumer** | Notification Service | `OrderCreated` event'ini dinler, bildirim oluşturur |

```text
Order Service ──(OrderCreated event)──> RabbitMQ ──> Notification Service
```

---

###  Event-Driven Mimaride Akış

**Senaryo: Sipariş oluşturma**

1. Order Service siparişi veritabanına kaydeder
2. `OrderCreated` event'ini RabbitMQ'ya gönderir
3. Notification Service bu event'i alır
4. Kullanıcıya "Siparişiniz alındı" bildirimi oluşturur

**Neden asenkron?** Order Service'in, Notification Service'in işini bitirmesini beklemesi gerekmez. Sipariş hemen onaylanır, bildirim arkada işlenir.

---

###  Queue ve Exchange

**Ne demek?**
- **Exchange:** Mesajları alır, kurallara göre kuyruklara yönlendirir
- **Queue:** Mesajların sırayla beklediği yerdir

---

###  Dayanıklılık (Durability)

**Ne demek?** RabbitMQ çökse bile mesajlar kaybolmaz, tekrar başlayınca işlenir.

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Producer | Order Service → `OrderCreated` event |
| Consumer | Notification Service → bildirim oluşturma |
| Asenkron | Order işlemi bildirimi beklemez |
| Dayanıklılık | Mesajlar kaybolmaz |

---

# 10.  API Gateway – Proje Örnekleri

---

## Önce Büyük Resim

API Gateway, dış dünyadan gelen tüm istekleri karşılayan **tek giriş noktasıdır**. REST isteklerini alır, gRPC'ye çevirir, ilgili servise iletir.

---

### Tek Giriş Noktası

```text
Tarayıcı/Mobil ──HTTP──> API Gateway (:8080) ──gRPC──> Auth Service (:50051)
                                │──────────gRPC──> User Service (:50052)
                                └──────────gRPC──> Product Service (:50053)
```

---

### Protokol Dönüşümü (HTTP → gRPC)

**Ne demek?** Dış dünya HTTP/REST konuşur, servislerimiz gRPC konuşur. API Gateway ikisi arasında çeviri yapar.

**Projedeki örnek (şu anki durum):** `api-gateway/cmd/api-gateway/main.go`

```go
// Şu an sadece health check var
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ok"))
})
```

**Olması gereken (planlanan):**
```go
// Register endpoint'i → Auth Service'e gRPC çağrısı
mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
    // HTTP isteğini gRPC'ye çevir
    resp, _ := authClient.Register(ctx, &authv1.RegisterRequest{
        Email:    r.FormValue("email"),
        Password: r.FormValue("password"),
    })
    // gRPC cevabını HTTP'ye çevir
    json.NewEncoder(w).Encode(resp)
})
```

---

### JWT Middleware (Planlanan)

**Ne demek?** Korumalı endpoint'lere gelen isteklerde JWT token kontrolü yapan katman.

```go
// Korunan endpoint'ler için JWT kontrolü
mux.Handle("/api/user/profile", authMiddleware(http.HandlerFunc(getProfile)))
```

---

##  Özet Tablo

| Alt Başlık | Projedeki Karşılığı |
|------------|---------------------|
| Tek Giriş Noktası | `:8080` portu |
| Protokol Dönüşümü | HTTP isteği → gRPC çağrısı |
| Mevcut Durum | Sadece `/health` endpoint'i |
| JWT Middleware | Korunan endpoint'ler için token kontrolü (planlandı) |

---

Bu, projendeki **10 ana konunun tamamının** alt başlıklarıyla birlikte proje örnekleridir. İstersen bir konuyu derinleştirelim, istersen API Gateway'i tamamlamaya geçelim.

---

# 11.  Hata Ayıklama ve Test Disiplini – Projeden Çıkan Gerçek Dersler

Bu bölüm, teoriden çok pratiği anlatır. Çünkü projeyi yazarken en çok öğrenilen şeyler genelde "çalışmayan yerler" sayesinde ortaya çıkar.

---

###  Önce Servisleri Durumlandırma

Bir problemi çözerken ilk iş, hangi servisin ayakta olduğunu ve hangisinin yanlış davrandığını anlamaktır.

**Uyguladığımız sıra:**
1. Tüm servisleri durdur
2. Tek tek yeniden ayağa kaldır
3. Önce health check
4. Sonra auth
5. Ardından ürün, sepet ve sipariş akışı

Bu yaklaşımın faydası şu:
- Hata hangi serviste hemen belli olur
- Loglar daha okunur hale gelir
- Bir değişiklik başka yeri bozduysa daha kolay yakalanır

---

###  Doğru Çalıştırma Komutu

Go projelerinde bazen sorun kodda değil, komutu yanlış dizinden çalıştırmaktır.

**Yaşadığımız örnek:**
```bash
go run ./api-gateway/cmd/api-gateway
```

Bu komut doğru klasörden çalıştırılmadığında yanlış path oluşabilir ve şu tarz hatalar görülebilir:
- `directory not found`
- `missing port in address`

Buradan çıkan ders:
- Her zaman önce bulunduğun dizini kontrol et
- Modül kökünden çalıştır
- Config dosyasında port ve adres formatlarını doğrula

---

###  Port ve Address Kontrolü

Bir servis "çalışıyor gibi" görünse bile yanlış port ile ayağa kalkmış olabilir.

**Kontrol etmen gerekenler:**
- `HTTP_PORT`
- `GRPC_PORT`
- `AUTH_ADDR`
- `USER_ADDR`
- `PRODUCT_ADDR`
- `CART_ADDR`
- `ORDER_ADDR`
- `PAYMENT_ADDR`
- `NOTIF_ADDR`

Eğer API Gateway ayakta değilse önce:
```bash
curl http://localhost:8080/health
```

Eğer gRPC servisleri görünmüyorsa:
```bash
grpcurl -plaintext localhost:50051 list
```

Bu iki komut, sorunun HTTP tarafında mı gRPC tarafında mı olduğunu hızlıca ayırır.

---

###  Field Adı Uyuşmazlığı

Bu projede en öğretici hatalardan biri, aynı anlamı taşıyan alan adlarının farklı katmanlarda farklı görünmesiydi.

**Örnek problem:**
- API Gateway `user_id` bekliyor
- Cart Service ya da proto tarafı `userId` / `UserId` / `user_id` gibi farklı isimler kullanıyor
- Sonuç olarak request doğru gibi görünse bile veri yanlış yere gidiyor

**Çözüm mantığı:**
- Dış istekten gelen alanları Gateway'de normalize et
- Kullanıcı kimliğini body'den almak yerine JWT'den çıkar
- Proto mesaj adlarını, Go struct alanlarını ve JSON field isimlerini ayrı ayrı düşün

**Bu projedeki doğru yaklaşım:**
- `user_id` body'den gelmesin
- Gateway JWT'den `sub` claim'ini alsın
- Cart request'i oluştururken `UserId` alanını gateway doldursun

Bu, hem güvenlik hem de veri bütünlüğü açısından daha doğru.

---

###  UUID ve Boş String Problemi

Sipariş akışında çok klasik ama çok can sıkıcı bir problem yaşanabilir:

- Veritabanında alan `UUID`
- Uygulama boş değer gönderiyor: `""`
- PostgreSQL bunu UUID'e çeviremiyor

**Belirti:**
```text
invalid input syntax for type uuid: ""
```

**Bu hata ne öğretir?**
- `""` ile `NULL` aynı şey değildir
- UUID alanına boş string gönderilmez
- SQL tarafında nullable alanlar açıkça modellemelidir

**Doğru yaklaşım:**
- Go tarafında `sql.NullString` ya da nullable yapıyı kullan
- Insert sırasında gerçekten değer yoksa `NULL` gönder
- Gereksiz boş string üretme

Bu proje bize şunu öğretti: veritabanı şeması ile Go struct'ı birebir aynı düşünülmez, arada dönüşüm katmanı gerekir.

---

###  Test Sırası

Test ederken her şeyi aynı anda denemek yerine sırayla gitmek çok daha güvenlidir.

**İyi test sırası:**
1. `health`
2. `register`
3. `login`
4. JWT ile `profile`
5. JWT ile `product list`
6. `add to cart`
7. `get cart`
8. `create order`
9. `get order`

Bu sırayı bozarsan, alttaki hata üstteki adımdan mı geliyor anlamak zorlaşır.

---

###  Sorun Çıktığında Bakılacak Yerler

Bir işlem bozulduğunda şu sırayla kontrol etmek çok iş gördü:

1. **Config**
   - Env var doğru mu?
   - Port doğru mu?
   - Secret doğru mu?

2. **Proto**
   - Field adı doğru mu?
   - Request/response alanı eşleşiyor mu?
   - Generate edilmiş kod güncel mi?

3. **Handler**
   - HTTP veya gRPC request doğru parse ediliyor mu?
   - JWT context doğru okunuyor mu?

4. **Use case**
   - İş kuralı doğru mu?
   - Validation doğru çalışıyor mu?

5. **Repository**
   - SQL sorgusu doğru mu?
   - UUID/nullable tipler doğru mu?

6. **Service bağımlılıkları**
   - Diğer servis açık mı?
   - Address doğru mu?

Bu kontrol listesi, hata ayıklamayı tahmin oyunundan çıkarıyor.

---

###  Projeden Öğrenilen En Önemli Alışkanlıklar

Bu projede teknik olarak çok şey öğrenildi ama birkaç alışkanlık daha değerli:

- Önce küçük parçayı çalıştır
- Sonra entegrasyonu ekle
- Hata mesajını dikkatle oku
- Proto ve DB şemasını ayrı ayrı doğrula
- Testi uçtan uca yapmadan "tamamdır" deme

---

###  Gerçek Hayat İçin Kısa Notlar

Bu proje şunu net gösterdi:

- Mikroservis kolay değil, ama düzenli olunca yönetilebilir
- gRPC hızlıdır, fakat sözleşme disiplini ister
- JWT güvenlik için güçlüdür, ama doğru middleware olmadan eksik kalır
- PostgreSQL ve Redis farklı amaçlar içindir, ikisini karıştırmamak gerekir
- API Gateway, dış dünya ile iç servisler arasında temiz bir sınır oluşturur

---

##  Final Öğrenme Özeti

Bu projeyi yazarken öğrendiğim ana şey şu:

> Her katman kendi işini yaparsa sistem büyürken bozulmaz.

Domain iş kuralını bilir, application akışı bilir, ports neye ihtiyaç olduğunu söyler, adapters dış dünyayla konuşur, transport isteği alır, gateway ise dış dünya ile servisler arasında çeviri yapar.

Bir hata çıktığında bu katmanlardan hangisinin kuralı bozduğunu bulmak, problemi çözmenin en doğru yoludur.

---

##  İleride Bakılacak Konular

Bu projeyi daha ileri taşımak için öğrenmeye devam edilebilecek başlıklar:

- gRPC interceptor yazımı
- gRPC client retry ve timeout stratejileri
- Redis cache invalidation
- RabbitMQ dead-letter queue
- Docker Compose ile servisleri tek komutta ayağa kaldırma
- CI/CD pipeline kurulumu
- OpenTelemetry ile distributed tracing
- SQL transaction ve saga desenleri

Bu not, sadece ne yaptığını değil, neden yaptığını da hatırlatmak için burada dursun.
