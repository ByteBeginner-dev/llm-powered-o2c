import { NodeProps } from '@xyflow/react'
import { Package } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function ProductNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Description', value: String(data.product_description || 'N/A') },
    { key: 'Unit', value: String(data.unit_of_measure || 'N/A') },
  ]

  return (
    <NodeCard
      type="Product"
      color="var(--node-product)"
      icon={<Package size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
