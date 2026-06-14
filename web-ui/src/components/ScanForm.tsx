import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Play, Square, ChevronDown, Loader2, Target, Hash, Clock, Zap, Settings2,
  User, Lock, Globe2, Ban, Activity, Wifi, CheckCircle2, XCircle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import type { ScanRequest, ScanPreset } from '@/lib/api';

interface ScanFormProps {
  formData: ScanRequest;
  onFormChange: (data: ScanRequest) => void;
  presets: ScanPreset[];
  isRunning: boolean;
  isStopping: boolean;
  loading: boolean;
  error: string | null;
  onStart: () => void;
  onStop: () => void;
}

export function ScanForm({
  formData,
  onFormChange,
  presets,
  isRunning,
  isStopping,
  loading,
  error,
  onStart,
  onStop,
}: ScanFormProps) {
  const { t, i18n } = useTranslation();
  const [showAdvanced, setShowAdvanced] = useState(false);

  const updateField = <K extends keyof ScanRequest>(key: K, value: ScanRequest[K]) => {
    onFormChange({ ...formData, [key]: value });
  };

  const applyPreset = (presetId: string) => {
    const preset = presets.find(p => p.id === presetId);
    if (preset) {
      onFormChange({
        ...formData,
        ports: preset.ports,
        scan_mode: preset.scan_mode,
        thread_num: preset.thread_num,
        timeout: preset.timeout,
      });
    }
  };

  const disabled = isRunning;

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <CardTitle className="flex items-center gap-2 text-base">
          <Target className="w-4 h-4 sm:w-5 sm:h-5 text-muted-foreground" />
          {t('scanTitle')}
        </CardTitle>
        {isRunning ? (
          <Button
            size="sm"
            variant="destructive"
            onClick={onStop}
            disabled={loading || isStopping}
            className="gap-2"
          >
            {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Square className="w-4 h-4" />}
            {t('scanStopBtn')}
          </Button>
        ) : (
          <Button
            size="sm"
            onClick={onStart}
            disabled={loading || !formData.host}
            className="gap-2"
          >
            {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Play className="w-4 h-4" />}
            {t('scanStartBtn')}
          </Button>
        )}
      </CardHeader>

      <CardContent className="space-y-4">
        {error && (
          <Alert variant="destructive">
            <XCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Row 1: Target */}
        <div className="space-y-1.5">
          <Label className="field-label inline-flex items-center gap-1.5">
            <Target className="w-3.5 h-3.5" />
            {t('scanTarget')}
          </Label>
          <Input
            placeholder={t('scanTargetPlaceholder')}
            value={formData.host}
            onChange={(e) => updateField('host', e.target.value)}
            disabled={disabled}
            className="field-input-mono"
          />
        </div>

        {/* Row 2: Ports + Preset + Threads + Timeout */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div className="space-y-1.5">
            <Label className="field-label inline-flex items-center gap-1.5">
              <Hash className="w-3.5 h-3.5" />
              {t('scanPorts')}
            </Label>
            <Input
              placeholder="1-65535"
              value={formData.ports}
              onChange={(e) => updateField('ports', e.target.value)}
              disabled={disabled}
              className="field-input-mono"
            />
          </div>
          <div className="space-y-1.5">
            <Label className="field-label inline-flex items-center gap-1.5">
              <Zap className="w-3.5 h-3.5" />
              {t('scanPreset')}
            </Label>
            <Select onValueChange={applyPreset} disabled={disabled}>
              <SelectTrigger>
                <SelectValue placeholder={t('scanPresetSelect')} />
              </SelectTrigger>
              <SelectContent>
                {presets.map(preset => (
                  <SelectItem key={preset.id} value={preset.id}>
                    {i18n.language === 'zh' ? preset.name : preset.name_en}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-1.5">
            <Label className="field-label inline-flex items-center gap-1.5">
              <Activity className="w-3.5 h-3.5" />
              {t('scanThreads')}
            </Label>
            <Input
              type="number"
              value={formData.thread_num}
              onChange={(e) => updateField('thread_num', parseInt(e.target.value) || 600)}
              disabled={disabled}
              className="field-input-mono"
            />
          </div>
          <div className="space-y-1.5">
            <Label className="field-label inline-flex items-center gap-1.5">
              <Clock className="w-3.5 h-3.5" />
              {t('scanTimeout')}
            </Label>
            <Input
              type="number"
              value={formData.timeout}
              onChange={(e) => updateField('timeout', parseInt(e.target.value) || 3)}
              disabled={disabled}
              className="field-input-mono"
            />
          </div>
        </div>

        {/* Advanced Options */}
        <Collapsible open={showAdvanced} onOpenChange={setShowAdvanced}>
          <CollapsibleTrigger asChild>
            <Button
              variant="ghost"
              size="sm"
              disabled={disabled}
              className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground p-0 h-auto font-normal"
            >
              <Settings2 className="w-4 h-4" />
              <ChevronDown className={`w-4 h-4 transition-transform ${showAdvanced ? 'rotate-180' : ''}`} />
              {t('scanAdvanced')}
            </Button>
          </CollapsibleTrigger>

          <CollapsibleContent className="mt-3">
            <div className="space-y-3 p-3 sm:p-4 rounded-lg bg-muted/50 border">
              {/* Switches */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-2 sm:gap-3">
                <div className="switch-row group">
                  <Label className="inline-flex items-center gap-2 cursor-pointer">
                    <Wifi className="w-4 h-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                    {t('scanDisablePing')}
                  </Label>
                  <Switch
                    checked={formData.disable_ping}
                    onCheckedChange={(checked) => updateField('disable_ping', checked)}
                    disabled={disabled}
                  />
                </div>
                <div className="switch-row group">
                  <Label className="inline-flex items-center gap-2 cursor-pointer">
                    <Lock className="w-4 h-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                    {t('scanDisableBrute')}
                  </Label>
                  <Switch
                    checked={formData.disable_brute}
                    onCheckedChange={(checked) => updateField('disable_brute', checked)}
                    disabled={disabled}
                  />
                </div>
                <div className="switch-row group">
                  <Label className="inline-flex items-center gap-2 cursor-pointer">
                    <CheckCircle2 className="w-4 h-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                    {t('scanAliveOnly')}
                  </Label>
                  <Switch
                    checked={formData.alive_only}
                    onCheckedChange={(checked) => updateField('alive_only', checked)}
                    disabled={disabled}
                  />
                </div>
              </div>

              {/* Credentials */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div className="space-y-1.5">
                  <Label className="field-label inline-flex items-center gap-1.5">
                    <User className="w-3.5 h-3.5" />
                    {t('scanUsername')}
                  </Label>
                  <Input
                    value={formData.username}
                    onChange={(e) => updateField('username', e.target.value)}
                    disabled={disabled}
                    className="field-input"
                  />
                </div>
                <div className="space-y-1.5">
                  <Label className="field-label inline-flex items-center gap-1.5">
                    <Lock className="w-3.5 h-3.5" />
                    {t('scanPassword')}
                  </Label>
                  <Input
                    type="password"
                    value={formData.password}
                    onChange={(e) => updateField('password', e.target.value)}
                    disabled={disabled}
                    className="field-input"
                  />
                </div>
                <div className="space-y-1.5">
                  <Label className="field-label inline-flex items-center gap-1.5">
                    <Globe2 className="w-3.5 h-3.5" />
                    {t('scanDomain')}
                  </Label>
                  <Input
                    value={formData.domain}
                    onChange={(e) => updateField('domain', e.target.value)}
                    disabled={disabled}
                    className="field-input"
                  />
                </div>
              </div>

              {/* Exclusions */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <Label className="field-label inline-flex items-center gap-1.5">
                    <Ban className="w-3.5 h-3.5" />
                    {t('scanExcludeHosts')}
                  </Label>
                  <Input
                    value={formData.exclude_hosts}
                    onChange={(e) => updateField('exclude_hosts', e.target.value)}
                    disabled={disabled}
                    className="field-input-mono"
                  />
                </div>
                <div className="space-y-1.5">
                  <Label className="field-label inline-flex items-center gap-1.5">
                    <Ban className="w-3.5 h-3.5" />
                    {t('scanExcludePorts')}
                  </Label>
                  <Input
                    value={formData.exclude_ports}
                    onChange={(e) => updateField('exclude_ports', e.target.value)}
                    disabled={disabled}
                    className="field-input-mono"
                  />
                </div>
              </div>
            </div>
          </CollapsibleContent>
        </Collapsible>
      </CardContent>
    </Card>
  );
}
