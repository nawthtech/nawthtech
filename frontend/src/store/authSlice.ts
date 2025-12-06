import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface User {
  id: number | string;
  username: string;
  email: string;
  full_name?: string;
  avatar?: string;
  phone?: string;
  created_at?: string;
  updated_at?: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  lastLogin: string | null;
  permissions: string[];
  roles: string[];
}

const initialState: AuthState = {
  user: null,
  token: localStorage.getItem('access_token') || null,
  isAuthenticated: !!localStorage.getItem('access_token'),
  isLoading: false,
  error: null,
  lastLogin: localStorage.getItem('last_login') || null,
  permissions: JSON.parse(localStorage.getItem('permissions') || '[]'),
  roles: JSON.parse(localStorage.getItem('roles') || '[]'),
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    // بدء تسجيل الدخول
    loginStart: (state) => {
      state.isLoading = true;
      state.error = null;
    },
    
    // نجاح تسجيل الدخول
    loginSuccess: (state, action: PayloadAction<{ 
      user: User; 
      token: string; 
      permissions?: string[]; 
      roles?: string[] 
    }>) => {
      state.isLoading = false;
      state.isAuthenticated = true;
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.permissions = action.payload.permissions || [];
      state.roles = action.payload.roles || [];
      state.lastLogin = new Date().toISOString();
      state.error = null;
      
      // حفظ في localStorage
      localStorage.setItem('access_token', action.payload.token);
      localStorage.setItem('last_login', state.lastLogin);
      localStorage.setItem('permissions', JSON.stringify(state.permissions));
      localStorage.setItem('roles', JSON.stringify(state.roles));
    },
    
    // فشل تسجيل الدخول
    loginFailure: (state, action: PayloadAction<string>) => {
      state.isLoading = false;
      state.isAuthenticated = false;
      state.user = null;
      state.token = null;
      state.permissions = [];
      state.roles = [];
      state.error = action.payload;
      
      // تنظيف localStorage
      localStorage.removeItem('access_token');
      localStorage.removeItem('last_login');
      localStorage.removeItem('permissions');
      localStorage.removeItem('roles');
    },
    
    // تسجيل الخروج
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      state.permissions = [];
      state.roles = [];
      state.error = null;
      state.lastLogin = null;
      
      // تنظيف localStorage
      localStorage.clear();
    },
    
    // تحديث بيانات المستخدم
    updateUser: (state, action: PayloadAction<Partial<User>>) => {
      if (state.user) {
        state.user = { ...state.user, ...action.payload };
      }
    },
    
    // تحديث الـ Token
    updateToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload;
      localStorage.setItem('access_token', action.payload);
    },
    
    // مسح الأخطاء
    clearError: (state) => {
      state.error = null;
    },
    
    // تحديث الصلاحيات
    updatePermissions: (state, action: PayloadAction<string[]>) => {
      state.permissions = action.payload;
      localStorage.setItem('permissions', JSON.stringify(action.payload));
    },
  },
});

export const {
  loginStart,
  loginSuccess,
  loginFailure,
  logout,
  updateUser,
  updateToken,
  clearError,
  updatePermissions,
} = authSlice.actions;

export default authSlice.reducer;