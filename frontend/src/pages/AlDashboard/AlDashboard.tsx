import React, { useState, useEffect } from 'react';
import {
  Box,
  Grid,
  Card,
  CardContent,
  Typography,
  LinearProgress,
  Button,
  Chip,
  Avatar,
  Divider,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  AutoAwesome,
  Image,
  TrendingUp,
  Description,
  Download,
  Share,
  Refresh,
  Settings,
} from '@mui/icons-material';
import AIContentGenerator from '../src/ai/components/AIContentGenerator/AIContentGenerator';
import AIMediaGenerator from '../src/ai/components/AIMediaGenerator/AIMediaGenerator';
import { useAI } from '../src/ai/hooks/useAI';

const AIDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'content' | 'media' | 'analysis' | 'strategy'>('content');
  const { getUsage, loading } = useAI();
  const [usage, setUsage] = useState<any>(null);
  
  useEffect(() => {
    loadUsage();
  }, []);
  
  const loadUsage = async () => {
    try {
      const data = await getUsage();
      setUsage(data);
    } catch (error) {
      console.error('Failed to load usage:', error);
    }
  };
  
  const tabs = [
    { id: 'content', label: 'توليد المحتوى', icon: <Description /> },
    { id: 'media', label: 'وسائط AI', icon: <Image /> },
    { id: 'analysis', label: 'التحليلات', icon: <TrendingUp /> },
    { id: 'strategy', label: 'التخطيط', icon: <AutoAwesome /> },
  ];
  
  const renderContent = () => {
    switch (activeTab) {
      case 'content':
        return <AIContentGenerator />;
      case 'media':
        return <AIMediaGenerator />;
      case 'analysis':
        return (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="h5" color="text.secondary">
              صفحة التحليلات - قريباً
            </Typography>
          </Box>
        );
      case 'strategy':
        return (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="h5" color="text.secondary">
              صفحة التخطيط الاستراتيجي - قريباً
            </Typography>
          </Box>
        );
      default:
        return null;
    }
  };
  
  return (
    <Box sx={{ bgcolor: 'background.default', minHeight: '100vh' }}>
      {/* Header */}
      <Box sx={{ bgcolor: '#7A3EF0', color: 'white', p: 3 }}>
        <Grid container alignItems="center" spacing={2}>
          <Grid item>
            <Avatar sx={{ bgcolor: '#00F6FF', width: 56, height: 56 }}>
              <AutoAwesome />
            </Avatar>
          </Grid>
          <Grid item xs>
            <Typography variant="h4" fontWeight="bold">
              NawthTech AI Studio
            </Typography>
            <Typography variant="body1" sx={{ opacity: 0.9 }}>
              منصة الذكاء الاصطناعي المتكاملة للنمو الرقمي
            </Typography>
          </Grid>
          <Grid item>
            <Tooltip title="إعدادات">
              <IconButton sx={{ color: 'white' }}>
                <Settings />
              </IconButton>
            </Tooltip>
          </Grid>
        </Grid>
      </Box>
      
      <Box sx={{ maxWidth: 1400, margin: '0 auto', p: 3 }}>
        <Grid container spacing={3}>
          {/* Sidebar */}
          <Grid item xs={12} md={3}>
            <Card sx={{ mb: 3 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  أدوات الذكاء الاصطناعي
                </Typography>
                <Divider sx={{ my: 2 }} />
                
                {tabs.map((tab) => (
                  <Button
                    key={tab.id}
                    fullWidth
                    startIcon={tab.icon}
                    onClick={() => setActiveTab(tab.id as any)}
                    sx={{
                      justifyContent: 'flex-start',
                      mb: 1,
                      bgcolor: activeTab === tab.id ? 'action.selected' : 'transparent',
                      '&:hover': {
                        bgcolor: 'action.hover',
                      },
                    }}
                  >
                    {tab.label}
                  </Button>
                ))}
              </CardContent>
            </Card>
            
            {/* Usage Stats */}
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  إحصائيات الاستخدام
                </Typography>
                
                {usage ? (
                  <>
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        المحتوى النصي
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={(usage.text_used / usage.text_limit) * 100}
                        sx={{ height: 8, borderRadius: 4 }}
                      />
                      <Typography variant="caption" color="text.secondary">
                        {usage.text_used} / {usage.text_limit} كلمة
                      </Typography>
                    </Box>
                    
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        الصور
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={(usage.images_used / usage.images_limit) * 100}
                        sx={{ height: 8, borderRadius: 4 }}
                      />
                      <Typography variant="caption" color="text.secondary">
                        {usage.images_used} / {usage.images_limit} صورة
                      </Typography>
                    </Box>
                    
                    <Box>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        الفيديوهات
                      </Typography>
                      <LinearProgress 
                        variant="determinate" 
                        value={(usage.videos_used / usage.videos_limit) * 100}
                        sx={{ height: 8, borderRadius: 4 }}
                      />
                      <Typography variant="caption" color="text.secondary">
                        {usage.videos_used} / {usage.videos_limit} فيديو
                      </Typography>
                    </Box>
                  </>
                ) : (
                  <Typography variant="body2" color="text.secondary">
                    جاري تحميل الإحصائيات...
                  </Typography>
                )}
                
                <Button
                  fullWidth
                  variant="outlined"
                  sx={{ mt: 3 }}
                  startIcon={<Refresh />}
                  onClick={loadUsage}
                  disabled={loading}
                >
                  تحديث الإحصائيات
                </Button>
              </CardContent>
            </Card>
            
            {/* Quick Actions */}
            <Card sx={{ mt: 3 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  إجراءات سريعة
                </Typography>
                <Grid container spacing={1}>
                  <Grid item xs={6}>
                    <Button
                      fullWidth
                      variant="contained"
                      size="small"
                      sx={{ bgcolor: '#7A3EF0' }}
                    >
                      مقال جديد
                    </Button>
                  </Grid>
                  <Grid item xs={6}>
                    <Button
                      fullWidth
                      variant="contained"
                      size="small"
                      sx={{ bgcolor: '#00F6FF', color: 'black' }}
                    >
                      تصميم جديد
                    </Button>
                  </Grid>
                  <Grid item xs={6}>
                    <Button
                      fullWidth
                      variant="outlined"
                      size="small"
                    >
                      تحليل سريع
                    </Button>
                  </Grid>
                  <Grid item xs={6}>
                    <Button
                      fullWidth
                      variant="outlined"
                      size="small"
                    >
                      تصدير جميع
                    </Button>
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
          
          {/* Main Content */}
          <Grid item xs={12} md={9}>
            {renderContent()}
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
};

export default AIDashboard;