/**
 * Admin API Types
 */

// Re-export main types
export type {
  DashboardStats,
  StoreMetrics,
  Order,
  UserActivity,
  SystemAlert,
  DashboardData,
  AnalyticsFilters,
  ExportOptions,
  AdminSettings,
} from '../admin';

// Additional admin-specific types
export interface AdminUser {
  id: string;
  name: string;
  email: string;
  role: string;
  permissions: string[];
  lastLogin?: string;
  createdAt: string;
  status: 'active' | 'inactive' | 'suspended';
  metadata?: Record<string, any>;
}

export interface AdminNotification {
  id: string;
  type: 'info' | 'warning' | 'error' | 'success';
  title: string;
  message: string;
  read: boolean;
  timestamp: string;
  action?: {
    label: string;
    url: string;
    method?: string;
  };
  priority: 'low' | 'medium' | 'high';
}

export interface SystemStatus {
  database: {
    status: 'connected' | 'disconnected' | 'degraded';
    latency: number;
    connections: number;
  };
  cache: {
    status: 'connected' | 'disconnected';
    hitRate: number;
    memoryUsage: number;
  };
  storage: {
    status: 'connected' | 'disconnected';
    totalSpace: number;
    usedSpace: number;
    availableSpace: number;
  };
  api: {
    status: 'healthy' | 'unhealthy';
    responseTime: number;
    errorRate: number;
  };
  overall: 'healthy' | 'degraded' | 'unhealthy';
}