#!/bin/bash

# ============================================
# NawthTech AI Models Installation Script
# ============================================
# هذا الملف لتثبيت جميع نماذج الذكاء الاصطناعي المجانية
# ============================================

set -e  # إيقاف عند حدوث خطأ

# ألوان للخروج
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# دالة طباعة
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# دالة التحقق من التثبيت
check_installed() {
    if command -v $1 &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# ============================================
# التثبيت الرئيسي
# ============================================

main() {
    clear
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════╗"
    echo "║      NawthTech AI Models Installer          ║"
    echo "║     منصة الذكاء الاصطناعي للنمو الرقمي      ║"
    echo "╚══════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # التحقق من أن المستخدم root
    if [ "$EUID" -eq 0 ]; then 
        print_warning "تشغيل كـ root! قد يكون هذا خطيراً."
        read -p "هل تريد المتابعة؟ (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # إنشاء مجلدات المشروع
    print_info "إنشاء هيكل المجلدات..."
    mkdir -p ./ai_models
    mkdir -p ./ai_models/text
    mkdir -p ./ai_models/image
    mkdir -p ./ai_models/video
    mkdir -p ./ai_models/audio
    mkdir -p ./data/ai/cache
    print_success "تم إنشاء هيكل المجلدات"
    
    # القسم 1: تثبيت Ollama (النماذج النصية)
    install_ollama
    
    # القسم 2: تثبيت نماذج Hugging Face
    install_huggingface_models
    
    # القسم 3: تثبيت Stable Diffusion للصور
    install_stable_diffusion
    
    # القسم 4: تثبيت نماذج الصوت
    install_audio_models
    
    # القسم 5: إعداد Environment Variables
    setup_environment
    
    # القسم 6: اختبار التثبيت
    test_installation
    
    print_success "✅ تم تثبيت جميع نماذج الذكاء الاصطناعي بنجاح!"
    echo ""
    print_info "🔧 الإعدادات المطلوبة:"
    echo "1. أضف مفاتيح API في ملف .env"
    echo "2. شغّل: docker-compose up -d"
    echo "3. افتح: http://localhost:3000"
    echo ""
    print_info "📁 هيكل المجلدات الجديد:"
    tree -L 2 ./ai_models
}

# ============================================
# 1. تثبيت Ollama
# ============================================

install_ollama() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "تثبيت Ollama للنماذج النصية المحلية..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    # التحقق من وجود Ollama
    if check_installed ollama; then
        print_success "Ollama مثبت مسبقاً"
    else
        print_info "جارٍ تثبيت Ollama..."
        
        # تثبيت Ollama (يدعم Linux و macOS)
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            # Linux
            curl -fsSL https://ollama.com/install.sh | sh
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            /bin/bash -c "$(curl -fsSL https://ollama.com/install.sh)"
        else
            print_error "نظام التشغيل غير مدعوم: $OSTYPE"
            exit 1
        fi
        
        print_success "تم تثبيت Ollama"
    fi
    
    # تشغيل خدمة Ollama
    print_info "تشغيل خدمة Ollama..."
    sudo systemctl enable ollama
    sudo systemctl start ollama
    
    # انتظار حتى يبدأ الخدمة
    sleep 5
    
    # تحميل النماذج النصية
    print_info "جارٍ تحميل النماذج النصية..."
    
    local models=(
        "llama3.2:3b"           # 3B parameters - سريع
        "mistral:7b"            # 7B - جيد للتوليد
        "qwen2.5:7b"           # 7B - دعم عربي ممتاز
        "phi3:mini"            # 3.8B - فعال
        "gemma:7b"             # من Google
        "llama3.2:1b"          # 1B - خفيف جداً
    )
    
    for model in "${models[@]}"; do
        print_info "تحميل: $model"
        ollama pull $model || print_warning "فشل تحميل $model، تخطي..."
    done
    
    print_success "تم تحميل النماذج النصية"
    
    # إنشاء نموذج مخصص لـ NawthTech
    print_info "إنشاء نموذج NawthTech المخصص..."
    cat > ./ai_models/nawthtech-model.Modelfile << 'EOF'
FROM llama3.2:3b

# System Prompt مخصص لـ NawthTech
SYSTEM """
أنت مساعد الذكاء الاصطناعي في NawthTech - منصة النمو الرقمي.
متخصص في:
1. التسويق الرقمي والنمو
2. استراتيجيات الأعمال
3. كتابة المحتوى بالعربية والإنجليزية
4. تحليل السوق والمنافسين
5. نصائح للشركات الناشئة

كن مفيداً، دقيقاً، ومركزاً على تقديم حلول عملية.
استخدم لغة واضحة واحترافية.
"""

PARAMETER temperature 0.7
PARAMETER top_p 0.9
PARAMETER num_ctx 4096
EOF
    
    ollama create nawthtech -f ./ai_models/nawthtech-model.Modelfile
    print_success "تم إنشاء نموذج NawthTech المخصص"
}

