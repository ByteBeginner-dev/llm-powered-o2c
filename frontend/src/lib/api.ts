import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 30000,
})

export default api

export interface Node {
  id: string
  type: string
  label: string
  properties: Record<string, any>
  position?: { x: number; y: number }
}

export interface Edge {
  id: string
  source: string
  target: string
  label?: string
  data?: Record<string, any>
}

export interface GraphResponse {
  nodes: Node[]
  edges: Edge[]
}

export interface NodeDetailResponse {
  node: Node
  neighbors: Node[]
  edges: Edge[]
}

export interface ChatRequest {
  query: string
}

export interface ChatResponse {
  answer: string
  sql?: string
  rows?: Array<Record<string, any>>
}

export const graphAPI = {
  getGraph: () => api.get<GraphResponse>('/api/graph'),
  getNode: (type: string, id: string) => 
    api.get<NodeDetailResponse>(`/api/node/${type}/${id}`),
  chat: (query: string) => 
    api.post<ChatResponse>('/api/chat', { query }),
}
