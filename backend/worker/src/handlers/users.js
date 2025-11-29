// worker/src/handlers/users.js
import { withDatabase } from '../utils/database.js'
import { ObjectId } from 'mongodb'

export const userHandlers = {
  getProfile: withDatabase(async (request, env) => {
    const { db, dbType } = request
    
    if (dbType !== 'mongodb') {
      throw new Error('This handler requires MongoDB')
    }

    const userId = request.user?.id
    
    if (!userId) {
      return Response.json({
        success: false,
        error: 'UNAUTHORIZED'
      }, { status: 401 })
    }

    const user = await db.collection('users').findOne({
      _id: new ObjectId(userId)
    })

    if (!user) {
      return Response.json({
        success: false,
        error: 'USER_NOT_FOUND'
      }, { status: 404 })
    }

    // إخفاء البيانات الحساسة
    const { password, ...safeUser } = user

    return Response.json({
      success: true,
      data: safeUser
    })
  }),

  // معالجات أخرى...
  getUsers: withDatabase(async (request, env) => {
    const { db } = request
    
    const users = await db.collection('users')
      .find({})
      .project({ password: 0 }) // استبعاد كلمة السر
      .limit(50)
      .toArray()

    return Response.json({
      success: true,
      data: users
    })
  })
}