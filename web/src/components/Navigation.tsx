import { faFileInvoice, faSearch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import classNames from "classnames";
import ky from "ky";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { Link, useLocation } from "react-router-dom";
import { User } from "../models/user";
import * as api from "../util/api";

export function Navigation() {
  const { data: user } = useQuery<User>("identity", api.getIdentity);
  const location = useLocation();
  const queryClient = useQueryClient();

  const logoutMutation = useMutation<void, ky.HTTPError>(api.getLogout, {
    onSettled: () => {
      queryClient.setQueryData("identity", null);
    },
  });

  const handleLogout = async () => {
    logoutMutation.mutate();
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
          <button onClick={handleLogout} className="px-3 py-2">
            {user?.email}
          </button>
        </div>
      </div>
    </nav>
  );
}
