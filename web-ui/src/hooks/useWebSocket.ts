import { useEffect, useRef, useCallback, useState } from 'react';

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

interface UseWebSocketOptions {
  onMessage?: (message: WSMessage) => void;
  onConnected?: () => void;
  onDisconnected?: () => void;
  reconnectInterval?: number;
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const {
    onMessage,
    onConnected,
    onDisconnected,
    reconnectInterval = 3000,
  } = options;

  const [isConnected, setIsConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      onConnected?.();
    };

    ws.onclose = () => {
      setIsConnected(false);
      onDisconnected?.();
      // Reconnect
      reconnectTimeoutRef.current = setTimeout(connect, reconnectInterval);
    };

    ws.onerror = () => {
      ws.close();
    };

    ws.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data);
        onMessage?.(message);
      } catch {
        console.error('Failed to parse WebSocket message');
      }
    };
  }, [onConnected, onDisconnected, onMessage, reconnectInterval]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    wsRef.current?.close();
    wsRef.current = null;
  }, []);

  const sendMessage = useCallback((type: MessageType, data?: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type, data, timestamp: Date.now() }));
    }
  }, []);

  useEffect(() => {
    connect();
    return () => disconnect();
  }, [connect, disconnect]);

  return {
    isConnected,
    sendMessage,
    connect,
    disconnect,
  };
}
