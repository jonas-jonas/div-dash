import { faBitcoin } from "@fortawesome/free-brands-svg-icons";
import {
  faChevronDown,
  faFileInvoice,
  faFileInvoiceDollar,
  faSearch,
  faSignOutAlt,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import classNames from "classnames";
import ky from "ky";
import { useCallback, useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { Link, useLocation } from "react-router-dom";
import { User } from "../models/user";
import * as api from "../util/api";

export function Navigation() {
  const { data: user } = useQuery<User>("identity", api.getIdentity);
  const location = useLocation();
  const queryClient = useQueryClient();

  const [showUserMenu, setShowUserMenu] = useState<boolean>(false);

  const logoutMutation = useMutation<void, ky.HTTPError>(api.getLogout, {
    onSettled: () => {
      queryClient.setQueryData("identity", null);
    },
  });

  const handleUserClick = useCallback(() => {
    setShowUserMenu((show) => !show);
  }, []);

  const handleLogout = async () => {
    logoutMutation.mutate();
  };

  const navItemClasses = (pathName: string) => {
    return classNames(
      "ml-4 h-full flex items-center border-b-4 border-transparent pt-1 px-3 transition-colors font-bold",
      {
        "border-gray-50 text-white": location.pathname.startsWith(pathName),
        "hover:border-gray-400 text-gray-500":
          !location.pathname.startsWith(pathName),
      }
    );
  };

  useEffect(() => {
    const escListener = function (event: KeyboardEvent) {
      if (event.key === "Escape") {
        setShowUserMenu(false);
      }
    };
    document.addEventListener("keydown", escListener);
    return () => {
      document.removeEventListener("keydown", escListener);
    };
  }, []);

  return (
    <nav
      className="bg-gray-900 w-100 h-16 shadow-md flex items-stretch text-white"
      role="navigation"
      aria-label="main navigation"
    >
      <div className="container mx-auto flex justify-between items-stretch">
        <div className="flex items-center">
          <Link className="navbar-item mr-16" to="/">
            <img src="/logo-white@2x.png" alt="Logo" width="112" height="28" />
          </Link>
          <Link className={navItemClasses("/accounts")} to="/accounts">
            <FontAwesomeIcon icon={faFileInvoice} className="mr-3" />
            <span className="tracking-wider">Accounts</span>
          </Link>
          <Link className={navItemClasses("/symbols/cs")} to="/symbols/cs">
            <FontAwesomeIcon icon={faFileInvoiceDollar} className="mr-3" />
            <span className="tracking-wider">Stocks</span>
          </Link>
          <Link
            className={navItemClasses("/symbols/crypto")}
            to="/symbols/crypto"
          >
            <FontAwesomeIcon icon={faBitcoin} className="mr-3" />
            <span className="tracking-wider">Cryptocurrencies</span>
          </Link>
        </div>
        <div className="flex items-center">
          <button className="px-3 py-2">
            <FontAwesomeIcon icon={faSearch} />
          </button>
          <span className="mx-5 border-l-2 border-white h-4/6 block"></span>
          <div className="relative">
            <button onClick={handleUserClick} className="px-3 py-2 focus:underline focus:outline-none">
              <span className="mr-2">{user?.email}</span>
              <FontAwesomeIcon icon={faChevronDown} size="xs" />
            </button>
            {showUserMenu && (
              <div className="absolute right-0 bg-white rounded shadow-xl border p-2 text-gray-800 flex flex-col w-full mt-4">
                <div className="absolute -top-2 w-4 h-4 right-0 -translate-x-1/2 transform rotate-45 bg-white border-l border-t"></div>
                <button className="px-4 text-left rounded hover:bg-gray-100 transition-colors relative flex items-center py-2">
                  <FontAwesomeIcon icon={faUser} className="absolute" />
                  <span className="ml-8 block font-bold">Profile</span>
                </button>
                <button
                  className="px-4 text-left rounded hover:bg-gray-100 transition-colors relative flex items-center py-2"
                  onClick={handleLogout}
                >
                  <FontAwesomeIcon icon={faSignOutAlt} className="absolute" />
                  <span className="ml-8 block font-bold">Logout</span>
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
