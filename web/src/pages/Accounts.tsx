import {
  faChartLine,
  faSpinner,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import ky from "ky";
import React, { useEffect, useState } from "react";
import ReactDOM from "react-dom";
import { useForm } from "react-hook-form";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { Link } from "react-router-dom";
import { AccountForm } from "../form/AccountForm";
import { Account } from "../models/account";
import * as api from "../util/api";
import { formatMoney } from "../util/formatter";

export function Accounts() {
  const { data: accounts, isLoading } = useQuery("accounts", api.getAccounts);
  const [creating, setCreating] = useState(false);

  const handleCreate = () => {
    setCreating(true);
  };

  return (
    <div className="container mx-auto py-8">
      <div className="flex justify-between mb-8">
        <h1 className="text-3xl">Your Accounts</h1>
        <button
          className="bg-gray-900 text-white py-2 px-4 rounded shadow"
          onClick={handleCreate}
        >
          + Account
        </button>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
        {!isLoading &&
          accounts &&
          accounts.map((account) => (
            <Link
              className="bg-white rounded px-6 py-4 transition-all border border-transparent hover:border-blue-600 hover:shadow"
              to={"/accounts/" + account.id}
              key={account.id}
            >
              <div className="flex flex-col mb-4 items-center w-full justify-between">
                <div>
                  <FontAwesomeIcon icon={faChartLine} size="7x" />
                </div>
                <h2 className="text-lg font-bold">{account.name}</h2>
                <span className="text-gray-700 text-sm mb-5">Stocks, ETFs</span>
                <h3 className="font-bold">{formatMoney(4232.23)}</h3>
              </div>
              <div className="flex justify-center mt-5">
                <button className="font-bold">View Account</button>
              </div>
            </Link>
          ))}
        {isLoading && (
          <>
            <AccountCardLoadingIndicator />
            <AccountCardLoadingIndicator />
            <AccountCardLoadingIndicator />
            <AccountCardLoadingIndicator />
          </>
        )}
      </div>

      {creating &&
        ReactDOM.createPortal(
          <CreateAccountModal
            close={() => setCreating(false)}
          ></CreateAccountModal>,
          document.body
        )}
    </div>
  );
}

function AccountCardLoadingIndicator() {
  return (
    <div className="bg-white rounded p-6 animate-pulse flex flex-col items-center px-6 py-4">
      <div className="bg-blue-100 w-24 h-24 mb-4"></div>
      <div className="bg-blue-100 w-20 h-5 mb-2"></div>
      <div className="bg-blue-100 w-16 h-4 mb-5"></div>
      <div className="bg-blue-100 w-24 h-6 mb-5"></div>
      <div className="bg-blue-100 w-28 h-5 mb-4"></div>
    </div>
  );
}

type CreateAccountModalProps = {
  close: () => void;
};

function CreateAccountModal({ close }: CreateAccountModalProps) {
  const { register, handleSubmit, formState, setError } =
    useForm<AccountForm>();

  const queryClient = useQueryClient();

  const createAccountMutation = useMutation<Account, ky.HTTPError, AccountForm>(
    api.postAccount,
    {
      onError: (error) => {
        setError("name", {
          message: error.message,
        });
      },
      onSuccess: (account) => {
        queryClient.setQueryData<Account[]>("accounts", (accounts) => {
          if (accounts) {
            return [...accounts, account];
          }
          return [account];
        });
        close();
      },
    }
  );

  const onSubmit = async (values: AccountForm) => {
    createAccountMutation.mutate(values);
  };

  useEffect(() => {
    const escKeyListener = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        close();
      }
    };

    document.addEventListener("keyup", escKeyListener);

    return () => {
      document.removeEventListener("keyup", escKeyListener);
    };
  }, [close]);

  return (
    <div className="top-0 fixed">
      <form
        className="w-96 mx-auto bg-white rounded fixed top-1/4 transform z-10 left-1/2 -translate-x-1/2 shadow"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="px-8 pt-8 flex justify-between items-start">
          <div className="mr-2">
            <h2 className="text-2xl font-bold">New Account</h2>
            <h3 className="text-gray-500">Lorem ipsum, dolor sit </h3>
          </div>
          <button
            onClick={close}
            className="flex flex-col items-center text-gray-500 hover:text-gray-900 transition-colors"
          >
            <FontAwesomeIcon icon={faTimes} />
            <span className="text-xs">ESC</span>
          </button>
        </div>
        <div className="px-8 py-8">
          <label className="block mb-8">
            <span>Name</span>
            <input
              type="text"
              className="bg-gray-50 block w-full px-3 py-2 focus:outline-none rounded border border-transparent focus:border-blue-700 transition-colors"
              placeholder="Enter account name"
              {...register("name", { required: true })}
            />
          </label>
          <div className="flex justify-items-stretch">
            <button
              className="bg-transparent text-gray-900 rounded py-2 hover:bg-gray-50 transition-colors focus:outline-none border border-gray-900 flex-1 mr-2"
              onClick={close}
              type="reset"
            >
              Cancel
            </button>
            <button
              className="bg-gray-900 text-white rounded py-2 hover:bg-gray-600 transition-colors focus:outline-none flex-1 ml-2 border border-gray-900"
              type="submit"
            >
              {formState.isSubmitting ? (
                <FontAwesomeIcon icon={faSpinner} spin />
              ) : (
                "Create Account"
              )}
            </button>
          </div>
        </div>
      </form>
      <div className="fixed w-full h-full top-0 z-0 backdrop-filter backdrop-blur-sm backdrop-opacity-45 backdrop-brightness-50"></div>
    </div>
  );
}
