import { Symbol } from "./symbol";

export type Balance = {
  symbol: Symbol;
  amount: number;
  costBasis: number;
  fiatAssetPrice: number;
  fiatValue: number;
  plAbsolute: number;
  plPercent: number;
};
