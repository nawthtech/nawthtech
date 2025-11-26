import api from './api';

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
}

export interface Order {
  id: string;
  user: string;
  service: string;
  amount: number;
  status: string;
  date: string;
  type: string;
  category: string;
}

export interface UserActivity {
  user: string;
  action: string;
  service?: string;
  time: string;
  ip: string;
  type: string;
}

export interface DashboardData {
  stats: DashboardStats;
  storeMetrics: StoreMetrics;
  recentOrders: Order[];
  userActivity: UserActivity[];
  systemAlerts: any[];
}

export const adminAPI = {
  getDashboardData: async (timeRange: string = 'month'): Promise<DashboardData> => {
    const response = await api.get(`/admin/dashboard?timeRange=${timeRange}`);
    return response.data;
  },

  getStoreMetrics: async (timeRange: string = 'month'): Promise<StoreMetrics> => {
    const response = await api.get(`/admin/store-metrics?timeRange=${timeRange}`);
    return response.data;
  },

  getRecentOrders: async (limit: number = 10): Promise<Order[]> => {
    const response = await api.get(`/admin/recent-orders?limit=${limit}`);
    return response.data;
  },

  getUserActivity: async (limit: number = 10): Promise<UserActivity[]> => {
    const response = await api.get(`/admin/user-activity?limit=${limit}`);
    return response.data;
  },

  exportReport: async (type: string, timeRange: string) => {
    const response = await api.get(`/admin/export-report?type=${type}&timeRange=${timeRange}`);
    return response.data;
  },

  updateOrderStatus: async (orderId: string, status: string) => {
    const response = await api.put(`/admin/orders/${orderId}/status`, { status });
    return response.data;
  },
};
