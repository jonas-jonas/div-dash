import {
  faChevronLeft,
  faChevronRight,
  faPlus,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import ky from "ky";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import { useRecoilValue } from "recoil";
import { tokenState } from "../state/authState";

type Transaction = {
  transactionId: string;
  symbol: string;
  type: string;
  transactionProvider: string;
  price: number;
  date: string;
  amount: number;
  side: "buy" | "sell";
};

type AccountParams = {
  accountId: string;
};

function formatDate(date: string) {
  const d = Date.parse(date);
  return new Intl.DateTimeFormat("de-DE", {
    dateStyle: "short",
  }).format(d);
}

function formatTime(date: string) {
  const d = Date.parse(date);
  return new Intl.DateTimeFormat("de-DE", {
    timeStyle: "medium",
  }).format(d);
}

function formatMoney(amount: number) {
  return new Intl.NumberFormat("de-DE", {
    style: "currency",
    currency: "EUR",
  }).format(amount);
}

export function Account() {
  const [transactions, setTransactions] = useState<Transaction[]>();
  const token = useRecoilValue(tokenState);

  const { accountId } = useParams<AccountParams>();

  useEffect(() => {
    const loadTransactions = async () => {
      try {
        const response = await ky.get(
          "/api/account/" + accountId + "/transaction",
          {
            headers: {
              Authorization: "Bearer " + token,
            },
          }
        );
        const transactions: Transaction[] = await response.json();
        setTransactions(transactions);
      } catch (error) {
        console.error(error);
      }
    };
    loadTransactions();
  }, [token, accountId]);

  return (
    <div className="container mx-auto pt-10">
      <div className="flex justify-between mb-4">
        <h1 className="text-2xl pl-4">Account</h1>
        <button className="bg-gray-900 text-white py-2 px-3 rounded shadow">
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
            {transactions?.map((transaction) => (
              <tr className="border-b border-gray-200">
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
                <td className="py-3 px-2 capitalize">{transaction.type}</td>
                <td className="py-3 px-2">{formatMoney(transaction.price)}</td>
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
          </tbody>
        </table>
        <div className="flex justify-end mt-4">
          <div className="bg-white rounded shadow p-2 text-gray-900">
            <button className="px-2 text-blue-700">
              <FontAwesomeIcon icon={faChevronLeft} size="sm" />
            </button>
            <button className="px-2 border border-blue-700 text-blue-700 rounded font-bold">1</button>
            <button className="px-2">2</button>
            <button className="px-2">3</button>
            <button className="px-2">...</button>
            <button className="px-2">22</button>
            <button className="px-2 text-blue-700">
              <FontAwesomeIcon icon={faChevronRight} size="sm" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
