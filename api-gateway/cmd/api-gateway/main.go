package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"GoCommerceX/api-gateway/internal/config"
	"GoCommerceX/api-gateway/internal/middleware"

	authv1 "GoCommerceX/proto/auth/v1"
	cartv1 "GoCommerceX/proto/cart/v1"
	orderv1 "GoCommerceX/proto/order/v1"
	productv1 "GoCommerceX/proto/product/v1"
	userv1 "GoCommerceX/proto/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	authClient    authv1.AuthServiceClient
	userClient    userv1.UserServiceClient
	productClient productv1.ProductServiceClient
	cartClient    cartv1.CartServiceClient
	orderClient   orderv1.OrderServiceClient
}

func main() {
	cfg := config.Load()

	authConn, _ := grpc.Dial(cfg.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	userConn, _ := grpc.Dial(cfg.UserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	productConn, _ := grpc.Dial(cfg.ProductAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cartConn, _ := grpc.Dial(cfg.CartAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	orderConn, _ := grpc.Dial(cfg.OrderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	defer authConn.Close()
	defer userConn.Close()
	defer productConn.Close()
	defer cartConn.Close()
	defer orderConn.Close()

	gw := &Gateway{
		authClient:    authv1.NewAuthServiceClient(authConn),
		userClient:    userv1.NewUserServiceClient(userConn),
		productClient: productv1.NewProductServiceClient(productConn),
		cartClient:    cartv1.NewCartServiceClient(cartConn),
		orderClient:   orderv1.NewOrderServiceClient(orderConn),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Açık endpoint'ler (JWT gerektirmez)
	mux.HandleFunc("/api/auth/register", gw.handleRegister)
	mux.HandleFunc("/api/auth/login", gw.handleLogin)

	// Korumalı endpoint'ler
	protected := http.NewServeMux()
	protected.HandleFunc("/api/user/profile", gw.handleGetProfile)
	protected.HandleFunc("/api/products/", gw.handleGetProduct) // ID ile
	protected.HandleFunc("/api/cart", gw.handleCart)
	protected.HandleFunc("/api/orders/", gw.handleGetOrder) // ID ile
	protected.HandleFunc("/api/orders/create", gw.handleCreateOrder)

	authMW := middleware.AuthMiddleware(cfg.JWTSecret)
	mux.Handle("/api/user/", authMW(protected))
	mux.Handle("/api/products/", authMW(protected))
	mux.Handle("/api/cart", authMW(protected))
	mux.Handle("/api/orders/", authMW(protected))
	mux.Handle("/api/orders/create", authMW(protected))

	handler := corsMiddleware(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("API Gateway listening on :%s\n", cfg.HTTPPort)
	log.Printf("  Auth: %s", cfg.AuthAddr)
	log.Printf("  User: %s", cfg.UserAddr)
	log.Printf("  Product: %s", cfg.ProductAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (gw *Gateway) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// Genişletilmiş register isteği
	var body struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}
	raw, _ := io.ReadAll(r.Body)
	json.Unmarshal(raw, &body)

	// Auth Service'e register isteği (first_name, last_name ile)
	authResp, err := gw.authClient.Register(r.Context(), &authv1.RegisterRequest{
		Email:     body.Email,
		Password:  body.Password,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Phone:     body.Phone,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, authResp)
}

func (gw *Gateway) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req authv1.LoginRequest
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	resp, err := gw.authClient.Login(r.Context(), &req)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// ---------- User Endpoints ----------

func (gw *Gateway) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	resp, err := gw.userClient.GetUser(r.Context(), &userv1.GetUserRequest{Id: userID})
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// ---------- Product Endpoints ----------

func (gw *Gateway) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// URL'den ID'yi al: /api/products/{id}
	id := r.URL.Path[len("/api/products/"):]
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "product id required"})
		return
	}

	resp, err := gw.productClient.GetProduct(r.Context(), &productv1.GetProductRequest{Id: id})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// ---------- Cart Endpoints ----------

func (gw *Gateway) handleCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	switch r.Method {
	case http.MethodGet:
		resp, err := gw.cartClient.GetCart(r.Context(), &cartv1.GetCartRequest{UserId: userID})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, resp)

	case http.MethodPost:
		var payload struct {
			ProductID string `json:"product_id"`
			Quantity  int32  `json:"quantity"`
		}
		raw, _ := io.ReadAll(r.Body)
		json.Unmarshal(raw, &payload)

		if payload.ProductID == "" || payload.Quantity <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "product_id and quantity are required"})
			return
		}

		req := cartv1.AddItemRequest{
			UserId:    userID,
			ProductId: payload.ProductID,
			Quantity:  payload.Quantity,
		}

		resp, err := gw.cartClient.AddItem(r.Context(), &req)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, resp)

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// ---------- Order Endpoints ----------

func (gw *Gateway) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	// URL'den ID'yi al: /api/orders/{id}
	id := r.URL.Path[len("/api/orders/"):]
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order id required"})
		return
	}

	resp, err := gw.orderClient.GetOrder(r.Context(), &orderv1.GetOrderRequest{Id: id})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Kullanıcı sadece kendi siparişini görebilir
	if resp.Order.UserId != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (gw *Gateway) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req orderv1.CreateOrderRequest
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)
	req.UserId = userID

	resp, err := gw.orderClient.CreateOrder(r.Context(), &req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// ---------- Yardımcı ----------

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func getContext(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