# ============================================
# 2. تثبيت نماذج Hugging Face
# ============================================

install_huggingface_models() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "تثبيت نماذج Hugging Face..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    # التحقق من وجود Python و pip
    if ! check_installed python3; then
        print_error "Python3 غير مثبت"
        print_info "جارٍ تثبيت Python3..."
        sudo apt-get update
        sudo apt-get install -y python3 python3-pip
    fi
    
    # تثبيت Hugging Face CLI
    print_info "تثبيت Hugging Face CLI..."
    pip3 install huggingface-hub
    
    # إنشاء مجلد النماذج
    mkdir -p ./ai_models/huggingface
    
    # تحميل النماذج المفيدة
    print_info "جارٍ تحميل النماذج من Hugging Face..."
    
    # نموذج الترجمة (عربي-إنجليزي)
    print_info "تحميل نموذج الترجمة..."
    huggingface-cli download \
        "Helsinki-NLP/opus-mt-ar-en" \
        --local-dir ./ai_models/huggingface/translation-ar-en \
        --local-dir-use-symlinks False
    
    # نموذج التلخيص
    print_info "تحميل نموذج التلخيص..."
    huggingface-cli download \
        "facebook/bart-large-cnn" \
        --local-dir ./ai_models/huggingface/summarization \
        --local-dir-use-symlinks False
    
    # نموذج التصنيف
    print_info "تحميل نموذج التصنيف..."
    huggingface-cli download \
        "distilbert-base-uncased-finetuned-sst-2-english" \
        --local-dir ./ai_models/huggingface/sentiment \
        --local-dir-use-symlinks False
    
    print_success "تم تحميل نماذج Hugging Face"
}

# ============================================
# 3. تثبيت Stable Diffusion
# ============================================

install_stable_diffusion() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "تثبيت Stable Diffusion للصور..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    # إنشاء مجلد Stable Diffusion
    mkdir -p ./ai_models/stable-diffusion
    
    print_info "تحميل Stable Diffusion XL..."
    huggingface-cli download \
        "stabilityai/stable-diffusion-xl-base-1.0" \
        --local-dir ./ai_models/stable-diffusion/sdxl \
        --local-dir-use-symlinks False \
        --exclude "*.safetensors" \
        --exclude "*.ckpt"
    
    # تحميل نموذج أصغر للاختبار
    print_info "تحميل Stable Diffusion 2.1 (أصغر)..."
    huggingface-cli download \
        "stabilityai/stable-diffusion-2-1" \
        --local-dir ./ai_models/stable-diffusion/sd2.1 \
        --local-dir-use-symlinks False \
        --exclude "*.safetensors" \
        --exclude "*.ckpt"
    
    # إنشاء Dockerfile لـ Stable Diffusion
    cat > ./ai_models/stable-diffusion/Dockerfile << 'EOF'
FROM pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

WORKDIR /app

