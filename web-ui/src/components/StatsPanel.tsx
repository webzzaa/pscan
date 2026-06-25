import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { PieChart, Pie, Cell } from 'recharts';
import { Activity, CheckCircle2, XCircle, Loader2 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { ChartContainer, ChartTooltip, ChartTooltipContent, type ChartConfig } from '@/components/ui/chart';
import { EmptyState } from '@/components/ui/empty-state';
import type { ScanStatus } from '@/lib/api';

interface StatsPanelProps {
  status: ScanStatus | null;
}

const CHART_COLORS = {
  hosts: 'hsl(var(--chart-1))',
  ports: 'hsl(var(--chart-2))',
  services: 'hsl(var(--chart-3))',
  vulns: 'hsl(var(--chart-4))',
} as const;

export function StatsPanel({ status }: StatsPanelProps) {
  const { t } = useTranslation();

  const isRunning = status?.state === 'running';
  const isStopping = status?.state === 'stopping';

  const statsChartData = useMemo(() => {
    if (!status?.stats) return [];
    return [
      { name: t('statsHosts'), value: status.stats.hosts_scanned || 0, fill: CHART_COLORS.hosts },
      { name: t('statsPorts'), value: status.stats.ports_scanned || 0, fill: CHART_COLORS.ports },
      { name: t('statsServices'), value: status.stats.services_found || 0, fill: CHART_COLORS.services },
      { name: t('statsVulns'), value: status.stats.vulns_found || 0, fill: CHART_COLORS.vulns },
    ].filter(d => d.value > 0);
  }, [status?.stats, t]);

  const chartConfig: ChartConfig = {
    hosts: { label: t('statsHosts'), color: CHART_COLORS.hosts },
    ports: { label: t('statsPorts'), color: CHART_COLORS.ports },
    services: { label: t('statsServices'), color: CHART_COLORS.services },
    vulns: { label: t('statsVulns'), color: CHART_COLORS.vulns },
  };

  const totalStats = useMemo(() => {
    if (!status?.stats) return 0;
    return (status.stats.hosts_scanned || 0) +
           (status.stats.ports_scanned || 0) +
           (status.stats.services_found || 0) +
           (status.stats.vulns_found || 0);
  }, [status?.stats]);

  const statsItems = [
    { key: 'hosts', label: t('statsHosts'), value: status?.stats.hosts_scanned || 0, color: CHART_COLORS.hosts },
    { key: 'ports', label: t('statsPorts'), value: status?.stats.ports_scanned || 0, color: CHART_COLORS.ports },
    { key: 'services', label: t('statsServices'), value: status?.stats.services_found || 0, color: CHART_COLORS.services },
    { key: 'vulns', label: t('statsVulns'), value: status?.stats.vulns_found || 0, color: CHART_COLORS.vulns },
  ];

  return (
    <Card className="h-full">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
        <CardTitle className="flex items-center gap-2 text-base">
          <Activity className="w-4 h-4 sm:w-5 sm:h-5 text-muted-foreground" />
          {t('resultsDistribution')}
        </CardTitle>
        <Badge
          variant={isRunning ? 'default' : isStopping ? 'secondary' : 'outline'}
          className="gap-1"
        >
          {isRunning ? <CheckCircle2 className="w-3 h-3" /> :
           isStopping ? <Loader2 className="w-3 h-3 animate-spin" /> :
           <XCircle className="w-3 h-3" />}
          {isRunning ? t('scanRunning') : isStopping ? t('statusStopping') : t('statusIdle')}
        </Badge>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Progress */}
        {isRunning && (
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-muted-foreground">{t('loading')}</span>
              <span className="font-mono text-primary text-lg">{status?.progress || 0}%</span>
            </div>
            <Progress value={status?.progress || 0} className="h-3" />
          </div>
        )}

        {/* Pie Chart */}
        {totalStats > 0 ? (
          <div className="flex justify-center py-4">
            <ChartContainer config={chartConfig} className="h-[180px] w-[180px] aspect-square">
              <PieChart>
                <Pie
                  data={statsChartData}
                  dataKey="value"
                  nameKey="name"
                  innerRadius={45}
                  outerRadius={75}
                  strokeWidth={3}
                  stroke="hsl(var(--background))"
                >
                  {statsChartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.fill} />
                  ))}
                </Pie>
                <ChartTooltip content={<ChartTooltipContent hideLabel />} />
              </PieChart>
            </ChartContainer>
          </div>
        ) : (
          <EmptyState
            icon={Activity}
            title={t('chartEmptyTitle')}
            description={t('chartEmptyDescription')}
            className="py-8"
          />
        )}

        {/* Stats List */}
        <div className="space-y-3">
          {statsItems.map((item) => (
            <div key={item.key} className="flex items-center justify-between p-3 rounded-lg bg-muted/50">
              <div className="flex items-center gap-3">
                <div className="w-4 h-4 rounded" style={{ backgroundColor: item.color }} />
                <span className="text-muted-foreground">{item.label}</span>
              </div>
              <span className="font-mono font-semibold text-lg">{item.value}</span>
            </div>
          ))}
        </div>

        {/* Total */}
        <div className="pt-4 border-t">
          <div className="flex items-center justify-between">
            <span className="text-muted-foreground font-medium">{t('items')}</span>
            <span className="font-mono font-bold text-2xl">{totalStats}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
