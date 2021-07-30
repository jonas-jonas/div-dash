import { Symbol } from "./symbol";

export type Account = {
  id: string;
  name: string;
  positions: AccountPosition[];
};

export type AccountPosition = {
  symbol: Symbol
  amount: number;
  buyIn: number;
  currentPrice: number;
  pnlRelative: number;
}