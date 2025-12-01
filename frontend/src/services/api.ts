/**
 * API service for making HTTP requests to the NawthTech backend
 * Compatible with Go backend in monorepo
 */

import { settings } from '../config';

// ==================== TYPES ====================
export interface ApiResponse<T = any> {
  data: T;
  status: number;
  success: boolean;
  message?: string;
  meta?: {
    pagination?: {
      page: number;
      limit: number;
      total: number;
      totalPages: number;
      hasNext: boolean;
      hasPrev: boolean;
    };
    [key: string]: any;
  };
}

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

export interface ErrorResponse {
  message: string;
  status: number;
  errors?: Record<string, string[]>;
  timestamp: string;
  path?: string;
}

export interface RequestConfig {
  headers?: Record<string, string>;
  params?: Record<string, any>;
  timeout?: number;
  signal?: AbortSignal;
  formData?: boolean;
  onUploadProgress?: UploadProgressCallback;
}

// ==================== API CLIENT ====================
class APIClient {
  private baseURL: string;
  private defaultHeaders: Record<string, string>;

  constructor() {
    this.baseURL = settings.api.baseURL;
    this.defaultHeaders = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      'Accept-Language': settings.localization.defaultLanguage,
    };
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
      // Access endpoints dynamically
      const endpoints = settings.api.endpoints as any;
      const endpointPath = endpoints[endpoint.category]?.[endpoint.endpoint];
      
      if (!endpointPath) {
        throw new Error(`Endpoint not found: ${endpoint.category}.${endpoint.endpoint}`);
      }
      
      url = `${this.baseURL}${endpointPath}`;
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

      throw {
        message: errorData.message || `Request failed with status ${status}`,
        status,
        errors: errorData.errors,
        timestamp: new Date().toISOString(),
        path: response.url,
      } as ErrorResponse;
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

    // Create abort controller for timeout
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);
    
    if (signal) {
      signal.addEventListener('abort', () => controller.abort());
    }

    try {
      const response = await fetch(url, {
        ...requestConfig,
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