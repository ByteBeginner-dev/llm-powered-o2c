import { useEffect, useState } from 'react'
import { Node, Edge } from '@/lib/api'

interface StatusBarProps {
  nodes: Node[]
  edges: Edge[]
}

export function StatusBar({ nodes, edges }: StatusBarProps) {
  const [lastUpdated, setLastUpdated] = useState<string>('just now')

  useEffect(() => {
    setLastUpdated('just now')
    const timeout = setTimeout(() => {
      const now = new Date()
      const minutes = Math.floor((now.getTime() - Date.now()) / 60000)
      if (minutes < 1) {
        setLastUpdated('just now')
      } else if (minutes < 60) {
        setLastUpdated(`${minutes}m ago`)
      } else {
        const hours = Math.floor(minutes / 60)
        setLastUpdated(`${hours}h ago`)
      }
    }, 60000)

    return () => clearTimeout(timeout)
  }, [nodes])

  return (
    <footer className="h-8 bg-bg-base border-t border-border flex items-center justify-between px-6 text-text-faint text-xs font-mono">
      <div>Nodes: {nodes.length} | Edges: {edges.length}</div>
      <div>SAP Order-to-Cash Dataset · FY2025</div>
      <div>Last updated: {lastUpdated}</div>
    </footer>
  )
}
