import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { graphAPI, ChatResponse } from '@/lib/api'
import axios from 'axios'

export interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
  sql?: string
  rows?: Array<Record<string, any>>
  highlightIds?: string[]
}

export function useChat() {
  const [messages, setMessages] = useState<Message[]>([])

  const sendMutation = useMutation({
    mutationFn: async (query: string) => {
      const response = await graphAPI.chat(query)
      return response.data
    },
    onSuccess: (data: ChatResponse) => {
      const assistantMessage: Message = {
        id: `ai-${Date.now()}`,
        role: 'assistant',
        content: data.answer,
        timestamp: new Date(),
        sql: data.sql,
        rows: data.rows,
        highlightIds: data.highlight_ids,
      }
      setMessages((prev) => [...prev, assistantMessage])
    },
    onError: (error: any) => {
      // Handle 400 and other errors from the API
      let errorMessage = 'Unable to process your question. Please try again.'
      
      if (axios.isAxiosError(error) && error.response?.data) {
        // Extract the answer from the error response (for 400 validation errors)
        const responseData = error.response.data as ChatResponse
        if (responseData.answer) {
          errorMessage = responseData.answer
        }
      }

      const assistantMessage: Message = {
        id: `ai-${Date.now()}`,
        role: 'assistant',
        content: errorMessage,
        timestamp: new Date(),
      }
      setMessages((prev) => [...prev, assistantMessage])
    },
  })

  const sendMessage = (query: string) => {
    // Add user message immediately
    const userMessage: Message = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: query,
      timestamp: new Date(),
    }
    setMessages((prev) => [...prev, userMessage])

    // Send to API
    sendMutation.mutate(query)
  }

  const clearMessages = () => {
    setMessages([])
  }

  return {
    messages,
    sendMessage,
    clearMessages,
    isLoading: sendMutation.isPending,
    error: sendMutation.error as Error | null,
  }
}