# تثبيت dependencies
RUN pip install --no-cache-dir \
    diffusers==0.24.0 \
    transformers==4.35.0 \
    accelerate==0.24.1 \
    torchvision==0.16.0 \
    pillow==10.1.0 \
    scipy==1.11.4 \
    flask==3.0.0

# نسخ النموذج المحلي
COPY sdxl/ /app/models/sdxl/

# إنشاء تطبيق Flask بسيط
COPY app.py /app/

EXPOSE 7860

CMD ["python", "app.py"]
EOF
    
    # إنشاء تطبيق Flask
    cat > ./ai_models/stable-diffusion/app.py << 'EOF'
from flask import Flask, request, jsonify
from diffusers import StableDiffusionXLPipeline
import torch
from PIL import Image
import io
import base64

app = Flask(__name__)

# تحميل النموذج
print("Loading Stable Diffusion XL model...")
pipe = StableDiffusionXLPipeline.from_pretrained(
    "/app/models/sdxl",
    torch_dtype=torch.float16,
    use_safetensors=True,
    variant="fp16"
)
pipe.to("cuda" if torch.cuda.is_available() else "cpu")
print("Model loaded successfully!")

@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "healthy", "model": "stable-diffusion-xl"})

@app.route('/generate', methods=['POST'])
def generate():
    try:
        data = request.json
        prompt = data.get('prompt', '')
        
        if not prompt:
            return jsonify({"error": "Prompt is required"}), 400
        
        # توليد الصورة
        image = pipe(
            prompt=prompt,
            num_inference_steps=25,
            guidance_scale=7.5
        ).images[0]
        
        # تحويل إلى base64
        buffered = io.BytesIO()
        image.save(buffered, format="PNG")
        img_str = base64.b64encode(buffered.getvalue()).decode()
        
        return jsonify({
            "success": True,
            "image": f"data:image/png;base64,{img_str}",
            "prompt": prompt
        })
        
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=7860)
EOF
    
    print_success "تم إعداد Stable Diffusion"
}

# ============================================
# 4. تثبيت نماذج الصوت
# ============================================

install_audio_models() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "تثبيت نماذج الصوت..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    mkdir -p ./ai_models/audio
    
    # تحميل Whisper للتعرف على الصوت
    print_info "تحميل Whisper (OpenAI)..."
    huggingface-cli download \
        "openai/whisper-medium" \
        --local-dir ./ai_models/audio/whisper \
        --local-dir-use-symlinks False
    
    # تحميل Bark لتوليد الصوت
    print_info "تحميل Bark (Suno AI)..."
    huggingface-cli download \
        "suno/bark" \
        --local-dir ./ai_models/audio/bark \
        --local-dir-use-symlinks False
    
    # تحميل XTTS للصوت متعدد اللغات
    print_info "تحميل XTTS-v2 (Coqui AI)..."
    huggingface-cli download \
        "coqui/XTTS-v2" \
        --local-dir ./ai_models/audio/xtts \
        --local-dir-use-symlinks False
    
    # إنشاء Dockerfile لخدمات الصوت
    cat > ./ai_models/audio/Dockerfile << 'EOF'
FROM python:3.10-slim

WORKDIR /app

