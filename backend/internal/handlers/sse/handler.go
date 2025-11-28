package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/models"
)

// SSEManager مدير اتصالات SSE
type SSEManager struct {
	clients    map[string]map[string]chan []byte
	mutex      sync.RWMutex
	broadcast  chan BroadcastMessage
	register   chan Client
	unregister chan Client
}

// Client عميل SSE
type Client struct {
	ID       string
	UserID   string
	Channels []string
	Messages chan []byte
}

// BroadcastMessage رسالة بث
type BroadcastMessage struct {
	Channels []string
	UserIDs  []string
	Data     interface{}
	Event    string
}

// Event هيكل الحدث
type Event struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Channel string      `json:"channel,omitempty"`
	Time    time.Time   `json:"time"`
}

var (
	manager *SSEManager
	once    sync.Once
)

// NewSSEManager إنشاء مدير SSE جديد
func NewSSEManager() *SSEManager {
	once.Do(func() {
		manager = &SSEManager{
			clients:    make(map[string]map[string]chan []byte),
			broadcast:  make(chan BroadcastMessage, 100),
			register:   make(chan Client, 100),
			unregister: make(chan Client, 100),
		}
		go manager.run()
	})
	return manager
}

// GetManager الحصول على مدير SSE
func GetManager() *SSEManager {
	if manager == nil {
		return NewSSEManager()
	}
	return manager
}

// تشغيل مدير SSE
func (m *SSEManager) run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			for _, channel := range client.Channels {
				if m.clients[channel] == nil {
					m.clients[channel] = make(map[string]chan []byte)
				}
				m.clients[channel][client.ID] = client.Messages
			}
			m.mutex.Unlock()

		case client := <-m.unregister:
			m.mutex.Lock()
			for _, channel := range client.Channels {
				if channelClients, exists := m.clients[channel]; exists {
					delete(channelClients, client.ID)
					if len(channelClients) == 0 {
						delete(m.clients, channel)
					}
				}
			}
			close(client.Messages)
			m.mutex.Unlock()

		case message := <-m.broadcast:
			m.broadcastMessage(message)
		}
	}
}

