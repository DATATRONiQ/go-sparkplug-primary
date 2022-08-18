import { useMemo } from "react";
import shallow from "zustand/shallow";
import { SparkplugState, useSparkplugStore } from "../store/store";
import {
  DataEntry,
  DeviceDataEntry,
  GroupDataEntry,
  MetricDataEntry,
  NodeDataEntry
} from "./tableData";

type SparkplugData = Pick<
  SparkplugState,
  "groups" | "nodes" | "devices" | "metrics"
>;

const createMetricData = (data: SparkplugData, metricId: string, groupId: string, nodeId: string, deviceId?: string): MetricDataEntry => {
  const metric = data.metrics[metricId];
  return {
    type: "metric",
    id: metricId,
    alias: metric.alias,
    dataType: metric.dataType,
    name: metric.name,
    value: metric.value,
    lastMessage: new Date(metric.timestamp),
    online: !metric.stale,
    groupId,
    nodeId,
    deviceId,
  };
};

const createDeviceData = (data: SparkplugData, deviceId: string): DeviceDataEntry => {
  const device = data.devices[deviceId];
  return {
    type: "device",
    id: deviceId,
    groupId: device.groupId,
    nodeId: device.nodeId,
    deviceId: device.id,
    lastMessage: new Date(device.lastMessageAt),
    online: device.online,
    _children: device.metricIds.map((metricId) =>
      createMetricData(data, metricId, device.groupId, device.nodeId, device.id)
    ),
  };
};

const createNodeData = (data: SparkplugData, nodeId: string): NodeDataEntry => {
  const node = data.nodes[nodeId];
  return {
    type: "node",
    id: nodeId,
    groupId: node.groupId,
    nodeId: node.id,
    lastMessage: new Date(node.lastMessageAt),
    online: node.online,
    _children: [
      ...node.deviceIds.map((deviceId) => createDeviceData(data, deviceId)),
      ...node.metricIds.map((metricId) => createMetricData(data, metricId, node.groupId, node.id)),
    ],
  };
};

const createGroupData = (
  data: SparkplugData,
  groupId: string
): GroupDataEntry => {
  const group = data.groups[groupId];
  return {
    type: "group",
    id: group.id,
    lastMessage: new Date(group.lastMessageAt),
    online: null,
    _children: group.nodeIds.map((nodeId) => createNodeData(data, nodeId)),
  };
};

export const useTableData = (): DataEntry[] => {
  const state: SparkplugData = useSparkplugStore(
    (state) => ({ groups: state.groups, nodes: state.nodes, devices: state.devices, metrics: state.metrics }),
    shallow
  );

  return useMemo(() => {
    const data = Object.keys(state.groups).map((groupId) => createGroupData(state, groupId));
    console.log("new data: ", data);
    return data;
  }, [state]);
};
