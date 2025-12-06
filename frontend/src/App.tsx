import { Fragment, useCallback, useEffect, useState, lazy, Suspense } from "react";
import { mc } from "./assets/mc";
import './App.css'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Provider } from 'react-redux';
import { store } from './store';
import { Box, CircularProgress, Typography } from '@mui/material';

// Theme
const theme = createTheme({
  palette: {
    primary: {
      main: '#7A3EF0',
    },
    secondary: {
      main: '#00F6FF',
    },
    background: {
      default: '#f8fafc',
    },
  },
  typography: {
    fontFamily: [
      'Noto Sans Arabic',
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      '"Helvetica Neue"',
      'Arial',
      'sans-serif',
    ].join(','),
    h1: {
      fontWeight: 700,
    },
    h2: {
      fontWeight: 600,
    },
    h3: {
      fontWeight: 600,
    },
  },
  direction: 'rtl',
});

// مكون تحميل
const LoadingFallback = () => (
  <Box sx={{ 
    display: 'flex', 
    flexDirection: 'column',
    alignItems: 'center', 
    justifyContent: 'center', 
    minHeight: '60vh',
    gap: 2
  }}>
    <CircularProgress size={60} />
    <Typography variant="h6" color="text.secondary">
      جاري التحميل...
    </Typography>
  </Box>
);

// مكونات ديناميكية مع معالجة الأخطاء
const AIDashboardComponent = lazy(() => 
  import('./pages/AIDashboard/AIDashboard')
    .then(module => ({ default: module.default }))
    .catch(() => ({ 
      default: () => (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h5" color="error">
            ⚠️ تعذر تحميل لوحة التحكم
          </Typography>
        </Box>
      )
    }))
);

const ContentGeneratorComponent = lazy(() => 
  import('./pages/ContentGenerator/ContentGenerator')
    .then(module => ({ default: module.default }))
    .catch(() => ({ 
      default: () => (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h5" color="error">
            ⚠️ تعذر تحميل مولد المحتوى
          </Typography>
        </Box>
      )
    }))
);

const MediaStudioComponent = lazy(() => 
  import('./pages/MediaStudio/MediaStudio')
    .then(module => ({ default: module.default }))
    .catch(() => ({ 
      default: () => (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h5" color="error">
            ⚠️ تعذر تحميل استوديو الوسائط
          </Typography>
        </Box>
      )
    }))
);

const StrategyPlannerComponent = lazy(() => 
  import('./pages/StrategyPlanner/StrategyPlanner')
    .then(module => ({ default: module.default }))
    .catch(() => ({ 
      default: () => (
        <Box sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h5" color="error">
            ⚠️ تعذر تحميل مخطط الاستراتيجيات
          </Typography>
        </Box>
      )
    }))
);

function App() {
  const [messages, setMessages] = useState<string[]>([]);
  const [isConnectionOpen, setIsConnectionOpen] = useState(false);

  const onToggleConnection = useCallback(() => {
    setIsConnectionOpen((isOpen) => !isOpen);
  }, []);

  useEffect(() => {
    if (!isConnectionOpen) return;

    const eventSource = new EventSource(import.meta.env.VITE_BACKEND_HOST + "/sse");

    eventSource.onopen = () => {
      console.log("[SSE] Connection established");
    };

    eventSource.onmessage = (event) => {
      setMessages((messages) => [...messages, event.data]);
    };

    eventSource.onerror = (event) => {
      console.error("[SSE] Error:", event);

      if (eventSource.readyState === EventSource.CLOSED) {
        console.log("[SSE] Connection closed because of an error");
        setIsConnectionOpen(false);
      }
    };

    const cleanup = () => {
      console.log("[SSE] Closing connection");
      eventSource.close();
      window.removeEventListener("beforeunload", cleanup);
    };

    window.addEventListener("beforeunload", cleanup);

    return cleanup;
  }, [isConnectionOpen]);

  useEffect(() => {
    window.scrollTo({
      top: document.body.scrollHeight,
      behavior: "smooth",
    });
  }, [messages]);

  return (
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Router>
          <div className="app-container">
            <Suspense fallback={<LoadingFallback />}>
              <Routes>
                <Route path="/" element={<Navigate to="/ai" />} />
                <Route path="/ai" element={<AIDashboardComponent />} />
                <Route path="/ai/content" element={<ContentGeneratorComponent />} />
                <Route path="/ai/media" element={<MediaStudioComponent />} />
                <Route path="/ai/strategy" element={<StrategyPlannerComponent />} />
              </Routes>
            </Suspense>
            
            {/* SSE Quotes Section */}
            <div className="sse-quotes" style={{ display: 'none' }}>
              <h1 className="text-4xl font-semibold">Here's some unnecessary quotes for you to read...</h1>

              {messages.map((message, index, elements) => (
                <Fragment key={index}>
                  <p className={mc("duration-200", index + 1 !== elements.length ? "opacity-40" : "scale-105 font-bold")}>{message}</p>
                </Fragment>
              ))}

              <button 
                className={mc("hover:opacity-75 duration-200 font-bold text-lg", isConnectionOpen ? "text-[#f06b6b]" : "text-[#6bf06b]")} 
                onClick={onToggleConnection}
              >
                {isConnectionOpen ? "Stop" : "Start"} Quotes
              </button>

              <div className="h-96 w-full" />
            </div>
          </div>
        </Router>
      </ThemeProvider>
    </Provider>
  );
}

export default App;