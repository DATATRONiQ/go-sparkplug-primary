import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Typography,
} from "@mui/material";
import React from "react";
import { useSparkplugStore } from "../store/store";

export const GroupTable: React.FC = () => {
  const groups = useSparkplugStore((state) => state.groups);
  return (
    <div>
      {Object.entries(groups).map(([groupId, group]) => (
        <Accordion key={groupId}>
          <AccordionSummary>
            <h3>{group.id}</h3>
          </AccordionSummary>
          <AccordionDetails>
            <Typography>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit.
              Suspendisse malesuada lacus ex, sit amet blandit leo lobortis
              eget.
            </Typography>
          </AccordionDetails>
        </Accordion>
      ))}
    </div>
  );
};
