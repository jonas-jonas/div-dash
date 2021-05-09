import ky from "ky";
import { Link } from "react-router-dom";
import { useRecoilState } from "recoil";
import { userState } from "../state/authState";

export function Navigation() {
  const [user, setUser] = useRecoilState(userState);

  const handleLogout = async () => {
    setUser(null);
    await ky.get("/api/auth/logout");
  };
  return (
    <nav
      className="bg-gray-900 w-100 h-16 shadow-md flex items-center text-white"
      role="navigation"
      aria-label="main navigation"
    >
      <div className="container mx-auto flex justify-between items-center">
        <div className="flex items-center">
          <Link className="navbar-item" to="/">
            <img src="/assets/logo.svg" alt="Logo" width="112" height="28" />
          </Link>
          <a className="ml-6" href="/">
            Documentation
          </a>
        </div>

        <div className="">
          {user && (
            <div className="navbar-item">
              <p className="mr-2">{user.email}</p>
              <div className="buttons">
                <button className="button is-primary" onClick={handleLogout}>
                  <strong>Logout</strong>
                </button>
              </div>
            </div>
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
