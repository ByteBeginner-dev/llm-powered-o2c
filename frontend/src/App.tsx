import { useState } from 'react'
import '@xyflow/react/dist/style.css'
import { useGraph } from '@/hooks/useGraph'
import { Header } from '@/components/layout/Header'
import { StatusBar } from '@/components/layout/StatusBar'
import { GraphPanel } from '@/components/graph/GraphPanel'
import { NodeDetailPanel } from '@/components/graph/NodeDetailPanel'
import { ChatPanel } from '@/components/chat/ChatPanel'
import { Node } from '@/lib/api'

function App() {
  const { nodes, edges, isLoading, error } = useGraph()
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [chatOpen, setChatOpen] = useState(true)

  const handleNodeClick = (nodeType: string, nodeId: string) => {
    const node = nodes.find(
      (n) => n.type === nodeType && n.label === nodeId
    )
    if (node) {
      setSelectedNode(node)
    }
  }

  return (
    <div className="flex flex-col h-screen bg-bg-base">
      {/* Header */}
      <Header chatOpen={chatOpen} onToggleChat={() => setChatOpen(o => !o)} />

      {/* Main Content */}
      <div className="flex flex-1 overflow-hidden">
        {/* Graph Panel */}
        <div className="flex-1 relative">
          {isLoading ? (
            <div className="flex items-center justify-center h-full">
              <div className="text-center">
                <div className="inline-block">
                  <div className="w-12 h-12 border-4 border-border border-t-accent rounded-full animate-spin" />
                </div>
                <p className="text-text-muted mt-4">Loading graph data...</p>
              </div>
            </div>
          ) : error ? (
            <div className="flex items-center justify-center h-full">
              <div className="text-center max-w-sm">
                <div className="text-4xl mb-4">⚠️</div>
                <h2 className="text-lg font-semibold text-error mb-2">
                  Failed to load graph
                </h2>
                <p className="text-text-muted text-sm">
                  {error.message || 'Make sure the backend is running on http://localhost:8089'}
                </p>
              </div>
            </div>
          ) : (
            <>
              <GraphPanel
                nodes={nodes}
                edges={edges}
                selectedNode={selectedNode}
                onSelectNode={setSelectedNode}
              />
              <NodeDetailPanel
                selectedNode={selectedNode}
                onClose={() => setSelectedNode(null)}
                onNavigate={(node) => setSelectedNode(node)}
              />
            </>
          )}
        </div>

        {/* Chat Panel */}
        {chatOpen && <ChatPanel onNodeClick={handleNodeClick} />}
      </div>

      {/* Status Bar */}
      <StatusBar nodes={nodes} edges={edges} />
    </div>
  )
}

export default App
