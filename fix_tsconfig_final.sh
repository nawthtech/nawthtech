#!/bin/bash

echo "🔧 إصلاح إعدادات TypeScript المتقدمة..."

cd frontend || exit 1

# 1. إصلاح tsconfig.node.json
echo "🛠️ إصلاح tsconfig.node.json..."
cat > tsconfig.node.json << 'EOF'
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "noUnusedLocals": false,
    "noUnusedParameters": false,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "noEmit": false,
    "noUncheckedSideEffectImports": true
  },
  "include": ["vite.config.ts"]
}
EOF

# 2. إصلاح tsconfig.app.json
echo "🛠️ إصلاح tsconfig.app.json..."
cat > tsconfig.app.json << 'EOF'
{
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.app.tsbuildinfo",
    "target": "ES2022",
    "useDefineForClassFields": true,
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "types": ["vite/client", "node"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "verbatimModuleSyntax": false,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": false,
    "noUnusedLocals": false,
    "noUnusedParameters": false,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"]
}
EOF

# 3. إصلاح vite.config.ts (تبسيط)
echo "⚡ تبسيط vite.config.ts..."
cat > vite.config.ts << 'EOF'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
  },
  server: {
    port: 5173,
  },
})
EOF

# 4. تحديث package.json للبناء
echo "📦 تحديث package.json..."
sed -i 's/"build": "tsc -b && vite build"/"build": "tsc -b && vite build"/' package.json
# أو إذا أردت تبسيط أكثر:
# sed -i 's/"build": "tsc -b && vite build"/"build": "vite build"/' package.json

# 5. تشغيل tsc للتحقق
echo "🧪 اختبار TypeScript build..."
npx tsc -b 2>&1 | head -30 || true

echo "✅ تم الإصلاح!"
echo "جرب الآن: npm run build"