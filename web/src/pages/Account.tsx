import {
  IconChevronLeft,
  IconChevronRight,
  IconFileUpload,
  IconLoader,
  IconPencil,
  IconTrash,
  IconX,
} from "@tabler/icons";
import classNames from "classnames";
import ky from "ky";
import { ChangeEvent, useEffect, useReducer, useState } from "react";
import ReactDOM from "react-dom";
import { Path, useForm, UseFormRegister } from "react-hook-form";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { useParams } from "react-router-dom";
import { TransactionForm } from "../form/TransactionForm";
import { Symbol, SymbolTypeLabels } from "../models/symbol";
import { Transaction } from "../models/transaction";
import * as api from "../util/api";
import {
  formatAmount,
  formatDate,
  formatMoney,
  formatTime,
} from "../util/formatter";

type AccountParams = {
  accountId: string;
};

type AccountImport = {
  file: File;
};

export function Account() {
  const [creating, setCreating] = useState(false);

  const queryClient = useQueryClient();

  const { accountId } = useParams<AccountParams>() as AccountParams;
  const { data: account } = useQuery(["account", accountId], () =>
    api.getAccount(accountId)
  );
  const { data: transactions, isLoading } = useQuery(
    ["account", accountId, "transactions"],
    () => api.getTransactions(accountId)
  );

  const accountImportMutation = useMutation<
    Transaction[],
    ky.HTTPError,
    AccountImport
  >((values) => api.postAccountImport(accountId, values), {
    onSuccess: (transactions) => {
      queryClient.setQueryData<Transaction[]>(
        ["account", accountId, "transactions"],
        transactions
      );
      queryClient.invalidateQueries("balance");
    },
    onError: (error, vars, ctx) => {
      // TODO: Add some kind of feedback here
      console.error(error);
    },
  });

  const getTransactionSideClasses = (side: "buy" | "sell") => {
    return classNames("font-bold py-3 pl-4 uppercase", {
      "text-green-600": side === "buy",
      "text-red-600": side === "sell",
    });
  };

  const handleImport = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      accountImportMutation.mutate({
        file: e.target.files[0],
      });
    }
  };

  return (
    <div className="container mx-auto pt-10">
      <div className="flex justify-between mb-4">
        <h1 className="text-2xl pl-4">{account?.name}</h1>
        <label className="p-2">
          <IconFileUpload />
          <input
            type="file"
            className="hidden"
            accept=".xls"
            onChange={handleImport}
          />
        </label>
      </div>
      <div className="w-full flex">
        <div className="flex-grow">
          <table className="table w-full text-left">
            <thead className="bg-white">
              <tr className="shadow">
                <th className="rounded-l py-3 pl-4">Direction</th>
                <th className="py-3 px-2">ID</th>
                <th className="py-3 px-2">Name</th>
                <th className="py-3 px-2">Type</th>
                <th className="py-3 px-2">Amount</th>
                <th className="py-3 px-2">Price</th>
                <th className="py-3 px-2">Total</th>
                <th className="py-3 px-2">Date</th>
                <th className="rounded-r py-3 px-3 text-right">
                  <button
                    className="bg-gray-900 text-white py-1 px-3 rounded "
                    onClick={() => setCreating(true)}
                  >
                    + Transaction
                  </button>
                </th>
              </tr>
            </thead>
            <tbody>
              {!isLoading &&
                transactions?.map((transaction) => (
                  <tr className="group" key={transaction.transactionId}>
                    <td className={getTransactionSideClasses(transaction.side)}>
                      {transaction.side}
                    </td>
                    <td className="py-3 px-2">
                      <span className="font-mono font-bold tracking-wider">
                        {transaction.transactionId}
                      </span>
                    </td>
                    <td className="py-3 px-2">
                      <div className=" flex items-center">
                        <span className="font-bold">{transaction.symbol}</span>
                        <span className="ml-1 bg-gray-900 text-white rounded px-1 text-sm">
                          {transaction.symbol}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-2 capitalize">
                      {SymbolTypeLabels[transaction.type]}
                    </td>
                    <td className="py-3 px-2">
                      {formatAmount(transaction.amount)}
                    </td>
                    <td className="py-3 px-2">
                      {formatMoney(transaction.price)}
                    </td>
                    <td className="py-3 px-2">
                      {formatMoney(transaction.amount * transaction.price)}
                    </td>
                    <td className="py-3 px-2 flex flex-col">
                      <span className="text-sm">
                        {formatDate(transaction.date)}
                      </span>
                      <span className="text-xs text-gray-700">
                        {formatTime(transaction.date)}
                      </span>
                    </td>
                    <td className="py-3 px-3 text-right">
                      <button className="text-gray-700 opacity-0 group-hover:opacity-100 px-2">
                        <IconPencil />
                      </button>
                      <button className="text-gray-700 opacity-0 group-hover:opacity-100 px-2">
                        <IconTrash />
                      </button>
                    </td>
                  </tr>
                ))}
              {isLoading && (
                <>
                  <TransactionRowLoadingIndicator />
                  <TransactionRowLoadingIndicator />
                  <TransactionRowLoadingIndicator />
                  <TransactionRowLoadingIndicator />
                </>
              )}
            </tbody>
          </table>
          {!isLoading && transactions?.length === 0 && (
            <div className="text-center mt-4">
              You have no transactions.{" "}
              <button
                className="text-blue-700"
                onClick={() => setCreating(true)}
              >
                Create one now
              </button>
            </div>
          )}
          {!isLoading && transactions && transactions.length > 0 && (
            <div className="flex justify-end mt-4">
              <div className="bg-white rounded shadow p-2 text-gray-900 flex items-center">
                <button className="px-2 text-blue-700">
                  <IconChevronLeft size={20} stroke={3} />
                </button>
                <button className="px-2 border border-blue-700 text-blue-700 rounded font-bold">
                  1
                </button>
                <button className="px-2">2</button>
                <button className="px-2">3</button>
                <button className="px-2">...</button>
                <button className="px-2">22</button>
                <button className="px-2 text-blue-700">
                  <IconChevronRight size={20} stroke={3} />
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
      {creating &&
        ReactDOM.createPortal(
          <CreateTransactionModal
            close={() => setCreating(false)}
            accountId={accountId}
          />,
          document.body
        )}
    </div>
  );
}

function TransactionRowLoadingIndicator() {
  return (
    <tr className="border-b border-gray-200 animate-pulse">
      <td className="py-3 px-2">
        <span className="h-5 w-20 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2 flex flex-col">
        <span className="h-4 w-16 bg-blue-100 block rounded mb-1"></span>
        <span className="h-3 w-14 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2">
        <span className="h-5 w-16 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2 capitalize">
        <span className="h-5 w-12 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2">
        <span className="h-5 w-12 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2">
        <span className="h-5 w-8 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2">
        <span className="h-5 w-12 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2 uppercase">
        <span className="h-5 w-10 bg-blue-100 block rounded"></span>
      </td>
      <td className="py-3 px-2">
        <span className="h-5 w-16 bg-blue-100 block rounded"></span>
      </td>
    </tr>
  );
}

type CreateTransactionModalProps = {
  close: () => void;
  accountId: string;
};

function CreateTransactionModal({
  close,
  accountId,
}: CreateTransactionModalProps) {
  const { register, handleSubmit, formState, setValue } =
    useForm<TransactionForm>({
      defaultValues: { side: "buy" },
    });
  const [error, setError] = useState<string>();
  const queryClient = useQueryClient();

  const createTransactionMutation = useMutation<
    Transaction,
    ky.HTTPError,
    TransactionForm
  >((values) => api.postTransaction(accountId, values), {
    onError: async (error) => {
      const json = await error.response.json();
      setError(json.message);
    },
    onSuccess: (transaction) => {
      queryClient.setQueryData<Transaction[]>(
        ["account", accountId, "transactions"],
        (transactions) => {
          if (transactions) {
            return [...transactions, transaction];
          }
          return [transaction];
        }
      );
      queryClient.invalidateQueries("balance");
      close();
    },
  });

  const onSubmit = async (values: TransactionForm) => {
    createTransactionMutation.mutate(values);
  };

  return (
    <div className="top-0 fixed">
      <form
        className="w-96 mx-auto bg-white rounded fixed top-1/4 transform z-10 left-1/2 -translate-x-1/2 shadow"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="px-8 pt-8 flex justify-between items-start">
          <div className="mr-2">
            <h2 className="text-2xl font-bold">New Transaction</h2>
            <h3 className="text-gray-500">Lorem ipsum, dolor sit </h3>
          </div>
          <button
            onClick={close}
            className="flex flex-col items-center text-gray-500 hover:text-gray-900 transition-colors"
            type="reset"
          >
            <IconX />
            <span className="text-xs">ESC</span>
          </button>
        </div>
        <div className="px-8 pt-4 pb-8">
          <div className="w-full flex mb-4 justify-around">
            <label className="text-center cursor-pointer hover:bg-green-50 py-3 transition-colors w-36 bg-gray-50 flex flex-col items-center">
              <input
                type="radio"
                value="buy"
                {...register("side", { required: true })}
                className="mb-1"
              />
              <span>Buy</span>
            </label>
            <label className="text-center cursor-pointer hover:bg-red-50 py-3 transition-colors w-36 bg-gray-50 flex flex-col items-center">
              <input
                type="radio"
                value="sell"
                {...register("side", { required: true })}
                className="mb-1"
              />
              <span>Sell</span>
            </label>
          </div>
          <label className="block mb-4">
            <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
              Symbol
            </span>
            <TypeAheadSymbolInput
              formKey="symbol"
              register={register}
              close={(symbol) => {
                setValue("symbol", symbol.symbolID);
                setValue("type", symbol.type);
              }}
              autoComplete="off"
            />
          </label>
          <label className="block mb-4">
            <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
              Type
            </span>
            <select
              {...register("type", { required: true })}
              className="block w-full px-4 py-2 focus:border-blue-700 transition-colors bg-gray-50 rounded"
            >
              {Object.entries(SymbolTypeLabels).map(([key, label]) => {
                return (
                  <option value={key} key={key}>
                    {label}
                  </option>
                );
              })}
            </select>
          </label>
          <div className="flex items-center justify-center">
            <label className="block mb-4 mr-2">
              <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
                Amount
              </span>
              <input
                type="number"
                step="0.000000001"
                className="block w-full px-4 py-2 focus:border-blue-700 transition-colors bg-gray-50 rounded"
                {...register("amount", { required: true })}
              />
            </label>
            <label className="block mb-4 ml-2">
              <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
                Price
              </span>
              <input
                type="number"
                step="0.000000001"
                className="block w-full px-4 py-2 focus:border-blue-700 transition-colors bg-gray-50 rounded"
                {...register("price", { required: true })}
              />
            </label>
          </div>
          <label className="block mb-4">
            <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
              Date
            </span>
            <input
              type="datetime-local"
              className="block w-full px-4 py-2 focus:border-blue-700 transition-colors bg-gray-50 rounded"
              {...register("date", { required: true })}
            />
          </label>
          {error && (
            <div className="bg-red-300 rounded shadow p-2 text-sm text-red-900 mb-4">
              {error}
            </div>
          )}
          <div className="flex justify-items-stretch mt-8">
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
                <IconLoader className="animate-spin-slow" />
              ) : (
                "Create"
              )}
            </button>
          </div>
        </div>
      </form>
      <div className="fixed w-full h-full top-0 z-0 backdrop-filter backdrop-blur-sm backdrop-opacity-45 backdrop-brightness-50"></div>
    </div>
  );
}

