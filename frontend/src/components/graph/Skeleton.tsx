interface SkeletonProps {
  className?: string
}

export function Skeleton({ className = '' }: SkeletonProps) {
  return (
    <div
      className={`bg-bg-elevated rounded-card animate-pulse ${className}`}
    />
  )
}
