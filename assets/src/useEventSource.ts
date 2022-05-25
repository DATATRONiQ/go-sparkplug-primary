import { useEffect } from "react";

export const useEventSource = <E>(url: string, onEvent: (e: E) => void) => {
  useEffect(() => {
    console.log("creating event source");
    const evtSource = new EventSource(url, {
      withCredentials: true,
    });
    evtSource.onmessage = (evt: MessageEvent<string>) => {
      console.log("message:", evt);
      const data: E = JSON.parse(evt.data);
      console.log("data:", data);
      onEvent(data);
    };
    evtSource.onopen = () => {
      console.log("opened");
    };
    evtSource.onerror = (evt) => {
      console.log("error:", evt);
    };
    return () => {
      evtSource.close();
    };
  }, [url, onEvent]);
};