type TypeAheadSymbolReducerState = {
  loading: boolean;
  searchResults: Symbol[];
  error?: string;
  show: boolean;
};

type TypeAheadSymbolReducerAction =
  | TypeAheadSymbolReducerActionLoad
  | TypeAheadSymbolReducerActionError
  | TypeAheadSymbolReducerActionFinished
  | TypeAheadSymbolReducerActionHide
  | TypeAheadSymbolReducerActionShow;

type TypeAheadSymbolReducerActionLoad = {
  type: "LOAD";
};
type TypeAheadSymbolReducerActionError = {
  type: "ERROR";
  error: string;
};
type TypeAheadSymbolReducerActionFinished = {
  type: "FINISHED";
  payload: Symbol[];
};
type TypeAheadSymbolReducerActionHide = {
  type: "HIDE";
};
type TypeAheadSymbolReducerActionShow = {
  type: "SHOW";
};

function TypeAheadSymbolReducer(
  state: TypeAheadSymbolReducerState,
  action: TypeAheadSymbolReducerAction
): TypeAheadSymbolReducerState {
  switch (action.type) {
    case "LOAD":
      return { show: true, loading: true, searchResults: [] };
    case "FINISHED":
      return { show: true, loading: false, searchResults: action.payload };
    case "ERROR":
      return {
        show: true,
        loading: false,
        error: action.error,
        searchResults: [],
      };
    case "HIDE":
      return {
        show: false,
        loading: state.loading,
        searchResults: state.searchResults,
        error: state.error,
      };
    case "SHOW":
      return {
        show: true,
        loading: state.loading,
        searchResults: state.searchResults,
        error: state.error,
      };
  }
}

