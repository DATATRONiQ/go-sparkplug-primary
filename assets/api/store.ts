interface BaseEntity {
  id: string;
  lastMessageAt: string;
}

export interface FetchedDevice extends BaseEntity {
  nodeId: string;
  groupId: string;
  online: boolean;
}

export interface FetchedNode extends BaseEntity {
  groupId: string;
  online: boolean;
  devices: FetchedDevice[];
}

export interface FetchedGroup extends BaseEntity {
  nodes: FetchedNode[];
}

export interface GetGroupsResponse {
  data: FetchedGroup[];
}
