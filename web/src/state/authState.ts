import { atom, DefaultValue, selector } from "recoil";
import { User } from "../models/user";

export const userState = atom<User | null>({
  key: "User",
  default: null,
});

export const loggedInState = selector({
  key: "loggedIn",
  get: ({ get }) => {
    const token = get(tokenState);
    return !!token;
  },
});

export const tokenState = atom<string | null>({
  key: "Token",
  default: null,
  effects_UNSTABLE: [
    ({ setSelf, onSet }) => {
      const savedValue = localStorage.getItem("token");
      if (savedValue != null) {
        setSelf(JSON.parse(savedValue));
      }

      onSet((newValue) => {
        if (newValue instanceof DefaultValue) {
          localStorage.removeItem("token");
        } else {
          localStorage.setItem("token", JSON.stringify(newValue));
        }
      });
    },
  ],
});
