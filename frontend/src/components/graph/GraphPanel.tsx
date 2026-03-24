import { useCallback, useEffect } from 'react'
import {
  ReactFlow,
  Controls,
  MiniMap,
  Background,
  useNodesState,
  useEdgesState,
  NodeTypes,
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

interface GraphPanelProps {
  nodes: Node[]
  edges: Edge[]
  selectedNode: Node | null
  onSelectNode: (node: Node | null) => void
}

export function GraphPanel({
  nodes,
  edges,
  selectedNode,
  onSelectNode,
}: GraphPanelProps) {
  const [reactFlowNodes, setReactFlowNodes] = useNodesState<any>([])
  const [reactFlowEdges, setReactFlowEdges] = useEdgesState<any>([])
  const { nodes: layoutedNodes, edges: layoutedEdges } = useLayoutedElements(
    nodes,
    edges
  )

  useEffect(() => {
    const flowNodes = layoutedNodes.map((node) => ({
      id: node.id,
      data: {
        id: node.id,
        label: node.label,
        ...node.properties,
      },
      position: node.position || { x: 0, y: 0 },
      type: node.type,
      selected: selectedNode?.id === node.id,
    }))

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
  }, [layoutedNodes, layoutedEdges, selectedNode, setReactFlowNodes, setReactFlowEdges])

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
