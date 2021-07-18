import {
  faChevronLeft,
  faChevronRight,
  faPlus,
  faSpinner,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import classNames from "classnames";
import ky from "ky";
import { ChangeEvent, useEffect, useReducer, useState } from "react";
import ReactDOM from "react-dom";
import { Path, useForm, UseFormRegister } from "react-hook-form";
import { useParams } from "react-router";
import { useRecoilState, useRecoilValue } from "recoil";
import { Symbol, SymbolType, SymbolTypeLabels } from "../models/symbol";
import { Transaction } from "../models/transaction";
import { accountByIdSelector } from "../state/accountState";
import { transactionsState } from "../state/transactionState";
import { formatDate, formatMoney, formatTime } from "../util/formatter";

type AccountParams = {
  accountId: string;
};

export function Account() {
  const [loading, setLoading] = useState(true);
  const [transactions, setTransactions] = useRecoilState(transactionsState);
  const [creating, setCreating] = useState(false);

  const { accountId } = useParams<AccountParams>();
  const account = useRecoilValue(accountByIdSelector(accountId));

  useEffect(() => {
    const loadTransactions = async () => {
      try {
        const response = await ky.get(
          "/api/account/" + accountId + "/transaction"
        );
        const transactions: Transaction[] = await response.json();
        setTransactions(transactions);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    };
    loadTransactions();
  }, [accountId, setTransactions]);

  return (
    <div className="container mx-auto pt-10">
      <div className="flex justify-between mb-4">
        <h1 className="text-2xl pl-4">{account?.name}</h1>
        <button
          className="bg-gray-900 text-white py-2 px-3 rounded shadow"
          onClick={() => setCreating(true)}
        >
          <FontAwesomeIcon icon={faPlus} />
        </button>
      </div>
      <div className="w-full">
        <table className="table w-full text-left">
          <thead className="bg-white">
            <tr className="shadow">
              <th className="rounded-l py-3 pl-4">Id</th>
              <th className="py-3 px-2">Date</th>
              <th className="py-3 px-2">Symbol</th>
              <th className="py-3 px-2">Type</th>
              <th className="py-3 px-2">Price</th>
              <th className="py-3 px-2">Amount</th>
              <th className="py-3 px-2">Total</th>
              <th className="py-3 px-2">Side</th>
              <th className="rounded-r py-3 px-2"></th>
            </tr>
          </thead>
          <tbody>
            {!loading &&
              transactions?.map((transaction) => (
                <tr
                  className="border-b border-gray-200"
                  key={transaction.transactionId}
                >
                  <td className="py-3 px-2">
                    <span className="font-mono font-bold tracking-wider text-blue-700">
                      {transaction.transactionId}
                    </span>
                  </td>
                  <td className="py-3 px-2 flex flex-col">
                    <span className="text-sm">
                      {formatDate(transaction.date)}
                    </span>
                    <span className="text-xs text-gray-700">
                      {formatTime(transaction.date)}
                    </span>
                  </td>
                  <td className="py-3 px-2">{transaction.symbol}</td>
                  <td className="py-3 px-2 capitalize">
                    {SymbolTypeLabels[transaction.type]}
                  </td>
                  <td className="py-3 px-2">
                    {formatMoney(transaction.price)}
                  </td>
                  <td className="py-3 px-2">{transaction.amount}</td>
                  <td className="py-3 px-2">
                    {formatMoney(transaction.amount * transaction.price)}
                  </td>
                  <td className="py-3 px-2 uppercase">{transaction.side}</td>
                  <td className="py-3 px-2">
                    <button className="text-blue-700 font-bold">Details</button>
                  </td>
                </tr>
              ))}
            {loading && (
              <>
                <TransactionRowLoadingIndicator />
                <TransactionRowLoadingIndicator />
                <TransactionRowLoadingIndicator />
                <TransactionRowLoadingIndicator />
              </>
            )}
          </tbody>
        </table>
        {!loading && transactions.length === 0 && (
          <div className="text-center mt-4">
            You have no transactions.{" "}
            <button className="text-blue-700" onClick={() => setCreating(true)}>
              Create one now
            </button>
          </div>
        )}
        {!loading && transactions.length > 0 && (
          <div className="flex justify-end mt-4">
            <div className="bg-white rounded shadow p-2 text-gray-900">
              <button className="px-2 text-blue-700">
                <FontAwesomeIcon icon={faChevronLeft} size="sm" />
              </button>
              <button className="px-2 border border-blue-700 text-blue-700 rounded font-bold">
                1
              </button>
              <button className="px-2">2</button>
              <button className="px-2">3</button>
              <button className="px-2">...</button>
              <button className="px-2">22</button>
              <button className="px-2 text-blue-700">
                <FontAwesomeIcon icon={faChevronRight} size="sm" />
              </button>
            </div>
          </div>
        )}
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

type CreateTransactionForm = {
  symbol: string;
  type: SymbolType;
  transactionProvider: string;
  price: string;
  date: string;
  amount: string;
  side: "buy" | "sell";
};

function CreateTransactionModal({
  close,
  accountId,
}: CreateTransactionModalProps) {
  const { register, handleSubmit, formState, setValue } =
    useForm<CreateTransactionForm>();
  const [error, setError] = useState<string>();
  const [, setTransactions] = useRecoilState(transactionsState);

  const onSubmit = async (values: CreateTransactionForm) => {
    const date = new Date(values.date);
    try {
      const response = await ky.post(
        "/api/account/" + accountId + "/transaction",
        {
          json: {
            ...values,
            date: date.toISOString(),
            amount: parseFloat(values.amount),
            price: parseFloat(values.price),
            transactionProvider: "binance",
          },
        }
      );
      const transaction: Transaction = await response.json();
      setTransactions((transactions) => [...transactions, transaction]);
      close();
    } catch (error) {
      if (error instanceof ky.HTTPError) {
        const json = await error.response.json();
        setError(json.message);
      }
    }
  };

  return (
    <div className="top-0 fixed">
      <form
        className="w-1/2 mx-auto bg-gray-50 rounded fixed top-1/4 transform z-10 left-1/2 -translate-x-1/2 shadow"
        onSubmit={handleSubmit(onSubmit)}
      >
        <div className="border-b border-gray-200 px-8 py-4 flex justify-between">
          <h2 className="text-xl">New Transaction</h2>
          <button onClick={close} className="px-2">
            <FontAwesomeIcon icon={faTimes} />
          </button>
        </div>
        <div className="px-8 py-4">
          <div className="rounded bg-white shadow w-full flex mb-4">
            <label className="w-full text-center cursor-pointer hover:bg-green-50 py-3 transition-colors">
              Buy
              <input
                type="radio"
                value="buy"
                {...register("side", { required: true })}
                className="ml-3"
              />
            </label>
            <label className="w-full text-center cursor-pointer hover:bg-red-50 py-3 transition-colors">
              Sell
              <input
                type="radio"
                value="sell"
                {...register("side", { required: true })}
                className="ml-3"
              />
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
              className="block w-full px-4 py-2 focus:bg-white rounded shadow focus:border-blue-700 transition-colors"
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
          <label className="block mb-4">
            <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
              Amount
            </span>
            <input
              type="number"
              step="0.00001"
              className="block w-full px-4 py-2 focus:bg-white rounded shadow focus:border-blue-700 transition-colors"
              {...register("amount", { required: true })}
            />
          </label>
          <label className="block mb-4">
            <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
              Price
            </span>
            <input
              type="number"
              step="0.00001"
              className="block w-full px-4 py-2 focus:bg-white rounded shadow focus:border-blue-700 transition-colors"
              {...register("price", { required: true })}
            />
          </label>
          <label className="block mb-4">
            <span className="text-xs text-gray-700 ml-4 font-bold tracking-wider">
              Date
            </span>
            <input
              type="datetime-local"
              className="block w-full px-4 py-2 focus:bg-white rounded shadow focus:border-blue-700 transition-colors"
              {...register("date", { required: true })}
            />
          </label>
          {error && (
            <div className="bg-red-300 rounded shadow p-2 text-sm text-red-900 mb-4">
              {error}
            </div>
          )}
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
        dispatch({ type: "ERROR", error: error });
      }
    }, 500);
    setSearchDebounce(debounceTimer);
  };

  const setValue = (result: any) => () => {
    close(result);
    dispatch({ type: "HIDE" });
  };

  useEffect(() => {
    function clickListener(e: MouseEvent) {
      if (
        e.target &&
        "id" in e.target &&
        (e.target as Element).id === "typeahead-symbol-input"
      ) {
        return;
      }
      dispatch({ type: "HIDE" });
    }
    document.addEventListener("click", clickListener);
    return () => {
      document.removeEventListener("click", clickListener);
    };
  }, []);

  const show = () => {
    dispatch({ type: "SHOW" });
  };

  return (
    <div className="relative">
      <input
        type="text"
        placeholder="BTC"
        className={classNames(
          "block w-full px-4 py-2 focus:bg-white shadow focus:border-blue-700 transition-colors",
          { "shadow-xl rounded-t": state.show, rounded: !state.show }
        )}
        {...register(formKey)}
        onChange={onSymbolChange}
        onClick={show}
        id="typeahead-symbol-input"
        {...rest}
      />
      {state.show && (
        <div className="absolute bg-white w-full shadow-xl rounded-b border-t px-2 py-3">
          {state.searchResults.map((result) => (
            <button
              className="px-2 py-2 flex justify-start items-center hover:bg-gray-100 w-full text-left rounded-lg transition-colors duration-75"
              key={result.symbolID}
              onClick={setValue(result)}
            >
              <span className="rounded px-1 bg-gray-300 text-sm font-bold mr-2 whitespace-nowrap">
                {SymbolTypeLabels[result.type]}
              </span>
              <span className="flex-shrink">
                {result.symbolID} - {result.symbolName}
              </span>
            </button>
          ))}
          {state.loading && (
            <div>
              <FontAwesomeIcon icon={faSpinner} spin />
            </div>
          )}
        </div>
      )}
    </div>
  );
}
