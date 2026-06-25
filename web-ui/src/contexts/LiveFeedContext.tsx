import { createContext, useContext, useState, useCallback, useEffect, useRef, useMemo, type ReactNode } from 'react';

export type MessageType =
  | 'scan_started'
  | 'scan_progress'
  | 'scan_result'
  | 'scan_completed'
  | 'scan_error'
  | 'connected'
  | 'ping'
  | 'pong';

export interface WSMessage {
  type: MessageType;
  timestamp: number;
  data?: unknown;
}

export interface LogEntry {
  id: number;
  time: string;
  type: string;
  target: string;
  status: string;
}

interface LiveFeedContextType {
  isConnected: boolean;
  logs: LogEntry[];
  clearLogs: () => void;
}

const LiveFeedContext = createContext<LiveFeedContextType | null>(null);

export function LiveFeedProvider({ children }: { children: ReactNode }) {
  const [isConnected, setIsConnected] = useState(false);
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const mountedRef = useRef(false);
  // 使用 Map 做可靠去重: key = type|target, value = LogEntry
  const logsMapRef = useRef<Map<string, LogEntry>>(new Map());
  // 递增计数器确保 id 唯一
  const counterRef = useRef(0);

  const clearLogs = useCallback(() => {
    logsMapRef.current.clear();
    counterRef.current = 0;
    setLogs([]);
  }, []);

  const connect = useCallback(() => {
    // 防止重复连接
    if (!mountedRef.current) return;
    if (wsRef.current) {
      const state = wsRef.current.readyState;
      if (state === WebSocket.OPEN || state === WebSocket.CONNECTING) {
        return;
      }
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      if (mountedRef.current) {
        setIsConnected(true);
      }
    };

    ws.onclose = () => {
      if (mountedRef.current) {
        setIsConnected(false);
        // 延迟重连
        reconnectTimeoutRef.current = setTimeout(() => {
          if (mountedRef.current) connect();
        }, 3000);
      }
    };

    ws.onerror = () => {
      ws.close();
    };

    ws.onmessage = (event) => {
      if (!mountedRef.current) return;
      try {
        const message: WSMessage = JSON.parse(event.data);
        if (message.type === 'scan_result' && message.data) {
          const data = message.data as Record<string, string>;
          // 使用服务端时间戳（如果有）格式化时间
          const serverTime = message.timestamp || Date.now();
          const newEntry: LogEntry = {
            id: ++counterRef.current, // 递增确保唯一且有序
            time: new Date(serverTime).toLocaleTimeString(),
            type: data.type || 'info',
            target: data.target || '',
            status: data.status || '',
          };

          const key = `${newEntry.type}|${newEntry.target}`;
          const existing = logsMapRef.current.get(key);

          // 去重逻辑：同一 type|target 只保留一条，优先保留详细信息
          if (existing) {
            const oldStatus = existing.status;
            // 只有新的更详细才更新（保留原始 id 以维持顺序）
            if ((oldStatus === 'identified' || oldStatus === 'open' || oldStatus === '') &&
                newEntry.status !== 'identified' && newEntry.status !== 'open' && newEntry.status !== '') {
              logsMapRef.current.set(key, { ...newEntry, id: existing.id });
            }
            // 已有信息，不重复添加
          } else {
            logsMapRef.current.set(key, newEntry);
          }

          // 从 Map 生成排序后的数组（按 id 倒序，最新在前）
          const sorted = Array.from(logsMapRef.current.values())
            .sort((a, b) => b.id - a.id)
            .slice(0, 100);
          setLogs(sorted);
        }
      } catch {
        // ignore parse errors
      }
    };
  }, []);

  useEffect(() => {
    mountedRef.current = true;
    connect();
    return () => {
      mountedRef.current = false;
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
        reconnectTimeoutRef.current = null;
      }
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [connect]);

  const value = useMemo(() => ({
    isConnected,
    logs,
    clearLogs,
  }), [isConnected, logs, clearLogs]);

  return (
    <LiveFeedContext.Provider value={value}>
      {children}
    </LiveFeedContext.Provider>
  );
}

export function useLiveFeed() {
  const context = useContext(LiveFeedContext);
  if (!context) {
    throw new Error('useLiveFeed must be used within a LiveFeedProvider');
  }
  return context;
}
