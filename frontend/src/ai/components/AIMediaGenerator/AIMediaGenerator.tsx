import React, { useState } from 'react';
import { Button, TextField, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel } from '@mui/material';
import { useAI } from '../../hooks/useAI';

const AIMediaGenerator: React.FC = () => {
  const [prompt, setPrompt] = useState('');
  const [style, setStyle] = useState('realistic');
  const { loading, error, generateImage } = useAI();

  const handleGenerate = async () => {
    if (!prompt.trim()) return;
    await generateImage(prompt, style);
  };

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        Generate Media
      </Typography>
      
      <TextField
        fullWidth
        multiline
        rows={3}
        label="Describe the image you want"
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        margin="normal"
      />
      
      <FormControl fullWidth margin="normal">
        <InputLabel>Style</InputLabel>
        <Select
          value={style}
          label="Style"
          onChange={(e) => setStyle(e.target.value)}
        >
          <MenuItem value="realistic">Realistic</MenuItem>
          <MenuItem value="cartoon">Cartoon</MenuItem>
          <MenuItem value="anime">Anime</MenuItem>
          <MenuItem value="painting">Painting</MenuItem>
          <MenuItem value="digital-art">Digital Art</MenuItem>
        </Select>
      </FormControl>
      
      <Button 
        variant="contained" 
        onClick={handleGenerate}
        disabled={loading}
        sx={{ mt: 2 }}
      >
        {loading ? <CircularProgress size={24} /> : 'Generate Image'}
      </Button>
      
      {error && (
        <Typography color="error" sx={{ mt: 2 }}>
          Error: {error}
        </Typography>
      )}
    </Box>
  );
};

export default AIMediaGenerator;
