import create from "zustand";
import {
  DeviceBirthEvent,
  DeviceDataEvent,
  DeviceDeathEvent,
  InitialEvent,
  NodeBirthEvent,
  NodeDataEvent,
  NodeDeathEvent,
} from "../api/event";
import { Group, Node, Device, Metric } from "../api/store";

interface StoreGroup extends Group {
  nodeIds: string[];
}

interface StoreNode extends Node {
  deviceIds: string[];
  metricIds: string[];
}

interface StoreDevice extends Device {
  metricIds: string[];
}

export interface SparkplugState {
  groups: { [groupId: string]: StoreGroup };
  nodes: { [nodeId: string]: StoreNode };
  devices: { [deviceId: string]: StoreDevice };
  metrics: { [metricId: string]: Metric };
  initialEvent: (event: InitialEvent) => void;
  nodeBirth: (event: NodeBirthEvent) => void;
  nodeDeath: (event: NodeDeathEvent) => void;
  nodeData: (event: NodeDataEvent) => void;
  deviceBirth: (event: DeviceBirthEvent) => void;
  deviceDeath: (event: DeviceDeathEvent) => void;
  deviceData: (event: DeviceDataEvent) => void;
}

const getNodeId = (node: Node): string => `${node.groupId}/${node.id}`;
const getDeviceId = (device: Device): string =>
  `${device.groupId}/${device.nodeId}/${device.id}`;
const getMetricId = (parentId: string, metric: Metric): string =>
  `${parentId}/${metric.alias}`;

type EventReducer<TEvent> = (
  state: SparkplugState,
  event: TEvent
) => Partial<SparkplugState>;

const handleInitialEvent: EventReducer<InitialEvent> = (state, event) => {
  const { groups } = event.data;
  const newState: Pick<SparkplugState, 'groups' | 'nodes' | 'devices' | 'metrics'> = {
    groups: {},
    nodes: {},
    devices: {},
    metrics: {},
  }

  groups.forEach((group) => {
    newState.groups[group.id] = {
      id: group.id,
      lastMessageAt: group.lastMessageAt,
      nodeIds: group.nodes.map(getNodeId),
    };

    group.nodes.forEach((node) => {
      const nodeId = getNodeId(node);
      newState.nodes[nodeId] = {
        id: node.id,
        groupId: group.id,
        lastMessageAt: node.lastMessageAt,
        deviceIds: node.devices.map(getDeviceId),
        metricIds: node.metrics.map((m) => getMetricId(nodeId, m)),
        online: node.online,
      };
      
      node.metrics.forEach((metric) => {
        const metricId = getMetricId(nodeId, metric);
        newState.metrics[metricId] = metric;
      });

      node.devices.forEach((device) => {
        const deviceId = getDeviceId(device);
        newState.devices[deviceId] = {
          id: device.id,
          groupId: group.id,
          nodeId: node.id,
          lastMessageAt: device.lastMessageAt,
          metricIds: device.metrics.map((m) => getMetricId(deviceId, m)),
          online: device.online,
        };

        device.metrics.forEach((metric) => {
          const metricId = getMetricId(deviceId, metric);
          newState.metrics[metricId] = metric;
        });
      });
    });
  });

  return newState;
};

const handleNodeBirth: EventReducer<NodeBirthEvent> = (state, event) => {
  const { node, nodeMetrics } = event.data;

  const groups = { ...state.groups };
  if (!state.groups[node.groupId]) {
    groups[node.groupId] = {
      id: node.groupId,
      lastMessageAt: node.lastMessageAt,
      nodeIds: [node.id],
    };
  }
  const nodes = { ...state.nodes };
  const nodeId = getNodeId(node);

  const metrics = { ...state.metrics };
  nodeMetrics.forEach((m) => {
    const metricId = getMetricId(nodeId, m);
    metrics[metricId] = m;
  });
  
  if (!nodes[nodeId]) {
    nodes[nodeId] = {
      id: node.id,
      groupId: node.groupId,
      lastMessageAt: node.lastMessageAt,
      deviceIds: [],
      metricIds: nodeMetrics.map((m) => getMetricId(nodeId, m)),
      online: node.online,
    };
  }

  return { groups, nodes, metrics };
};

const handleNodeData: EventReducer<NodeDataEvent> = (state, event) => {
  const { node, nodeMetrics } = event.data;
  if (!state.groups[node.groupId]) {
    return state;
  }
  const nodeId = getNodeId(node);
  if (!state.nodes[nodeId]) {
    return state;
  }
  const groups = { ...state.groups };
  const nodes = { ...state.nodes };
  groups[node.groupId].lastMessageAt = node.lastMessageAt;
  nodes[nodeId].lastMessageAt = node.lastMessageAt;
  nodes[nodeId].online = node.online;
  const metrics = { ...state.metrics };
  nodeMetrics.forEach((m) => {
    const metricId = getMetricId(nodeId, m);
    metrics[metricId] = m;
  });
  return { groups, nodes, metrics };
};

