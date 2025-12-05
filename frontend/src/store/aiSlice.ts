import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { aiService } from '../ai/services/api';
import { AIModel, AIUsage, ContentHistoryItem, MediaItem } from '../ai/types/ai';

interface AIState {
  models: AIModel[];
  usage: AIUsage | null;
  history: ContentHistoryItem[];
  mediaLibrary: MediaItem[];
  activeModel: string;
  config: {
    language: string;
    autoSave: boolean;
    notifications: boolean;
  };
  loading: boolean;
  error: string | null;
}

const initialState: AIState = {
  models: [],
  usage: null,
  history: [],
  mediaLibrary: [],
  activeModel: 'gemini-2.0-flash',
  config: {
    language: 'ar',
    autoSave: true,
    notifications: true,
  },
  loading: false,
  error: null,
};

// Async Thunks
export const fetchModels = createAsyncThunk(
  'ai/fetchModels',
  async () => {
    const response = await aiService.getAvailableModels();
    return response.data.models;
  }
);

export const fetchUsage = createAsyncThunk(
  'ai/fetchUsage',
  async () => {
    const response = await aiService.getUsage();
    return response.data;
  }
);

const aiSlice = createSlice({
  name: 'ai',
  initialState,
  reducers: {
    setActiveModel: (state, action: PayloadAction<string>) => {
      state.activeModel = action.payload;
    },
    setLanguage: (state, action: PayloadAction<string>) => {
      state.config.language = action.payload;
    },
    addToHistory: (state, action: PayloadAction<ContentHistoryItem>) => {
      state.history.unshift(action.payload);
      // احتفظ بـ 50 عنصراً فقط
      if (state.history.length > 50) {
        state.history.pop();
      }
    },
    addToMediaLibrary: (state, action: PayloadAction<MediaItem>) => {
      state.mediaLibrary.unshift(action.payload);
      // احتفظ بـ 30 عنصراً فقط
      if (state.mediaLibrary.length > 30) {
        state.mediaLibrary.pop();
      }
    },
    clearHistory: (state) => {
      state.history = [];
    },
    clearMediaLibrary: (state) => {
      state.mediaLibrary = [];
    },
    setConfig: (state, action: PayloadAction<Partial<AIState['config']>>) => {
      state.config = { ...state.config, ...action.payload };
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchModels.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchModels.fulfilled, (state, action) => {
        state.loading = false;
        state.models = action.payload;
      })
      .addCase(fetchModels.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch models';
      })
      .addCase(fetchUsage.fulfilled, (state, action) => {
        state.usage = action.payload;
      });
  },
});

export const {
  setActiveModel,
  setLanguage,
  addToHistory,
  addToMediaLibrary,
  clearHistory,
  clearMediaLibrary,
  setConfig,
} = aiSlice.actions;

export default aiSlice.reducer;