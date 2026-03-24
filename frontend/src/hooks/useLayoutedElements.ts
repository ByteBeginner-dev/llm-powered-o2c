import { useMemo } from 'react'
import * as dagre from 'dagre'
import { Node, Edge } from '@/lib/api'

const NODE_WIDTH = 192
const NODE_HEIGHT = 120
const RANK_SEP = 350
const NODE_SEP = 50

export function useLayoutedElements(nodes: Node[], edges: Edge[]) {
  const { nodes: layoutedNodes, edges: layoutedEdges } = useMemo(() => {
    if (nodes.length === 0) return { nodes: [], edges: [] }

    const g = new dagre.graphlib.Graph({ compound: false })
    g.setGraph({
      rankdir: 'LR',
      ranksep: RANK_SEP,
      nodesep: NODE_SEP,
    })
    g.setDefaultEdgeLabel(() => ({}))

    // Add nodes
    nodes.forEach((node) => {
      g.setNode(node.id, { width: NODE_WIDTH, height: NODE_HEIGHT })
    })

    // Add edges
    edges.forEach((edge) => {
      g.setEdge(edge.source, edge.target)
    })

    // Run layout
    dagre.layout(g)

    // Extract positions
    const layoutedNodes = nodes.map((node) => {
      const nodeWithPosition = g.node(node.id)
      return {
        ...node,
        position: {
          x: nodeWithPosition.x - NODE_WIDTH / 2,
          y: nodeWithPosition.y - NODE_HEIGHT / 2,
        },
      }
    })

    return { nodes: layoutedNodes, edges }
  }, [nodes, edges])

  return { nodes: layoutedNodes, edges: layoutedEdges }
}
