import React, { useState } from 'react';
import { useContentGeneration } from '../../hooks/useContentGeneration';
import { 
  Box, 
  Button, 
  TextField, 
  Select, 
  MenuItem, 
  FormControl, 
  InputLabel,
  Card,
  CardContent,
  CardActions,
  Typography,
  CircularProgress,
  Alert,
  Chip,
  Grid,
  Paper,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  ContentCopy,
  Download,
  Save,
  History,
  Refresh,
} from '@mui/icons-material';

const AIContentGenerator: React.FC = () => {
  const [topic, setTopic] = useState('');
  const [contentType, setContentType] = useState<'blog_post' | 'social_media' | 'email' | 'ad_copy' | 'product_description'>('blog_post');
  const [language, setLanguage] = useState<'ar' | 'en'>('ar');
  const [tone, setTone] = useState<'professional' | 'casual' | 'persuasive' | 'informative'>('professional');
  const [length, setLength] = useState<'short' | 'medium' | 'long'>('medium');
  
  const {
    generatedContent,
    isGenerating,
    history,
    generateAndSave,
    editContent,
    saveContent,
    loadFromHistory,
    exportContent,
    copyToClipboard,
    clearContent,
  } = useContentGeneration();
  
  const handleGenerate = async () => {
    if (!topic.trim()) return;
    
    await generateAndSave(contentType, topic, {
      language,
      tone,
      length,
    });
  };
  
  const handleSave = () => {
    const title = `${contentType} - ${topic}`;
    saveContent(title, [contentType, language]);
  };
  
  const handleExport = (format: 'txt' | 'md' | 'html') => {
    exportContent(format);
  };
  
  return (
    <Box sx={{ maxWidth: 1200, margin: '0 auto', p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ color: '#7A3EF0', mb: 4 }}>
        🚀 NawthTech AI Content Generator
      </Typography>
      
      <Grid container spacing={3}>
        {/* لوحة التحكم */}
        <Grid item xs={12} md={4}>
          <Card sx={{ mb: 3, bgcolor: 'background.paper' }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                إعدادات المحتوى
              </Typography>
              
              <TextField
                fullWidth
                label="الموضوع أو الفكرة"
                value={topic}
                onChange={(e) => setTopic(e.target.value)}
                sx={{ mb: 2 }}
                placeholder="مثال: استراتيجيات النمو الرقمي للشركات الناشئة"
              />
              
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>نوع المحتوى</InputLabel>
                <Select
                  value={contentType}
                  label="نوع المحتوى"
                  onChange={(e) => setContentType(e.target.value as any)}
                >
                  <MenuItem value="blog_post">مقال مدونة</MenuItem>
                  <MenuItem value="social_media">منشور وسائط اجتماعية</MenuItem>
                  <MenuItem value="email">نص بريد إلكتروني</MenuItem>
                  <MenuItem value="ad_copy">نص إعلاني</MenuItem>
                  <MenuItem value="product_description">وصف منتج</MenuItem>
                </Select>
              </FormControl>
              
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>اللغة</InputLabel>
                <Select
                  value={language}
                  label="اللغة"
                  onChange={(e) => setLanguage(e.target.value as any)}
                >
                  <MenuItem value="ar">العربية</MenuItem>
                  <MenuItem value="en">English</MenuItem>
                </Select>
              </FormControl>
              
              <FormControl fullWidth sx={{ mb: 3 }}>
                <InputLabel>النبرة</InputLabel>
                <Select
                  value={tone}
                  label="النبرة"
                  onChange={(e) => setTone(e.target.value as any)}
                >
                  <MenuItem value="professional">مهنية</MenuItem>
                  <MenuItem value="casual">غير رسمية</MenuItem>
                  <MenuItem value="persuasive">إقناعية</MenuItem>
                  <MenuItem value="informative">إعلامية</MenuItem>
                </Select>
              </FormControl>
              
              <Button
                fullWidth
                variant="contained"
                onClick={handleGenerate}
                disabled={isGenerating || !topic.trim()}
                sx={{
                  bgcolor: '#7A3EF0',
                  '&:hover': { bgcolor: '#6A2EE0' },
                  mb: 2,
                }}
                startIcon={isGenerating ? <CircularProgress size={20} color="inherit" /> : <Refresh />}
              >
                {isGenerating ? 'جاري التوليد...' : 'توليد المحتوى'}
              </Button>
              
              <Alert severity="info" sx={{ mt: 2 }}>
                يستخدم الذكاء الاصطناعي لتوليد محتوى فريد ومخصص لعلامتك التجارية
              </Alert>
            </CardContent>
          </Card>
          
          {/* التاريخ */}
          {history.length > 0 && (
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  المحتوى السابق
                </Typography>
                {history.slice(0, 5).map((item) => (
                  <Paper
                    key={item.id}
                    sx={{
                      p: 2,
                      mb: 1,
                      cursor: 'pointer',
                      '&:hover': { bgcolor: 'action.hover' },
                    }}
                    onClick={() => loadFromHistory(item.id)}
                  >
                    <Typography variant="body2" noWrap>
                      {item.content.substring(0, 50)}...
                    </Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                      <Chip label={item.type} size="small" />
                      <Typography variant="caption" color="text.secondary">
                        {new Date(item.timestamp).toLocaleTimeString('ar-EG')}
                      </Typography>
                    </Box>
                  </Paper>
                ))}
              </CardContent>
            </Card>
          )}
        </Grid>
        
        {/* محرر المحتوى */}
        <Grid item xs={12} md={8}>
          <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <CardContent sx={{ flexGrow: 1, overflow: 'auto' }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                <Typography variant="h6">
                  المحتوى المولد
                </Typography>
                
                <Box>
                  <Tooltip title="نسخ">
                    <IconButton onClick={copyToClipboard} disabled={!generatedContent}>
                      <ContentCopy />
                    </IconButton>
                  </Tooltip>
                  
                  <Tooltip title="حفظ">
                    <IconButton onClick={handleSave} disabled={!generatedContent}>
                      <Save />
                    </IconButton>
                  </Tooltip>
                  
                  <Tooltip title="مسح">
                    <IconButton onClick={clearContent} disabled={!generatedContent}>
                      <History />
                    </IconButton>
                  </Tooltip>
                </Box>
              </Box>
              
              {generatedContent ? (
                <TextField
                  fullWidth
                  multiline
                  rows={20}
                  value={generatedContent}
                  onChange={(e) => editContent(e.target.value)}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      fontFamily: language === 'ar' ? "'Noto Sans Arabic', sans-serif" : 'inherit',
                      direction: language === 'ar' ? 'rtl' : 'ltr',
                      textAlign: language === 'ar' ? 'right' : 'left',
                    },
                  }}
                />
              ) : (
                <Box
                  sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    height: 400,
                    border: '2px dashed #ddd',
                    borderRadius: 1,
                    p: 3,
                  }}
                >
                  <Typography variant="body1" color="text.secondary" gutterBottom>
                    {isGenerating ? 'جاري توليد المحتوى...' : 'لم يتم توليد محتوى بعد'}
                  </Typography>
                  {isGenerating && <CircularProgress sx={{ mt: 2 }} />}
                </Box>
              )}
            </CardContent>
            
            <CardActions sx={{ justifyContent: 'space-between', p: 2 }}>
              <Box>
                <Typography variant="caption" color="text.secondary">
                  الطول: {generatedContent?.length || 0} حرف
                </Typography>
              </Box>
              
              <Box>
                <Button
                  size="small"
                  onClick={() => handleExport('txt')}
                  disabled={!generatedContent}
                  sx={{ mr: 1 }}
                >
                  TXT
                </Button>
                <Button
                  size="small"
                  onClick={() => handleExport('md')}
                  disabled={!generatedContent}
                  sx={{ mr: 1 }}
                >
                  Markdown
                </Button>
                <Button
                  size="small"
                  onClick={() => handleExport('html')}
                  disabled={!generatedContent}
                >
                  HTML
                </Button>
              </Box>
            </CardActions>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default AIContentGenerator;