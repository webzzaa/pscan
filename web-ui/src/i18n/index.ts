import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

const resources = {
  en: {
    translation: {
      // App
      appTitle: 'Fscan Web UI',
      appDescription: 'Network Security Scanner',

      // Navigation
      navScan: 'Scan',
      navResults: 'Results',
      navSettings: 'Settings',

      // Scan Page
      scanTitle: 'New Scan',
      scanTarget: 'Target',
      scanTargetPlaceholder: 'IP, IP range, domain (e.g., 192.168.1.0/24)',
      scanPorts: 'Ports',
      scanPortsPlaceholder: 'Port range (e.g., 1-1000,3306,8080)',
      scanPreset: 'Preset',
      scanPresetSelect: 'Select preset...',
      scanMode: 'Scan Mode',
      scanModeAll: 'All',
      scanModeIcmp: 'ICMP Only',
      scanThreads: 'Threads',
      scanTimeout: 'Timeout (s)',
      scanAdvanced: 'Advanced Options',
      scanDisablePing: 'Disable Ping',
      scanDisableBrute: 'Disable Brute Force',
      scanAliveOnly: 'Alive Only',
      scanUsername: 'Username',
      scanPassword: 'Password',
      scanDomain: 'Domain',
      scanExcludeHosts: 'Exclude Hosts',
      scanExcludePorts: 'Exclude Ports',
      scanStartBtn: 'Start Scan',
      scanStopBtn: 'Stop Scan',
      scanRunning: 'Scan Running...',

      // Status
      statusIdle: 'Idle',
      statusRunning: 'Running',
      statusStopping: 'Stopping',

      // Stats
      statsHosts: 'Hosts',
      statsPorts: 'Ports',
      statsServices: 'Services',
      statsVulns: 'Vulnerabilities',

      // Results
      resultsTitle: 'Scan Results',
      resultsDistribution: 'Results Distribution',
      chartEmptyTitle: 'No data available',
      chartEmptyDescription: 'Statistics will be displayed here after scanning',
      resultsExport: 'Export',
      resultsClear: 'Clear',
      resultsEmpty: 'No results yet',
      resultsEmptyDescription: 'Results will appear here after scanning',
      resultsFilterAll: 'All',
      resultsFilterHosts: 'Hosts',
      resultsFilterPorts: 'Ports',
      resultsFilterServices: 'Services',
      resultsFilterVulns: 'Vulnerabilities',

      // Live Feed
      liveFeed: 'Live Feed',
      liveFeedConnected: 'Connected',
      liveFeedDisconnected: 'Disconnected',
      liveFeedEmptyDescription: 'Start a scan to see real-time results',

      // Settings
      settingsTitle: 'Settings',
      settingsLanguage: 'Language',
      settingsTheme: 'Theme',
      settingsThemeLight: 'Light',
      settingsThemeDark: 'Dark',
      settingsThemeSystem: 'System',

      // Common
      loading: 'Loading...',
      error: 'Error',
      success: 'Success',
      cancel: 'Cancel',
      confirm: 'Confirm',
      close: 'Close',
      items: 'items',
      refresh: 'Refresh',
      clearAll: 'Clear all',
      export: 'Export',
      exportJson: 'Export JSON',
      exportCsv: 'Export CSV',

      // Validation & Errors
      targetRequired: 'Target is required',
      startScanFailed: 'Failed to start scan',
      stopScanFailed: 'Failed to stop scan',
      clearConfirmTitle: 'Clear Results',
      clearConfirm: 'Are you sure you want to clear all results? This action cannot be undone.',

      // Theme
      lightMode: 'Light mode',
      darkMode: 'Dark mode',

      // Result Types
      typeHost: 'host',
      typePort: 'port',
      typeService: 'service',
      typeVuln: 'vuln',
    },
  },
  zh: {
    translation: {
      // App
      appTitle: 'Fscan Web UI',
      appDescription: '网络安全扫描器',

      // Navigation
      navScan: '扫描',
      navResults: '结果',
      navSettings: '设置',

      // Scan Page
      scanTitle: '新建扫描',
      scanTarget: '目标',
      scanTargetPlaceholder: 'IP、IP段、域名 (如: 192.168.1.0/24)',
      scanPorts: '端口',
      scanPortsPlaceholder: '端口范围 (如: 1-1000,3306,8080)',
      scanPreset: '预设',
      scanPresetSelect: '选择预设...',
      scanMode: '扫描模式',
      scanModeAll: '全部',
      scanModeIcmp: '仅ICMP',
      scanThreads: '线程数',
      scanTimeout: '超时(秒)',
      scanAdvanced: '高级选项',
      scanDisablePing: '禁用Ping',
      scanDisableBrute: '禁用爆破',
      scanAliveOnly: '仅存活检测',
      scanUsername: '用户名',
      scanPassword: '密码',
      scanDomain: '域名',
      scanExcludeHosts: '排除主机',
      scanExcludePorts: '排除端口',
      scanStartBtn: '开始扫描',
      scanStopBtn: '停止扫描',
      scanRunning: '扫描进行中...',

      // Status
      statusIdle: '空闲',
      statusRunning: '运行中',
      statusStopping: '停止中',

      // Stats
      statsHosts: '主机',
      statsPorts: '端口',
      statsServices: '服务',
      statsVulns: '漏洞',

      // Results
      resultsTitle: '扫描结果',
      resultsDistribution: '结果分布',
      chartEmptyTitle: '暂无数据',
      chartEmptyDescription: '扫描后将在此显示统计图表',
      resultsExport: '导出',
      resultsClear: '清空',
      resultsEmpty: '暂无结果',
      resultsEmptyDescription: '扫描结果将在此显示',
      resultsFilterAll: '全部',
      resultsFilterHosts: '主机',
      resultsFilterPorts: '端口',
      resultsFilterServices: '服务',
      resultsFilterVulns: '漏洞',

      // Live Feed
      liveFeed: '实时动态',
      liveFeedConnected: '已连接',
      liveFeedDisconnected: '已断开',
      liveFeedEmptyDescription: '开始扫描后将在此显示实时结果',

      // Settings
      settingsTitle: '设置',
      settingsLanguage: '语言',
      settingsTheme: '主题',
      settingsThemeLight: '浅色',
      settingsThemeDark: '深色',
      settingsThemeSystem: '跟随系统',

      // Common
      loading: '加载中...',
      error: '错误',
      success: '成功',
      cancel: '取消',
      confirm: '确认',
      close: '关闭',
      items: '条',
      refresh: '刷新',
      clearAll: '清空全部',
      export: '导出',
      exportJson: '导出 JSON',
      exportCsv: '导出 CSV',

      // Validation & Errors
      targetRequired: '请输入扫描目标',
      startScanFailed: '启动扫描失败',
      stopScanFailed: '停止扫描失败',
      clearConfirmTitle: '清空结果',
      clearConfirm: '确定要清空所有结果吗？此操作不可撤销。',

      // Theme
      lightMode: '浅色模式',
      darkMode: '深色模式',

      // Result Types
      typeHost: '主机',
      typePort: '端口',
      typeService: '服务',
      typeVuln: '漏洞',
    },
  },
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: 'zh',
    fallbackLng: 'zh',
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
