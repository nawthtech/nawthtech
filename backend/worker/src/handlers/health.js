// worker/src/handlers/health.js
import { getDatabaseManager } from '../utils/database.js'

export const healthHandlers = {
  async check(request, env) {
    const dbManager = getDatabaseManager(env)
    const dbHealth = await dbManager.healthCheck()

    const healthData = {
      status: dbHealth.status,
      database: dbHealth.type,
      timestamp: new Date().toISOString(),
      environment: env.ENVIRONMENT,
      version: env.API_VERSION,
      service: 'nawthtech-worker'
    }

    return Response.json({
      success: true,
      message: `Service is ${dbHealth.status}`,
      data: healthData
    })
  },

  async ready(request, env) {
    const dbManager = getDatabaseManager(env)
    const dbHealth = await dbManager.healthCheck()
    
    if (dbHealth.status !== 'healthy') {
      return Response.json({
        success: false,
        error: 'SERVICE_NOT_READY',
        message: 'Database is not ready'
      }, { status: 503 })
    }

    return Response.json({
      success: true,
      message: 'Service is ready',
      data: {
        status: 'ready',
        database: dbHealth.type,
        timestamp: new Date().toISOString()
      }
    })
  }
}