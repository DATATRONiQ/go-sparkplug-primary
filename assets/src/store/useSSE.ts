import { useCallback, useEffect, useState } from "react";
import { useSparkplugStore } from "./store";
import { GetGroupsResponse } from "../api/store";
import { isValidEventType, SparkplugEvent } from "../api/event";

export const useSSE = () => {
  const [initialized, setInitialized] = useState(false);

  const initialEvent = useSparkplugStore((state) => state.initialEvent);
  const nodeBirth = useSparkplugStore((state) => state.nodeBirth);
  const nodeDeath = useSparkplugStore((state) => state.nodeDeath);
  const nodeData = useSparkplugStore((state) => state.nodeData);
  const deviceBirth = useSparkplugStore((state) => state.deviceBirth);
  const deviceDeath = useSparkplugStore((state) => state.deviceDeath);
  const deviceData = useSparkplugStore((state) => state.deviceData);

  const handleEvent = useCallback(
    (event: SparkplugEvent) => {
      switch (event.type) {
        case "NBIRTH":
          nodeBirth(event);
          break;
        case "NDATA":
          nodeData(event);
          break;
        case "NDEATH":
          nodeDeath(event);
          break;
        case "DBIRTH":
          deviceBirth(event);
          break;
        case "DDATA":
          deviceData(event);
          break;
        case "DDEATH":
          deviceDeath(event);
          break;
        default:
          console.log("Unhandled event: ", event);
          break;
      }
    },
    [
      initialEvent,
      nodeBirth,
      nodeDeath,
      nodeData,
      deviceBirth,
      deviceDeath,
      deviceData,
    ]
  );

  // initial import
  useEffect(() => {
    if (initialized) {
      return;
    }
    setInitialized(true);
    fetch("/api/groups")
      .then((res) => res.json())
      .then((resJson: GetGroupsResponse) =>
        initialEvent({
          type: "INITIAL",
          timestamp: new Date().toISOString(),
          data: {
            groups: resJson,
          },
        })
      );
  }, [initialized, initialEvent]);

  // update on new events
  useEffect(() => {
    if (!initialized) {
      return;
    }
    const eventSource = new EventSource("/api/groups/stream");
    eventSource.onmessage = (event: MessageEvent) => {
      const parsedEvent: SparkplugEvent = JSON.parse(event.data);
      console.log("useSSE onmessage:", parsedEvent);
      if (!isValidEventType(parsedEvent.type)) {
        return;
      }
      handleEvent(parsedEvent);
    };
    eventSource.onopen = () => {
      console.log("useSSE opened");
    };
    eventSource.onerror = (event: Event) => {
      console.log("useSSE error:", event);
    };
    return () => {
      console.log("useSSE closing");
      eventSource.close();
    };
  }, [initialized, handleEvent]);
};
