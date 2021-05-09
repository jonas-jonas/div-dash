import ky from "ky";
import React, { useEffect, useState } from "react";
import { Redirect, Route, Switch } from "react-router-dom";
import { useRecoilState, useRecoilValue } from "recoil";
import { Navigation } from "./components/Navigation";
import { Home } from "./pages/Home";
import { Login } from "./pages/Login";
import { loggedInState, tokenState, userState } from "./state/authState";

function App() {
  const [loading, setLoading] = useState<boolean>(true);

  const [, setUser] = useRecoilState(userState);
  const [token, setToken] = useRecoilState(tokenState);
  const isLoggedIn = useRecoilValue(loggedInState);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = ky.get("/api/auth/identity", {
          headers: {
            Authorization: "Bearer " + token,
          },
        });
        setUser(await response.json());
      } catch (error) {
        setUser(null);
        setToken(null);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [setUser, token, setToken]);

  if (loading) {
    return <p>Loading data...</p>;
  } else {
    return (
      <div>
        <Navigation />
        <Switch>
          <Route
            path="/login"
            render={(props) =>
              isLoggedIn ? (
                <Redirect
                  to={{ pathname: "/", state: { from: props.location } }}
                />
              ) : (
                <Login />
              )
            }
          />
        </Switch>
        <Route path="/">
          <Home></Home>
        </Route>
      </div>
    );
  }
}

export default App;