const handleNodeDeath: EventReducer<NodeDeathEvent> = (state, event) => {
  const { node } = event.data;

  if (!state.groups[node.groupId]) {
    return state;
  }
  if (!state.nodes[getNodeId(node)]) {
    return state;
  }
  const groups = { ...state.groups };
  const nodes = { ...state.nodes };
  const devices = { ...state.devices };
  const metrics = { ...state.metrics };

  const nodeId = getNodeId(node);
  groups[node.groupId].lastMessageAt = node.lastMessageAt;
  nodes[nodeId].lastMessageAt = node.lastMessageAt;
  nodes[nodeId].online = false;
  nodes[nodeId].deviceIds.forEach((deviceId) => {
    devices[deviceId].online = false;
    devices[deviceId].metricIds.forEach((metricId) => {
      metrics[metricId].stale = true;
    });
  });
  nodes[nodeId].metricIds.forEach((metricId) => {
    metrics[metricId].stale = true;
  });
  return { groups, nodes, devices, metrics };
};

const handleDeviceBirth: EventReducer<DeviceBirthEvent> = (state, event) => {
  const { device, deviceMetrics, node } = event.data;

  if (!state.groups[node.groupId]) {
    return state;
  }
  const nodeId = getNodeId(node);
  if (!state.nodes[nodeId]) {
    return state;
  }
  const deviceId = getDeviceId(device);
  const groups = { ...state.groups };
  const nodes = { ...state.nodes };
  const devices = { ...state.devices };
  const metrics = { ...state.metrics };
  groups[node.groupId].lastMessageAt = node.lastMessageAt;

  nodes[nodeId].lastMessageAt = node.lastMessageAt;
  if (!nodes[nodeId].deviceIds.includes(deviceId)) {
    nodes[nodeId].deviceIds.push(deviceId);
  }

  devices[deviceId] = {
    ...device,
    metricIds: deviceMetrics.map((m) => getMetricId(deviceId, m)),
  };

  deviceMetrics.forEach((m) => {
    const metricId = getMetricId(deviceId, m);
    metrics[metricId] = m;
  });

  return { groups, nodes, devices, metrics };
};

const handleDeviceData: EventReducer<DeviceDataEvent> = (state, event) => {
  const { device, deviceMetrics, node } = event.data;
  if (!state.groups[node.groupId]) {
    return state;
  }
  const nodeId = getNodeId(node);
  if (!state.nodes[nodeId]) {
    return state;
  }
  const deviceId = getDeviceId(device);
  if (!state.devices[deviceId]) {
    return state;
  }

  const groups = { ...state.groups };
  const nodes = { ...state.nodes };
  const devices = { ...state.devices };
  const metrics = { ...state.metrics };

  groups[node.groupId].lastMessageAt = node.lastMessageAt;

  nodes[nodeId].lastMessageAt = node.lastMessageAt;

  devices[deviceId].lastMessageAt = device.lastMessageAt;

  deviceMetrics.forEach((m) => {
    const metricId = getMetricId(deviceId, m);
    metrics[metricId] = m;
  });

  return { groups, nodes, devices, metrics };
};

const handleDeviceDeath: EventReducer<DeviceDeathEvent> = (state, event) => {
  const { device, node } = event.data;

  if (!state.groups[node.groupId]) {
    return state;
  }
  const nodeId = getNodeId(node);
  if (!state.nodes[nodeId]) {
    return state;
  }
  const deviceId = getDeviceId(device);
  if (!state.devices[deviceId]) {
    return state;
  }

  const groups = { ...state.groups };
  const nodes = { ...state.nodes };
  const devices = { ...state.devices };
  const metrics = { ...state.metrics };

  groups[node.groupId].lastMessageAt = node.lastMessageAt;

  nodes[nodeId].lastMessageAt = node.lastMessageAt;

  devices[deviceId].lastMessageAt = device.lastMessageAt;
  devices[deviceId].online = false;
  devices[deviceId].metricIds.forEach((metricId) => {
    metrics[metricId].stale = true;
  });
  return { groups, nodes, devices, metrics };
};

export const useSparkplugStore = create<SparkplugState>((set) => ({
  groups: {},
  nodes: {},
  devices: {},
  metrics: {},
  initialEvent: (event) => set((state) => handleInitialEvent(state, event)),
  nodeBirth: (event) => set((state) => handleNodeBirth(state, event)),
  nodeDeath: (event) => set((state) => handleNodeDeath(state, event)),
  nodeData: (event) => set((state) => handleNodeData(state, event)),
  deviceBirth: (event) => set((state) => handleDeviceBirth(state, event)),
  deviceDeath: (event) => set((state) => handleDeviceDeath(state, event)),
  deviceData: (event) => set((state) => handleDeviceData(state, event)),
}));
