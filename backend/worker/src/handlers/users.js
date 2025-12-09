// worker/src/handlers/users.js
export const userHandlers = {
  getProfile: async (request, env) => {
    const userId = request.user?.id

    if (!userId) {
      return new Response(
        JSON.stringify({ success: false, error: 'UNAUTHORIZED' }),
        { status: 401, headers: { 'Content-Type': 'application/json' } }
      )
    }

    // الاتصال بقاعدة D1
    const db = env.D1('NAWTHTECH_DB')

    try {
      const result = await db.prepare(`
        SELECT id, username, email, created_at
        FROM users
        WHERE id = ?
      `).bind(userId).all()

      const user = result.results[0]

      if (!user) {
        return new Response(
          JSON.stringify({ success: false, error: 'USER_NOT_FOUND' }),
          { status: 404, headers: { 'Content-Type': 'application/json' } }
        )
      }

      return new Response(
        JSON.stringify({ success: true, data: user }),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      )

    } catch (err) {
      return new Response(
        JSON.stringify({ success: false, error: 'DB_ERROR', message: err.message }),
        { status: 500, headers: { 'Content-Type': 'application/json' } }
      )
    }
  },

  getUsers: async (request, env) => {
    const db = env.D1('NAWTHTECH_DB')

    try {
      const result = await db.prepare(`
        SELECT id, username, email, created_at
        FROM users
        LIMIT 50
      `).all()

      return new Response(
        JSON.stringify({ success: true, data: result.results }),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      )

    } catch (err) {
      return new Response(
        JSON.stringify({ success: false, error: 'DB_ERROR', message: err.message }),
        { status: 500, headers: { 'Content-Type': 'application/json' } }
      )
    }
  }
}