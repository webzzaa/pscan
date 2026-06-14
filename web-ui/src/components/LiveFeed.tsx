import { useTranslation } from 'react-i18next';
import {
  Wifi, WifiOff, Server, Network, Shield, AlertTriangle, CircleDot, Inbox, Activity
} from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { EmptyState } from '@/components/ui/empty-state';
import { useLiveFeed } from '@/contexts/LiveFeedContext';

interface LiveFeedProps {
  compact?: boolean;
  showTypeLabel?: boolean;
}

const TYPE_ICONS = {
  host: Server,
  port: Network,
  service: Shield,
  vuln: AlertTriangle,
} as const;

export function LiveFeed({ compact = false, showTypeLabel = false }: LiveFeedProps) {
  const { t } = useTranslation();
  const { isConnected, logs } = useLiveFeed();

  const getTypeIcon = (type: string) => {
    const key = type?.toLowerCase() as keyof typeof TYPE_ICONS;
    return TYPE_ICONS[key] || CircleDot;
  };

  const getTypeLabel = (type: string) => {
    switch (type?.toLowerCase()) {
      case 'host': return t('typeHost');
      case 'port': return t('typePort');
      case 'service': return t('typeService');
      case 'vuln': return t('typeVuln');
      default: return type;
    }
  };

  return (
    <Card className={compact ? '' : 'flex-1 flex flex-col'}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
        <CardTitle className="flex items-center gap-2 text-base">
          <Activity className="w-4 h-4 sm:w-5 sm:h-5 text-muted-foreground" />
          {t('liveFeed')}
        </CardTitle>
        <div className="flex items-center gap-2">
          {isConnected ? (
            <Badge variant="default" className="gap-1">
              <Wifi className="w-3 h-3" />
              <span className="hidden sm:inline">{t('liveFeedConnected')}</span>
            </Badge>
          ) : (
            <Badge variant="destructive" className="gap-1">
              <WifiOff className="w-3 h-3" />
              <span className="hidden sm:inline">{t('liveFeedDisconnected')}</span>
            </Badge>
          )}
          <Badge variant="outline" className="font-mono">{logs.length}/100</Badge>
        </div>
      </CardHeader>
      <CardContent className={compact ? 'pt-0' : 'pt-0 flex-1 min-h-0'}>
        <ScrollArea className={compact ? 'h-52 lg:h-56' : 'h-full'}>
          {logs.length === 0 ? (
            <EmptyState
              icon={Inbox}
              title={t('resultsEmpty')}
              description={t('liveFeedEmptyDescription')}
              className={compact ? 'py-6' : 'py-8'}
            />
          ) : (
            <div className="space-y-0.5">
              {logs.map((log) => {
                const TypeIcon = getTypeIcon(log.type);
                return (
                  <div key={log.id} className="log-line animate-fade-in group">
                    <span className="log-time">{log.time}</span>
                    <Badge
                      variant={log.type?.toLowerCase() as 'host' | 'port' | 'service' | 'vuln'}
                      className="gap-1 text-xs"
                    >
                      <TypeIcon className="w-3 h-3" />
                      {showTypeLabel && getTypeLabel(log.type)}
                    </Badge>
                    <span className="log-target">{log.target}</span>
                    <span className="text-muted-foreground truncate ml-auto text-xs">
                      {log.status}
                    </span>
                  </div>
                );
              })}
            </div>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  );
}
