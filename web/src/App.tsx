import ky from "ky";
import React, { useEffect, useState } from "react";
import { Redirect, Route, Switch } from "react-router-dom";
import { useRecoilState, useRecoilValue } from "recoil";
import { Navigation } from "./components/Navigation";
import { Account } from "./pages/Account";
import { Accounts } from "./pages/Accounts";
import { Home } from "./pages/Home";
import { Login } from "./pages/Login";
import { SymbolPage } from "./pages/Symbol";
import { loggedInState, userState } from "./state/authState";

function App() {
  const [loading, setLoading] = useState<boolean>(true);

  const [, setUser] = useRecoilState(userState);
  const isLoggedIn = useRecoilValue(loggedInState);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = ky.get("/api/auth/identity");
        setUser(await response.json());
      } catch (error) {
        setUser(null);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [setUser]);

  if (loading) {
    return <p>Loading data...</p>;
  } else if (isLoggedIn) {
    return (
      <div>
        <Navigation />
        <Switch>
          <Route path="/accounts" exact>
            <Accounts></Accounts>
          </Route>
          <Route path="/account/:accountId" exact>
            <Account></Account>
          </Route>
          <Route path="/symbol/:symbolId" exact>
            <SymbolPage></SymbolPage>
          </Route>
          <Route path="/" exact>
            <Home></Home>
          </Route>
        </Switch>
      </div>
    );
  } else {
    return (
      <div className="h-full">
        <Switch>
          <Route path="/login" exact>
            <Login />
          </Route>
          <Route>
            <Redirect to="/login" />
          </Route>
        </Switch>
      </div>
    );
  }
}

export default App;
