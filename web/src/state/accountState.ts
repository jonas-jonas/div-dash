import { atom } from "recoil";
import { Account } from "../models/account";

export const accountsState = atom<Account[]>({
  key: "Accounts",
  default: [],
});
