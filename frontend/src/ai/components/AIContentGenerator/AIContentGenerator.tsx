import React, { useState } from 'react';
import { Button, TextField, Box, Typography, CircularProgress } from '@mui/material';
import { useContentGeneration } from '../../../hooks/useContentGeneration';

const AIContentGenerator: React.FC = () => {
  const [prompt, setPrompt] = useState('');
  const [length] = useState<'short' | 'medium' | 'long'>('medium');
  const [tone, setTone] = useState<'professional' | 'casual' | 'persuasive' | 'informative'>('professional');
  const { loading, content, error, generateContent } = useContentGeneration();

  const handleGenerate = async () => {
    if (!prompt.trim()) return;
    await generateContent(prompt, { length, tone });
  };

  const handleItemClick = (item: any) => {
    setPrompt(item);
  };

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Generate Content
      </Typography>
      
      <TextField
        fullWidth
        multiline
        rows={3}
        label="Enter your prompt"
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        margin="normal"
      />
      
      <Button 
        variant="contained" 
        onClick={handleGenerate}
        disabled={loading}
        sx={{ mt: 2 }}
      >
        {loading ? <CircularProgress size={24} /> : 'Generate'}
      </Button>
      
      {error && (
        <Typography color="error" sx={{ mt: 2 }}>
          Error: {error}
        </Typography>
      )}
      
      {content && (
        <Box sx={{ mt: 3, p: 2, bgcolor: 'background.default', borderRadius: 1 }}>
          <Typography variant="subtitle1" gutterBottom>
            Generated Content:
          </Typography>
          <Typography variant="body1">
            {content}
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default AIContentGenerator;
