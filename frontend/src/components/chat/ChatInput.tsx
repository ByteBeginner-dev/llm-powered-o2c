import { useState, useRef, useEffect } from 'react'
import { Send, Loader } from 'lucide-react'

interface ChatInputProps {
  onSend: (message: string) => void
  isLoading: boolean
}

export function ChatInput({ onSend, isLoading }: ChatInputProps) {
  const [input, setInput] = useState('')
  const [rows, setRows] = useState(1)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto'
      const scrollHeight = textareaRef.current.scrollHeight
      const lineHeight = parseInt(
        window.getComputedStyle(textareaRef.current).lineHeight
      )
      const newRows = Math.min(Math.ceil(scrollHeight / lineHeight), 3)
      setRows(newRows)
      textareaRef.current.style.height = `${scrollHeight}px`
    }
  }, [input])

  const handleSend = () => {
    if (input.trim() && !isLoading) {
      onSend(input)
      setInput('')
      setRows(1)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const charCount = input.length

  return (
    <div className="p-4 border-t border-border">
      <div className="space-y-2">
        <textarea
          ref={textareaRef}
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Ask about orders, deliveries, billing, payments..."
          rows={rows}
          disabled={isLoading}
          className="w-full px-4 py-2 bg-bg-base border border-border rounded-card text-text-primary placeholder-text-faint focus:outline-none focus:ring-2 focus:ring-accent resize-none disabled:opacity-50"
        />

        <div className="flex items-center justify-between">
          {charCount > 200 && (
            <div className="text-xs text-text-muted">
              {charCount} characters
            </div>
          )}

          <div className="flex-1" />

          <button
            onClick={handleSend}
            disabled={isLoading || !input.trim()}
            className="ml-2 px-4 py-2 bg-accent hover:bg-accent/90 disabled:bg-text-faint disabled:opacity-50 text-white rounded-card font-medium text-sm flex items-center gap-2 transition-smooth"
          >
            {isLoading ? (
              <>
                <Loader size={16} className="animate-spin" />
              </>
            ) : (
              <>
                <Send size={16} />
                Send
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  )
}
