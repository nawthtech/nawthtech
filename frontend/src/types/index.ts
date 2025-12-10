export interface Service {
  id: string;
  name: string;
  category: string;
  price: number;
  unit: string;
  minQuantity: number;
  maxQuantity: number;
  step: number;
  platform: string;
  featured: boolean;
  aiRecommended: boolean;
  trending: boolean;
  aiPowered: boolean;
  description: string;
  features: string[];
  stats: {
    orders: number;
    successRate: number;
    deliveryTime: string;
  };
}

export interface CartItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  unit: string;
  platform: string;
}

export interface ViewTab {
  id: string;
  icon: React.ReactNode;
  label: string;
}

export interface Category {
  id: string;
  icon: React.ReactNode;
  label: string;
}

export interface AIRecommendation {
  icon: React.ReactNode;
  title: string;
  description: string;
  tooltip: string;
}