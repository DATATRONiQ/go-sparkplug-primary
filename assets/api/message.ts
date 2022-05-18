export type MessageType =
  | "NBIRTH"
  | "NDEATH"
  | "NDATA"
  | "NCMD"
  | "DBIRTH"
  | "DDEATH"
  | "DDATA"
  | "DCMD";

export interface FetchedMessage {
  groupId: string;
  nodeId: string;
  deviceId: string;
  type: MessageType;
  metricAmount: number;
  receivedAt: string;
}

export interface GetMessagesResponse {
    data: FetchedMessage[];
}