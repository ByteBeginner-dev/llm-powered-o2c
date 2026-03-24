import { NodeProps } from '@xyflow/react'
import { Truck } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function DeliveryNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Date', value: data.creation_date ? new Date(String(data.creation_date)).toLocaleDateString('en-US', { day: '2-digit', month: 'short', year: 'numeric' }) : 'N/A' },
    { key: 'Status', value: String(data.overall_delivery_status || 'N/A') },
  ]

  return (
    <NodeCard
      type="Delivery"
      color="var(--node-delivery)"
      icon={<Truck size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
