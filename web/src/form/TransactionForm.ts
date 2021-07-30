import { SymbolType } from "../models/symbol";

export type TransactionForm = {
  symbol: string;
  type: SymbolType;
  transactionProvider: string;
  price: string;
  date: string;
  amount: string;
  side: "buy" | "sell";
  accountId: string;
};
