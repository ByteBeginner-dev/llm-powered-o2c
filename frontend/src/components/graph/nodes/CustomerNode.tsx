import { NodeProps } from '@xyflow/react'
import { Users } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function CustomerNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Name', value: String(data.customer_name || 'N/A') },
    { key: 'Country', value: String(data.country || 'N/A') },
  ]

  return (
    <NodeCard
      type="Customer"
      color="var(--node-customer)"
      icon={<Users size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
