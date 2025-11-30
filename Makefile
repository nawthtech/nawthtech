# Makefile for NawthTech Platform

.PHONY: help install dev build test deploy clean

# الألوان للواجهة
ORANGE = \033[38;5;208m
BLUE = \033[38;5;19m
WHITE = \033[0m
RESET = \033[0m

# المساعدة
help:
	@echo "$(ORANGE)🚀 NawthTech Platform Commands$(RESET)"
	@echo ""
	@echo "$(BLUE)Development:$(RESET)"
	@echo "  make install    - تثبيت جميع الاعتمادات"
	@echo "  make dev        - تشغيل بيئة التطوير الكاملة"
	@echo "  make backend    - تشغيل الخادم الخلفي فقط"
	@echo "  make frontend   - تشغيل الواجهة الأمامية فقط"
	@echo ""
	@echo "$(BLUE)Building:$(RESET)"
	@echo "  make build      - بناء جميع المكونات"
	@echo "  make build-prod - بناء للإنتاج"
	@echo ""
	@echo "$(BLUE)Testing:$(RESET)"
	@echo "  make test       - تشغيل جميع الاختبارات"
	@echo "  make test-backend - اختبارات الخادم الخلفي"
	@echo "  make test-frontend - اختبارات الواجهة الأمامية"
	@echo ""
	@echo "$(BLUE)Deployment:$(RESET)"
	@echo "  make deploy     - النشر إلى البيئة المستهدفة"
	@echo "  make docker     - بناء صور Docker"
	@echo ""
	@echo "$(BLUE)Maintenance:$(RESET)"
	@echo "  make clean      - تنظيف الملفات المبنية"
	@echo "  make logs       - عرض السجلات"
	@echo ""

# التثبيت
install:
	@echo "$(ORANGE)📦 تثبيت اعتمادات NawthTech...$(RESET)"
	@cd backend && make deps
	@cd frontend && npm install
	@echo "$(ORANGE)✅ تم التثبيت بنجاح$(RESET)"

# بيئة التطوير
dev:
	@echo "$(ORANGE)🚀 تشغيل بيئة التطوير الكاملة...$(RESET)"
	docker-compose up --build

# الخادم الخلفي فقط
backend:
	@echo "$(ORANGE)🔧 تشغيل الخادم الخلفي...$(RESET)"
	@cd backend && make run

# الواجهة الأمامية فقط
frontend:
	@echo "$(ORANGE)🎨 تشغيل الواجهة الأمامية...$(RESET)"
	@cd frontend && npm run dev

# البناء
build:
	@echo "$(ORANGE)🏗️ بناء جميع المكونات...$(RESET)"
	@cd backend && make build
	@cd frontend && npm run build

# البناء للإنتاج
build-prod:
	@echo "$(ORANGE)🏗️ بناء للإنتاج...$(RESET)"
	@cd backend && make build
	@cd frontend && npm run build:prod

# الاختبارات
test:
	@echo "$(ORANGE)🧪 تشغيل جميع الاختبارات...$(RESET)"
	@cd backend && make test
	@cd frontend && npm run test

test-backend:
	@cd backend && make test

test-frontend:
	@cd frontend && npm run test

# النشر
deploy:
	@echo "$(ORANGE)🚀 النشر إلى الإنتاج...$(RESET)"
	@echo "هذا الأمر سينشر التطبيق إلى البيئة المستهدفة"

# بناء Docker
docker:
	@echo "$(ORANGE)🐳 بناء صور Docker...$(RESET)"
	docker-compose build

# التنظيف
clean:
	@echo "$(ORANGE)🧹 تنظيف الملفات المبنية...$(RESET)"
	@cd backend && make clean
	@cd frontend && npm run clean
	@docker-compose down -v

# السجلات
logs:
	@echo "$(ORANGE)📋 عرض سجلات التطبيق...$(RESET)"
	docker-compose logs -f

# فحص الصحة
health:
	@echo "$(ORANGE)🔍 فحص صحة النظام...$(RESET)"
	@cd backend && make health

# قاعدة البيانات
migrate:
	@echo "$(ORANGE)🗄️ ترحيل قاعدة البيانات...$(RESET)"
	@cd backend && make migrate