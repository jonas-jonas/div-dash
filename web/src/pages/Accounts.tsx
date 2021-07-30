import {
  faChartLine,
  faSpinner,
  faTimes
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import ky from "ky";
import React, { useState } from "react";
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
              to={"/account/" + account.id}
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
  return (
    <div className="top-0 fixed">
      <form
        className="w-96 mx-auto bg-white rounded fixed top-1/4 transform z-10 left-1/2 -translate-x-1/2 shadow"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="border-b border-gray-200 px-8 py-4 flex justify-between">
          <h2 className="text-xl">New Account</h2>
          <button onClick={close}>
            <FontAwesomeIcon icon={faTimes} />
          </button>
        </div>
        <div className="px-8 py-4">
          <label className="block mb-4">
            <span className="text-xs text-gray-600 ml-3">Name</span>
            <input
              type="text"
              className="bg-gray-100 block w-full px-3 py-2 focus:outline-none rounded-md border border-gray-400 focus:border-blue-700 transition-colors"
              {...register("name", { required: true })}
            />
          </label>
          <button
            className="mx-auto block bg-gray-900 text-white rounded px-6 py-2 shadow hover:bg-gray-600 transition-colors focus:outline-none"
            type="submit"
          >
            {formState.isSubmitting ? (
              <FontAwesomeIcon icon={faSpinner} spin />
            ) : (
              "Create"
            )}
          </button>
        </div>
      </form>
      <div className="bg-gray-600 opacity-40 fixed w-full h-full top-0 z-0"></div>
    </div>
  );
}
