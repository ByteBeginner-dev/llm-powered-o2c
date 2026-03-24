import { NodeProps } from '@xyflow/react'
import { DollarSign } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function PaymentNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Amount', value: data.payment_amount ? `₹${parseFloat(String(data.payment_amount)).toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : 'N/A' },
    { key: 'Date', value: data.document_date ? new Date(String(data.document_date)).toLocaleDateString('en-US', { day: '2-digit', month: 'short', year: 'numeric' }) : 'N/A' },
  ]

  return (
    <NodeCard
      type="Payment"
      color="var(--node-payment)"
      icon={<DollarSign size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
