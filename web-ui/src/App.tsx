import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Terminal, Languages, Moon, Sun, Radar, BarChart3, Github, ExternalLink } from 'lucide-react';
import { ScanPage } from '@/pages/ScanPage';
import { ResultsPage } from '@/pages/ResultsPage';
import { LiveFeedProvider } from '@/contexts/LiveFeedContext';
import { Button } from '@/components/ui/button';
import './i18n';

type Page = 'scan' | 'results';

function App() {
  const { t, i18n } = useTranslation();
  const [currentPage, setCurrentPage] = useState<Page>('scan');
  const [isDark, setIsDark] = useState(() => {
    if (typeof window !== 'undefined') {
      const stored = localStorage.getItem('theme');
      if (stored) return stored === 'dark';
      return window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    return false;
  });

  useEffect(() => {
    document.documentElement.classList.toggle('dark', isDark);
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
  }, [isDark]);

  return (
    <LiveFeedProvider>
    <div className="min-h-screen bg-background flex flex-col">
      {/* Header */}
      <header className="sticky top-0 z-50 border-b bg-background/95 backdrop-blur-sm">
        <div className="container h-14 sm:h-16 flex items-center justify-between">
          {/* Left: Logo + Nav */}
          <div className="flex items-center gap-4 sm:gap-6">
            {/* Brand */}
            <a
              href="https://github.com/shadow1ng/fscan"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-foreground hover:opacity-80 transition-opacity"
            >
              <div className="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-foreground flex items-center justify-center">
                <Terminal className="w-4 h-4 sm:w-5 sm:h-5 text-background" />
              </div>
              <span className="font-semibold text-base sm:text-lg tracking-tight">fscan</span>
              <span className="text-xs text-muted-foreground font-mono px-1.5 py-0.5 rounded bg-muted hidden sm:inline">v2.1</span>
            </a>

            <div className="h-5 w-px bg-border hidden sm:block" />

            {/* Navigation */}
            <nav className="flex items-center gap-1">
              <Button
                variant={currentPage === 'scan' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setCurrentPage('scan')}
                className="gap-2"
              >
                <Radar className="w-4 h-4" />
                <span className="hidden sm:inline">{t('navScan')}</span>
              </Button>
              <Button
                variant={currentPage === 'results' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setCurrentPage('results')}
                className="gap-2"
              >
                <BarChart3 className="w-4 h-4" />
                <span className="hidden sm:inline">{t('navResults')}</span>
              </Button>
            </nav>
          </div>

          {/* Right: Actions */}
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              asChild
            >
              <a
                href="https://github.com/shadow1ng/fscan"
                target="_blank"
                rel="noopener noreferrer"
                title="GitHub"
              >
                <Github className="w-5 h-5" />
              </a>
            </Button>
            <div className="h-5 w-px bg-border mx-1" />
            <Button
              variant="ghost"
              size="icon"
              onClick={() => i18n.changeLanguage(i18n.language === 'zh' ? 'en' : 'zh')}
              title={i18n.language === 'zh' ? 'English' : '中文'}
            >
              <Languages className="w-5 h-5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setIsDark(!isDark)}
              title={isDark ? t('lightMode') : t('darkMode')}
            >
              {isDark ? (
                <Sun className="w-5 h-5" />
              ) : (
                <Moon className="w-5 h-5" />
              )}
            </Button>
          </div>
        </div>
      </header>

      {/* Main */}
      <main className="container flex-1 py-3 sm:py-4 lg:py-5">
        {currentPage === 'scan' ? <ScanPage /> : <ResultsPage />}
      </main>

      {/* Footer */}
      <footer className="border-t py-3 sm:py-4">
        <div className="container flex items-center justify-center gap-2 text-xs sm:text-sm text-muted-foreground">
          <Terminal className="w-4 h-4" />
          <span>{t('appDescription')}</span>
          <span className="opacity-40">·</span>
          <a
            href="https://github.com/shadow1ng/fscan"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 hover:text-foreground transition-colors"
          >
            GitHub
            <ExternalLink className="w-3.5 h-3.5" />
          </a>
        </div>
      </footer>
    </div>
    </LiveFeedProvider>
  );
}

export default App;
