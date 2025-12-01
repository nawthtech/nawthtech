/**
 * Admin API Service for NawthTech Dashboard
 * Integrates with Go backend in monorepo
 */

import { api } from './api';
import type { PaginationParams } from './api';
import type { ApiResponse } from './types';

// ==================== TYPES ====================
export interface DashboardStats {
  totalUsers: number;
  totalOrders: number;
  totalRevenue: number;
  activeServices: number;
  pendingOrders: number;
  supportTickets: number;
  conversionRate: number;
  bounceRate: number;
  storeVisits: number;
  newCustomers: number;
  growthRate: number;
  averageOrderValue: number;
  customerSatisfaction: number;
  monthlyRecurringRevenue?: number;
  churnRate?: number;
  lifetimeValue?: number;
}

export interface StoreMetrics {
  totalProducts: number;
  lowStockItems: number;
  storeRevenue: number;
  storeOrders: number;
  averageOrderValue: number;
  topSellingCategory: string;
  customerSatisfaction: number;
  returnRate: number;
  inventoryValue: number;
  bestSellingProduct?: string;
  worstSellingProduct?: string;
  revenueTrend: 'up' | 'down' | 'stable';
}

export interface Order {
  id: string;
  user: {
    id: string;
    name: string;
    email: string;
    avatar?: string;
  };
  service: {
    id: string;
    name: string;
    category: string;
  };
  amount: number;
  status: 'pending' | 'processing' | 'completed' | 'cancelled' | 'refunded';
  date: string;
  type: string;
  category: string;
  paymentMethod?: string;
  paymentStatus?: 'paid' | 'pending' | 'failed';
  notes?: string;
  attachments?: string[];
}

export interface UserActivity {
  id: string;
  user: {
    id: string;
    name: string;
    email: string;
    avatar?: string;
  };
  action: string;
  service?: {
    id: string;
    name: string;
  };
  time: string;
  ip: string;
  type: 'login' | 'logout' | 'purchase' | 'view' | 'update' | 'delete' | 'create';
  details?: Record<string, any>;
  userAgent?: string;
  location?: {
    city?: string;
    country?: string;
    coordinates?: {
      lat: number;
      lng: number;
    };
  };
}

export interface SystemAlert {
  id: string;
  type: 'error' | 'warning' | 'info' | 'success';
  title: string;
  message: string;
  timestamp: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  resolved: boolean;
  actionRequired: boolean;
  service?: string;
  metadata?: Record<string, any>;
}

export interface UserReport {
  id: string;
  name: string;
  email: string;
  role: string;
  status: 'active' | 'inactive' | 'suspended' | 'pending';
  joinedDate: string;
  lastLogin: string;
  totalOrders: number;
  totalSpent: number;
  subscriptionPlan?: string;
  tags?: string[];
}

export interface RevenueReport {
  period: string;
  revenue: number;
  orders: number;
  averageOrderValue: number;
  growth: number;
  expenses?: number;
  profit?: number;
  margin?: number;
}

export interface DashboardData {
  stats: DashboardStats;
  storeMetrics: StoreMetrics;
  recentOrders: Order[];
  userActivity: UserActivity[];
  systemAlerts: SystemAlert[];
  revenueTrend: RevenueReport[];
  topUsers: UserReport[];
  performanceMetrics?: {
    responseTime: number;
    uptime: number;
    errorRate: number;
    serverLoad: number;
  };
}

export interface AnalyticsFilters {
  timeRange: 'today' | 'yesterday' | 'week' | 'month' | 'quarter' | 'year' | 'custom';
  startDate?: string;
  endDate?: string;
  category?: string;
  service?: string;
  userGroup?: string;
  status?: string;
}

export interface ExportOptions {
  format: 'csv' | 'excel' | 'pdf' | 'json';
  includeCharts: boolean;
  includeDetails: boolean;
  timeZone: string;
  language: string;
}

export interface AdminSettings {
  siteMaintenance: boolean;
  registrationEnabled: boolean;
  emailNotifications: boolean;
  autoBackup: boolean;
  backupFrequency: 'daily' | 'weekly' | 'monthly';
  maxLoginAttempts: number;
  sessionTimeout: number;
  apiRateLimit: number;
  cacheEnabled: boolean;
  cacheDuration: number;
  securityLevel: 'low' | 'medium' | 'high' | 'strict';
}

