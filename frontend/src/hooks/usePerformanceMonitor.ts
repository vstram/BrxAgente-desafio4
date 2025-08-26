import { useEffect, useRef, useState } from 'react';

interface PerformanceMetrics {
  renderTime: number;
  memoryUsage: number;
  componentMounts: number;
  rerenders: number;
  lastRenderTime: Date;
}

interface PerformanceMonitorOptions {
  enabled?: boolean;
  logInterval?: number;
  componentName?: string;
}

export function usePerformanceMonitor(options: PerformanceMonitorOptions = {}) {
  const {
    enabled = true,
    logInterval = 5000, // 5 segundos
    componentName = 'Unknown Component'
  } = options;

  const renderCountRef = useRef(0);
  const mountCountRef = useRef(0);
  const lastRenderTimeRef = useRef<number>();
  const intervalRef = useRef<number>();
  
  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    renderTime: 0,
    memoryUsage: 0,
    componentMounts: 0,
    rerenders: 0,
    lastRenderTime: new Date()
  });

  // Contar renders
  useEffect(() => {
    if (!enabled) return;
    
    const now = performance.now();
    renderCountRef.current++;
    
    if (lastRenderTimeRef.current) {
      const renderTime = now - lastRenderTimeRef.current;
      // console.log(`🎭 [${componentName}] Render #${renderCountRef.current} took ${renderTime.toFixed(2)}ms`);
      
      if (renderTime > 100) {
        // console.warn(`⚠️ [${componentName}] SLOW RENDER: ${renderTime.toFixed(2)}ms`);
      }
    }
    
    lastRenderTimeRef.current = now;
  });

  // Contar montagens
  useEffect(() => {
    if (!enabled) return;
    
    mountCountRef.current++;
    // console.log(`🏗️ [${componentName}] Component mounted (#${mountCountRef.current})`);
    
    if (mountCountRef.current > 1) {
      // console.warn(`⚠️ [${componentName}] REMOUNTED ${mountCountRef.current} times - possible memory leak!`);
    }
    
    return () => {
      // console.log(`🧹 [${componentName}] Component unmounted`);
    };
  }, []);

  // Monitoramento contínuo
  useEffect(() => {
    if (!enabled) return;

    intervalRef.current = setInterval(() => {
      const memoryInfo = (performance as any).memory;
      const memory = memoryInfo ? memoryInfo.usedJSHeapSize / 1024 / 1024 : 0; // MB
      
      const currentMetrics: PerformanceMetrics = {
        renderTime: lastRenderTimeRef.current || 0,
        memoryUsage: memory,
        componentMounts: mountCountRef.current,
        rerenders: renderCountRef.current,
        lastRenderTime: new Date()
      };
      
      setMetrics(currentMetrics);
      
      // console.log(`📊 [${componentName}] Performance Report:`, {
      //   totalRenders: renderCountRef.current,
      //   totalMounts: mountCountRef.current,
      //   memoryUsage: `${memory.toFixed(2)}MB`,
      //   avgRenderTime: lastRenderTimeRef.current ? `${(lastRenderTimeRef.current).toFixed(2)}ms` : 'N/A'
      // });
      
      // Alertas
      if (renderCountRef.current > 50) {
        // console.warn(`⚠️ [${componentName}] HIGH RENDER COUNT: ${renderCountRef.current}`);
      }
      
      if (memory > 100) {
        // console.warn(`⚠️ [${componentName}] HIGH MEMORY USAGE: ${memory.toFixed(2)}MB`);
      }
      
    }, logInterval);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [enabled, logInterval, componentName]);

  return {
    metrics,
    renderCount: renderCountRef.current,
    mountCount: mountCountRef.current
  };
}