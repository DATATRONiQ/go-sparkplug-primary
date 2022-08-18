import React from "react";
import "react-tabulator/lib/css/materialize/tabulator_materialize.min.css";
import {
  ReactTabulator,
  ColumnDefinition,
  ReactTabulatorOptions,
} from "react-tabulator";
import { Tabulator } from "react-tabulator/lib/types/TabulatorTypes";
import { DataEntry } from "./tableData";
import { useTableData } from "./useTableData";

const idFormatter = (cell: Tabulator.CellComponent): string => {
  const data = cell.getData() as DataEntry;
  switch (data.type) {
    case "group":
      return data.id;
    case "node":
      return data.nodeId;
    case "device":
      return data.deviceId;
    case "metric":
      return `${data.alias}-${data.name}`;
    default:
      return "";
  }
};

const lastMessageFormatter = (cell: Tabulator.CellComponent): string => {
  const data = cell.getData() as DataEntry;
  return data.lastMessage.toISOString();
};

const valueFormatter = (cell: Tabulator.CellComponent): string => {
  const data = cell.getData() as DataEntry;
  if (data.type !== "metric") {
    return "";
  }
  if (data.value === null) {
    return "null";
  }

  switch (data.dataType) {
    case "Int8":
    case "Int16":
    case "Int32":
    case "Int64":
    case "UInt8":
    case "UInt16":
    case "UInt32":
    case "UInt64":
    case "Float":
    case "Double":
      return data.value;
    case "Boolean":
      return data.value ? "true" : "false";
    case "String":
    case "Text":
    case "UUID":
      return data.value;
    default:
      return "";
  }
};

const columns: ColumnDefinition[] = [
  { title: "ID", field: "id", formatter: idFormatter },
  { title: "Online", field: "online", formatter: "tickCross" },
  { title: "Data-Type", field: "dataType" },
  { title: "Value", field: "value", formatter: valueFormatter },
  {
    title: "Last Message",
    field: "lastMessage",
    formatter: lastMessageFormatter,
  },
];

const options: ReactTabulatorOptions = {
  dataTree: true,
  dataTreeStartExpanded: true,
  layout: "fitDataStretch",
  reactiveData: true,
};

export const MetricTable: React.FC = () => {
  const data = useTableData();
  console.log("data", data);
  return <ReactTabulator data={data} columns={columns} options={options} />;
};
