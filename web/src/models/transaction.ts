export type Transaction = {
  transactionId: string;
  symbol: string;
  type: string;
  transactionProvider: string;
  price: number;
  date: string;
  amount: number;
  side: "buy" | "sell";
};
