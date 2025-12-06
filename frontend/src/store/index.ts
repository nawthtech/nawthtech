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
  devTools: process.env.NODE_ENV !== 'production',
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;