type TypeAheadSymbolInputProps<T> = {
  formKey: Path<T>;
  register: UseFormRegister<T>;
  close: (value: Symbol) => void;
} & React.InputHTMLAttributes<HTMLInputElement>;

function TypeAheadSymbolInput<T>({
  formKey,
  register,
  close,
  ...rest
}: TypeAheadSymbolInputProps<T>) {
  const [searchDebounce, setSearchDebounce] = useState(-1);
  const [state, dispatch] = useReducer(TypeAheadSymbolReducer, {
    loading: false,
    show: false,
    searchResults: [],
  });
  const [selectedEntry, setSelectedEntry] = useState(0);

  const onSymbolChange = async (evt: ChangeEvent<HTMLInputElement>) => {
    if (searchDebounce >= 0) {
      clearTimeout(searchDebounce);
      setSearchDebounce(-1);
    }
    if (!evt.target.value) {
      return;
    }
    dispatch({ type: "LOAD" });
    const debounceTimer = window.setTimeout(async () => {
      try {
        const resp = await ky.get("/api/symbol/search", {
          searchParams: {
            query: evt.target.value || "",
            count: 5,
          },
        });
        const symbols: Symbol[] = await resp.json();
        dispatch({ type: "FINISHED", payload: symbols });
      } catch (error) {
        if (error instanceof ky.HTTPError) {
          dispatch({ type: "ERROR", error: error.message });
        }
      }
    }, 500);
    setSearchDebounce(debounceTimer);
  };

  const setValue = (result: any) => () => {
    dispatch({ type: "HIDE" });
    close(result);
  };

  useEffect(() => {
    function keyListener(e: KeyboardEvent) {
      switch (e.key) {
        case "ArrowUp":
          setSelectedEntry((selectedEntry) => selectedEntry - 1);
          break;
        case "ArrowDown":
          setSelectedEntry((selectedEntry) => selectedEntry + 1);
          break;
        case "Escape":
          dispatch({ type: "HIDE" });
          break;
      }
    }
    document.addEventListener("keyup", keyListener);
    return () => {
      document.removeEventListener("keyup", keyListener);
    };
  }, [state.searchResults, close, selectedEntry]);

  return (
    <div className="relative">
      <input
        type="text"
        placeholder="Search..."
        className={classNames(
          "block w-full px-4 py-2 focus:border-blue-700 transition-colors bg-gray-50 rounded"
        )}
        {...register(formKey)}
        onChange={onSymbolChange}
        id="typeahead-symbol-input"
        {...rest}
      />
      {state.show && (
        <div className="absolute bg-white w-full shadow-xl mt-1">
          {state.searchResults.map((result, i) => (
            <button
              className={classNames(
                "px-2 py-1 flex justify-start items-center hover:bg-gray-50 w-full text-left transition-colors duration-75 hover:border-blue-700 border-l-4 border-transparent",
                { "bg-gray-100": selectedEntry === i }
              )}
              key={result.symbolID}
              onClick={setValue(result)}
            >
              <span className="flex-shrink">{result.symbolName}</span>
              <span className="ml-1 bg-gray-900 text-white rounded px-1 text-sm">
                {result.symbolID}
              </span>
            </button>
          ))}
          {state.loading && (
            <div className="text-center py-2">
              <IconLoader className="animate-spin-slow" />
            </div>
          )}
          <button
            onClick={() => dispatch({ type: "HIDE" })}
            className="w-full text-center font-bold mt-2 py-1 flex justify-center"
          >
            <IconX className="mr-2" />
            Close
          </button>
        </div>
      )}
    </div>
  );
}
