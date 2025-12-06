import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import aiReducer from './aiSlice';
import storeReducer from './storeSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    ai: aiReducer,
    store: storeReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }),
  // استخدام متغيرات بيئة Vite
  devTools: import.meta.env.DEV || import.meta.env.MODE === 'development',
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;