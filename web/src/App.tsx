import ky from "ky";
import React, { useEffect, useState } from "react";
import { Redirect, Route, Switch } from "react-router-dom";
import { useRecoilState, useRecoilValue } from "recoil";
import { Navigation } from "./components/Navigation";
import { Home } from "./pages/Home";
import { Login } from "./pages/Login";
import { Accounts } from "./pages/Accounts";
import { loggedInState, tokenState, userState } from "./state/authState";
import { Account } from "./pages/Account";

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
            exact
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
          <Route path="/accounts" exact>
            <Accounts></Accounts>
          </Route>
          <Route path="/account/:accountId" exact>
            <Account></Account>
          </Route>
          <Route path="/" exact>
            <Home></Home>
          </Route>
        </Switch>
      </div>
    );
  }
}

export default App;
