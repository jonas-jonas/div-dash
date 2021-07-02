import { atom } from "recoil";
import { Balance } from "../models/balance";

export const balancesState = atom<Balance | null>({
  key: "Balances",
  default: null,
});
