import { useState, useEffect, useCallback } from 'react';
import { serviceManager } from '../services/api';

export const useServices = () => {
  const [services, setServices] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadServices = useCallback(async () => {
    try {
      setLoading(true);
      const response = await serviceManager.getServiceStatus();
      setServices(response);
      setError(null);
    } catch (err) {
      setError('Failed to load services');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  const initializeServices = useCallback(async () => {
    try {
      setLoading(true);
      await serviceManager.initializeServices();
      await loadServices();
    } catch (err) {
      setError('Failed to initialize services');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [loadServices]);

  const getServiceHealth = useCallback(async () => {
    try {
      const response = await serviceManager.healthCheck();
      return response;
    } catch (err) {
      console.error('Health check failed:', err);
      throw err;
    }
  }, []);

  useEffect(() => {
    loadServices();
    
    // Auto-refresh every 30 seconds
    const interval = setInterval(loadServices, 30000);
    return () => clearInterval(interval);
  }, [loadServices]);

  return {
    services,
    loading,
    error,
    loadServices,
    initializeServices,
    getServiceHealth,
  };
};