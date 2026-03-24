import { useState, useEffect } from 'react'
import { Moon, Sun, MessageSquare, X } from 'lucide-react'

interface HeaderProps {
  chatOpen: boolean
  onToggleChat: () => void
}

export function Header({ chatOpen, onToggleChat }: HeaderProps) {
  const [theme, setTheme] = useState<'light' | 'dark'>('dark')
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    // Check connection on mount and every 5 minutes
    const checkConnection = async () => {
      try {
        const response = await fetch('http://localhost:8089/health', {
          mode: 'no-cors',
        })
        setIsConnected(response.ok || response.status === 0)
      } catch {
        setIsConnected(false)
      }
    }

    checkConnection()
    const interval = setInterval(checkConnection, 300000) // 5 minutes
    return () => clearInterval(interval)
  }, [])

  const toggleTheme = () => {
    const newTheme = theme === 'dark' ? 'light' : 'dark'
    setTheme(newTheme)
    document.documentElement.setAttribute('data-theme', newTheme)
    localStorage.setItem('theme', newTheme)
  }

  return (
    <header className="h-12 bg-bg-surface border-b border-border flex items-center justify-between px-6">
      <div className="flex items-center gap-3">
        {/* Hexagon Logo */}
        <svg
          width="28"
          height="28"
          viewBox="0 0 28 28"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          className="text-accent"
        >
          <path
            d="M14 2L22.5 6.5V15.5L14 20L5.5 15.5V6.5L14 2Z"
            stroke="currentColor"
            strokeWidth="1.5"
            fill="none"
          />
          <path
            d="M9 8L17 12"
            stroke="currentColor"
            strokeWidth="1"
            opacity="0.5"
          />
        </svg>

        <div>
          <h1 className="text-sm font-semibold text-text-primary">O2C Graph</h1>
          <p className="text-xs text-text-muted">Order-to-Cash Intelligence</p>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {/* Connection Status */}
        <div className="flex items-center gap-2 px-3 py-1 bg-bg-elevated rounded-card">
          <div
            className={`w-2 h-2 rounded-full ${
              isConnected ? 'bg-success' : 'bg-error'
            }`}
          />
          <span className="text-xs text-text-muted">
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
        </div>

        {/* Chat Toggle */}
        <button
          onClick={onToggleChat}
          className="p-2 hover:bg-bg-elevated rounded-card transition-smooth text-text-muted hover:text-text-primary"
          title={chatOpen ? 'Hide AI Chat' : 'Show AI Chat'}
        >
          {chatOpen ? <X size={18} /> : <MessageSquare size={18} />}
        </button>

        {/* Theme Toggle */}
        <button
          onClick={toggleTheme}
          className="p-2 hover:bg-bg-elevated rounded-card transition-smooth text-text-muted hover:text-text-primary"
          title="Toggle theme"
        >
          {theme === 'dark' ? (
            <Sun size={18} />
          ) : (
            <Moon size={18} />
          )}
        </button>
      </div>
    </header>
  )
}
