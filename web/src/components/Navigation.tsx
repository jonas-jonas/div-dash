import classNames from "classnames";
import { Link, useLocation } from "react-router-dom";
import { useRecoilState, useResetRecoilState } from "recoil";
import { tokenState, userState } from "../state/authState";

export function Navigation() {
  const [user, setUser] = useRecoilState(userState);
  const resetToken = useResetRecoilState(tokenState);
  const location = useLocation();

  const handleLogout = async () => {
    setUser(null);
    resetToken();
  };

  const navItemClasses = (pathName: string) => {
    return classNames(
      "ml-6 h-full flex items-center border-b-4 border-transparent pt-1 px-3 transition-colors",
      {
        "border-gray-50": location.pathname === pathName,
        "hover:border-gray-400": location.pathname !== pathName,
      }
    );
  };

  return (
    <nav
      className="bg-gray-900 w-100 h-16 shadow-md flex items-stretch text-white"
      role="navigation"
      aria-label="main navigation"
    >
      <div className="container mx-auto flex justify-between items-stretch">
        <div className="flex items-center">
          <Link className="navbar-item" to="/">
            <img src="/assets/logo.svg" alt="Logo" width="112" height="28" />
          </Link>
          <Link className={navItemClasses("/portfolios")} to="/portfolios">
            Portfolios
          </Link>
        </div>

        <div className="flex items-center">
          {user && (
            <>
              <p className="mr-2">{user.email}</p>
              <div className="buttons">
                <button
                  className="px-4 py-2 bg-white shadow rounded text-gray-900"
                  onClick={handleLogout}
                >
                  <strong>Logout</strong>
                </button>
              </div>
            </>
          )}
          {!user && (
            <>
              <Link className="mr-6" to="/login">
                Log in
              </Link>
              <Link
                className="px-4 py-2 bg-white shadow rounded text-gray-900"
                to="/register"
              >
                <strong>Sign up</strong>
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
