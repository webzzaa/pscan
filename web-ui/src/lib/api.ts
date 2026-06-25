const API_BASE = '/api';

export interface ScanRequest {
  host: string;
  ports?: string;
  exclude_hosts?: string;
  exclude_ports?: string;
  scan_mode?: string;
  thread_num?: number;
  timeout?: number;
  module_thread_num?: number;
  disable_ping?: boolean;
  disable_brute?: boolean;
  alive_only?: boolean;
  username?: string;
  password?: string;
  domain?: string;
  poc_path?: string;
  poc_name?: string;
  poc_full?: boolean;
  disable_poc?: boolean;
}

export interface ScanStatus {
  state: 'idle' | 'running' | 'stopping';
  start_time?: string;
  progress: number;
  stats: ScanStats;
}

export interface ScanStats {
  hosts_scanned: number;
  ports_scanned: number;
  services_found: number;
  vulns_found: number;
}

export interface ResultItem {
  id: number;
  time: string;
  type: string;
  target: string;
  status: string;
  details?: Record<string, unknown>;
}

export interface ScanPreset {
  id: string;
  name: string;
  name_en: string;
  description: string;
  description_en: string;
  ports: string;
  scan_mode: string;
  thread_num: number;
  timeout: number;
}

export interface PluginInfo {
  name: string;
  type: string;
  description: string;
  description_en: string;
  enabled: boolean;
}

// API functions
export async function startScan(request: ScanRequest): Promise<{ status: string; start_time: string }> {
  const response = await fetch(`${API_BASE}/scan/start`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to start scan');
  }
  return response.json();
}

export async function stopScan(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/scan/stop`, {
    method: 'POST',
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to stop scan');
  }
  return response.json();
}

export async function getScanStatus(): Promise<ScanStatus> {
  const response = await fetch(`${API_BASE}/scan/status`);
  if (!response.ok) {
    throw new Error('Failed to get scan status');
  }
  return response.json();
}

export async function getResults(type?: string): Promise<{ items: ResultItem[]; total: number; stats: ScanStats }> {
  const url = type ? `${API_BASE}/results?type=${type}` : `${API_BASE}/results`;
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error('Failed to get results');
  }
  return response.json();
}

export async function exportResults(format: 'json' | 'csv'): Promise<Blob> {
  const response = await fetch(`${API_BASE}/results/export?format=${format}`);
  if (!response.ok) {
    throw new Error('Failed to export results');
  }
  return response.blob();
}

export async function clearResults(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/results/clear`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to clear results');
  }
  return response.json();
}

export async function getPresets(): Promise<ScanPreset[]> {
  const response = await fetch(`${API_BASE}/config/presets`);
  if (!response.ok) {
    throw new Error('Failed to get presets');
  }
  return response.json();
}

export async function getPlugins(): Promise<PluginInfo[]> {
  const response = await fetch(`${API_BASE}/config/plugins`);
  if (!response.ok) {
    throw new Error('Failed to get plugins');
  }
  return response.json();
}
