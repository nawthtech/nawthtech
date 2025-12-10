import React from 'react';
import { Mail, Globe, Instagram, Twitter, Linkedin, Github, Youtube, Facebook, Heart } from 'lucide-react';

export const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();
  
  const socialLinks = [
    { icon: <Instagram size={20} />, url: 'https://instagram.com/nawthtech', label: 'Instagram' },
    { icon: <Twitter size={20} />, url: 'https://twitter.com/nawthtech', label: 'Twitter' },
    { icon: <Linkedin size={20} />, url: 'https://linkedin.com/company/nawthtech', label: 'LinkedIn' },
    { icon: <Github size={20} />, url: 'https://github.com/nawthtech', label: 'GitHub' },
    { icon: <Youtube size={20} />, url: 'https://youtube.com/@nawthtech', label: 'YouTube' },
    { icon: <Facebook size={20} />, url: 'https://facebook.com/nawthtech', label: 'Facebook' },
  ];

  const platformBadges = [
    { platform: 'Instagram', color: '#E4405F' },
    { platform: 'TikTok', color: '#000000' },
    { platform: 'Twitter', color: '#1DA1F2' },
    { platform: 'YouTube', color: '#FF0000' },
    { platform: 'Facebook', color: '#1877F2' },
    { platform: 'Twitch', color: '#9146FF' },
    { platform: 'Telegram', color: '#26A5E4' },
    { platform: 'Discord', color: '#5865F2' },
  ];

  const footerLinks = [
    { label: 'الشروط والأحكام', url: '/terms' },
    { label: 'سياسة الخصوصية', url: '/privacy' },
    { label: 'سياسة الاسترجاع', url: '/refund' },
    { label: 'الدعم الفني', url: '/support' },
    { label: 'الأسئلة الشائعة', url: '/faq' },
    { label: 'عن الشركة', url: '/about' },
  ];

  return (
    <footer className="store-footer">
      {/* قسم الاتصال */}
      <div className="footer-contact">
        <h3 className="footer-title">📞 تواصل معنا</h3>
        <div className="contact-info">
          <div className="contact-item">
            <Mail size={18} />
            <a href="mailto:support@nawthtech.com">support@nawthtech.com</a>
          </div>
          <div className="contact-item">
            <Globe size={18} />
            <a href="https://www.nawthtech.com" target="_blank" rel="noopener noreferrer">
              www.nawthtech.com
            </a>
          </div>
        </div>
      </div>

      {/* روابط التواصل الاجتماعي */}
      <div className="footer-social">
        <h3 className="footer-title">تابعنا على</h3>
        <div className="social-icons">
          {socialLinks.map((social, index) => (
            <a
              key={index}
              href={social.url}
              target="_blank"
              rel="noopener noreferrer"
              className="social-icon"
              aria-label={social.label}
            >
              {social.icon}
            </a>
          ))}
        </div>
      </div>

      {/* المنصات المدعومة */}
      <div className="footer-platforms">
        <h3 className="footer-title">المنصات المدعومة</h3>
        <div className="platform-badges">
          {platformBadges.map((badge, index) => (
            <span
              key={index}
              className="platform-badge"
              style={{ backgroundColor: badge.color }}
            >
              {badge.platform}
            </span>
          ))}
        </div>
      </div>

      {/* روابط سريعة */}
      <div className="footer-links">
        <h3 className="footer-title">روابط سريعة</h3>
        <div className="links-grid">
          {footerLinks.map((link, index) => (
            <a key={index} href={link.url} className="footer-link">
              {link.label}
            </a>
          ))}
        </div>
      </div>

      {/* حقوق النشر */}
      <div className="footer-copyright">
        <div className="copyright-text">
          <p>
            © {currentYear} NawthTech. جميع الحقوق محفوظة.
            <br />
            تم التطوير بـ <Heart size={14} style={{ display: 'inline-block', margin: '0 4px' }} /> في السعودية
          </p>
        </div>
        
        <div className="footer-logo">
          <svg width="100" height="30" viewBox="0 0 120 40">
            <defs>
              <linearGradient id="footerGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style={{ stopColor: '#bc8cff', stopOpacity: 0.8 }} />
                <stop offset="100%" style={{ stopColor: '#7c3aed', stopOpacity: 0.8 }} />
              </linearGradient>
            </defs>
            <text x="0" y="28" fontFamily="'Segoe UI', 'Inter', sans-serif" fontWeight="800" fontSize="20" fill="url(#footerGradient)">
              NawthTech
            </text>
          </svg>
        </div>
      </div>

      {/* رسالة شكر */}
      <div className="footer-message">
        <p className="thank-you-message">
          ✨ شكراً لاختيار NawthTech شريكاً لنموك الرقمي ✨
        </p>
        <div className="ai-badge">
          <span className="ai-label">مدعوم بالذكاء الاصطناعي</span>
        </div>
      </div>
    </footer>
  );
};