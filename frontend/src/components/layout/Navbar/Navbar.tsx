import React from 'react';
import { Boxes, Star, Robot, TrendingUp, Zap } from 'lucide-react';

interface NavbarProps {
  activeView: string;
  onViewChange: (view: string) => void;
}

export const Navbar: React.FC<NavbarProps> = ({
  activeView,
  onViewChange
}) => {
  const views = [
    { id: 'all', icon: <Boxes size={16} />, label: 'جميع الخدمات' },
    { id: 'featured', icon: <Star size={16} />, label: 'المميزة' },
    { id: 'ai', icon: <Robot size={16} />, label: 'الذكاء الاصطناعي' },
    { id: 'trending', icon: <TrendingUp size={16} />, label: 'الرائجة' },
    { id: 'new', icon: <Zap size={16} />, label: 'الجديدة' },
  ];

  return (
    <nav className="navbar">
      <div className="views-container">
        {views.map((view) => (
          <button
            key={view.id}
            className={`view-tab ${activeView === view.id ? 'active' : ''}`}
            onClick={() => onViewChange(view.id)}
            aria-label={`عرض ${view.label}`}
          >
            {view.icon}
            <span>{view.label}</span>
          </button>
        ))}
      </div>
    </nav>
  );
};