// بث الرسالة
func (m *SSEManager) broadcastMessage(msg BroadcastMessage) {
	jsonData, err := json.Marshal(Event{
		Type:    msg.Event,
		Data:    msg.Data,
		Time:    time.Now(),
		Channel: msg.Event,
	})
	if err != nil {
		slog.Error("فشل في ترميز رسالة البث", "error", err)
		return
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// البث للقنوات المحددة
	for _, channel := range msg.Channels {
		if channelClients, exists := m.clients[channel]; exists {
			for clientID, messageChan := range channelClients {
				// إذا تم تحديد مستخدمين معينين، تحقق من أن العميل هو أحدهم
				if len(msg.UserIDs) > 0 {
					clientUserID := extractUserID(clientID)
					if !contains(msg.UserIDs, clientUserID) {
						continue
					}
				}

				select {
				case messageChan <- jsonData:
				default:
					// تجنب الانسداد - تجاهل إذا كانت القناة ممتلئة
					slog.Warn("قناة العميل ممتلئة، تجاهل الرسالة", "client_id", clientID)
				}
			}
		}
	}
}

// Broadcast بث رسالة للقنوات
func (m *SSEManager) Broadcast(channels []string, userIDs []string, data interface{}, eventType string) {
	message := BroadcastMessage{
		Channels: channels,
		UserIDs:  userIDs,
		Data:     data,
		Event:    eventType,
	}

	select {
	case m.broadcast <- message:
	default:
		slog.Warn("قناة البث ممتلئة، تجاهل الرسالة")
	}
}

// Handler معالج SSE الرئيسي
func Handler(c *gin.Context) {
	// التحقق من دعم الـ SSE
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "التدفق غير مدعوم",
			"success": false,
		})
		return
	}

	// إعداد الاتصال لـ SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=UTF-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	// إرسال الرؤوس إلى العميل
	c.Writer.WriteHeader(http.StatusOK)
	flusher.Flush()

	// الحصول على معرف المستخدم من السياق (مباشرة من gin)
	userID := ""
	if userIDVal, exists := c.Get("userID"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	// الحصول على القنوات المطلوبة
	channels := c.QueryArray("channels")
	if len(channels) == 0 {
		channels = []string{"notifications", "updates"}
	}

	// إنشاء عميل جديد
	client := Client{
		ID:       generateClientID(userID),
		UserID:   userID,
		Channels: channels,
		Messages: make(chan []byte, 10),
	}

	// تسجيل العميل
	manager := GetManager()
	manager.register <- client
	defer func() {
		manager.unregister <- client
	}()

	// إعداد المسجل
	requestID := getRequestID(c.Request.Context())
	eventLogger := logger.Stdout.With(slog.String("request_id", requestID))

	// تسجيل طلب SSE
	eventLogger.Info("عميل SSE متصل", 
		"user_id", userID, 
		"client_id", client.ID, 
		"channels", channels)

	// إرسال حدث الاتصال
	sendEvent(c.Writer, flusher, Event{
		Type: "connected",
		Data: gin.H{
			"message":   "تم الاتصال بنجاح",
			"client_id": client.ID,
			"channels":  channels,
			"timestamp": time.Now(),
		},
		Time: time.Now(),
	})

	// حلقة الاستماع للرسائل
	for {
		select {
		case message, ok := <-client.Messages:
			if !ok {
				eventLogger.Info("قناة العميل مغلقة", "client_id", client.ID)
				return
			}

			// إرسال الرسالة إلى العميل
			fmt.Fprintf(c.Writer, "data: %s\n\n", string(message))
			flusher.Flush()

			eventLogger.Debug("تم إرسال حدث SSE", "client_id", client.ID)

		case <-c.Request.Context().Done():
			eventLogger.Info("عميل SSE انقطع", 
				"client_id", client.ID, 
				"user_id", userID)
			return

		case <-time.After(30 * time.Second):
			// إرسال نبضة قلب للحفاظ على الاتصال
			sendEvent(c.Writer, flusher, Event{
				Type: "heartbeat",
				Data: gin.H{
					"timestamp": time.Now(),
				},
				Time: time.Now(),
			})
		}
	}
}

// NotificationHandler معالج الإشعارات عبر SSE
func NotificationHandler(c *gin.Context) {
	// هذا يمكن دمجه مع المعالج الرئيسي أو استخدامه بشكل منفصل
	Handler(c)
}

// AdminHandler معالج SSE للمسؤولين
func AdminHandler(c *gin.Context) {
	// التحقق من صلاحيات المسؤول
	userRole := ""
	if roleVal, exists := c.Get("userRole"); exists {
		if role, ok := roleVal.(string); ok {
			userRole = role
		}
	}
	
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "غير مصرح",
			"success": false,
		})
		return
	}

	// إضافة قنوات المسؤول
	channels := c.QueryArray("channels")
	channels = append(channels, "admin", "system", "monitoring")
	c.Request.URL.RawQuery = "channels="
	for i, channel := range channels {
		if i > 0 {
			c.Request.URL.RawQuery += "&channels="
		}
		c.Request.URL.RawQuery += channel
	}
	
	Handler(c)
}

// إرسال حدث فردي
func sendEvent(w http.ResponseWriter, flusher http.Flusher, event Event) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		slog.Error("فشل في ترميز الحدث", "error", err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
	flusher.Flush()
}

// الدوال المساعدة

