import { NodeProps } from '@xyflow/react'
import { Building2 } from 'lucide-react'
import { NodeCard } from './NodeCard'

interface NodeData {
  id: string
  label: string
  [key: string]: any
}

export function PlantNode({ data: rawData, selected }: NodeProps) {
  const data = rawData as NodeData
  const properties = [
    { key: 'Name', value: String(data.plant_name || 'N/A') },
    { key: 'Country', value: String(data.country || 'N/A') },
  ]

  return (
    <NodeCard
      type="Plant"
      color="var(--node-plant)"
      icon={<Building2 size={16} />}
      label={data.label || data.id}
      properties={properties}
      isSelected={selected}
    />
  )
}
