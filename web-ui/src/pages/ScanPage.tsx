import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { TooltipProvider } from '@/components/ui/tooltip';
import { ScanForm } from '@/components/ScanForm';
import { LiveFeed } from '@/components/LiveFeed';
import { StatsPanel } from '@/components/StatsPanel';
import { startScan, stopScan, getScanStatus, getPresets, type ScanRequest, type ScanStatus, type ScanPreset } from '@/lib/api';
import { useLiveFeed } from '@/contexts/LiveFeedContext';

const DEFAULT_FORM: ScanRequest = {
  host: '',
  ports: '',
  scan_mode: 'all',
  thread_num: 600,
  timeout: 3,
  disable_ping: false,
  disable_brute: false,
  alive_only: false,
  username: '',
  password: '',
  domain: '',
  exclude_hosts: '',
  exclude_ports: '',
};

export function ScanPage() {
  const { t } = useTranslation();
  const { clearLogs } = useLiveFeed();
  const [status, setStatus] = useState<ScanStatus | null>(null);
  const [presets, setPresets] = useState<ScanPreset[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<ScanRequest>(DEFAULT_FORM);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [statusData, presetsData] = await Promise.all([
          getScanStatus(),
          getPresets(),
        ]);
        setStatus(statusData);
        setPresets(presetsData);
      } catch (err) {
        console.error('Failed to fetch data:', err);
      }
    };
    fetchData();

    const interval = setInterval(async () => {
      try {
        const statusData = await getScanStatus();
        setStatus(statusData);
      } catch {
        // ignore
      }
    }, 2000);

    return () => clearInterval(interval);
  }, []);

  const handleStart = async () => {
    if (!formData.host) {
      setError(t('targetRequired'));
      return;
    }
    setLoading(true);
    setError(null);
    clearLogs();
    try {
      await startScan(formData);
    } catch (err) {
      setError(err instanceof Error ? err.message : t('startScanFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleStop = async () => {
    setLoading(true);
    try {
      await stopScan();
    } catch (err) {
      setError(err instanceof Error ? err.message : t('stopScanFailed'));
    } finally {
      setLoading(false);
    }
  };

  const isRunning = status?.state === 'running';
  const isStopping = status?.state === 'stopping';

  return (
    <TooltipProvider>
      <div className="grid grid-cols-1 lg:grid-cols-10 gap-4 h-full">
        {/* Left Panel - 7 cols */}
        <div className="lg:col-span-7 flex flex-col gap-4 min-h-0">
          <ScanForm
            formData={formData}
            onFormChange={setFormData}
            presets={presets}
            isRunning={isRunning}
            isStopping={isStopping}
            loading={loading}
            error={error}
            onStart={handleStart}
            onStop={handleStop}
          />
          <LiveFeed />
        </div>

        {/* Right Panel - 3 cols */}
        <div className="lg:col-span-3">
          <StatsPanel status={status} />
        </div>
      </div>
    </TooltipProvider>
  );
}
