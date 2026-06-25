import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { BarChart, Bar, XAxis, YAxis, PieChart, Pie, Cell } from 'recharts';
import { BarChart3 } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ChartContainer, ChartTooltip, ChartTooltipContent, type ChartConfig } from '@/components/ui/chart';
import { EmptyState } from '@/components/ui/empty-state';
import type { ResultItem } from '@/lib/api';

interface ResultsChartProps {
  results: ResultItem[];
}

const CHART_COLORS = {
  host: 'hsl(var(--chart-1))',
  port: 'hsl(var(--chart-2))',
  service: 'hsl(var(--chart-3))',
  vuln: 'hsl(var(--chart-4))',
} as const;

export function ResultsChart({ results }: ResultsChartProps) {
  const { t } = useTranslation();

  const chartData = useMemo(() => {
    const counts = { host: 0, port: 0, service: 0, vuln: 0 };
    results.forEach(r => {
      const type = r.type?.toLowerCase() as keyof typeof counts;
      if (type in counts) counts[type]++;
    });
    return [
      { name: t('typeHost'), value: counts.host, fill: CHART_COLORS.host },
      { name: t('typePort'), value: counts.port, fill: CHART_COLORS.port },
      { name: t('typeService'), value: counts.service, fill: CHART_COLORS.service },
      { name: t('typeVuln'), value: counts.vuln, fill: CHART_COLORS.vuln },
    ];
  }, [results, t]);

  const chartConfig: ChartConfig = {
    host: { label: t('typeHost'), color: CHART_COLORS.host },
    port: { label: t('typePort'), color: CHART_COLORS.port },
    service: { label: t('typeService'), color: CHART_COLORS.service },
    vuln: { label: t('typeVuln'), color: CHART_COLORS.vuln },
  };

  const hasResults = results.length > 0;

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="flex items-center gap-2 text-base">
          <BarChart3 className="w-4 h-4 sm:w-5 sm:h-5 text-muted-foreground" />
          {t('resultsDistribution')}
        </CardTitle>
        <Badge variant="secondary" className="font-mono">{results.length}</Badge>
      </CardHeader>
      <CardContent>
        {hasResults ? (
          <div className="space-y-4">
            {/* Pie Chart */}
            <div className="flex justify-center">
              <ChartContainer config={chartConfig} className="h-[160px] w-[160px] aspect-square">
                <PieChart>
                  <Pie
                    data={chartData.filter(d => d.value > 0)}
                    dataKey="value"
                    nameKey="name"
                    innerRadius={35}
                    outerRadius={60}
                    strokeWidth={2}
                    stroke="hsl(var(--background))"
                  >
                    {chartData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.fill} />
                    ))}
                  </Pie>
                  <ChartTooltip content={<ChartTooltipContent hideLabel />} />
                </PieChart>
              </ChartContainer>
            </div>

            {/* Legend */}
            <div className="space-y-2">
              {chartData.map((item, index) => (
                <div key={index} className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <div
                      className="w-3 h-3 rounded-sm shrink-0"
                      style={{ backgroundColor: item.fill }}
                    />
                    <span className="text-muted-foreground">{item.name}</span>
                  </div>
                  <span className="font-mono font-medium">{item.value}</span>
                </div>
              ))}
            </div>

            {/* Bar Chart */}
            <ChartContainer config={chartConfig} className="h-[120px] w-full">
              <BarChart data={chartData} layout="vertical" margin={{ left: 0, right: 8 }}>
                <XAxis type="number" hide />
                <YAxis
                  type="category"
                  dataKey="name"
                  tickLine={false}
                  axisLine={false}
                  width={50}
                  tick={{ fontSize: 11 }}
                />
                <ChartTooltip content={<ChartTooltipContent />} />
                <Bar dataKey="value" radius={3}>
                  {chartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.fill} />
                  ))}
                </Bar>
              </BarChart>
            </ChartContainer>
          </div>
        ) : (
          <EmptyState
            icon={BarChart3}
            title={t('chartEmptyTitle')}
            description={t('chartEmptyDescription')}
            className="py-6"
          />
        )}
      </CardContent>
    </Card>
  );
}
