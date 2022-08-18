import { Node, Device, Metric, FullGroup } from "./store";

interface BaseEvent {
  type: string;
  timestamp: string;
  data: unknown;
}

export interface InitialEvent extends BaseEvent {
  type: "INITIAL";
  data: {
    groups: FullGroup[];
  };
}

export interface NodeBirthEvent extends BaseEvent {
  type: "NBIRTH";
  data: {
    node: Node;
    nodeMetrics: Metric[];
  };
}
export interface NodeDataEvent extends BaseEvent {
  type: "NDATA";
  data: {
    node: Node;
    nodeMetrics: Metric[];
  };
}
export interface NodeDeathEvent extends BaseEvent {
  type: "NDEATH";
  data: {
    node: Node;
  };
}
export interface DeviceBirthEvent extends BaseEvent {
  type: "DBIRTH";
  data: {
    node: Node;
    device: Device;
    deviceMetrics: Metric[];
  };
}
export interface DeviceDataEvent extends BaseEvent {
  type: "DDATA";
  data: {
    node: Node;
    device: Device;
    deviceMetrics: Metric[];
  };
}
export interface DeviceDeathEvent extends BaseEvent {
  type: "DDEATH";
  data: {
    node: Node;
    device: Device;
  };
}

export type SparkplugEvent =
  | InitialEvent
  | NodeBirthEvent
  | NodeDataEvent
  | NodeDeathEvent
  | DeviceBirthEvent
  | DeviceDataEvent
  | DeviceDeathEvent;

export type EventType = SparkplugEvent["type"];

/**
 * a helper record with two purposes:
 * 1. for the isValidEventType function
 * 2. I want TypeScript to alarm me if I forget to add a new event type
 */
const EVENT_TYPES: Record<EventType, true> = {
  INITIAL: true,
  NBIRTH: true,
  NDATA: true,
  NDEATH: true,
  DBIRTH: true,
  DDATA: true,
  DDEATH: true,
};

export const isValidEventType = (event: string): event is EventType =>
  event in EVENT_TYPES;
