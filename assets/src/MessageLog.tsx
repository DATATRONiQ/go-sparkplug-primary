import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from "@mui/material";
import { FC } from "react";
import { Todo } from "./util";

interface Message {
  receivedAt: string;
  groupId: string;
  nodeId: string;
  messageType: string;
  deviceId: string;
  payload: Todo;
}

interface Props {
  messages: Message[];
}

export const MessageLog: FC<Props> = ({ messages }) => (
  <div style={{
  }}>
    <h1>Message Log</h1>
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }}>
        <TableHead>
          <TableCell>Received At</TableCell>
          <TableCell>Group Id</TableCell>
          <TableCell>Node Id</TableCell>
          <TableCell>Message Type</TableCell>
          <TableCell>Device Id</TableCell>
        </TableHead>
        <TableBody>
          {messages.map((message) => (
            <TableRow key={message.receivedAt}>
              <TableCell>{message.receivedAt}</TableCell>
              <TableCell>{message.groupId}</TableCell>
              <TableCell>{message.nodeId}</TableCell>
              <TableCell>{message.messageType}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  </div>
);
