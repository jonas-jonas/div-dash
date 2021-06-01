import { SymbolType } from "./symbol";

export type Transaction = {
  transactionId: string;
  symbol: string;
  type: SymbolType;
  transactionProvider: string;
  price: number;
  date: string;
  amount: number;
  side: "buy" | "sell";
};
