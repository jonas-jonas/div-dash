import ky from "ky";
import React, { useEffect, useState } from "react";
import { Route, Switch } from "react-router-dom";
import { useRecoilState } from "recoil";
import { Navigation } from "./components/Navigation";
import { Home } from "./pages/Home";
import { userState } from "./state/authState";

function App() {
  const [loading, setLoading] = useState<boolean>(true);

  const [, setUser] = useRecoilState(userState);

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
  } else {
    return (
      <div>
        <Navigation />
        <Switch>
          <Route path="/">
            <Home></Home>
          </Route>
        </Switch>
      </div>
    );
  }
}

export default App;
