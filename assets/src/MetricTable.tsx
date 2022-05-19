import React, { useMemo } from "react";
import "react-tabulator/lib/css/materialize/tabulator_materialize.min.css";
import {
  ReactTabulator,
  ColumnDefinition,
  ReactTabulatorOptions,
} from "react-tabulator";
import {
  FetchedDevice,
  FetchedGroup,
  FetchedMetric,
  FetchedNode,
} from "../api/store";
import { Tabulator } from "react-tabulator/lib/types/TabulatorTypes";

interface Props {
  groups: FetchedGroup[];
}

interface MetricDataEntry {
  id: string;
  name: string;
  alias: number;
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

type DataEntry =
  | GroupDataEntry
  | NodeDataEntry
  | DeviceDataEntry
  | MetricDataEntry;

const createMetricData = (
  metric: FetchedMetric,
  groupId: string,
  nodeId: string,
  deviceId?: string
): MetricDataEntry => {
  const id =
    deviceId === undefined
      ? `${groupId}/${nodeId}/${metric.alias}`
      : `${groupId}/${nodeId}/${deviceId}/${metric.alias}`;
  return {
    id,
    name: metric.name,
    alias: metric.alias,
    type: "metric",
    groupId,
    nodeId,
    deviceId,
    online: metric.stale === false,
    lastMessage: new Date(metric.timestamp),
  };
};

const createDeviceData = (devices: FetchedDevice[]): DeviceDataEntry[] =>
  devices.map((device) => ({
    id: `${device.groupId}/${device.nodeId}/${device.id}`,
    type: "device",
    groupId: device.groupId,
    nodeId: device.nodeId,
    deviceId: device.id,
    online: device.online,
    lastMessage: new Date(device.lastMessageAt),
    _children: device.metrics.map((metric) => createMetricData(metric, device.groupId, device.nodeId, device.id)),
  }));

const createNodeData = (nodes: FetchedNode[]): NodeDataEntry[] =>
  nodes.map((node) => ({
    id: `${node.groupId}/${node.id}`,
    type: "node",
    groupId: node.groupId,
    nodeId: node.id,
    online: node.online,
    lastMessage: new Date(node.lastMessageAt),
    _children: [
      ...createDeviceData(node.devices),
      ...node.metrics.map((metric) =>
        createMetricData(metric, node.groupId, node.id)
      ),
    ],
  }));

const createGroupData = (groups: FetchedGroup[]): GroupDataEntry[] =>
  groups.map((group) => ({
    id: group.id,
    type: "group",
    online: null,
    lastMessage: new Date(group.lastMessageAt),
    _children: createNodeData(group.nodes),
  }));

const idFormatter = (cell: Tabulator.CellComponent): string => {
  const data = cell.getData() as DataEntry;
  switch (data.type) {
    case "group":
      return data.id;
    case "node":
      return data.nodeId;
    case "device":
      return data.deviceId;
    case "metric":
      return `${data.alias}-${data.name}`;
    default:
      return "";
  }
};

const lastMessageFormatter = (cell: Tabulator.CellComponent): string => {
  const data = cell.getData() as DataEntry;
  return data.lastMessage.toISOString();
};

const columns: ColumnDefinition[] = [
  { title: "ID", field: "id", formatter: idFormatter },
  { title: "Online", field: "online", formatter: "tickCross" },
  {
    title: "Last Message",
    field: "lastMessage",
    formatter: lastMessageFormatter,
  },
];

const options: ReactTabulatorOptions = {
  dataTree: true,
  dataTreeStartExpanded: true,
  layout: "fitDataStretch",
};

export const MetricTable: React.FC<Props> = ({ groups }) => {
  const data = useMemo(() => createGroupData(groups), [groups]);
  return <ReactTabulator data={data} columns={columns} options={options} />;
};
