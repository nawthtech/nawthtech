import React from 'react';
import { Box, Typography, Container } from '@mui/material';

const StrategyPlanner: React.FC = () => {
  return (
    <Container maxWidth="lg">
      <Box sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom>
          📊 مُخطِّط الاستراتيجيات
        </Typography>
        <Typography variant="body1" color="text.secondary">
          خطط واستراتيجيات ذكية لتحقيق أهدافك الرقمية
        </Typography>
        
        <Box sx={{ mt: 4, p: 3, bgcolor: 'background.paper', borderRadius: 2 }}>
          <Typography variant="h6" gutterBottom>
            🔄 قيد التطوير
          </Typography>
          <Typography>
            هذه الصفحة قيد التطوير وسيتم إطلاقها قريباً.
          </Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default StrategyPlanner;