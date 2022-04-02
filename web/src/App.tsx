import { useMemo } from "react";
import { useQuery } from "react-query";
import { Route, Routes } from "react-router-dom";
import { Navigation } from "./components/Navigation";
import { Account } from "./pages/Account";
import { Accounts } from "./pages/Accounts";
import { Home } from "./pages/Home";
import { Login } from "./pages/Login";
import { SymbolPage } from "./pages/Symbol";
import { SymbolListPage } from "./pages/SymbolList";
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
        <Routes>
          <Route path="/accounts" element={<Accounts />} />
          <Route path="/accounts/:accountId" element={<Account />} />
          <Route path="/symbols/:type" element={<SymbolListPage />} />

          <Route path="/symbols/:type/:symbolId" element={<SymbolPage />} />
          <Route path="/" element={<Home />} />
        </Routes>
      </div>
    );
  } else {
    return (
      <div className="h-full">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/*" element={<Login />} />
        </Routes>
      </div>
    );
  }
}

export default App;
