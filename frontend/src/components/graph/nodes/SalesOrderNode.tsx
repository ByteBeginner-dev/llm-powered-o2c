import { NodeProps } from '@xyflow/react'
import { ShoppingCart } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function SalesOrderNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Amount', value: data.total_net_amount ? `₹${parseFloat(String(data.total_net_amount)).toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : 'N/A' },
    { key: 'Status', value: String(data.overall_delivery_status || 'N/A') },
  ]

  return (
    <NodeCard
      type="Sales Order"
      color="var(--node-order)"
      icon={<ShoppingCart size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
