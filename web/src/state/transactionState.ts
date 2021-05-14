import { atom } from "recoil";
import { Transaction } from "../models/transaction";

export const transactionsState = atom<Transaction[]>({
  key: "Transactions",
  default: [],
});
