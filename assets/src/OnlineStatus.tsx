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

interface DeviceStatus {
  id: string;
  online: boolean;
}

interface NodeStatus {
  id: string;
  online: boolean;
  devices: DeviceStatus[];
}

interface GroupStatus {
  id: string;
  nodes: NodeStatus[];
}

interface Props {
  groups: GroupStatus[];
}

const GroupRow: FC<{ group: GroupStatus }> = ({ group }) => {
  const [open, setOpen] = React.useState(false);

  return (
    <React.Fragment>
      <TableRow sx={{ "& > *": { borderBottom: "unset" } }}>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={() => setOpen(!open)}
          >
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {group.id}
        </TableCell>
        <TableCell align="left">{group.nodes.length}</TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                Nodes
              </Typography>
              <Table size="small" aria-label="purchases">
                <TableHead>
                  <TableRow>
                    <TableCell>ID</TableCell>
                    <TableCell>Online</TableCell>
                    <TableCell align="right">Devices</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {group.nodes.map((node) => (
                    <TableRow key={node.id}>
                      <TableCell component="th" scope="row">
                        {node.id}
                      </TableCell>
                      <TableCell>
                        {node.online ? "ONLINE" : "OFFLINE"}
                      </TableCell>
                      <TableCell align="right">{node.devices.length}</TableCell>
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

export const OnlineStatus: FC<Props> = ({ groups }) => (
  <div>
    <h1>Online Status</h1>
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
