import { createContext } from "react";
import { ApiClient } from "./api";

export const ApiContext = createContext<ApiClient | null>(null);
