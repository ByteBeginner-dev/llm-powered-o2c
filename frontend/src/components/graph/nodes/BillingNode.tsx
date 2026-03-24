import { NodeProps } from '@xyflow/react'
import { FileText } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function BillingNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Amount', value: data.billing_amount ? `₹${parseFloat(String(data.billing_amount)).toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}` : 'N/A' },
    { key: 'Type', value: String(data.billing_document_type || 'N/A') },
  ]

  return (
    <NodeCard
      type="Billing"
      color="var(--node-billing)"
      icon={<FileText size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
