import { useQuery } from '@tanstack/react-query'
import { graphAPI, NodeDetailResponse } from '@/lib/api'

export function useNodeDetail(type: string | null, id: string | null) {
  const { data, isLoading, error } = useQuery<NodeDetailResponse | null>({
    queryKey: ['node', type, id],
    queryFn: async () => {
      if (!type || !id) return null
      const response = await graphAPI.getNode(type, id)
      return response.data
    },
    enabled: !!type && !!id,
  })

  return {
    nodeDetail: data || null,
    isLoading,
    error: error as Error | null,
  }
}
