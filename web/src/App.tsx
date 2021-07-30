import React, { useMemo } from "react";
import { useQuery } from "react-query";
import { Redirect, Route, Switch } from "react-router-dom";
import { Navigation } from "./components/Navigation";
import { Account } from "./pages/Account";
import { Accounts } from "./pages/Accounts";
import { Home } from "./pages/Home";
import { Login } from "./pages/Login";
import { SymbolPage } from "./pages/Symbol";
import { getIdentity } from "./util/api";

function App() {
  const { isLoading, data, error } = useQuery("identity", getIdentity, {
    retry: false,
  });
  const isLoggedIn = useMemo(() => !error && data, [error, data]);

  if (isLoading) {
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
          <Route path="/">
            <Redirect to="/" />
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
