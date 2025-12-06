#!/bin/bash

echo "🚀 إعداد commit نهائي..."

# إضافة جميع الملفات
git add -A

# إنشاء commit
cat > /tmp/commit_msg.txt << 'EOF'
refactor: overhaul frontend configuration and AI system

Configuration Changes:
- Simplify package.json with working build script
- Remove complex tsconfig files (app.json, node.json)
- Clean up TypeScript configuration
- Remove duplicate vite.config.js

AI System Additions:
- Create AIContentGenerator component
- Create AIMediaGenerator component  
- Fix and update useAI hook
- Update AI services (api.ts, content.ts, media.ts)
- Create useContentGeneration hook

This commit resolves build issues and establishes
a clean, working frontend foundation.
EOF

git commit -F /tmp/commit_msg.txt
rm -f /tmp/commit_msg.txt

echo "✅ تم commit بنجاح!"
echo "جاري الدفع إلى GitHub..."
git push origin main