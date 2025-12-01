import React, { useState } from 'react';
import { 
  Eye, 
  EyeOff, 
  LogIn, 
  Mail, 
  Lock, 
  AlertCircle
  // إزالة User و Smartphone لأنهم غير مستخدمين
} from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { api } from '../../services/api';
import { settings, isFeatureEnabled } from '../../config';
// إزالة getApiEndpoint لأنه غير مستخدم
import './AuthForms.css';

// Types
interface LoginFormData {
  email: string;
  password: string;
  rememberMe: boolean;
}

interface LoginErrors {
  email?: string;
  password?: string;
  submit?: string;
}

interface LoginResponse {
  data: {
    token: string;
    user: {
      id: string;
      name: string;
      email: string;
      role: string;
      avatar?: string;
      permissions?: string[];
      subscription?: {
        plan: string;
        expiresAt: string;
        status: string;
      };
    };
    refreshToken?: string;
    expiresIn: number;
  };
}

const LoginForm: React.FC = () => {
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
    rememberMe: false
  });
  
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [errors, setErrors] = useState<LoginErrors>({});
  const [showDemo, setShowDemo] = useState<boolean>(settings.development.debug);
  
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target;
    
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
    
    // Clear error when user starts typing
    if (errors[name as keyof LoginErrors]) {
      setErrors(prev => ({
        ...prev,
        [name]: undefined
      }));
    }
  };

  const validateForm = (): boolean => {
    const newErrors: LoginErrors = {};
    
    // Email validation
    if (!formData.email.trim()) {
      newErrors.email = 'البريد الإلكتروني مطلوب';
    } else if (!/\S+@\S+\.\S+/.test(formData.email) && !/^[0-9+\-\s()]+$/.test(formData.email)) {
      newErrors.email = 'البريد الإلكتروني أو رقم الهاتف غير صالح';
    }
    
    // Password validation
    if (!formData.password) {
      newErrors.password = 'كلمة المرور مطلوبة';
    } else if (formData.password.length < settings.security.password.minLength) {
      newErrors.password = `كلمة المرور يجب أن تكون ${settings.security.password.minLength} أحرف على الأقل`;
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleLoginSuccess = (response: LoginResponse) => {
    const { token, user, refreshToken, expiresIn } = response.data;
    
    // Store auth data based on remember me choice
    const storage = formData.rememberMe ? localStorage : sessionStorage;
    
    storage.setItem('nawthtech_auth_token', token);
    storage.setItem('nawthtech_user', JSON.stringify(user));
    
    if (refreshToken) {
      storage.setItem('nawthtech_refresh_token', refreshToken);
    }
    
    // Set cookie for Go backend (if needed)
    document.cookie = `auth_token=${token}; path=/; max-age=${expiresIn}; secure=${settings.app.environment === 'production'}; SameSite=Strict`;
    
    // Set last login time
    localStorage.setItem('nawthtech_last_login', new Date().toISOString());
    
    // Redirect based on user role
    if (user.role === 'admin') {
      navigate('/admin/dashboard');
    } else {
      navigate('/dashboard');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return;
    
    setLoading(true);
    setErrors({});

    try {
      // Determine if input is email or phone
      const identifier = formData.email.includes('@') 
        ? { email: formData.email }
        : { phone: formData.email.replace(/\D/g, '') };
      
      const response = await api.post<LoginResponse>(
        { category: 'auth', endpoint: 'login' },
        {
          ...identifier,
          password: formData.password,
          rememberMe: formData.rememberMe,
          deviceInfo: {
            userAgent: navigator.userAgent,
            platform: navigator.platform,
            language: navigator.language,
            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
          }
        }
      );

      handleLoginSuccess(response.data);
      
    } catch (error: any) {
      console.error('Login error:', error);
      
      let errorMessage = 'حدث خطأ أثناء تسجيل الدخول. يرجى المحاولة مرة أخرى.';
      
      if (error.status === 401) {
        errorMessage = 'البريد الإلكتروني أو كلمة المرور غير صحيحة';
      } else if (error.status === 429) {
        errorMessage = 'لقد تجاوزت عدد المحاولات المسموح بها. يرجى المحاولة مرة أخرى لاحقاً';
      } else if (error.status === 403) {
        errorMessage = 'حسابك غير مفعل أو محظور. يرجى التواصل مع الدعم';
      } else if (error.message) {
        errorMessage = error.message;
      }
      
      setErrors({ submit: errorMessage });
    } finally {
      setLoading(false);
    }
  };

  const handleSocialLogin = async (provider: string) => {
    try {
      // Get OAuth URL from backend
      const response = await api.get<{ url: string }>(
        { category: 'auth', endpoint: `social/${provider}/url` }
      );
      
      // Redirect to OAuth provider
      window.location.href = response.data.url;
      
    } catch (error: any) {
      setErrors({ 
        submit: `تعذر الاتصال بخدمة ${provider}. يرجى المحاولة لاحقاً` 
      });
    }
  };

  const handleDemoLogin = async () => {
    if (!settings.development.debug) return;
    
    setLoading(true);
    
    try {
      const response = await api.post<LoginResponse>(
        { category: 'auth', endpoint: 'demo' },
        {
          role: 'user',
          deviceInfo: {
            userAgent: navigator.userAgent,
            platform: navigator.platform
          }
        }
      );

      handleLoginSuccess(response.data);
      
    } catch (error: any) {
      setErrors({ submit: 'تعذر تسجيل الدخول التجريبي' });
      setLoading(false);
    }
  };

  // حل مشكلة social platforms
  const socialPlatforms = settings.social.platforms as any;

  return (
    <div className="auth-form-container">
      <div className="auth-header">
        <div className="auth-logo">
          {settings.app.environment === 'development' && (
            <div className="env-badge dev">تطوير</div>
          )}
          <img 
            src="/assets/logo.png" 
            alt={settings.app.name} 
            onError={(e: React.SyntheticEvent<HTMLImageElement, Event>) => {
              e.currentTarget.src = '/assets/logo-default.png';
            }}
          />
          <h2>{settings.app.name}</h2>
          <span className="app-version">v{settings.app.version}</span>
        </div>
        <h1>مرحباً بعودتك!</h1>
        <p>سجل الدخول إلى حسابك للمتابعة</p>
      </div>

      <form onSubmit={handleSubmit} className="auth-form">
        {errors.submit && (
          <div className="error-message">
            <AlertCircle size={18} />
            <span>{errors.submit}</span>
          </div>
        )}

        <div className="form-group">
          <label htmlFor="email" className="form-label">
            <Mail size={18} />
            البريد الإلكتروني أو رقم الهاتف
          </label>
          <input
            type="text"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className={`form-input ${errors.email ? 'error' : ''}`}
            placeholder="أدخل بريدك الإلكتروني أو رقم هاتفك"
            dir="auto"
            autoComplete="email"
            disabled={loading}
          />
          {errors.email && (
            <span className="error-text">{errors.email}</span>
          )}
        </div>

        <div className="form-group">
          <label htmlFor="password" className="form-label">
            <Lock size={18} />
            كلمة المرور
          </label>
          <div className="password-input-container">
            <input
              type={showPassword ? 'text' : 'password'}
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className={`form-input ${errors.password ? 'error' : ''}`}
              placeholder="أدخل كلمة المرور"
              dir="auto"
              autoComplete="current-password"
              disabled={loading}
            />
            <button
              type="button"
              className="password-toggle"
              onClick={() => setShowPassword(!showPassword)}
              disabled={loading}
            >
              {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
          {errors.password && (
            <span className="error-text">{errors.password}</span>
          )}
        </div>

        <div className="form-options">
          <label className="checkbox-label">
            <input
              type="checkbox"
              name="rememberMe"
              checked={formData.rememberMe}
              onChange={handleChange}
              disabled={loading}
            />
            <span className="checkmark"></span>
            تذكرني على هذا الجهاز
          </label>
          
          <Link 
            to="/auth/forgot-password" 
            className="forgot-password"
            onClick={(e: React.MouseEvent) => loading && e.preventDefault()}
          >
            نسيت كلمة المرور؟
          </Link>
        </div>

        <button 
          type="submit" 
          className={`auth-button primary ${loading ? 'loading' : ''}`}
          disabled={loading}
        >
          {loading ? (
            <div className="loading-spinner"></div>
          ) : (
            <>
              <LogIn size={18} />
              تسجيل الدخول
            </>
          )}
        </button>

        {isFeatureEnabled('socialMediaIntegration') && (
          <>
            <div className="auth-divider">
              <span>أو</span>
            </div>

            <div className="social-login">
              {socialPlatforms.google?.enabled && (
                <button 
                  type="button" 
                  className="social-button google"
                  onClick={() => handleSocialLogin('google')}
                  disabled={loading}
                >
                  <img 
                    src="/assets/icons/google.svg" 
                    alt="Google"
                    onError={(e: React.SyntheticEvent<HTMLImageElement, Event>) => {
                      e.currentTarget.src = '/assets/icons/default.svg';
                    }}
                  />
                  متابعة مع جوجل
                </button>
              )}

              {socialPlatforms.apple?.enabled && (
                <button 
                  type="button" 
                  className="social-button apple"
                  onClick={() => handleSocialLogin('apple')}
                  disabled={loading}
                >
                  <img 
                    src="/assets/icons/apple.svg" 
                    alt="Apple"
                    onError={(e: React.SyntheticEvent<HTMLImageElement, Event>) => {
                      e.currentTarget.src = '/assets/icons/default.svg';
                    }}
                  />
                  متابعة مع آبل
                </button>
              )}
            </div>
          </>
        )}

        <div className="auth-footer">
          <p>
            ليس لديك حساب؟{' '}
            <Link 
              to="/auth/register" 
              className="auth-link"
              onClick={(e: React.MouseEvent) => loading && e.preventDefault()}
            >
              إنشاء حساب جديد
            </Link>
          </p>
        </div>
      </form>

      {/* Demo Credentials - Only shown in development */}
      {settings.development.debug && settings.development.mockData && (
        <div className="demo-credentials">
          <div className="demo-header">
            <h4>بيانات تجريبية:</h4>
            <button 
              type="button" 
              className="demo-toggle"
              onClick={() => setShowDemo(!showDemo)}
            >
              {showDemo ? 'إخفاء' : 'إظهار'}
            </button>
          </div>
          
          {showDemo && (
            <>
              <div className="credential-item">
                <strong>البريد:</strong> demo@nawthtech.com
              </div>
              <div className="credential-item">
                <strong>كلمة المرور:</strong> 123456
              </div>
              <button 
                type="button" 
                className="demo-login-button"
                onClick={handleDemoLogin}
                disabled={loading}
              >
                {loading ? 'جاري التسجيل...' : 'تسجيل تجريبي سريع'}
              </button>
            </>
          )}
        </div>
      )}

      {/* Security Notice */}
      <div className="security-notice">
        <Lock size={14} />
        <span>اتصال آمن عبر SSL/TLS</span>
      </div>
    </div>
  );
};

export default LoginForm;