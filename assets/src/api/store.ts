interface BaseEntity {
  id: string;
  lastMessageAt: string;
}

export interface Group extends BaseEntity {}

export interface FullGroup extends Group {
  nodes: FullNode[];
}

export interface Node extends BaseEntity {
  groupId: string;
  online: boolean;
}

export interface FullNode extends Node {
  devices: FullDevice[];
  metrics: Metric[];
}

export interface Device extends BaseEntity {
  nodeId: string;
  groupId: string;
  online: boolean;
}

export interface FullDevice extends Device {
  metrics: Metric[];
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
  | "UUID";

export interface Metric {
  name: string;
  alias: number;
  stale: boolean;
  dataType: DataType;
  timestamp: string;
  isNull: boolean;
  value: any;
}

export type GetGroupsResponse = FullGroup[];
