import React, { useCallback, useEffect, useState } from "react";
import "./App.css";
import { MessageLog } from "./MessageLog";
import { FetchedMessage, GetMessagesResponse } from "../api/message";
import { FetchedGroup, GetGroupsResponse } from "../api/store";
import { Grid } from "@mui/material";
import { MetricTable } from "./MetricTable";

export const App: React.FC = () => {
  const [messages, setMessages] = useState<FetchedMessage[]>([]);
  const [groups, setGroups] = useState<FetchedGroup[]>([]);

  const refresh = useCallback(async () => {
    fetch("/api/messages")
      .then((res) => res.json())
      .then((resJson: GetMessagesResponse) => setMessages(resJson.data));
    fetch("/api/groups")
      .then((res) => res.json())
      .then((resJson: GetGroupsResponse) => setGroups(resJson.data));
  }, []);

  // TODO: Solve with server sent events (SSE) or websocket

  useEffect(() => {
    const interval = setInterval(refresh, 10000);
    return () => clearInterval(interval);
  }, [refresh]);

  return (
    <div className="App">
      <Grid container spacing={2}>
        <Grid item xs={12} lg={8}>
          <MetricTable groups={groups} />
        </Grid>
        <Grid item xs={12} lg={4}>
          <MessageLog messages={messages} />
        </Grid>
      </Grid>
    </div>
  );
}
