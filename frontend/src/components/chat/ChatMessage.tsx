import { useState } from 'react'
import { Message } from '@/hooks/useChat'
import { ChevronDown } from 'lucide-react'

interface ChatMessageProps {
  message: Message
  onNodeClick?: (type: string, id: string) => void
}

function extractNodeReferences(text: string, onNodeClick?: (type: string, id: string) => void): (string | JSX.Element)[] {
  // Match patterns like "BD: 90504248" or "SO: 100000001" etc
  const pattern = /([A-Z]{1,3}):\s?([0-9]+)/g
  const parts: (string | JSX.Element)[] = []
  let lastIndex = 0
  let match

  const typeMap: Record<string, string> = {
    'SO': 'SalesOrder',
    'BD': 'Billing',
    'DE': 'Delivery',
    'PA': 'Payment',
    'CU': 'Customer',
    'PR': 'Product',
    'PL': 'Plant',
  }

  while ((match = pattern.exec(text)) !== null) {
    // Add text before match
    if (match.index > lastIndex) {
      parts.push(text.substring(lastIndex, match.index))
    }

    const typeCode = match[1]
    const id = match[2]
    const fullType = typeMap[typeCode]

    if (fullType) {
      parts.push(
        <button
          key={`${typeCode}-${id}`}
          onClick={() => onNodeClick?.(fullType, id)}
          className="inline-flex items-center gap-1 px-2 py-0.5 rounded-card bg-accent/20 border border-accent text-accent hover:bg-accent/30 transition-smooth font-mono text-xs"
        >
          {match[0]} →
        </button>
      )
    } else {
      parts.push(match[0])
    }

    lastIndex = match.index + match[0].length
  }

  if (lastIndex < text.length) {
    parts.push(text.substring(lastIndex))
  }

  return parts.length > 0 ? parts : [text]
}

export function ChatMessage({
  message,
  onNodeClick,
}: ChatMessageProps) {
  const [showSql, setShowSql] = useState(false)

  const isUser = message.role === 'user'

  return (
    <div
      className={`flex gap-3 mb-4 animate-in fade-in slide-in-from-bottom-2 ${
        isUser ? 'flex-row-reverse' : ''
      }`}
    >
      {/* Avatar */}
      <div
        className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-semibold flex-shrink-0 ${
          isUser
            ? 'bg-accent/20 text-accent'
            : 'bg-bg-elevated text-text-muted'
        }`}
      >
        {isUser ? 'You' : 'AI'}
      </div>

      {/* Message bubble */}
      <div
        className={`flex-1 max-w-xs ${
          isUser ? 'flex justify-end' : ''
        }`}
      >
        <div
          className={`px-4 py-2 rounded-card ${
            isUser
              ? 'bg-accent/20 border border-accent/40 text-text-primary'
              : 'bg-bg-elevated border border-border text-text-primary'
          }`}
        >
          <p className="text-sm leading-relaxed">
            {extractNodeReferences(message.content, onNodeClick)}
          </p>

          {/* SQL Block */}
          {message.sql && (
            <div className="mt-3 pt-3 border-t border-border">
              <button
                onClick={() => setShowSql(!showSql)}
                className="text-xs text-text-muted hover:text-accent transition-smooth flex items-center gap-1"
              >
                <ChevronDown
                  size={14}
                  className={`transition-transform ${
                    showSql ? 'rotate-180' : ''
                  }`}
                />
                View generated SQL
              </button>

              {showSql && (
                <pre className="mt-2 p-2 bg-bg-base rounded text-xs font-mono text-text-muted overflow-x-auto">
                  {message.sql}
                </pre>
              )}
            </div>
          )}

          {/* Data Table */}
          {message.rows && message.rows.length > 0 && message.rows.length <= 20 && (
            <div className="mt-3 pt-3 border-t border-border overflow-x-auto">
              <table className="text-xs w-full">
                <thead>
                  <tr>
                    {Object.keys(message.rows[0]).map((key) => (
                      <th
                        key={key}
                        className="text-left text-text-muted font-semibold px-2 py-1 border-b border-border"
                      >
                        {key}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {message.rows.map((row, idx) => (
                    <tr key={idx} className="hover:bg-bg-elevated/50">
                      {Object.values(row).map((value, idx) => (
                        <td
                          key={idx}
                          className="px-2 py-1 text-text-muted font-mono"
                        >
                          {String(value)}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Timestamp */}
        <div className={`text-xs text-text-faint mt-1 ${isUser ? 'text-right' : ''}`}>
          {message.timestamp.toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
          })}
        </div>
      </div>
    </div>
  )
}
