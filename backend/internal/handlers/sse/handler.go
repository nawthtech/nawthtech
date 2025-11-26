package sse

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"github.com/nawthtech/nawthtech/backend/internal/quote"
)

// Handler معالج SSE (Server-Sent Events)
func Handler(w http.ResponseWriter, r *http.Request) {
	// التحقق من دعم الـ SSE
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusBadRequest)
		return
	}

	// إعداد الاتصال لـ SSE
	w.Header().Set("Content-Type", "text/event-stream; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// إرسال الرؤوس إلى العميل
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	// إعداد المسجل
	requestID := getRequestID(r.Context())
	eventLogger := logger.Stdout.With(slog.String("request_id", requestID))

	var eventCount uint64

	// دالة ملائمة للاستدعاء عند الاتصال وفي الحلقة
	sendData := func() {
		quoteData := quote.GetRandom()
		fmt.Fprintf(w, "data: %s\n\n", quoteData)

		flusher.Flush()

		eventCount++
		eventLogger.Debug("sent sse event", "quote", quoteData, "event_count", eventCount)
	}

	// تسجيل طلب SSE
	eventLogger.Info("عميل SSE متصل")

	// إرسال بيانات عند الاتصال
	sendData()

	// إعداد المجدول لإرسال بيانات كل ثانية
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{}, 1)

	// انتظار انقطاع العميل ثم إغلاق قناة done
	go func() {
		<-r.Context().Done()

		eventLogger.Info("عميل SSE انقطع", slog.Uint64("event_count", eventCount))
		close(done)
	}()

	// انتظار tick أو قناة done، إرسال بيانات كل tick
	for {
		select {
		case <-ticker.C:
			sendData()
		case <-done:
			return
		}
	}
}

// getRequestID استخراج معرف الطلب من السياق
func getRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}
	return "unknown"
}