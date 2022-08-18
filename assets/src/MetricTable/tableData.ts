import {
  Metric
} from "../api/store";

export interface MetricDataEntry
  extends Pick<Metric, "name" | "alias" | "dataType" | "value"> {
  id: string;
  type: "metric";
  groupId: string;
  nodeId: string;
  deviceId?: string;
  online: boolean;
  lastMessage: Date;
}

export interface DeviceDataEntry {
  id: string;
  type: "device";
  groupId: string;
  nodeId: string;
  deviceId: string;
  online: boolean;
  lastMessage: Date;
  _children: MetricDataEntry[];
}

export interface NodeDataEntry {
  id: string;
  type: "node";
  groupId: string;
  nodeId: string;
  online: boolean;
  lastMessage: Date;
  _children: (DeviceDataEntry | MetricDataEntry)[];
}

export interface GroupDataEntry {
  id: string;
  type: "group";
  online: null;
  lastMessage: Date;
  _children: NodeDataEntry[];
}

export type DataEntry =
  | GroupDataEntry
  | NodeDataEntry
  | DeviceDataEntry
  | MetricDataEntry;
