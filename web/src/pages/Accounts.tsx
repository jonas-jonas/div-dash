import {
  faChartLine,
  faPlus,
  faSpinner,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import ky from "ky";
import React, { useEffect, useState } from "react";
import ReactDOM from "react-dom";
import { useForm } from "react-hook-form";
import { Link } from "react-router-dom";
import { useRecoilState, useRecoilValue } from "recoil";
import { Account } from "../models/account";
import { tokenState } from "../state/authState";
import { accountsState } from "../state/accountState";

export function Accounts() {
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [accounts, setAccounts] = useRecoilState(accountsState);
  const token = useRecoilValue(tokenState);

  useEffect(() => {
    const loadAccounts = async () => {
      try {
        const response = await ky.get("/api/account", {
          headers: {
            Authorization: "Bearer " + token,
          },
        });
        const accounts: Account[] = await response.json();
        setAccounts(accounts);
      } catch (error) {
        if (error instanceof ky.HTTPError) {
        }
      } finally {
        setLoading(false);
      }
    };
    loadAccounts();
  }, [token, setAccounts]);

  return (
    <div className="container mx-auto py-8">
      <div className="flex justify-between mb-8">
        <h1 className="text-3xl">Accounts</h1>
        <button className="bg-gray-900 text-white py-2 px-4 rounded shadow">
          +
        </button>
      </div>
      <div className="grid grid-cols-4 gap-8">
        {!loading &&
          accounts.map((account) => (
            <Link
              className="bg-white rounded-lg shadow p-6 border border-gray-200 hover:shadow-lg transition-shadow"
              to={"/account/" + account.id}
              key={account.id}
            >
              <div className="flex justify-between mb-4">
                <FontAwesomeIcon icon={faChartLine} size="lg" />
                <span className="bg-gray-300 uppercase font-bold text-xs rounded p-1">
                  Default
                </span>
              </div>
              <h2 className="text-lg font-bold">{account.name}</h2>
            </Link>
          ))}
        {!loading && (
          <button
            className="rounded-lg shadow p-6 border border-gray-200 flex items-center justify-center text-gray-400"
            onClick={() => setCreating(true)}
          >
            <FontAwesomeIcon icon={faPlus} size="2x" />
          </button>
        )}
        {loading && (
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
    <div className="bg-white rounded-lg shadow p-6 border border-gray-200 hover:shadow-lg transition-shadow animate-pulse">
      <div className="flex justify-between mb-4 items-center">
        <span className="bg-blue-100 uppercase font-bold text-xs rounded w-8 h-6"></span>
        <span className="bg-blue-100 uppercase font-bold text-xs rounded h-4 w-16"></span>
      </div>
      <div className="h-6 bg-blue-100"></div>
    </div>
  );
}

type CreateAccountModalProps = {
  close: () => void;
};

type CreateAccountForm = {
  name: string;
};

function CreateAccountModal({ close }: CreateAccountModalProps) {
  const { register, handleSubmit, formState, setError } =
    useForm<CreateAccountForm>();
  const token = useRecoilValue(tokenState);
  const [, setAccounts] = useRecoilState(accountsState);

  const onSubmit = async (values: CreateAccountForm) => {
    try {
      const response = await ky.post("/api/account", {
        json: values,
        headers: {
          Authorization: "Bearer " + token,
        },
      });

      const account: Account = await response.json();
      setAccounts((accounts) => [...accounts, account]);
      close();
    } catch (error) {
      if (error instanceof ky.HTTPError) {
        setError("name", {
          message: error.message,
        });
      }
    }
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
