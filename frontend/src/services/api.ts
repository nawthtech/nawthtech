/**
 * API service for making HTTP requests to the NawthTech backend
 * Compatible with Go backend in monorepo
 */

import { settings, getApiEndpoint } from '../config';
import type { ApiResponse, RequestConfig, ErrorResponse } from './types';

// ==================== TYPES ====================
export interface PaginationParams {
  page?: number;
  limit?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  search?: string;
}

export interface UploadProgressEvent {
  loaded: number;
  total: number;
  percentage: number;
}

export type UploadProgressCallback = (progress: UploadProgressEvent) => void;

// ==================== API CLIENT ====================
class APIClient {
  private baseURL: string;
  private defaultHeaders: Record<string, string>;
  private requestInterceptor?: (config: RequestConfig) => RequestConfig;
  private responseInterceptor?: (response: Response) => Response | Promise<Response>;
  private errorInterceptor?: (error: ErrorResponse) => ErrorResponse | Promise<ErrorResponse>;

  constructor() {
    this.baseURL = settings.api.baseURL;
    this.defaultHeaders = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      'Accept-Language': settings.localization.defaultLanguage,
    };
  }

  // ==================== INTERCEPTORS ====================
  setRequestInterceptor(interceptor: (config: RequestConfig) => RequestConfig): void {
    this.requestInterceptor = interceptor;
  }

  setResponseInterceptor(interceptor: (response: Response) => Response | Promise<Response>): void {
    this.responseInterceptor = interceptor;
  }

  setErrorInterceptor(interceptor: (error: ErrorResponse) => ErrorResponse | Promise<ErrorResponse>): void {
    this.errorInterceptor = interceptor;
  }

  // ==================== TOKEN MANAGEMENT ====================
  private getAuthToken(): string | null {
    // Try different storage locations
    const token = 
      localStorage.getItem('nawthtech_auth_token') ||
      sessionStorage.getItem('nawthtech_auth_token') ||
      document.cookie
        .split('; ')
        .find(row => row.startsWith('auth_token='))
        ?.split('=')[1] || 
      null;

    return token;
  }

  private getAuthHeaders(): Record<string, string> {
    const token = this.getAuthToken();
    const headers: Record<string, string> = {};

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    // Add CSRF token if enabled
    if (settings.security.csrf.enabled) {
      const csrfToken = this.getCsrfToken();
      if (csrfToken) {
        headers[settings.security.csrf.headerName] = csrfToken;
      }
    }

    return headers;
  }

  private getCsrfToken(): string | null {
    return document.cookie
      .split('; ')
      .find(row => row.startsWith('csrf_token='))
      ?.split('=')[1] || null;
  }

  // ==================== REQUEST BUILDING ====================
  private buildUrl(
    endpoint: string | { category: string; endpoint: string },
    params?: Record<string, any>
  ): string {
    let url: string;

    if (typeof endpoint === 'string') {
      url = endpoint.startsWith('http') ? endpoint : `${this.baseURL}${endpoint}`;
    } else {
      url = getApiEndpoint(endpoint.category, endpoint.endpoint);
    }

    // Add query parameters
    if (params && Object.keys(params).length > 0) {
      const queryParams = new URLSearchParams();
      
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            value.forEach(v => queryParams.append(`${key}[]`, String(v)));
          } else if (typeof value === 'object') {
            queryParams.append(key, JSON.stringify(value));
          } else {
            queryParams.append(key, String(value));
          }
        }
      });

      const queryString = queryParams.toString();
      if (queryString) {
        url += `${url.includes('?') ? '&' : '?'}${queryString}`;
      }
    }

    return url;
  }

  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    // Apply response interceptor if exists
    if (this.responseInterceptor) {
      const interceptorResult = this.responseInterceptor(response);
      response = interceptorResult instanceof Promise ? await interceptorResult : interceptorResult;
    }

    const contentType = response.headers.get('content-type');
    const isJson = contentType?.includes('application/json');
    const status = response.status;

    // Handle errors
    if (!response.ok) {
      let errorData: any;

      try {
        errorData = isJson ? await response.json() : await response.text();
      } catch {
        errorData = { message: `HTTP ${status}: ${response.statusText}` };
      }

      const error: ErrorResponse = {
        message: errorData.message || `Request failed with status ${status}`,
        status,
        errors: errorData.errors,
        timestamp: new Date().toISOString(),
        path: response.url,
      };

      // Apply error interceptor if exists
      if (this.errorInterceptor) {
        const interceptorResult = this.errorInterceptor(error);
        throw interceptorResult instanceof Promise ? await interceptorResult : interceptorResult;
      }

      throw error;
    }

    // Handle successful responses
    if (status === 204 || response.headers.get('content-length') === '0') {
      return {
        data: null as T,
        status,
        success: true,
        message: 'No content',
      };
    }

    const data = isJson ? await response.json() : await response.text();

    return {
      data,
      status,
      success: true,
      message: data?.message,
      meta: data?.meta,
    };
  }

  // ==================== CORE REQUEST METHOD ====================
  async request<T = any>(
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE',
    endpoint: string | { category: string; endpoint: string },
    data?: any,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const {
      headers = {},
      params,
      timeout = settings.api.timeout,
      signal,
      formData = false,
      onUploadProgress,
    } = config || {};

    // Build request configuration
    const authHeaders = this.getAuthHeaders();
    const requestHeaders: Record<string, string> = {
      ...this.defaultHeaders,
      ...authHeaders,
      ...headers,
    };

    // Remove Content-Type for FormData
    if (formData && requestHeaders['Content-Type']) {
      delete requestHeaders['Content-Type'];
    }

    // Build URL
    const url = this.buildUrl(endpoint, params);

    // Prepare request config
    const requestConfig: RequestInit = {
      method,
      headers: requestHeaders,
      credentials: 'include', // Important for cookies with Go backend
      signal,
    };

    // Handle request body
    if (data && method !== 'GET' && method !== 'DELETE') {
      if (formData) {
        requestConfig.body = data;
      } else {
        requestConfig.body = JSON.stringify(data);
      }
    }

    // Handle upload progress
    if (onUploadProgress && data instanceof FormData) {
      const xhr = new XMLHttpRequest();
      
      return new Promise((resolve, reject) => {
        xhr.open(method, url);
        
        // Set headers
        Object.entries(requestHeaders).forEach(([key, value]) => {
          xhr.setRequestHeader(key, value);
        });

        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable && onUploadProgress) {
            onUploadProgress({
              loaded: event.loaded,
              total: event.total,
              percentage: Math.round((event.loaded / event.total) * 100),
            });
          }
        };

        xhr.onload = () => {
          const response = new Response(xhr.responseText, {
            status: xhr.status,
            statusText: xhr.statusText,
            headers: new Headers(
              xhr.getAllResponseHeaders()
                .split('\r\n')
                .filter(Boolean)
                .reduce((acc: Record<string, string>, line) => {
                  const [key, value] = line.split(': ');
                  if (key && value) {
                    acc[key] = value;
                  }
                  return acc;
                }, {})
            ),
          });

          this.handleResponse<T>(response).then(resolve).catch(reject);
        };

        xhr.onerror = () => {
          reject({
            message: 'Network error',
            status: 0,
            timestamp: new Date().toISOString(),
          } as ErrorResponse);
        };

        xhr.ontimeout = () => {
          reject({
            message: 'Request timeout',
            status: 408,
            timestamp: new Date().toISOString(),
          } as ErrorResponse);
        };

        xhr.timeout = timeout;
        xhr.send(requestConfig.body as any);
      });
    }

    // Apply request interceptor if exists
    let finalConfig = { ...requestConfig };
    if (this.requestInterceptor) {
      finalConfig = this.requestInterceptor(finalConfig);
    }

    // Create abort controller for timeout
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);
    
    if (signal) {
      signal.addEventListener('abort', () => controller.abort());
    }

    try {
      const response = await fetch(url, {
        ...finalConfig,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);
      return await this.handleResponse<T>(response);
    } catch (error) {
      clearTimeout(timeoutId);

      if (error instanceof DOMException && error.name === 'AbortError') {
        throw {
          message: 'Request timeout',
          status: 408,
          timestamp: new Date().toISOString(),
        } as ErrorResponse;
      }

      throw {
        message: error instanceof Error ? error.message : 'Network error',
        status: 0,
        timestamp: new Date().toISOString(),
      } as ErrorResponse;
    }
  }

  // ==================== HTTP METHODS ====================
  get<T = any>(
    endpoint: string | { category: string; endpoint: string },
    config?: Omit<RequestConfig, 'formData' | 'onUploadProgress'>
  ): Promise<ApiResponse<T>> {
    return this.request<T>('GET', endpoint, undefined, config);
  }

  post<T = any>(
    endpoint: string | { category: string; endpoint: string },
    data?: any,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    return this.request<T>('POST', endpoint, data, config);
  }

  put<T = any>(
    endpoint: string | { category: string; endpoint: string },
    data?: any,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    return this.request<T>('PUT', endpoint, data, config);
  }

  patch<T = any>(
    endpoint: string | { category: string; endpoint: string },
    data?: any,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    return this.request<T>('PATCH', endpoint, data, config);
  }

  delete<T = any>(
    endpoint: string | { category: string; endpoint: string },
    config?: Omit<RequestConfig, 'formData' | 'onUploadProgress'>
  ): Promise<ApiResponse<T>> {
    return this.request<T>('DELETE', endpoint, undefined, config);
  }

  // ==================== SPECIALIZED METHODS ====================
  async uploadFile<T = any>(
    file: File,
    additionalData?: Record<string, any>,
    onProgress?: UploadProgressCallback
  ): Promise<ApiResponse<T>> {
    const formData = new FormData();
    formData.append('file', file);

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          formData.append(key, value);
        }
      });
    }

    return this.post<T>(
      { category: 'media', endpoint: 'upload' },
      formData,
      {
        formData: true,
        onUploadProgress: onProgress,
      }
    );
  }

  async uploadFiles<T = any>(
    files: File[],
    additionalData?: Record<string, any>,
    onProgress?: UploadProgressCallback
  ): Promise<ApiResponse<T>> {
    const formData = new FormData();

    files.forEach((file, index) => {
      formData.append(`files[${index}]`, file);
    });

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          formData.append(key, value);
        }
      });
    }

    return this.post<T>(
      { category: 'media', endpoint: 'upload' },
      formData,
      {
        formData: true,
        onUploadProgress: onProgress,
      }
    );
  }

  async downloadFile(
    endpoint: string | { category: string; endpoint: string },
    filename?: string,
    config?: Omit<RequestConfig, 'formData' | 'onUploadProgress'>
  ): Promise<void> {
    const response = await this.get<Blob>(endpoint, {
      ...config,
      headers: {
        ...config?.headers,
        'Accept': 'application/octet-stream',
      },
    });

    const blob = new Blob([response.data], { type: response.data.type });
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = filename || 'download';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(downloadUrl);
  }

  async getPaginated<T = any>(
    endpoint: string | { category: string; endpoint: string },
    params: PaginationParams = {},
    config?: Omit<RequestConfig, 'formData' | 'onUploadProgress'>
  ): Promise<ApiResponse<T[]>> {
    const { page = 1, limit = 10, sortBy, sortOrder, search, ...restParams } = params;

    const queryParams: Record<string, any> = {
      page,
      limit,
      ...restParams,
    };

    if (sortBy) {
      queryParams.sort_by = sortBy;
    }

    if (sortOrder) {
      queryParams.sort_order = sortOrder;
    }

    if (search) {
      queryParams.search = search;
    }

    return this.get<T[]>(endpoint, {
      ...config,
      params: queryParams,
    });
  }
}

