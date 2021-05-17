import { Asset } from "./asset";

export type Balance = {
  asset: Asset;
  amount: number;
  costBasis: number;
  fiatAssetPrice: number;
  fiatValue: number;
};