# تثبيت system dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    libsndfile1 \
    && rm -rf /var/lib/apt/lists/*

# تثبيت Python packages
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# نسخ النماذج
COPY whisper/ /app/models/whisper/
COPY bark/ /app/models/bark/
COPY xtts/ /app/models/xtts/

# نسخ تطبيق الصوت
COPY audio_app.py /app/

EXPOSE 7861

CMD ["python", "audio_app.py"]
EOF
    
    # إنشاء requirements.txt للصوت
    cat > ./ai_models/audio/requirements.txt << 'EOF'
openai-whisper==20231117
TTS==0.22.0
torch==2.1.0
torchaudio==2.1.0
flask==3.0.0
numpy==1.24.3
scipy==1.11.4
soundfile==0.12.1
EOF
    
    # إنشاء تطبيق الصوت
    cat > ./ai_models/audio/audio_app.py << 'EOF'
from flask import Flask, request, jsonify
import whisper
import torch
import io
import base64
from TTS.api import TTS

app = Flask(__name__)

# تحميل Whisper
print("Loading Whisper model...")
whisper_model = whisper.load_model("/app/models/whisper")
print("Whisper loaded!")

# تحميل TTS
print("Loading TTS models...")
tts = TTS(model_name="tts_models/multilingual/multi-dataset/xtts_v2", progress_bar=False)
print("TTS loaded!")

@app.route('/health', methods=['GET'])
def health():
    return jsonify({"status": "healthy"})

@app.route('/transcribe', methods=['POST'])
def transcribe():
    try:
        if 'audio' not in request.files:
            return jsonify({"error": "No audio file"}), 400
        
        audio_file = request.files['audio']
        
        # تحويل الصوت إلى نص
        result = whisper_model.transcribe(audio_file)
        
        return jsonify({
            "success": True,
            "text": result["text"],
            "language": result["language"]
        })
        
    except Exception as e:
        return jsonify({"error": str(e)}), 500

@app.route('/tts', methods=['POST'])
def text_to_speech():
    try:
        data = request.json
        text = data.get('text', '')
        language = data.get('language', 'en')
        
        if not text:
            return jsonify({"error": "Text is required"}), 400
        
        # توليد الصوت
        audio_path = "/tmp/output.wav"
        tts.tts_to_file(
            text=text,
            speaker_wav=None,  # استخدام الصوت الافتراضي
            language=language,
            file_path=audio_path
        )
        
        # قراءة الملف وإرجاعه
        with open(audio_path, 'rb') as f:
            audio_data = f.read()
        
        audio_b64 = base64.b64encode(audio_data).decode()
        
        return jsonify({
            "success": True,
            "audio": f"data:audio/wav;base64,{audio_b64}",
            "text": text,
            "language": language
        })
        
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=7861)
EOF
    
    print_success "تم إعداد نماذج الصوت"
}

# ============================================
# 5. إعداد Environment Variables
# ============================================

setup_environment() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "إعداد Environment Variables..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    # إنشاء ملف .env
    cat > .env << 'EOF'
# ============================================
# NawthTech AI Environment Configuration
# ============================================

# Ollama Configuration
OLLAMA_HOST=http://localhost:11434
OLLAMA_MODEL=nawthtech
OLLAMA_KEEP_ALIVE=5m

# Hugging Face (احصل على Token من: https://huggingface.co/settings/tokens)
HUGGINGFACE_TOKEN=your_huggingface_token_here

# Google Gemini (مجاني: https://makersuite.google.com/app/apikey)
GEMINI_API_KEY=your_gemini_api_key_here

# Stability AI (25 صورة مجانية/شهر: https://platform.stability.ai/)
STABILITY_API_KEY=your_stability_key_here

# Model Paths
AI_MODELS_PATH=./ai_models
AI_CACHE_PATH=./data/ai/cache
AI_DATA_PATH=./data/ai

# Service Ports
OLLAMA_PORT=11434
STABLE_DIFFUSION_PORT=7860
AUDIO_SERVICE_PORT=7861
BACKEND_PORT=8080
FRONTEND_PORT=3000

# AI Configuration
DEFAULT_TEXT_MODEL=gemini-2.0-flash
DEFAULT_IMAGE_MODEL=stable-diffusion-xl
DEFAULT_AUDIO_MODEL=whisper-medium
DEFAULT_TRANSLATION_MODEL=opus-mt-ar-en

# Rate Limits (Free Tier)
MAX_REQUESTS_PER_MINUTE=30
MAX_REQUESTS_PER_DAY=1000
MAX_IMAGES_PER_DAY=10
MAX_VIDEOS_PER_DAY=3

# User Quotas (Free Tier)
FREE_USER_QUOTA_TEXT=10000      # كلمات/شهر
FREE_USER_QUOTA_IMAGES=10       # صور/شهر
FREE_USER_QUOTA_VIDEOS=3        # فيديوهات/شهر
FREE_USER_QUOTA_AUDIO=30        # دقائق/شهر

# Cache Settings
AI_CACHE_TTL=24h
AI_CACHE_MAX_SIZE=10GB

# Logging
AI_LOG_LEVEL=info
AI_LOG_PATH=./logs/ai.log

# Monitoring
ENABLE_AI_METRICS=true
AI_METRICS_PORT=9090
EOF
    
    print_success "تم إنشاء ملف .env"
    print_warning "⚠️  يرجى تعديل ملف .env وإضافة مفاتيح API الخاصة بك"
    
    # إنشاء ملف docker-compose.yml للـ AI services
    cat > docker-compose.ai.yml << 'EOF'
version: '3.8'

services:
  # Ollama Service
  ollama:
    image: ollama/ollama:latest
    container_name: nawthtech-ollama
    ports:
      - "11434:11434"
    volumes:
      - ./ai_models/ollama:/root/.ollama
      - ./ai_models/nawthtech-model.Modelfile:/root/nawthtech.Modelfile
    environment:
      - OLLAMA_HOST=0.0.0.0
      - OLLAMA_KEEP_ALIVE=5m
    restart: unless-stopped
    networks:
      - ai-network
    command: >
      sh -c "
        ollama serve &
        sleep 10 &&
        ollama create nawthtech -f /root/nawthtech.Modelfile
        wait
      "
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu]

  # Stable Diffusion Service
  stable-diffusion:
    build: ./ai_models/stable-diffusion
    container_name: nawthtech-sd
    ports:
      - "7860:7860"
    volumes:
      - ./ai_models/stable-diffusion/sdxl:/app/models/sdxl
      - ./data/ai/cache:/root/.cache
    environment:
      - HF_TOKEN=${HUGGINGFACE_TOKEN}
      - MODEL_PATH=/app/models/sdxl
    restart: unless-stopped
    networks:
      - ai-network
    depends_on:
      - ollama
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]

  # Audio Service
  audio-service:
    build: ./ai_models/audio
    container_name: nawthtech-audio
    ports:
      - "7861:7861"
    volumes:
      - ./ai_models/audio:/app/models
      - ./data/ai/cache:/root/.cache
    environment:
      - HF_TOKEN=${HUGGINGFACE_TOKEN}
    restart: unless-stopped
    networks:
      - ai-network

  # AI Gateway (Reverse Proxy)
  ai-gateway:
    image: nginx:alpine
    container_name: nawthtech-ai-gateway
    ports:
      - "8000:80"
    volumes:
      - ./nginx/ai-gateway.conf:/etc/nginx/conf.d/default.conf
    restart: unless-stopped
    networks:
      - ai-network
    depends_on:
      - ollama
      - stable-diffusion
      - audio-service

networks:
  ai-network:
    driver: bridge

volumes:
  ollama_data:
  sd_cache:
  audio_cache:
EOF
    
    # إنشاء مجلد nginx للإعدادات
    mkdir -p nginx
    
    cat > nginx/ai-gateway.conf << 'EOF'
server {
    listen 80;
    server_name localhost;
    
    # Health check endpoint
    location /health {
        return 200 '{"status": "healthy"}';
        add_header Content-Type application/json;
    }
    
    # Ollama API
    location /api/ollama/ {
        proxy_pass http://ollama:11434/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # زيادة timeouts للنماذج الكبيرة
        proxy_read_timeout 300s;
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
    }
    
    # Stable Diffusion API
    location /api/sd/ {
        proxy_pass http://stable-diffusion:7860/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # Audio Service API
    location /api/audio/ {
        proxy_pass http://audio-service:7861/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # Default
    location / {
        return 404 '{"error": "Not found"}';
        add_header Content-Type application/json;
    }
}
EOF
    
    print_success "تم إعداد Docker Compose للـ AI Services"
}

# ============================================
# 6. اختبار التثبيت
# ============================================

test_installation() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "اختبار التثبيت..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    # اختبار Ollama
    print_info "اختبار Ollama..."
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        print_success "Ollama يعمل بشكل صحيح"
    else
        print_warning "Ollama غير قيد التشغيل. جارٍ البدء..."
        sudo systemctl start ollama
        sleep 5
    fi
    
    # اختبار Python packages
    print_info "اختبار Python packages..."
    if python3 -c "import huggingface_hub" 2>/dev/null; then
        print_success "Hugging Face Hub مثبت"
    else
        print_warning "جارٍ تثبيت Hugging Face Hub..."
        pip3 install huggingface-hub
    fi
    
    # اختبار Docker
    print_info "اختبار Docker..."
    if docker --version > /dev/null 2>&1; then
        print_success "Docker مثبت"
    else
        print_error "Docker غير مثبت. يرجى تثبيته أولاً."
        print_info "تعليمات التثبيت: https://docs.docker.com/engine/install/"
        exit 1
    fi
    
    # اختبار Docker Compose
    if docker-compose --version > /dev/null 2>&1; then
        print_success "Docker Compose مثبت"
    else
        print_warning "Docker Compose غير مثبت. جارٍ التثبيت..."
        sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
            -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
    fi
    
    # التحقق من مساحة التخزين
    print_info "التحقق من مساحة التخزين..."
    free_space=$(df -h . | awk 'NR==2 {print $4}')
    print_info "المساحة الحرة: $free_space"
    
    if [[ ${free_space%G} -lt 20 ]]; then
        print_warning "⚠️  مساحة التخزين منخفضة! تحتاج 20GB على الأقل للنماذج."
    fi
    
    # التحقق من الذاكرة RAM
    print_info "التحقق من الذاكرة..."
    total_ram=$(free -g | awk 'NR==2 {print $2}')
    print_info "الذاكرة الكلية: ${total_ram}GB"
    
    if [[ $total_ram -lt 8 ]]; then
        print_warning "⚠️  الذاكرة منخفضة! تحتاج 8GB على الأقل لتشغيل النماذج."
    fi
    
    # إنشاء تقرير التثبيت
    create_installation_report
}

# ============================================
# إنشاء تقرير التثبيت
# ============================================

create_installation_report() {
    echo ""
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    print_info "إنشاء تقرير التثبيت..."
    echo -e "${BLUE}══════════════════════════════════════════════${NC}"
    
    cat > INSTALLATION_REPORT.md << 'EOF'
# NawthTech AI Models Installation Report

## 📅 تاريخ التثبيت
'$(date)'

## ✅ المكونات المثبتة

### 1. Ollama (النماذج النصية)
- ✅ llama3.2:3b
- ✅ mistral:7b  
- ✅ qwen2.5:7b
- ✅ phi3:mini
- ✅ gemma:7b
- ✅ nawthtech (مخصص)

### 2. Hugging Face Models
- ✅ Helsinki-NLP/opus-mt-ar-en (ترجمة)
- ✅ facebook/bart-large-cnn (تلخيص)
- ✅ distilbert-base-uncased-finetuned-sst-2-english (تصنيف)

### 3. Stable Diffusion (الصور)
- ✅ stabilityai/stable-diffusion-xl-base-1.0
- ✅ stabilityai/stable-diffusion-2-1

### 4. Audio Models (الصوت)
- ✅ openai/whisper-medium (تعرف على الكلام)
- ✅ suno/bark (توليد صوت)
- ✅ coqui/XTTS-v2 (نص إلى صوت)

## 🚀 كيفية التشغيل

### الطريقة 1: استخدام Docker Compose
```bash
# تشغيل جميع خدمات AI
docker-compose -f docker-compose.ai.yml up -d

# عرض السجلات
docker-compose -f docker-compose.ai.yml logs -f