import {
  Box,
  Collapse,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import React, { FC } from "react";
import { FetchedGroup, FetchedNode } from "../api/store";

interface Props {
  groups: FetchedGroup[];
}

const ExpandButton: FC<{ open: boolean; onClick: () => void }> = ({
  open,
  onClick,
}) => (
  <IconButton aria-label="expand row" size="small" onClick={onClick}>
    {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
  </IconButton>
);

const NodeRow: FC<{ node: FetchedNode }> = ({ node }) => {
  const { id, online, devices } = node;
  const [open, setOpen] = React.useState(false);

  return (
    <React.Fragment>
      <TableRow sx={{ "& > *": { borderBottom: "unset" } }}>
        <TableCell>
          <ExpandButton open={open} onClick={() => setOpen(!open)} />
        </TableCell>
        <TableCell component="th" scope="row">
          {id}
        </TableCell>
        <TableCell>
          {online ? "ONLINE" : "OFFLINE"}
        </TableCell>
        <TableCell align="left">{devices.length}</TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Devices
              </Typography>
              <Table size="small" aria-label="devices">
                <TableHead>
                  <TableRow>
                    <TableCell>ID</TableCell>
                    <TableCell>Online</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {devices.map((device) => (
                    <TableRow key={device.id}>
                      <TableCell component="th" scope="row">
                        {device.id}
                      </TableCell>
                      <TableCell>
                        {device.online ? "ONLINE" : "OFFLINE"}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
};

const GroupRow: FC<{ group: FetchedGroup }> = ({ group }) => {
  const { id, nodes } = group;
  const [open, setOpen] = React.useState(false);

  return (
    <React.Fragment>
      <TableRow sx={{ "& > *": { borderBottom: "unset" } }}>
        <TableCell>
          <ExpandButton open={open} onClick={() => setOpen(!open)} />
        </TableCell>
        <TableCell component="th" scope="row">
          {id}
        </TableCell>
        <TableCell>{nodes.length}</TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Nodes
              </Typography>
              <Table size="small" aria-label="nodes">
                <TableHead>
                  <TableRow>
                    <TableCell />
                    <TableCell>ID</TableCell>
                    <TableCell>Online</TableCell>
                    <TableCell>Devices</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {nodes.map((node) => (
                    <NodeRow key={node.id} node={node} />
                  ))}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
};

export const OnlineStatus: FC<Props> = ({ groups }) => (
  <div>
    <h2>Groups</h2>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }}>
          <TableHead>
            <TableCell />
            <TableCell>Group Id</TableCell>
            <TableCell>Nodes</TableCell>
          </TableHead>
          <TableBody>
            {groups.map((group) => (
              <GroupRow key={group.id} group={group} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
  </div>
);
