import { Fragment, useCallback, useEffect, useState } from "react";
import { mc } from "./assets/mc";
import './App.css'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Provider } from 'react-redux';
import { store } from './store';

// مكونات بسيطة مع data-testid للاختبارات
const AIDashboard = () => (
  <div data-testid="ai-dashboard" style={{ padding: '2rem' }}>
    <h2>لوحة تحكم الذكاء الاصطناعي</h2>
    <p>الصفحة الرئيسية لأدوات الذكاء الاصطناعي</p>
  </div>
);

const ContentGenerator = () => (
  <div data-testid="content-generator" style={{ padding: '2rem' }}>
    <h2>مولد المحتوى</h2>
    <p>أداة توليد المحتوى باستخدام الذكاء الاصطناعي</p>
  </div>
);

const MediaStudio = () => (
  <div data-testid="media-studio" style={{ padding: '2rem' }}>
    <h2>استوديو الوسائط</h2>
    <p>أداة إنشاء الوسائط باستخدام الذكاء الاصطناعي</p>
  </div>
);

const StrategyPlanner = () => (
  <div data-testid="strategy-planner" style={{ padding: '2rem' }}>
    <h2>مخطط الاستراتيجيات</h2>
    <p>أداة تخطيط الاستراتيجيات باستخدام الذكاء الاصطناعي</p>
  </div>
);

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
            <Routes>
              <Route path="/" element={<Navigate to="/ai" />} />
              <Route path="/ai" element={<AIDashboard />} />
              <Route path="/ai/content" element={<ContentGenerator />} />
              <Route path="/ai/media" element={<MediaStudio />} />
              <Route path="/ai/strategy" element={<StrategyPlanner />} />
            </Routes>
            
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
                data-testid="quotes-button"
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