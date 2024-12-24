package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"wallet/internal/service"

	jwt "github.com/dgrijalva/jwt-go"
)

// As the application is handling money values in the wallet, it is important to keep the precision of the values.
// The floatConversion variable is used to convert the float values to integers by multiplying them by 10000.
// This way, the application can handle the values with the required precision.
var floatConversion = 10000

type WalletHandler struct {
	service service.WalletServiceInterface
}

func NewWalletHandler(service service.WalletServiceInterface) *WalletHandler {
	return &WalletHandler{service: service}
}

// WalletChangeRequest is a struct to represent the request payload for deposit and withdraw operations.
type WalletChangeRequest struct {
	Currency string  `json:"currency" `
	Amount   float32 `json:"amount"`
}

// WalletChangeResponse is a struct to represent the response payload for deposit and withdraw operations.
type WalletChangeResponse struct {
	Messsage    string             `json:"message" `
	New_balance map[string]float32 `json:"new_balances"`
}

// ExchangeRateRequest is a struct to represent the request payload for getting the exchange rate between two currencies.
type ExchangeRateRequest struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
}

// ExchangeRequest is a struct to represent the request payload for exchanging money between two currencies.
type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float32 `json:"amount"`
}

// RegisterUserRequest is a struct to represent the request payload for registering a new user.
type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pw       string `json:"pw"`
}

// RegisterUserResponse is a struct to represent the response payload for registering a new user.
type RegisterUserResponse struct {
	Message string `json:"message"`
}

// LoginRequest is a struct to represent the request payload for logging in a user.
type LoginRequest struct {
	Username string `json:"username"`
	Pw       string `json:"pw"`
}

// WalletDeposit is an HTTP handler to deposit money into the wallet.
func (h *WalletHandler) WalletDeposit(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header from the request.
	auth := r.Header.Get(
		"Authorization",
	)

	// Verify the token and extract the user ID and username from the claims.
	uid, username, err := verifyTokenWithClaims(auth)
	if uid == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req WalletChangeRequest
	var res WalletChangeResponse
	// Decode the request payload into the req variable.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Deposit the amount into the wallet.
	err = h.service.Deposit(r.Context(), uid, floatToIntConversion(req.Amount), req.Currency)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated balance after the deposit operation.
	balances, err := h.service.GetBalance(r.Context(), username)
	balancesFloat32 := make(map[string]float32)
	balancesFloat32 = intMapToFloatMapConversion(balances)

	res.Messsage = "Deposit successful"
	res.New_balance = balancesFloat32
	json.NewEncoder(w).Encode(res)
}

// WalletWithdraw is an HTTP handler to withdraw money from the wallet.
func (h *WalletHandler) WalletWithdraw(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header from the request.
	auth := r.Header.Get(
		"Authorization",
	)

	// Verify the token and extract the user ID and username from the claims.
	uid, username, err := verifyTokenWithClaims(auth)
	var req WalletChangeRequest
	var res WalletChangeResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Withdraw the amount from the wallet.
	err = h.service.Withdraw(r.Context(), uid, floatToIntConversion(req.Amount), req.Currency)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated balance after the withdraw operation.
	balances, err := h.service.GetBalance(r.Context(), username)
	res.Messsage = "Withdraw successful"

	balancesFloat32 := make(map[string]float32)
	balancesFloat32 = intMapToFloatMapConversion(balances)

	res.New_balance = balancesFloat32
	json.NewEncoder(w).Encode(res)
}

// GetBalance is an HTTP handler to get the balance of the wallet.
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Extract the Authorization header from the request.
	auth := r.Header.Get(
		"Authorization",
	)
	// Verify the token and extract the user ID and username from the claims.
	_, username, err := verifyTokenWithClaims(auth)
	if err != nil {
		http.Error(w, "Invalid wallet ID", http.StatusBadRequest)
		return
	}

	// Get the balance of the wallet.
	balances, err := h.service.GetBalance(r.Context(), username)
	balancesFloat32 := make(map[string]float32)
	balancesFloat32 = intMapToFloatMapConversion(balances)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(balancesFloat32)
}

// GetExchangeRates is an HTTP handler to get the exchange rates between different currencies.
func (h *WalletHandler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get(
		"Authorization",
	)
	err := verifyToken(auth)
	log.Println(err)
	if err != nil {
		http.Error(w, "Authorization invalid", http.StatusUnauthorized)
		return
	}
	rates, err := h.service.GetExchangeRates(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rates)
}