// ==================== ADMIN API SERVICE ====================
export const adminAPI = {
  // ==================== DASHBOARD ====================
  getDashboardData: async (
    filters: AnalyticsFilters = { timeRange: 'month' }
  ): Promise<ApiResponse<DashboardData>> => {
    return api.get<DashboardData>(
      { category: 'admin', endpoint: 'dashboard' },
      {
        params: filters,
      }
    );
  },

  getDashboardStats: async (
    filters: AnalyticsFilters = { timeRange: 'month' }
  ): Promise<ApiResponse<DashboardStats>> => {
    return api.get<DashboardStats>(
      { category: 'admin', endpoint: 'dashboard/stats' },
      {
        params: filters,
      }
    );
  },

  getStoreMetrics: async (
    filters: AnalyticsFilters = { timeRange: 'month' }
  ): Promise<ApiResponse<StoreMetrics>> => {
    return api.get<StoreMetrics>(
      { category: 'admin', endpoint: 'dashboard/store-metrics' },
      {
        params: filters,
      }
    );
  },

  // ==================== ORDERS ====================
  getRecentOrders: async (
    limit: number = 10,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<Order[]>> => {
    return api.get<Order[]>(
      { category: 'admin', endpoint: 'orders/recent' },
      {
        params: {
          limit,
          ...filters,
        },
      }
    );
  },

  getAllOrders: async (
    params: PaginationParams & {
      status?: string;
      dateFrom?: string;
      dateTo?: string;
      userId?: string;
      serviceId?: string;
    } = {}
  ): Promise<ApiResponse<Order[]>> => {
    return api.getPaginated<Order>(
      { category: 'admin', endpoint: 'orders' },
      params
    );
  },

  getOrderDetails: async (orderId: string): Promise<ApiResponse<Order>> => {
    return api.get<Order>(
      { category: 'admin', endpoint: `orders/${orderId}` }
    );
  },

  updateOrderStatus: async (
    orderId: string,
    status: Order['status'],
    notes?: string
  ): Promise<ApiResponse<Order>> => {
    return api.put<Order>(
      { category: 'admin', endpoint: `orders/${orderId}/status` },
      { status, notes }
    );
  },

  refundOrder: async (
    orderId: string,
    amount?: number,
    reason?: string
  ): Promise<ApiResponse<Order>> => {
    return api.post<Order>(
      { category: 'admin', endpoint: `orders/${orderId}/refund` },
      { amount, reason }
    );
  },

  // ==================== USERS ====================
  getUserActivity: async (
    limit: number = 10,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<UserActivity[]>> => {
    return api.get<UserActivity[]>(
      { category: 'admin', endpoint: 'users/activity' },
      {
        params: {
          limit,
          ...filters,
        },
      }
    );
  },

  getAllUsers: async (
    params: PaginationParams & {
      status?: string;
      role?: string;
      dateFrom?: string;
      dateTo?: string;
    } = {}
  ): Promise<ApiResponse<UserReport[]>> => {
    return api.getPaginated<UserReport>(
      { category: 'admin', endpoint: 'users' },
      params
    );
  },

  getUserDetails: async (userId: string): Promise<ApiResponse<UserReport>> => {
    return api.get<UserReport>(
      { category: 'admin', endpoint: `users/${userId}` }
    );
  },

  updateUserStatus: async (
    userId: string,
    status: UserReport['status'],
    reason?: string
  ): Promise<ApiResponse<UserReport>> => {
    return api.put<UserReport>(
      { category: 'admin', endpoint: `users/${userId}/status` },
      { status, reason }
    );
  },

  updateUserRole: async (
    userId: string,
    role: string,
    permissions?: string[]
  ): Promise<ApiResponse<UserReport>> => {
    return api.put<UserReport>(
      { category: 'admin', endpoint: `users/${userId}/role` },
      { role, permissions }
    );
  },

  impersonateUser: async (userId: string): Promise<ApiResponse<{ token: string }>> => {
    return api.post<{ token: string }>(
      { category: 'admin', endpoint: `users/${userId}/impersonate` }
    );
  },

  // ==================== ANALYTICS & REPORTS ====================
  getRevenueReport: async (
    filters: AnalyticsFilters
  ): Promise<ApiResponse<RevenueReport[]>> => {
    return api.get<RevenueReport[]>(
      { category: 'admin', endpoint: 'analytics/revenue' },
      {
        params: filters,
      }
    );
  },

  getUserGrowthReport: async (
    filters: AnalyticsFilters
  ): Promise<ApiResponse<{ period: string; newUsers: number; activeUsers: number; churnedUsers: number }[]>> => {
    return api.get(
      { category: 'admin', endpoint: 'analytics/user-growth' },
      {
        params: filters,
      }
    );
  },

  getServiceAnalytics: async (
    serviceId?: string,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<any>> => {
    return api.get(
      { category: 'admin', endpoint: 'analytics/services' },
      {
        params: {
          serviceId,
          ...filters,
        },
      }
    );
  },

  // ==================== SYSTEM & ALERTS ====================
  getSystemAlerts: async (
    params: PaginationParams & {
      severity?: SystemAlert['severity'];
      resolved?: boolean;
      type?: SystemAlert['type'];
    } = {}
  ): Promise<ApiResponse<SystemAlert[]>> => {
    return api.getPaginated<SystemAlert>(
      { category: 'admin', endpoint: 'system/alerts' },
      params
    );
  },

  resolveAlert: async (alertId: string, notes?: string): Promise<ApiResponse<SystemAlert>> => {
    return api.put<SystemAlert>(
      { category: 'admin', endpoint: `system/alerts/${alertId}/resolve` },
      { notes }
    );
  },

  acknowledgeAlert: async (alertId: string): Promise<ApiResponse<SystemAlert>> => {
    return api.put<SystemAlert>(
      { category: 'admin', endpoint: `system/alerts/${alertId}/acknowledge` }
    );
  },

  getSystemMetrics: async (): Promise<ApiResponse<{
    cpuUsage: number;
    memoryUsage: number;
    diskUsage: number;
    activeConnections: number;
    responseTime: number;
    uptime: number;
  }>> => {
    return api.get(
      { category: 'admin', endpoint: 'system/metrics' }
    );
  },

  // ==================== SETTINGS ====================
  getAdminSettings: async (): Promise<ApiResponse<AdminSettings>> => {
    return api.get<AdminSettings>(
      { category: 'admin', endpoint: 'settings' }
    );
  },

  updateAdminSettings: async (
    settings: Partial<AdminSettings>
  ): Promise<ApiResponse<AdminSettings>> => {
    return api.put<AdminSettings>(
      { category: 'admin', endpoint: 'settings' },
      settings
    );
  },

  // ==================== EXPORT & BACKUP ====================
  exportReport: async (
    type: 'orders' | 'users' | 'revenue' | 'analytics' | 'all',
    options: ExportOptions,
    filters?: AnalyticsFilters
  ): Promise<ApiResponse<{ url: string; filename: string; expiresAt: string }>> => {
    return api.post<{ url: string; filename: string; expiresAt: string }>(
      { category: 'admin', endpoint: 'export' },
      {
        type,
        options,
        filters,
      }
    );
  },

  createBackup: async (): Promise<ApiResponse<{ backupId: string; createdAt: string; size: number }>> => {
    return api.post<{ backupId: string; createdAt: string; size: number }>(
      { category: 'admin', endpoint: 'backup' }
    );
  },

  restoreBackup: async (backupId: string): Promise<ApiResponse<{ message: string }>> => {
    return api.post<{ message: string }>(
      { category: 'admin', endpoint: `backup/${backupId}/restore` }
    );
  },

  // ==================== BULK OPERATIONS ====================
  bulkUpdateOrders: async (
    orderIds: string[],
    updates: Partial<{
      status: Order['status'];
      category: string;
      assignedTo: string;
    }>
  ): Promise<ApiResponse<{ updated: number; failed: number; details: any[] }>> => {
    return api.post<{ updated: number; failed: number; details: any[] }>(
      { category: 'admin', endpoint: 'orders/bulk-update' },
      { orderIds, updates }
    );
  },

  bulkUpdateUsers: async (
    userIds: string[],
    updates: Partial<{
      status: UserReport['status'];
      role: string;
      subscriptionPlan: string;
    }>
  ): Promise<ApiResponse<{ updated: number; failed: number; details: any[] }>> => {
    return api.post<{ updated: number; failed: number; details: any[] }>(
      { category: 'admin', endpoint: 'users/bulk-update' },
      { userIds, updates }
    );
  },

  sendBulkNotifications: async (
    userIds: string[],
    notification: {
      title: string;
      message: string;
      type: 'email' | 'push' | 'both';
      data?: Record<string, any>;
    }
  ): Promise<ApiResponse<{ sent: number; failed: number }>> => {
    return api.post<{ sent: number; failed: number }>(
      { category: 'admin', endpoint: 'notifications/bulk' },
      { userIds, notification }
    );
  },

  // ==================== UTILITIES ====================
  clearCache: async (cacheType?: 'all' | 'data' | 'images' | 'api'): Promise<ApiResponse<{ cleared: string[] }>> => {
    return api.post<{ cleared: string[] }>(
      { category: 'admin', endpoint: 'cache/clear' },
      { cacheType }
    );
  },

  sendTestEmail: async (
    email: string,
    template?: string
  ): Promise<ApiResponse<{ message: string }>> => {
    return api.post<{ message: string }>(
      { category: 'admin', endpoint: 'test/email' },
      { email, template }
    );
  },

  checkSystemHealth: async (): Promise<ApiResponse<{
    status: 'healthy' | 'degraded' | 'unhealthy';
    components: {
      database: boolean;
      redis: boolean;
      storage: boolean;
      email: boolean;
      api: boolean;
      auth: boolean;
    };
    issues: Array<{
      component: string;
      issue: string;
      severity: string;
    }>;
    lastCheck: string;
  }>> => {
    return api.get(
      { category: 'admin', endpoint: 'health' }
    );
  },

  // ==================== LOGS & AUDIT ====================
  getAuditLogs: async (
    params: PaginationParams & {
      userId?: string;
      action?: string;
      startDate?: string;
      endDate?: string;
      ip?: string;
    } = {}
  ): Promise<ApiResponse<UserActivity[]>> => {
    return api.getPaginated<UserActivity>(
      { category: 'admin', endpoint: 'audit-logs' },
      params
    );
  },

  getErrorLogs: async (
    params: PaginationParams & {
      level?: 'error' | 'warning' | 'info';
      startDate?: string;
      endDate?: string;
      service?: string;
    } = {}
  ): Promise<ApiResponse<Array<{
    id: string;
    timestamp: string;
    level: string;
    message: string;
    service: string;
    stackTrace?: string;
    userId?: string;
    ip?: string;
  }>>> => {
    return api.getPaginated(
      { category: 'admin', endpoint: 'error-logs' },
      params
    );
  },
};

// ==================== HELPER FUNCTIONS ====================
export const adminHelpers = {
  /**
   * Format dashboard data for charts
   */
  formatChartData: (
    data: DashboardData,
    chartType: 'line' | 'bar' | 'pie' | 'donut'
  ): any[] => {
    switch (chartType) {
      case 'line':
        return data.revenueTrend.map(item => ({
          period: item.period,
          revenue: item.revenue,
          orders: item.orders,
        }));
      case 'bar':
        return [
          { name: 'Users', value: data.stats.totalUsers },
          { name: 'Orders', value: data.stats.totalOrders },
          { name: 'Revenue', value: data.stats.totalRevenue },
          { name: 'Services', value: data.stats.activeServices },
        ];
      case 'pie':
        return Object.entries(data.systemAlerts.reduce((acc, alert) => {
          acc[alert.type] = (acc[alert.type] || 0) + 1;
          return acc;
        }, {} as Record<string, number>)).map(([type, count]) => ({ type, count }));
      default:
        return [];
    }
  },

  /**
   * Calculate dashboard metrics changes
   */
  calculateMetricsChange: (
    current: DashboardStats,
    previous?: DashboardStats
  ): Record<string, { value: number; change: number; trend: 'up' | 'down' | 'stable' }> => {
    if (!previous) return {};

    const changes: Record<string, { value: number; change: number; trend: 'up' | 'down' | 'stable' }> = {};
    const keys = Object.keys(current) as Array<keyof DashboardStats>;

    keys.forEach(key => {
      const currentVal = current[key] as number;
      const previousVal = previous[key] as number;

      if (typeof currentVal === 'number' && typeof previousVal === 'number' && previousVal !== 0) {
        const change = ((currentVal - previousVal) / previousVal) * 100;
        changes[key] = {
          value: currentVal,
          change: parseFloat(change.toFixed(2)),
          trend: change > 0 ? 'up' : change < 0 ? 'down' : 'stable',
        };
      }
    });

    return changes;
  },

  /**
   * Generate export filename
   */
  generateExportFilename: (
    type: string,
    timeRange: string = 'all'
  ): string => {
    const timestamp = new Date().toISOString().split('T')[0];
    return `nawthtech-${type}-${timeRange}-${timestamp}`;
  },

  /**
   * Check if admin has permission
   */
  hasPermission: (
    requiredPermission: string,
    userPermissions?: string[]
  ): boolean => {
    if (!userPermissions) return false;
    
    // Check direct permission
    if (userPermissions.includes(requiredPermission)) {
      return true;
    }

    // Check wildcard permissions
    if (userPermissions.includes('*') || userPermissions.includes('admin.*')) {
      return true;
    }

    // Check category permission
    const [category] = requiredPermission.split('.');
    if (userPermissions.includes(`${category}.*`)) {
      return true;
    }

    return false;
  },

  /**
   * Validate admin filters
   */
  validateFilters: (filters: AnalyticsFilters): string[] => {
    const errors: string[] = [];

    if (filters.timeRange === 'custom') {
      if (!filters.startDate || !filters.endDate) {
        errors.push('Custom time range requires both start and end dates');
      } else {
        const start = new Date(filters.startDate);
        const end = new Date(filters.endDate);

        if (start > end) {
          errors.push('Start date must be before end date');
        }

        // Limit custom range to 1 year
        const maxRange = 365 * 24 * 60 * 60 * 1000; // 1 year in milliseconds
        if (end.getTime() - start.getTime() > maxRange) {
          errors.push('Custom time range cannot exceed 1 year');
        }
      }
    }

    return errors;
  },
};

// ==================== DEFAULT EXPORT ====================
export default adminAPI;