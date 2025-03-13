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
  topicCount: number;
  brokerCount: number;
  consumerGroups: string[];
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
  page: number;
  pageSize: number;
  search?: string;
}

export interface KafkaTopicCreate {
  name: string;
  partitions: number;
  replicas: number;
}

export interface KafkaTopicUpdate {
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
  partition: number;
  offset: number;
  count: number;
  keyFilter?: string;
  valueFilter?: string;
  groupId?: string;
}

export function queryTopicList(clusterId: number, params: KafkaTopicQuery) {
  return request.get(`/kafka/clusters/${clusterId}/topics`, { params });
}

export function createTopic(clusterId: number, data: KafkaTopicCreate) {
  return request.post(`/kafka/clusters/${clusterId}/topics`, data);
}

export function alterTopic(clusterId: number, topicName: string, data: KafkaTopicUpdate) {
  return request.put(`/kafka/clusters/${clusterId}/topics/${topicName}`, data);
}

export function deleteTopic(clusterId: number, topicName: string) {
  return request.delete(`/kafka/clusters/${clusterId}/topics/${topicName}`);
}

export function getTopicPartitions(clusterId: number, topicName: string) {
  return request.get(`/kafka/clusters/${clusterId}/topics/${topicName}/partitions`);
}

export function consumeMessages(clusterId: number, topicName: string, params: KafkaConsumeParams) {
  return request.get(`/kafka/clusters/${clusterId}/topics/${topicName}/messages`, { params });
}
