import request from '@/utils/request';

// Cluster APIs
export interface KafkaCluster {
  id: number;
  name: string;
  brokerServers: string;
  securityProtocol: string;
  saslMechanism: string;
  username: string;
  password: string;
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface KafkaClusterQuery {
  page: number;
  pageSize: number;
  search?: string;
}

export interface KafkaClusterCreate {
  name: string;
  brokerServers: string;
  securityProtocol: string;
  saslMechanism: string;
  username: string;
  password: string;
  description: string;
}

export interface KafkaClusterUpdate extends KafkaClusterCreate {
  id: number;
}

export function queryClusterList(params: KafkaClusterQuery) {
  return request.get('/kafka/clusters', { params });
}

export function createCluster(data: KafkaClusterCreate) {
  return request.post('/kafka/clusters', data);
}

export function updateCluster(data: KafkaClusterUpdate) {
  return request.put(`/kafka/clusters/${data.id}`, data);
}

export function deleteCluster(id: number) {
  return request.delete(`/kafka/clusters/${id}`);
}

// Topic APIs
export interface KafkaTopic {
  name: string;
  partitions: number;
  replicas: number;
  avgLogSize: number;
  logSize: number;
}

export interface KafkaTopicQuery {
  clusterId: number | string;
  page: number;
  pageSize: number;
  search?: string;
}

export interface KafkaTopicCreate {
  clusterId: number | string;
  name: string;
  partitions: number;
  replicas: number;
}

export interface KafkaTopicUpdate {
  clusterId: number | string;
  name: string;
  partitions: number;
}

export interface KafkaMessage {
  partition: number;
  offset: number;
  key: string;
  value: string;
  timestamp: string;
}

export interface KafkaConsumeParams {
  clusterId: number | string;
  name: string;
  partition: number;
  offset: number;
  count: number;
  keyFilter?: string;
  valueFilter?: string;
}

export function queryTopicList(params: KafkaTopicQuery) {
  return request.get('/kafka/topics', { params });
}

export function createTopic(data: KafkaTopicCreate) {
  return request.post('/kafka/topics', data);
}

export function alterTopic(data: KafkaTopicUpdate) {
  return request.put(`/kafka/topics/${data.name}`, data);
}

export function deleteTopic(params: { clusterId: number | string; name: string }) {
  return request.delete(`/kafka/topics/${params.name}`, { params });
}

export function getTopicPartitions(params: { clusterId: number | string; name: string }) {
  return request.get(`/kafka/topics/${params.name}/partitions`, { params });
}

export function consumeMessages(params: KafkaConsumeParams) {
  return request.get(`/kafka/topics/${params.name}/messages`, { params });
}
