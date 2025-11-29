import { getMongoDB } from '../utils/mongodb.js'

export function withMongoDB(handler) {
  return async (request, env, ...args) => {
    try {
      // الحصول على رابط الاتصال من الأسرار
      const databaseUrl = request.secrets?.databaseUrl || env.DATABASE_URL
      
      if (!databaseUrl) {
        throw new Error('MongoDB connection string is required')
      }

      console.log('🔗 Connecting to MongoDB...')
      const { db, client } = await getMongoDB(databaseUrl)
      
      // إضافة اتصال DB إلى request
      request.db = db
      request.mongoClient = client

      const result = await handler(request, env, ...args)
      
      return result
    } catch (error) {
      console.error('❌ MongoDB Middleware Error:', error)
      
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