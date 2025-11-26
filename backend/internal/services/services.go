package services

// خدمات أساسية (سيتم تطويرها لاحقاً)
type AdminService struct{}
type UserService struct{}
type AuthService struct{}
type StoreService struct{}
type CartService struct{}
type PaymentService struct{}
type AIService struct{}
type EmailService struct{}
type UploadService struct{}

func NewAdminService() *AdminService   { return &AdminService{} }
func NewUserService() *UserService     { return &UserService{} }
func NewAuthService() *AuthService     { return &AuthService{} }
func NewStoreService() *StoreService   { return &StoreService{} }
func NewCartService() *CartService     { return &CartService{} }
func NewPaymentService() *PaymentService { return &PaymentService{} }
func NewAIService() *AIService         { return &AIService{} }
func NewEmailService() *EmailService   { return &EmailService{} }
func NewUploadService() *UploadService { return &UploadService{} }