// GetExchangeRate is an HTTP handler to get the exchange rate between two currencies.
func (h *WalletHandler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	var req ExchangeRateRequest
	auth := r.Header.Get(
		"Authorization",
	)
	err := verifyToken(auth)
	if err != nil {
		http.Error(w, "Authorization invalid", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rates, err := h.service.GetExchangeRate(r.Context(), req.FromCurrency, req.ToCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rates)
}

// Exchange is an HTTP handler to exchange money between two currencies.
func (h *WalletHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	var req ExchangeRequest
	auth := r.Header.Get(
		"Authorization",
	)
	uid, username, err := verifyTokenWithClaims(auth)
	if err != nil {
		http.Error(w, "Authorization invalid", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get the exchange rate between the two currencies.
	rate, err := h.service.GetExchangeRate(r.Context(), req.FromCurrency, req.ToCurrency)
	// Get the balance of the wallet.
	balance, err := h.service.GetBalance(r.Context(), username)

	// Convert the balance from int to float.
	balanceFloat32 := map[string]float32{}
	balanceFloat32 = intMapToFloatMapConversion(balance)
	// subtract the amount from the fromCurrency balance and add the amount to the toCurrency balance.
	balanceFloat32[req.FromCurrency] -= req.Amount
	balanceFloat32[req.ToCurrency] += req.Amount * rate
	// Update the balance in the wallet.
	err = h.service.Withdraw(r.Context(), uid, floatToIntConversion(req.Amount), req.FromCurrency)
	err = h.service.Deposit(r.Context(), uid, floatToIntConversion(req.Amount*rate), req.ToCurrency)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(balanceFloat32)
}

// RegisterUser is an HTTP handler to register a new user.
func (h *WalletHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	var res RegisterUserResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.service.RegisterUser(r.Context(), req.Username, req.Email, req.Pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Message = "User registered successfully"
	json.NewEncoder(w).Encode(res)
}

// Login is an HTTP handler to log in a user.
func (h *WalletHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := h.service.Login(r.Context(), req.Username, req.Pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(token)
}

// verifyToken is a helper function to verify the token in the Authorization header.
func verifyToken(auth string) error {
	if auth == "" {
		return errors.New("Authorization header is required")
	}

	parts := strings.Split(auth, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("Authorization header format must be Bearer {token}")
	}

	tokenString := parts[1]

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	if token.Valid {
		return nil
	} else {
		return errors.New("Invalid token")
	}
}

// verifyTokenWithClaims is a helper function to verify the token in the Authorization header and extract the claims.
func verifyTokenWithClaims(auth string) (int32, string, error) {
	if auth == "" {
		return 0, "", errors.New("Authorization header is required")
	}

	parts := strings.Split(auth, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, "", errors.New("Authorization header format must be Bearer {token}")
	}

	tokenString := parts[1]

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	if !token.Valid {
		return 0, "", errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["uid"].(float64)
		if !ok {
			return 0, "", errors.New("uid not found in token")
		}
		username, ok := claims["username"].(string)
		if !ok {
			return 0, "", errors.New("username not found in token")
		}
		return int32(uid), username, nil
	}

	return 0, "", errors.New("Invalid token")
}

// intToFloatConversion is a helper function to convert an integer value to a float value.
func intToFloatConversion(value int32) float32 {
	return float32(value) / float32(floatConversion)
}

// floatToIntConversion is a helper function to convert a float value to an integer value.
func floatToIntConversion(value float32) int32 {
	return int32(value * float32(floatConversion))
}

// intMapToFloatMapConversion is a helper function to convert a map of integer values to a map of float values.
func intMapToFloatMapConversion(m map[string]int32) map[string]float32 {
	mf := make(map[string]float32)
	for k, v := range m {
		mf[k] = intToFloatConversion(v)
	}
	return mf
}

// floatMapToIntMapConversion is a helper function to convert a map of float values to a map of integer values.
func floatMapToIntMapConversion(m map[string]float32) map[string]int32 {
	mf := make(map[string]int32)
	for k, v := range m {
		mf[k] = floatToIntConversion(v)
	}
	return mf
}
