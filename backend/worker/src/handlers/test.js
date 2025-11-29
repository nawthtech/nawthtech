import { withMongoDB } from '../middleware/mongodb.js'

export const testHandlers = {
  testConnection: withMongoDB(async (request, env) => {
    const { db } = request
    
    // اختبار الاتصال بمحاولة جلب قائمة collections
    const collections = await db.listCollections().toArray()
    
    return Response.json({
      success: true,
      message: 'MongoDB connection successful!',
      data: {
        database: db.databaseName,
        collections: collections.map(c => c.name),
        connected: true
      }
    })
  })
}