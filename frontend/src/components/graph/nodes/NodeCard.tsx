import { Handle, Position } from '@xyflow/react'
import { ReactNode } from 'react'

interface NodeCardProps {
  type: string
  color: string
  icon: ReactNode
  label: string
  properties: { key: string; value: string }[]
  isSelected?: boolean
}

export function NodeCard({
  type,
  color,
  icon,
  label,
  properties,
  isSelected,
}: NodeCardProps) {
  return (
    <div
      className={`w-48 rounded-card bg-bg-surface border transition-smooth ${
        isSelected
          ? `border-accent ring-accent glow-accent`
          : 'border-border hover:border-accent'
      }`}
      style={{
        borderLeft: `4px solid ${color}`,
      }}
    >
      <div className="p-3">
        {/* Header */}
        <div className="flex items-center gap-2 mb-2">
          <div className="text-accent">{icon}</div>
          <span className="text-xs text-text-muted uppercase tracking-wide">
            {type}
          </span>
        </div>

        {/* ID */}
        <div className="font-mono text-sm text-text-primary font-medium truncate mb-2">
          {label}
        </div>

        {/* Properties */}
        {properties.length > 0 && (
          <div className="space-y-1 mb-3 border-t border-border pt-2">
            {properties.slice(0, 2).map((prop) => (
              <div key={prop.key} className="text-xs">
                <span className="text-text-faint">{prop.key}:</span>
                <span className="text-text-muted ml-1 truncate">
                  {prop.value}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Handles */}
      <Handle type="target" position={Position.Left} />
      <Handle type="source" position={Position.Right} />
    </div>
  )
}
