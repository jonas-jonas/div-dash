import { atom, selectorFamily } from "recoil";
import { Account } from "../models/account";

export const accountsState = atom<Account[]>({
  key: "Accounts",
  default: [],
});

export const accountByIdSelector = selectorFamily({
  key: "AccountByIdSelector",
  get:
    (id) =>
    ({ get }) => {
      const accounts = get(accountsState);
      for (const account of accounts) {
        if (account.id === id) {
          return account;
        }
      }
    },
});
