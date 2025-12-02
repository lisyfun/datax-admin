import { defineStore } from 'pinia'
import { listConnections, type RedisConnectionInfo } from '@/api/redis'

export const useRedisStore = defineStore('redis', {
  state: () => ({
    connections: [] as (RedisConnectionInfo & { favorite?: boolean })[],
    currentConnId: undefined as number | undefined,
    loading: false,
    keyword: '' as string,
  }),
  actions: {
    async fetchConnections() {
      this.loading = true
      try {
        const { data } = await listConnections({ page: 1, page_size: 100 })
        this.connections = (data.items || []).map((c: any) => ({ ...c, favorite: c.favorite || false }))
        if (!this.currentConnId && this.connections.length) {
          this.currentConnId = this.connections[0].id
        }
      } finally {
        this.loading = false
      }
    },
    setCurrent(id?: number) { this.currentConnId = id },
    setKeyword(kw: string) { this.keyword = kw },
    toggleFavorite(id: number) {
      const idx = this.connections.findIndex(c => c.id === id)
      if (idx >= 0) this.connections[idx].favorite = !this.connections[idx].favorite
    },
    filtered(): (RedisConnectionInfo & { favorite?: boolean })[] {
      const kw = this.keyword.trim().toLowerCase()
      const list = kw ? this.connections.filter(c => c.name.toLowerCase().includes(kw) || c.host.toLowerCase().includes(kw)) : this.connections
      return [...list].sort((a, b) => Number(!!b.favorite) - Number(!!a.favorite))
    },
  }
})
