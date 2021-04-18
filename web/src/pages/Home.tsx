import { useRecoilValue } from "recoil";
import { userState } from "../state/authState";

export function Home() {
  const user = useRecoilValue(userState);
  return <div className="container">Hi {user?.email || "Not logged in"}</div>;
}