func generateClientID(userID string) string {
	if userID == "" {
		return fmt.Sprintf("anonymous_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%s_%d", userID, time.Now().UnixNano())
}

func extractUserID(clientID string) string {
	// استخراج معرف المستخدم من معرف العميل
	// التنسيق: userID_timestamp
	for i := len(clientID) - 1; i >= 0; i-- {
		if clientID[i] == '_' {
			return clientID[:i]
		}
	}
	return clientID
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// ================================
// دوال البث الجاهزة
// ================================

// BroadcastNotification بث إشعار
func BroadcastNotification(notification models.Notification, userIDs []string) {
	manager := GetManager()
	manager.Broadcast([]string{"notifications"}, userIDs, notification, "notification")
}

// BroadcastSystemAlert بث تنبيه نظام
func BroadcastSystemAlert(alert SystemAlert, channels []string) {
	manager := GetManager()
	if len(channels) == 0 {
		channels = []string{"system", "admin"}
	}
	manager.Broadcast(channels, []string{}, alert, "system_alert")
}

// BroadcastOrderUpdate بث تحديث طلب
func BroadcastOrderUpdate(order models.Order, userID string) {
	manager := GetManager()
	manager.Broadcast([]string{"orders"}, []string{userID}, order, "order_update")
}

// BroadcastServiceUpdate بث تحديث خدمة
func BroadcastServiceUpdate(service models.Service, channels []string) {
	manager := GetManager()
	if len(channels) == 0 {
		channels = []string{"services", "updates"}
	}
	manager.Broadcast(channels, []string{}, service, "service_update")
}

// BroadcastAdminStats بث إحصائيات للمسؤولين
func BroadcastAdminStats(stats DashboardStats) {
	manager := GetManager()
	manager.Broadcast([]string{"admin", "stats"}, []string{}, stats, "admin_stats")
}

// BroadcastUserStatus بث حالة المستخدم
func BroadcastUserStatus(userID string, status string) {
	manager := GetManager()
	manager.Broadcast([]string{"users", "status"}, []string{userID}, gin.H{
		"user_id": userID,
		"status":  status,
		"timestamp": time.Now(),
	}, "user_status")
}

// BroadcastPaymentUpdate بث تحديث دفع
func BroadcastPaymentUpdate(payment models.Payment, userID string) {
	manager := GetManager()
	manager.Broadcast([]string{"payments"}, []string{userID}, payment, "payment_update")
}

// BroadcastCartUpdate بث تحديث سلة التسوق
func BroadcastCartUpdate(cart models.Cart, userID string) {
	manager := GetManager()
	manager.Broadcast([]string{"cart"}, []string{userID}, cart, "cart_update")
}

// ================================
// هياكل البيانات الإضافية
// ================================

// SystemAlert تنبيه النظام
type SystemAlert struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // error, warning, info
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"` // low, medium, high, critical
	Timestamp time.Time `json:"timestamp"`
	Data      gin.H     `json:"data,omitempty"`
}

// DashboardStats إحصائيات لوحة التحكم
type DashboardStats struct {
	TotalUsers     int       `json:"total_users"`
	TotalOrders    int       `json:"total_orders"`
	TotalRevenue   float64   `json:"total_revenue"`
	ActiveServices int       `json:"active_services"`
	PendingOrders  int       `json:"pending_orders"`
	Timestamp      time.Time `json:"timestamp"`
}

// ConnectionInfo معلومات الاتصال
type ConnectionInfo struct {
	ClientID  string   `json:"client_id"`
	UserID    string   `json:"user_id"`
	Channels  []string `json:"channels"`
	Connected bool     `json:"connected"`
	Timestamp time.Time `json:"timestamp"`
}

// GetConnectionStats الحصول على إحصائيات الاتصال
func GetConnectionStats() gin.H {
	manager := GetManager()
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	stats := gin.H{
		"total_channels": len(manager.clients),
		"total_clients":  0,
		"channels":       gin.H{},
		"timestamp":      time.Now(),
	}

	for channel, clients := range manager.clients {
		stats["total_clients"] = stats["total_clients"].(int) + len(clients)
		stats["channels"].(gin.H)[channel] = len(clients)
	}

	return stats
}

// DisconnectClient فصل عميل
func DisconnectClient(clientID string) {
	manager := GetManager()
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for channel, clients := range manager.clients {
		if _, exists := clients[clientID]; exists {
			delete(clients, clientID)
			if len(clients) == 0 {
				delete(manager.clients, channel)
			}
		}
	}
}