// ==================== HELPER FUNCTIONS ====================
export const apiHelpers = {
  /**
   * Fetch with retry logic
   */
  async fetchWithRetry<T = any>(
    fn: () => Promise<ApiResponse<T>>,
    maxRetries: number = settings.api.retryAttempts,
    initialDelay: number = 1000
  ): Promise<ApiResponse<T>> {
    let lastError: ErrorResponse;
    
    for (let attempt = 0; attempt < maxRetries; attempt++) {
      try {
        return await fn();
      } catch (error: unknown) {
        const typedError = error as ErrorResponse;
        lastError = typedError;
        
        // Don't retry on 4xx errors (except 429 - Too Many Requests)
        if (typedError.status >= 400 && typedError.status < 500 && typedError.status !== 429) {
          throw error;
        }

        // Wait before retrying (with exponential backoff)
        if (attempt < maxRetries - 1) {
          const delay = initialDelay * Math.pow(2, attempt);
          await new Promise(resolve => setTimeout(resolve, delay));
        }
      }
    }
    
    throw lastError!;
  },

  /**
   * Create cancelable request
   */
  createCancelableRequest() {
    const controller = new AbortController();
    
    return {
      signal: controller.signal,
      cancel: () => controller.abort(),
    };
  },

  /**
   * Debounced API call
   */
  createDebouncedApiCall<T = any>(
    fn: (...args: any[]) => Promise<ApiResponse<T>>,
    delay: number = settings.performance.debounce.search
  ) {
    let timeoutId: ReturnType<typeof setTimeout>;
    let abortController: AbortController | null = null;
    
    return (...args: any[]): Promise<ApiResponse<T>> => {
      return new Promise((resolve, reject) => {
        if (abortController) {
          abortController.abort();
        }
        
        abortController = new AbortController();
        
        clearTimeout(timeoutId);
        
        timeoutId = setTimeout(async () => {
          try {
            const result = await fn(...args);
            abortController = null;
            resolve(result);
          } catch (error) {
            abortController = null;
            reject(error);
          }
        }, delay);
      });
    };
  },

  /**
   * Validate file before upload
   */
  validateFile(file: File): { valid: boolean; errors: string[] } {
    const errors: string[] = [];
    const maxSize = settings.upload.maxFileSize;
    const allowedTypes = [
      ...settings.upload.allowedImageTypes,
      ...settings.upload.allowedDocumentTypes,
      ...settings.upload.allowedVideoTypes,
    ];

    // Check file size
    if (file.size > maxSize) {
      errors.push(`File size exceeds ${maxSize / 1024 / 1024}MB limit`);
    }

    // Check file type
    if (!allowedTypes.includes(file.type)) {
      errors.push(`File type ${file.type} is not allowed`);
    }

    return {
      valid: errors.length === 0,
      errors,
    };
  },
};

