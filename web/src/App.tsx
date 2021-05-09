import ky from "ky";
import React, { useEffect, useState } from "react";
import { Link, Route, Switch } from "react-router-dom";
import { useRecoilState } from "recoil";
import "./App.scss";
import { Home } from "./pages/Home";
import { userState } from "./state/authState";

function App() {
  const [loading, setLoading] = useState<boolean>(true);

  const [user, setUser] = useRecoilState(userState);

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

  const handleLogout = async () => {
    setUser(null);
    await ky.get("/api/auth/logout");
  };

  if (loading) {
    return <p>Loading data...</p>;
  } else {
    return (
      <div>
        <nav
          className="navbar container"
          role="navigation"
          aria-label="main navigation"
        >
          <div className="navbar-brand">
            <Link className="navbar-item" to="/">
              <img
                src="/assets/logo.svg"
                alt="Bulma: Free, open source, and modern CSS framework based on Flexbox"
                width="112"
                height="28"
              />
            </Link>

            <button
              className="navbar-burger"
              aria-label="menu"
              aria-expanded="false"
            >
              <span aria-hidden="true"></span>
              <span aria-hidden="true"></span>
              <span aria-hidden="true"></span>
            </button>
          </div>

          <div id="navbarBasicExample" className="navbar-menu">
            <div className="navbar-start">
              <a className="navbar-item" href="/">
                Documentation
              </a>

              <div className="navbar-item has-dropdown is-hoverable">
                <a className="navbar-link" href="/">
                  More
                </a>

                <div className="navbar-dropdown">
                  <a className="navbar-item" href="/">
                    About
                  </a>
                  <a className="navbar-item" href="/">
                    Jobs
                  </a>
                  <a className="navbar-item" href="/">
                    Contact
                  </a>
                  <hr className="navbar-divider" />
                  <a className="navbar-item" href="/">
                    Report an issue
                  </a>
                </div>
              </div>
            </div>

            <div className="navbar-end">
              {user && (
                <div className="navbar-item">
                  <p className="mr-2">{user.email}</p>
                  <div className="buttons">
                    <button
                      className="button is-primary"
                      onClick={handleLogout}
                    >
                      <strong>Logout</strong>
                    </button>
                  </div>
                </div>
              )}
              {!user && (
                <div className="navbar-item">
                  <div className="buttons">
                    <a className="button is-primary" href="/register">
                      <strong>Sign up</strong>
                    </a>
                    <a className="button is-light" href="/login">
                      Log in
                    </a>
                  </div>
                </div>
              )}
            </div>
          </div>
        </nav>
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
