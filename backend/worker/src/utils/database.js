// worker/src/utils/database.js - MongoDB فقط
import { MongoClient } from 'mongodb'

export class DatabaseManager {
  constructor(env) {
    this.env = env
    this.mongoClient = null
    this.mongoDb = null
  }

  // الاتصال بـ MongoDB
  async connect() {
    const connectionString = this.env.DATABASE_URL
    
    if (!connectionString) {
      throw new Error('DATABASE_URL is required')
    }

    try {
      const options = {
        maxPoolSize: 5,
        minPoolSize: 1,
        maxIdleTimeMS: 30000,
        serverSelectionTimeoutMS: 10000,
        socketTimeoutMS: 45000,
        compressors: 'zlib',
        retryWrites: true,
        w: 'majority'
      }

      this.mongoClient = new MongoClient(connectionString, options)
      await this.mongoClient.connect()
      
      // استخراج اسم قاعدة البيانات من الرابط أو استخدام الافتراضي
      const dbName = this.extractDatabaseName(connectionString) || 'nawthtech'
      this.mongoDb = this.mongoClient.db(dbName)
      
      console.log('✅ Connected to MongoDB successfully!')
      return this.mongoDb
    } catch (error) {
      console.error('❌ MongoDB connection failed:', error)
      throw error
    }
  }

  // استخراج اسم قاعدة البيانات من رابط الاتصال
  extractDatabaseName(connectionString) {
    try {
      // معالجة رابط MongoDB
      if (connectionString.includes('mongodb+srv://')) {
        // تنسيق SRV
        const url = new URL(connectionString.replace('mongodb+srv://', 'https://'))
        const pathname = url.pathname
        return pathname && pathname !== '/' ? pathname.replace('/', '') : null
      } else if (connectionString.includes('mongodb://')) {
        // تنسيق عادي
        const url = new URL(connectionString.replace('mongodb://', 'http://'))
        const pathname = url.pathname
        return pathname && pathname !== '/' ? pathname.replace('/', '') : null
      }
      
      return null
    } catch (error) {
      console.warn('Could not parse database name from connection string')
      return null
    }
  }

  // الحصول على اتصال قاعدة البيانات
  getConnection() {
    if (this.mongoDb) {
      return {
        type: 'mongodb',
        db: this.mongoDb,
        client: this.mongoClient
      }
    }
    
    throw new Error('No database connection available')
  }

  // إغلاق الاتصال
  async disconnect() {
    if (this.mongoClient) {
      await this.mongoClient.close()
      console.log('🔌 Disconnected from MongoDB')
      this.mongoClient = null
      this.mongoDb = null
    }
  }

  // فحص صحة الاتصال
  async healthCheck() {
    try {
      if (this.mongoDb) {
        await this.mongoDb.command({ ping: 1 })
        return { status: 'healthy', type: 'mongodb' }
      }
      return { status: 'disconnected', type: 'none' }
    } catch (error) {
      return { status: 'unhealthy', type: 'mongodb', error: error.message }
    }
  }
}

// اتصال مخبأ عالمي
let cachedDatabaseManager = null

// إنشاء أو استرجاع مدير قاعدة البيانات
export function getDatabaseManager(env) {
  if (cachedDatabaseManager) {
    return cachedDatabaseManager
  }

  cachedDatabaseManager = new DatabaseManager(env)
  return cachedDatabaseManager
}

// وسيط موحد لقاعدة البيانات
export function withDatabase(handler) {
  return async (request, env, ...args) => {
    const dbManager = getDatabaseManager(env)
    
    try {
      // الاتصال إذا لم يكن متصلاً
      if (!dbManager.mongoDb) {
        await dbManager.connect()
      }

      // إضافة اتصال قاعدة البيانات إلى request
      const connection = dbManager.getConnection()
      request.db = connection.db
      request.dbType = connection.type
      request.dbClient = connection.client

      const result = await handler(request, env, ...args)
      return result

    } catch (error) {
      console.error('Database middleware error:', error)
      
      return new Response(
        JSON.stringify({
          success: false,
          error: 'DATABASE_CONNECTION_FAILED',
          message: 'Unable to connect to database'
        }),
        { 
          status: 503,
          headers: { 'Content-Type': 'application/json' }
        }
      )
    }
  }
}

// للاستخدام بدون وسيط (يدوي)
export async function createDatabaseConnection(env) {
  const dbManager = new DatabaseManager(env)
  await dbManager.connect()
  return dbManager
}

// دالة مساعدة للتعامل مع ObjectId
export function toObjectId(id) {
  if (!id) return null
  
  try {
    const { ObjectId } = require('mongodb')
    return new ObjectId(id)
  } catch (error) {
    console.error('Invalid ObjectId:', id)
    return null
  }
}