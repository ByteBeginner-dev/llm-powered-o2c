import { useEffect, useRef } from 'react'
import { MessageSquare, Zap } from 'lucide-react'
import { useChat } from '@/hooks/useChat'
import { ChatMessage } from './ChatMessage'
import { ChatInput } from './ChatInput'

interface ChatPanelProps {
  onNodeClick?: (type: string, id: string) => void
}

const EXAMPLE_QUERIES = [
  'Which products have the most billing documents?',
  'Show me sales orders with missing deliveries',
  'Trace the full flow of billing document 90504248',
]

export function ChatPanel({ onNodeClick }: ChatPanelProps) {
  const { messages, sendMessage, isLoading } = useChat()
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleExampleQuery = (query: string) => {
    sendMessage(query)
  }

  return (
    <div className="w-2/5 bg-bg-surface border-l border-border flex flex-col">
      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4">
        {messages.length === 0 ? (
          <div className="h-full flex flex-col items-center justify-center text-center">
            <MessageSquare
              size={48}
              className="text-text-faint mb-4 opacity-50"
            />
            <h2 className="text-lg font-semibold text-text-primary mb-2">
              Ask anything about your O2C data
            </h2>
            <p className="text-sm text-text-muted mb-6 max-w-sm">
              Query sales orders, deliveries, billing documents, payments, and more
            </p>

            {/* Example Queries */}
            <div className="space-y-2 w-full">
              {EXAMPLE_QUERIES.map((query) => (
                <button
                  key={query}
                  onClick={() => handleExampleQuery(query)}
                  className="w-full px-4 py-2 text-sm text-text-primary bg-bg-elevated hover:bg-bg-elevated/80 border border-border rounded-card transition-smooth text-left flex items-start gap-2"
                >
                  <Zap size={14} className="text-accent mt-0.5 flex-shrink-0" />
                  {query}
                </button>
              ))}
            </div>
          </div>
        ) : (
          <div>
            {messages.map((message) => (
              <ChatMessage
                key={message.id}
                message={message}
                onNodeClick={onNodeClick}
              />
            ))}

            {/* Loading indicator */}
            {isLoading && (
              <div className="flex gap-3 mb-4">
                <div className="w-6 h-6 rounded-full flex items-center justify-center bg-bg-elevated text-text-muted flex-shrink-0">
                  AI
                </div>
                <div className="flex items-center gap-1 px-4 py-2 rounded-card bg-bg-elevated">
                  <span className="w-2 h-2 bg-accent rounded-full animate-bounce" />
                  <span className="w-2 h-2 bg-accent rounded-full animate-bounce delay-100" />
                  <span className="w-2 h-2 bg-accent rounded-full animate-bounce delay-200" />
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>
        )}
      </div>

      {/* Input */}
      <ChatInput onSend={sendMessage} isLoading={isLoading} />
    </div>
  )
}
