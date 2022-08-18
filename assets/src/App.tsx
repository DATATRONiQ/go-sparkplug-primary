import React, { useEffect, useState } from "react";
import "./App.css";
import { MessageLog } from "./MessageLog";
import { FetchedMessage, GetMessagesResponse } from "./api/message";
import { Grid } from "@mui/material";
import { MetricTable } from "./MetricTable/MetricTable";

export const App: React.FC = () => {
  const [messages, setMessages] = useState<FetchedMessage[]>([]);

  useEffect(() => {
    fetch("/api/messages")
      .then((res) => res.json())
      .then((resJson: GetMessagesResponse) => setMessages(resJson));
  }, []);

  return (
    <div className="App">
      <Grid container spacing={2}>
        <Grid item xs={12} lg={8}>
          <MetricTable />
        </Grid>
        <Grid item xs={12} lg={4}>
          <MessageLog messages={messages} />
        </Grid>
      </Grid>
    </div>
  );
};
