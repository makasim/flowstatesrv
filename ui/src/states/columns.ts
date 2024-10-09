import { ColumnDef } from "@tanstack/react-table";

// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type StateData = {
  id: string;
  stateId: string;
  rev: bigint;
  transition: string;
};

export const columns: ColumnDef<StateData>[] = [
  {
    accessorKey: "stateId",
    header: "ID",
  },
  {
    accessorKey: "rev",
    header: "REV",
  },
  {
    accessorKey: "transition",
    header: "Transtion",
  },
];
