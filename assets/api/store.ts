interface BaseEntity {
  id: string;
  lastMessageAt: string;
}

export type DataType =
| "Int8"
| "Int16"
| "Int32"
| "Int64"
| "UInt8"
| "UInt16"
| "UInt32"
| "UInt64"
| "Float"
| "Double"
| "Boolean"
| "String"
| "Text"
| "UUID"

export interface FetchedMetric {
  name: string;
  alias: number;
  stale: boolean;
  dataType: DataType;
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
