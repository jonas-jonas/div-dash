import { Symbol } from "./symbol";

export type BalanceItem = {
  symbol: Symbol;
  amount: number;
  costBasis: number;
  fiatAssetPrice: number;
  pnl: PNL;
}

export type PNL = {
  pnl: number;
  pnlPercent: number;
}

export type Balance = {
  symbols: BalanceItem[];
  fiatValue: number;
  costBasis: number;
  pnl: PNL;
};
