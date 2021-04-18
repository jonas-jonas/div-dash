import { atom, selector } from "recoil";
import { User } from "../models/user";

export const userState = atom<User | null>({
  key: "User",
  default: null
})

export const loggedInState = selector({
  key: "loggedIn",
  get: ({ get }) => {
    const user = get(userState);
    return !!user;
  },
  set: ({ set }) => {
    set(userState, null);
  }
})