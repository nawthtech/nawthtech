#!/bin/bash

echo "🧹 تنظيف وإعادة تعيين كاملة..."

# 1. تنظيف ملفات git الزائدة
echo "🗑️  حذف ملفات التعارض..."
rm -f frontend/package_*.json 2>/dev/null || true
rm -f frontend/*.tsbuildinfo 2>/dev/null || true
rm -f frontend/vite.config.js 2>/dev/null || true

# 2. إعادة تعيين package.json
echo "📦 إعادة إنشاء package.json..."
cd frontend || exit 1

cat > package.json << 'EOF'
{
  "name": "nawthtech-client",
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "lint": "eslint . --ext ts,tsx --max-warnings 5",
    "preview": "vite preview",
    "test": "vitest run --passWithNoTests"
  },
  "dependencies": {
    "@mui/material": "^5.15.0",
    "@mui/icons-material": "^5.15.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    "axios": "^1.6.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@types/node": "^20.0.0",
    "@vitejs/plugin-react": "^4.0.0",
    "typescript": "^5.0.0",
    "vite": "^5.0.0"
  }
}
EOF

# 3. تبسيط vite.config.ts
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

# 4. تبسيط tsconfig
echo "📝 تبسيط tsconfig..."
cat > tsconfig.json << 'EOF'
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "node",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": false,
    "noUnusedLocals": false,
    "noUnusedParameters": false
  },
  "include": ["src"],
  "exclude": ["node_modules", "dist"]
}
EOF

# 5. حذف ملفات tsconfig المعقدة
rm -f tsconfig.app.json tsconfig.node.json 2>/dev/null || true

cd ..

# 6. إعادة تعيين git
echo "🔄 إعادة تعيين git..."
git add -A
git status

echo "✅ تم التنظيف!"
echo ""
echo "📋 الخطوات التالية:"
echo "1. git commit -m 'chore: clean up and reset configuration'"
echo "2. git push origin main"
echo "3. cd frontend && npm install && npm run build"