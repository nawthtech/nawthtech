import React from 'react';
import { 
  Filter, Brain, Target, TrendingUp, 
  Instagram, Twitter, Youtube, Facebook, 
  Package, Users, MessageSquare, PieChart,
  CheckCircle
} from 'lucide-react';

interface SidebarProps {
  activeCategory: string;
  onCategorySelect: (category: string) => void;
}

export const Sidebar: React.FC<SidebarProps> = ({
  activeCategory,
  onCategorySelect
}) => {
  const categories = [
    { id: 'all', icon: <Package size={18} />, label: 'جميع الخدمات' },
    { id: 'instagram', icon: <Instagram size={18} />, label: 'إنستغرام' },
    { id: 'tiktok', icon: '🎵', label: 'تيك توك' },
    { id: 'twitter', icon: <Twitter size={18} />, label: 'تويتر' },
    { id: 'youtube', icon: <Youtube size={18} />, label: 'يوتيوب' },
    { id: 'facebook', icon: <Facebook size={18} />, label: 'فيسبوك' },
    { id: 'followers', icon: <Users size={18} />, label: 'المتابعين' },
    { id: 'likes', icon: '❤️', label: 'الإعجابات' },
    { id: 'comments', icon: <MessageSquare size={18} />, label: 'التعليقات' },
    { id: 'analytics', icon: <PieChart size={18} />, label: 'التحليلات' },
  ];

  const aiRecommendations = [
    {
      icon: <Target size={18} />,
      title: 'مخصص لك',
      description: 'خدمات تناسب أهدافك',
      tooltip: 'بناءً على سجل مشترياتك'
    },
    {
      icon: <TrendingUp size={18} />,
      title: 'الأكثر طلباً',
      description: 'خدمات شائعة بين العملاء',
      tooltip: 'الأكثر طلباً هذا الأسبوع'
    }
  ];

  return (
    <aside className="store-sidebar">
      {/* رأس الشريط الجانبي */}
      <div className="sidebar-header">
        <div className="logo-icon">
          <svg width="40" height="40" viewBox="0 0 120 40">
            <defs>
              <linearGradient id="sidebarGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style={{ stopColor: '#bc8cff', stopOpacity: 1 }} />
                <stop offset="100%" style={{ stopColor: '#7c3aed', stopOpacity: 1 }} />
              </linearGradient>
            </defs>
            <rect x="5" y="5" width="30" height="30" rx="8" fill="url(#sidebarGradient)" />
            <text x="20" y="25" fontFamily="'Segoe UI', 'Inter', sans-serif" fontWeight="900" fontSize="14" textAnchor="middle" fill="#ffffff">NT</text>
          </svg>
        </div>
        <div className="sidebar-title">
          <span>NawthTech</span> متجر
        </div>
      </div>

      {/* قسم التصنيفات */}
      <div className="sidebar-section">
        <h3>
          <Filter size={14} />
          التصنيفات
        </h3>
        <div className="categories-list">
          {categories.map((category) => (
            <button
              key={category.id}
              className={`category-item ${activeCategory === category.id ? 'active' : ''}`}
              onClick={() => onCategorySelect(category.id)}
              title={category.label}
            >
              {category.icon}
              <span>{category.label}</span>
            </button>
          ))}
        </div>
      </div>

      {/* قسم التوصيات الذكية */}
      <div className="sidebar-section">
        <h3>
          <Brain size={14} />
          التوصيات الذكية
        </h3>
        <div className="ai-recommendations">
          {aiRecommendations.map((rec, index) => (
            <div
              key={index}
              className="recommendation-item"
              title={rec.tooltip}
              role="button"
              tabIndex={0}
            >
              <div className="rec-icon">
                {rec.icon}
              </div>
              <div className="rec-content">
                <div className="rec-title">{rec.title}</div>
                <div className="rec-desc">{rec.description}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* مؤشر حالة النظام */}
      <div className="sidebar-section">
        <div className="status-indicator">
          <span className="status-dot"></span>
          <span>النظام نشط</span>
        </div>
      </div>
    </aside>
  );
};