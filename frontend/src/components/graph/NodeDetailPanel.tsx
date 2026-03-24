import { useNodeDetail } from '@/hooks/useNodeDetail'
import { Node } from '@/lib/api'
import { X, ArrowRight } from 'lucide-react'
import { Skeleton } from './Skeleton'

interface NodeDetailPanelProps {
  selectedNode: Node | null
  onClose: () => void
  onNavigate: (node: Node) => void
}

function getNodeColor(type: string): string {
  const colors: Record<string, string> = {
    SalesOrder: 'var(--node-order)',
    Delivery: 'var(--node-delivery)',
    Billing: 'var(--node-billing)',
    Payment: 'var(--node-payment)',
    Customer: 'var(--node-customer)',
    Product: 'var(--node-product)',
    Plant: 'var(--node-plant)',
  }
  return colors[type] || 'var(--accent)'
}

function formatValue(value: any): string {
  if (value === null || value === undefined) return 'N/A'
  if (typeof value === 'number') {
    // Format as currency if it looks like an amount
    if (value > 0 && value < 10000000) {
      return `₹${value.toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
    }
    return value.toString()
  }
  if (typeof value === 'string') {
    // Try to format as date
    if (/^\d{4}-\d{2}-\d{2}/.test(value)) {
      return new Date(value).toLocaleDateString('en-US', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
      })
    }
    return value
  }
  return String(value)
}

export function NodeDetailPanel({
  selectedNode,
  onClose,
  onNavigate,
}: NodeDetailPanelProps) {
  const { nodeDetail, isLoading } = useNodeDetail(
    selectedNode?.type || null,
    selectedNode?.id || null
  )

  if (!selectedNode) return null

  return (
    <div className="absolute inset-y-0 right-0 w-80 bg-bg-surface border-l border-border flex flex-col shadow-lg animate-in slide-in-from-right transition-all duration-200">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-border">
        <div className="flex items-center gap-2">
          <div
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: getNodeColor(selectedNode.type) }}
          />
          <div>
            <div className="text-xs text-text-muted uppercase tracking-wide">
              {selectedNode.type}
            </div>
            <div className="font-mono text-sm font-medium truncate">
              {selectedNode.label}
            </div>
          </div>
        </div>
        <button
          onClick={onClose}
          className="p-1 hover:bg-bg-elevated rounded text-text-muted hover:text-text-primary transition-smooth"
        >
          <X size={18} />
        </button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        {isLoading ? (
          <div className="p-4 space-y-4">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-3/4" />
          </div>
        ) : (
          <>
            {/* Properties */}
            {selectedNode.properties && Object.keys(selectedNode.properties).length > 0 && (
              <div className="p-4">
                <h3 className="text-xs font-semibold text-text-primary uppercase tracking-wide mb-3">
                  Properties
                </h3>
                <div className="space-y-2">
                  {Object.entries(selectedNode.properties).map(([key, value]) => (
                    <div key={key} className="text-xs">
                      <div className="text-text-faint">{key}</div>
                      <div className="font-mono text-text-muted mt-1 break-all">
                        {formatValue(value)}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Connected To */}
            {nodeDetail?.neighbors && nodeDetail.neighbors.length > 0 && (
              <div className="p-4 border-t border-border">
                <h3 className="text-xs font-semibold text-text-primary uppercase tracking-wide mb-3">
                  Connected To
                </h3>
                <div className="space-y-2 flex flex-wrap gap-2">
                  {nodeDetail.neighbors.map((neighbor) => (
                    <button
                      key={neighbor.id}
                      onClick={() => onNavigate(neighbor)}
                      className="text-xs px-2 py-1 rounded-card border border-border hover:bg-bg-elevated transition-smooth"
                      style={{
                        backgroundColor: getNodeColor(neighbor.type) + '1a',
                        borderColor: getNodeColor(neighbor.type),
                        color: 'var(--text-primary)',
                      }}
                    >
                      {neighbor.type}: {neighbor.label}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Footer */}
      <div className="p-4 border-t border-border">
        <button
          onClick={() => {
            const elem = document.querySelector('[title="Fit view"]') as HTMLButtonElement
            elem?.click()
          }}
          className="w-full py-2 px-3 bg-accent hover:bg-accent/90 text-white rounded-card text-xs font-medium flex items-center justify-center gap-2 transition-smooth"
        >
          <ArrowRight size={14} />
          View in Graph
        </button>
      </div>
    </div>
  )
}
