import { useQuery } from '@tanstack/react-query'
import { graphAPI, GraphResponse } from '@/lib/api'

export function useGraph() {
  const { data, isLoading, error } = useQuery<GraphResponse>({
    queryKey: ['graph'],
    queryFn: async () => {
      const response = await graphAPI.getGraph()
      return response.data
    },
  })

  return {
    nodes: data?.nodes || [],
    edges: data?.edges || [],
    isLoading,
    error: error as Error | null,
  }
}
