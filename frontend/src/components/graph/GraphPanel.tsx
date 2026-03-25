import { useCallback, useEffect } from 'react'
import {
  ReactFlow,
  Controls,
  MiniMap,
  Background,
  useNodesState,
  useEdgesState,
  NodeTypes,
  useReactFlow,
} from '@xyflow/react'
import { useLayoutedElements } from '@/hooks/useLayoutedElements'
import { SalesOrderNode } from './nodes/SalesOrderNode'
import { DeliveryNode } from './nodes/DeliveryNode'
import { BillingNode } from './nodes/BillingNode'
import { PaymentNode } from './nodes/PaymentNode'
import { CustomerNode } from './nodes/CustomerNode'
import { ProductNode } from './nodes/ProductNode'
import { PlantNode } from './nodes/PlantNode'
import { Node, Edge } from '@/lib/api'
import { Zap } from 'lucide-react'

const nodeTypes: NodeTypes = {
  SalesOrder: SalesOrderNode,
  Delivery: DeliveryNode,
  BillingDocument: BillingNode,
  Payment: PaymentNode,
  Customer: CustomerNode,
  Product: ProductNode,
  Plant: PlantNode,
}

// Separate component to use useReactFlow inside ReactFlow context
function HighlightAutoPan({ 
  highlightedIds, 
  reactFlowNodes 
}: { 
  highlightedIds: string[], 
  reactFlowNodes: any[] 
}) {
  const { setCenter } = useReactFlow()

  useEffect(() => {
    if (highlightedIds.length === 0) return

    // Find the first matching node
    const target = reactFlowNodes.find(n =>
      highlightedIds.includes(n.id) || highlightedIds.includes(n.data?.label)
    )

    if (target?.position) {
      // Pan to that node
      setCenter(target.position.x, target.position.y, {
        zoom: 1.5,
        duration: 800,   // smooth animation
      })
    }
  }, [highlightedIds, reactFlowNodes, setCenter])

  return null
}

interface GraphPanelProps {
  nodes: Node[]
  edges: Edge[]
  selectedNode: Node | null
  onSelectNode: (node: Node | null) => void
  highlightedIds: string[]
}

export function GraphPanel({
  nodes,
  edges,
  selectedNode,
  onSelectNode,
  highlightedIds,
}: GraphPanelProps) {
  const [reactFlowNodes, setReactFlowNodes] = useNodesState<any>([])
  const [reactFlowEdges, setReactFlowEdges] = useEdgesState<any>([])
  const { nodes: layoutedNodes, edges: layoutedEdges } = useLayoutedElements(
    nodes,
    edges
  )

  useEffect(() => {
    const flowNodes = layoutedNodes.map((node) => {
      const isHighlighted = highlightedIds.includes(node.id) ||
                            highlightedIds.includes(node.label)

      return {
        id: node.id,
        data: {
          id: node.id,
          label: node.label,
          ...node.properties,
        },
        position: node.position || { x: 0, y: 0 },
        type: node.type,
        selected: selectedNode?.id === node.id,
        style: {
          // Glow effect on highlighted nodes
          boxShadow: isHighlighted
            ? '0 0 0 3px #4f98a3, 0 0 20px rgba(79,152,163,0.6)'
            : undefined,
          opacity: highlightedIds.length > 0 && !isHighlighted ? 0.3 : 1,
          transition: 'all 0.4s ease',
          zIndex: isHighlighted ? 1000 : undefined,
        }
      }
    })

    const flowEdges = layoutedEdges.map((edge) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      label: edge.label,
      animated:
        selectedNode &&
        (selectedNode.id === edge.source || selectedNode.id === edge.target),
    }))

    setReactFlowNodes(flowNodes)
    setReactFlowEdges(flowEdges)
  }, [layoutedNodes, layoutedEdges, selectedNode, setReactFlowNodes, setReactFlowEdges, highlightedIds])

  const handleNodeClick = useCallback(
    (event: any, node: any) => {
      event.preventDefault()
      const selectedNodeData = nodes.find((n) => n.id === node.id)
      onSelectNode(selectedNodeData || null)
    },
    [nodes, onSelectNode]
  )

  const handlePaneClick = useCallback(() => {
    onSelectNode(null)
  }, [onSelectNode])

  return (
    <div className="relative flex-1 h-full bg-bg-base">
      <ReactFlow
        nodes={reactFlowNodes}
        edges={reactFlowEdges}
        nodeTypes={nodeTypes}
        onNodeClick={handleNodeClick}
        onPaneClick={handlePaneClick}
      >
        <HighlightAutoPan highlightedIds={highlightedIds} reactFlowNodes={reactFlowNodes} />
        <Background />
        <Controls position="bottom-right" />
        <MiniMap
          position="bottom-left"
          style={{
            backgroundColor: 'var(--bg-surface)',
          }}
        />

        {/* Node Legend */}
        <div className="absolute top-4 left-4 bg-bg-surface border border-border rounded-card p-4 text-xs">
          <div className="font-semibold text-text-primary mb-2">Node Types</div>
          <div className="space-y-1">
            {[
              { label: 'Sales Order', color: 'var(--node-order)' },
              { label: 'Delivery', color: 'var(--node-delivery)' },
              { label: 'Billing', color: 'var(--node-billing)' },
              { label: 'Payment', color: 'var(--node-payment)' },
              { label: 'Customer', color: 'var(--node-customer)' },
              { label: 'Product', color: 'var(--node-product)' },
              { label: 'Plant', color: 'var(--node-plant)' },
            ].map((item) => (
              <div key={item.label} className="flex items-center gap-2">
                <div
                  className="w-2 h-5 rounded-sm"
                  style={{ backgroundColor: item.color }}
                />
                <span className="text-text-muted">{item.label}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Fit View Button */}
        <button
          onClick={() => {
            const fitViewButton = document.querySelector(
              '[title="Fit view"]'
            ) as HTMLButtonElement
            fitViewButton?.click()
          }}
          className="absolute top-4 right-4 bg-accent hover:bg-accent/90 text-white px-3 py-2 rounded-card text-xs font-medium flex items-center gap-1 transition-smooth"
        >
          <Zap size={14} />
          Fit View
        </button>
      </ReactFlow>
    </div>
  )
}
