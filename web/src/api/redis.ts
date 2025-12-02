import request from '@/utils/request'

export interface RedisConnectionInfo {
  id: number
  name: string
  host: string
  port: number
  username: string
  db: number
  use_tls: boolean
}

export async function listConnections(params: { page?: number; page_size?: number; keyword?: string }) {
  return request.get('/redis/connections', { params })
}

export async function createConnection(data: any) {
  return request.post('/redis/connections', data)
}

export async function updateConnection(id: number, data: any) {
  return request.put(`/redis/connections/${id}`, data)
}

export async function deleteConnection(id: number) {
  return request.delete(`/redis/connections/${id}`)
}

export async function testConnection(id: number) {
  return request.post(`/redis/connections/${id}/test`)
}

export async function selectDb(id: number, db: number) {
  return request.post(`/redis/connections/${id}/select`, { db })
}

export async function listKeys(params: { conn_id: number; db?: number; pattern?: string; cursor?: number; count?: number; type?: string }) {
  return request.get('/redis/keys', { params })
}

export async function getKey(params: { conn_id: number; db?: number; key: string }) {
  return request.post('/redis/key', params)
}

export async function setKey(data: { conn_id: number; db?: number; type: string; key: string; value: any; ttl_seconds?: number }) {
  return request.post('/redis/keys', data)
}

export async function deleteKey(data: { conn_id: number; db?: number; key: string }) {
  return request.post('/redis/key/delete', data)
}

export async function getTTL(params: { conn_id: number; db?: number; key: string }) {
  return request.post('/redis/key/ttl', params)
}

export async function setTTL(data: { conn_id: number; db?: number; key: string; seconds: number }) {
  return request.post('/redis/key/ttl/set', data)
}

export async function renameKey(data: { conn_id: number; db?: number; key: string; new_key: string }) {
  return request.post('/redis/keys/rename', data)
}

export async function copyKey(data: { conn_id: number; db?: number; key: string; dest_key: string; replace?: boolean; dest_db?: number }) {
  return request.post('/redis/keys/copy', data)
}

export async function moveKey(data: { conn_id: number; db?: number; key: string; dest_db: number }) {
  return request.post('/redis/keys/move', data)
}

export async function exportKeys(params: { conn_id: number; db?: number; pattern?: string; count?: number }) {
  return request.get('/redis/export', { params })
}

export async function countKeys(params: { conn_id: number; db?: number; pattern?: string; batch?: number }) {
  return request.get('/redis/keys/count', { params })
}

export async function execCommand(data: { conn_id: number; db?: number; line: string }) {
  return request.post('/redis/cli', data)
}

export async function execBulkCommand(data: { conn_id: number; db?: number; lines: string[] }) {
  return request.post('/redis/cli/bulk', data)
}
