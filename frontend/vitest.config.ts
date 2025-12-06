import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  
  // يجب أن يتطابق مع vite.config.ts
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@pages': path.resolve(__dirname, './src/pages'),
      '@components': path.resolve(__dirname, './src/components'),
      '@ai': path.resolve(__dirname, './src/ai'),
      '@store': path.resolve(__dirname, './src/store'),
      '@services': path.resolve(__dirname, './src/services'),
      '@hooks': path.resolve(__dirname, './src/hooks'),
      '@utils': path.resolve(__dirname, './src/utils'),
      '@assets': path.resolve(__dirname, './src/assets'),
    },
    extensions: ['.mjs', '.js', '.mts', '.ts', '.jsx', '.tsx', '.json'],
  },
  
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    
    // حل مشكلة استيراد الملفات في الاختبارات
    server: {
      deps: {
        inline: [
          /@mui/,
          /@testing-library/,
          /@reduxjs\/toolkit/,
          /react-redux/,
        ],
      },
    },
    
    // Mock للأنظمة غير المتوفرة في بيئة الاختبار
    environmentOptions: {
      jsdom: {
        resources: 'usable',
      },
    },
    
    // الملفات التي يتم تضمينها في الاختبارات
    include: [
      'src/**/*.{test,spec}.{js,jsx,ts,tsx}',
      'src/**/__tests__/**/*.{js,jsx,ts,tsx}',
    ],
    
    // الملفات المستثناة
    exclude: [
      'node_modules/',
      'dist/',
      'build/',
      '**/*.config.*',
      '**/.{idea,git,cache,output,temp}/**',
      '**/{karma,rollup,webpack,vite,vitest,babel,postcss}.config.*',
    ],
    
    // مقلدات افتراضية للملفات
    mockReset: true,
    restoreMocks: true,
    
    // التغطية
    coverage: {
      provider: 'v8',
      enabled: true,
      reporter: ['text', 'json', 'html'],
      reportsDirectory: './coverage',
      exclude: [
        'node_modules/',
        'src/test/',
        '**/*.d.ts',
        '**/*.config.*',
        '**/types.ts',
        '**/index.ts',
        '**/vite-env.d.ts',
        '**/*.stories.{js,jsx,ts,tsx}',
        '**/*.story.{js,jsx,ts,tsx}',
        'public/',
        'dist/',
        'coverage/',
      ],
      include: ['src/**/*.{js,jsx,ts,tsx}'],
      thresholds: {
        lines: 0,
        functions: 0,
        branches: 0,
        statements: 0,
      },
    },
    
    // تهيئة إضافية
    deps: {
      optimizer: {
        web: {
          include: [
            'react',
            'react-dom',
            'react-router-dom',
            '@mui/material',
            '@mui/icons-material',
          ],
        },
      },
      interopDefault: true,
    },
    
    // تحديثات فورية في وضع المراقبة
    watch: false,
    
    // إعدادات النوع
    typecheck: {
      tsconfig: './tsconfig.json',
    },
    
    // إعدادات UI إذا استخدمت واجهة مستخدم vitest
    ui: false,
    
    // إعدادات الأداء
    maxWorkers: '50%',
    minWorkers: 1,
    
    // إعدادات الإخراج
    silent: false,
    reporters: ['default'],
    
    // حل مشكلة اختبارات الملفات الكبيرة
    testTimeout: 10000,
    hookTimeout: 10000,
    
    // إعدادات إضافية لتحسين الأداء
    cache: {
      dir: './node_modules/.vitest',
    },
    
    // إعدادات العزل
    isolate: true,
    
    // إعدادات XML للتقارير (إذا احتجت)
    outputFile: {
      json: './test-results.json',
    },
  },
})