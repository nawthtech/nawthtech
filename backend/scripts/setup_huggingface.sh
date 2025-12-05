#!/bin/bash

# NawthTech Hugging Face Setup Script

echo "🤖 إعداد Hugging Face لـ NawthTech"
echo "=================================="

# 1. التحقق من Token
if [ -z "$HUGGINGFACE_TOKEN" ]; then
    echo "❌ HUGGINGFACE_TOKEN غير محدد"
    echo ""
    echo "📝 اتبع هذه الخطوات:"
    echo "1. اذهب إلى: https://huggingface.co/settings/tokens"
    echo "2. اضغط على 'New token'"
    echo "3. اختر 'Fine-grained'"
    echo "4. أدخل الاسم: 'nawthtech-ai-platform'"
    echo "5. اختر الأذونات المطلوبة (انظر الوثائق)"
    echo "6. انسخ Token"
    echo "7. أضف إلى .env: HUGGINGFACE_TOKEN=your_token_here"
    exit 1
fi

echo "✅ Token موجود"

# 2. اختبار الاتصال
echo "🔍 اختبار الاتصال بـ Hugging Face..."
curl -s -H "Authorization: Bearer $HUGGINGFACE_TOKEN" \
    https://huggingface.co/api/whoami | python3 -m json.tool

# 3. تحميل النماذج الأساسية
echo "📥 تحميل النماذج الأساسية..."

MODELS=(
    "google/flan-t5-xl"
    "mistralai/Mistral-7B-Instruct-v0.2"
    "Qwen/Qwen2.5-7B-Instruct"
)

for MODEL in "${MODELS[@]}"; do
    echo "جارٍ تحميل: $MODEL"
    huggingface-cli download $MODEL \
        --local-dir "./models/$MODEL" \
        --local-dir-use-symlinks False \
        --resume-download || echo "⚠️ فشل تحميل $MODEL"
done

# 4. إنشاء ملف التكوين
echo "📁 إنشاء ملف التكوين..."

cat > huggingface_config.json << EOF
{
  "token": "$HUGGINGFACE_TOKEN",
  "models": {
    "text": [
      "google/flan-t5-xl",
      "mistralai/Mistral-7B-Instruct-v0.2",
      "Qwen/Qwen2.5-7B-Instruct"
    ],
    "image": [
      "stabilityai/stable-diffusion-xl-base-1.0"
    ],
    "audio": [
      "openai/whisper-large-v3"
    ]
  },
  "rate_limit": 30,
  "cache_dir": "./cache/huggingface"
}
EOF

echo "✅ تم الإعداد بنجاح!"
echo ""
echo "📋 الخطوات التالية:"
echo "1. اختبر التطبيق: go run cmd/test_huggingface/main.go"
echo "2. ابدأ الخدمة: docker-compose up -d"
echo "3. افتح: http://localhost:3000"