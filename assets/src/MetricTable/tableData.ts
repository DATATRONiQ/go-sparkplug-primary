import {
  FullGroup,
  Group,
  Node,
  Metric,
  Device,
  FullNode,
  FullDevice,
} from "../api/store";

interface MetricDataEntry
  extends Pick<Metric, "name" | "alias" | "dataType" | "value"> {
  id: string;
  type: "metric";
  groupId: string;
  nodeId: string;
  deviceId?: string;
  online: boolean;
  lastMessage: Date;
}

interface DeviceDataEntry {
  id: string;
  type: "device";
  groupId: string;
  nodeId: string;
  deviceId: string;
  online: boolean;
  lastMessage: Date;
  _children: MetricDataEntry[];
}

interface NodeDataEntry {
  id: string;
  type: "node";
  groupId: string;
  nodeId: string;
  online: boolean;
  lastMessage: Date;
  _children: (DeviceDataEntry | MetricDataEntry)[];
}

interface GroupDataEntry {
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

const groupToRowId = (group: Group): string => group.id;
const nodeToRowId = (node: Node): string => `${node.groupId}/${node.id}`;
const deviceToRowId = (device: Device): string =>
  `${device.groupId}/${device.nodeId}/${device.id}`;

const metricToDataEntry = (
  metric: Metric,
  groupId: string,
  nodeId: string,
  deviceId?: string
): MetricDataEntry => {
  const id =
    deviceId === undefined
      ? `${groupId}/${nodeId}/${metric.alias}`
      : `${groupId}/${nodeId}/${deviceId}/${metric.alias}`;
  return {
    id: id,
    type: "metric",
    groupId: groupId,
    nodeId: nodeId,
    deviceId: deviceId,
    name: metric.name,
    alias: metric.alias,
    dataType: metric.dataType,
    value: metric.value,
    online: !metric.stale,
    lastMessage: new Date(metric.timestamp),
  };
};

const fullDeviceToDataEntry = (device: FullDevice): DeviceDataEntry => ({
  id: deviceToRowId(device),
  type: "device",
  groupId: device.groupId,
  nodeId: device.nodeId,
  deviceId: device.id,
  online: device.online,
  lastMessage: new Date(device.lastMessageAt),
  _children: device.metrics.map((m) => metricToDataEntry(m, device.groupId, device.nodeId, device.id)),
});

const fullNodeToDataEntry = (node: FullNode): NodeDataEntry => ({
  id: nodeToRowId(node),
  type: "node",
  groupId: node.groupId,
  nodeId: node.id,
  online: node.online,
  lastMessage: new Date(node.lastMessageAt),
  _children: [
      ...node.metrics.map((m) => metricToDataEntry(m, node.groupId, node.id)),
      ...node.devices.map(fullDeviceToDataEntry),
  ],
});

export const fullGroupToDataEntry = (group: FullGroup): GroupDataEntry => ({
  id: groupToRowId(group),
  type: "group",
  online: null,
  lastMessage: new Date(group.lastMessageAt),
  _children: group.nodes.map(fullNodeToDataEntry),
});