// ==================== API INSTANCE ====================
export const api = new APIClient();

// Set up default interceptors
api.setRequestInterceptor((config) => {
  // Add request timestamp
  const headers = config.headers || {};
  headers['X-Request-Timestamp'] = Date.now().toString();
  
  return {
    ...config,
    headers,
  };
});

api.setResponseInterceptor(async (response) => {
  // Handle rate limiting
  const remaining = response.headers.get('X-RateLimit-Remaining');
  const reset = response.headers.get('X-RateLimit-Reset');
  
  if (remaining && parseInt(remaining, 10) < 10) {
    console.warn(`Rate limit low: ${remaining} requests remaining`);
  }
  
  return response;
});

api.setErrorInterceptor(async (error) => {
  // Log errors in development
  if (settings.development.debug) {
    console.error('API Error:', error);
  }

  // Handle specific error cases
  if (error.status === 401) {
    // Token expired, clear auth
    localStorage.removeItem('nawthtech_auth_token');
    sessionStorage.removeItem('nawthtech_auth_token');
    
    // Redirect to login if not already there
    if (!window.location.pathname.includes('/auth/login')) {
      window.location.href = '/auth/login?expired=true';
    }
  }

  if (error.status === 429) {
    error.message = 'Too many requests. Please try again later.';
  }

  return error;
});

// ==================== EXPORT ====================
export default api;