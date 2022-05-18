import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import { FC } from "react";
import { FetchedMessage } from "../api/message";

interface Props {
  messages: FetchedMessage[];
}

export const MessageLog: FC<Props> = ({ messages }) => (
  <div style={{}}>
    <h2>Message Log</h2>
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
              <TableCell>{message.type}</TableCell>
              <TableCell>{message.deviceId}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  </div>
);
