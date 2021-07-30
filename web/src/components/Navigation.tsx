import { faFileInvoice, faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import classNames from "classnames";
import ky from "ky";
import { Link, useLocation } from "react-router-dom";
import { useRecoilState } from "recoil";
import { userState } from "../state/authState";

export function Navigation() {
  const [user, setUser] = useRecoilState(userState);
  const location = useLocation();

  const handleLogout = async () => {
    setUser(null);
    await ky.get("/api/auth/logout");
  };

  const navItemClasses = (pathName: string) => {
    return classNames(
      "ml-20 h-full flex items-center border-b-4 border-transparent pt-1 px-3 transition-colors font-bold",
      {
        "border-gray-50 text-white": location.pathname === pathName,
        "hover:border-gray-400 text-gray-500": location.pathname !== pathName,
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
            <img src="/logo-white@2x.png" alt="Logo" width="112" height="28" />
          </Link>
          <Link className={navItemClasses("/accounts")} to="/accounts">
            <FontAwesomeIcon icon={faFileInvoice} className="mr-3" />
            <span className="tracking-wider">Accounts</span>
          </Link>
        </div>
        <div className="flex items-center">
          <button className="px-3 py-2">
            <FontAwesomeIcon icon={faSearch} />
          </button>
          <span className="mx-5 border-l-2 border-white h-4/6 block"></span>
          <button onClick={handleLogout} className="px-3 py-2">{user?.email}</button>
        </div>
      </div>
    </nav>
  );
}
