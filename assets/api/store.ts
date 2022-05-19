interface BaseEntity {
  id: string;
  lastMessageAt: string;
}

export interface FetchedMetric {
  name: string;
  alias: number;
  stale: boolean;
  dataType: string;
  timestamp: string;
  isNull: boolean;
  value: any;
}

export interface FetchedDevice extends BaseEntity {
  nodeId: string;
  groupId: string;
  online: boolean;
  metrics: FetchedMetric[];
}

export interface FetchedNode extends BaseEntity {
  groupId: string;
  online: boolean;
  devices: FetchedDevice[];
  metrics: FetchedMetric[];
}

export interface FetchedGroup extends BaseEntity {
  nodes: FetchedNode[];
}

export interface GetGroupsResponse {
  data: FetchedGroup[];
}
