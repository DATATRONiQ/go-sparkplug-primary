import { useEffect, useMemo, useReducer, useState } from "react";
import {
  DeviceBirthEvent,
  DeviceDataEvent,
  DeviceDeathEvent,
  isValidEventType,
  NodeBirthEvent,
  NodeDataEvent,
  NodeDeathEvent,
  SparkplugEvent,
} from "../api/event";
import { FullGroup, GetGroupsResponse } from "../api/store";
import { DataEntry, fullGroupToDataEntry } from "./tableData";

type GroupsReducer<T extends SparkplugEvent> = (
  groups: FullGroup[],
  event: T
) => FullGroup[];

const handleNodeBirth: GroupsReducer<NodeBirthEvent> = (groups, event) => {
  const currGroup = groups.find((g) => g.id === event.data.node.groupId);
  if (!currGroup) {
    return groups;
  }
  const currNode = currGroup.nodes.find((n) => n.id === event.data.node.id);
  const newNodes = currGroup.nodes.filter((n) => n.id !== event.data.node.id);
  newNodes.push({
    id: event.data.node.id,
    groupId: currGroup.id,
    lastMessageAt: event.data.node.lastMessageAt,
    metrics: event.data.nodeMetrics,
    online: event.data.node.online,
    devices: currNode ? currNode.devices : [],
  });
  currGroup.nodes = newNodes;
  currGroup.lastMessageAt = event.data.node.lastMessageAt;
  return [...groups];
};

const handleNodeData: GroupsReducer<NodeDataEvent> = (groups, event) => {
  const currGroup = groups.find((g) => g.id === event.data.node.groupId);
  if (!currGroup) {
    return groups;
  }
  const currNode = currGroup.nodes.find((n) => n.id === event.data.node.id);
  if (!currNode) {
    return groups;
  }
  event.data.nodeMetrics.forEach((m) => {
    const currMetric = currNode.metrics.find((nm) => nm.alias === m.alias);
    if (currMetric) {
      currMetric.value = m.value;
      currMetric.stale = m.stale;
      currMetric.timestamp = m.timestamp;
    }
  });

  currNode.lastMessageAt = event.data.node.lastMessageAt;
  currNode.online = event.data.node.online;
  return [...groups];
};

const handleNodeDeath: GroupsReducer<NodeDeathEvent> = (groups, event) => {
  const currGroup = groups.find((g) => g.id === event.data.node.groupId);
  if (!currGroup) {
    return groups;
  }
  const currNode = currGroup.nodes.find((n) => n.id === event.data.node.id);
  if (!currNode) {
    return groups;
  }
  currNode.online = false;
  currGroup.lastMessageAt = event.data.node.lastMessageAt;
  currNode.metrics.forEach((m) => (m.stale = true));
  currNode.devices.forEach((d) => {
    d.online = false;
    d.metrics.forEach((m) => (m.stale = true));
  });
  return [...groups];
};

const handleDeviceBirth: GroupsReducer<DeviceBirthEvent> = (groups, event) => {
  const currGroup = groups.find((g) => g.id === event.data.node.groupId);
  if (!currGroup) {
    return groups;
  }
  const currNode = currGroup.nodes.find((n) => n.id === event.data.node.id);
  if (!currNode) {
    return groups;
  }
  const currDevice = currNode.devices.find(
    (d) => d.id === event.data.device.id
  );
  if (!currDevice) {
    currNode.devices.push({
      id: event.data.device.id,
      nodeId: currNode.id,
      groupId: currGroup.id,
      lastMessageAt: event.data.device.lastMessageAt,
      online: event.data.device.online,
      metrics: event.data.deviceMetrics,
    });
  } else {
    currDevice.online = true;
    currDevice.metrics = event.data.deviceMetrics;
    currDevice.lastMessageAt = event.data.device.lastMessageAt;
  }
  currGroup.lastMessageAt = event.data.node.lastMessageAt;
  return [...groups];
};

const handleDeviceData: GroupsReducer<DeviceDataEvent> = (groups, event) => {
  const currGroup = groups.find((g) => g.id === event.data.node.groupId);
  if (!currGroup) {
    return groups;
  }
  const currNode = currGroup.nodes.find((n) => n.id === event.data.node.id);
  if (!currNode) {
    return groups;
  }
  const currDevice = currNode.devices.find(
    (d) => d.id === event.data.device.id
  );
  if (!currDevice) {
    return groups;
  }
  event.data.deviceMetrics.forEach((m) => {
    const currMetric = currDevice.metrics.find((dm) => dm.alias === m.alias);
    if (currMetric) {
      currMetric.value = m.value;
      currMetric.stale = false;
      currMetric.timestamp = m.timestamp;
    }
  });
  currDevice.lastMessageAt = event.data.device.lastMessageAt;
  currGroup.lastMessageAt = event.data.node.lastMessageAt;
  return [...groups];
};

const handleDeviceDeath: GroupsReducer<DeviceDeathEvent> = (groups, event) => {
  const currGroup = groups.find((g) => g.id === event.data.node.groupId);
  if (!currGroup) {
    return groups;
  }
  const currNode = currGroup.nodes.find((n) => n.id === event.data.node.id);
  if (!currNode) {
    return groups;
  }
  const currDevice = currNode.devices.find(
    (d) => d.id === event.data.device.id
  );
  if (!currDevice) {
    return groups;
  }
  currDevice.online = false;
  currDevice.metrics.forEach((m) => (m.stale = true));
  currGroup.lastMessageAt = event.data.node.lastMessageAt;
  return [...groups];
};

const useFullGroupsReducer = () =>
  useReducer<React.Reducer<FullGroup[], SparkplugEvent>>((groups, event) => {
    console.log("Reducing event: ", event);
    switch (event.type) {
      case "INITIAL":
        return event.data.groups;
      case "NBIRTH":
        return handleNodeBirth(groups, event);
      case "NDATA":
        return handleNodeData(groups, event);
      case "NDEATH":
        return handleNodeDeath(groups, event);
      case "DBIRTH":
        return handleDeviceBirth(groups, event);
      case "DDATA":
        return handleDeviceData(groups, event);
      case "DDEATH":
        return handleDeviceDeath(groups, event);
      default:
        console.log("Unhandled event: ", event);
        return groups;
    }
  }, []);

export const useTableData = (): DataEntry[] => {
  const [initialized, setInitialized] = useState(false);
  const [groups, dispatch] = useFullGroupsReducer();

  // initial import
  useEffect(() => {
    if (initialized) {
      return;
    }
    setInitialized(true);
    fetch("/api/groups")
      .then((res) => res.json())
      .then((resJson: GetGroupsResponse) =>
        dispatch({
          type: "INITIAL",
          timestamp: new Date().toISOString(),
          data: {
            groups: resJson,
          },
        })
      );
  }, [initialized]);

  // update on new events
  useEffect(() => {
    if (!initialized) {
      return;
    }
    const eventSource = new EventSource("/api/groups/stream");
    eventSource.onmessage = (event: MessageEvent) => {
      const parsedEvent: SparkplugEvent = JSON.parse(event.data);
      console.log("Received event: ", parsedEvent);
      if (!isValidEventType(parsedEvent.type)) {
        return;
      }
      dispatch(parsedEvent);
    };
    eventSource.onopen = () => {
      console.log("opened");
    }
    eventSource.onerror = (event: Event) => {
      console.log(event);
    };
    return () => {
      console.log("closing");
      eventSource.close();
    };
  }, [initialized]);

  const data = useMemo(() => groups.map(fullGroupToDataEntry), [groups]);
  console.log("New data: ", data);
  return data;
};
