import { faSadTear, faSpinner } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import classNames from "classnames";
import ky from "ky";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useRecoilState } from "recoil";
import { Balance } from "../models/balance";
import { balancesState } from "../state/balanceState";
import { formatMoney, formatPercent } from "../util/formatter";

export function PortfolioBalance() {
  const [balance, setBalance] = useRecoilState(balancesState);
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const loadBalance = async () => {
      try {
        const response = await ky.get("/api/balance");
        const balance: Balance = await response.json();
        setBalance(balance);
      } catch (error) {
        if (error instanceof ky.HTTPError) {
          setError(error.message);
        } else {
          setError("General error");
        }
      } finally {
        setLoading(false);
      }
    };
    loadBalance();
  }, [setBalance]);

  return (
    <div className="col-span-2">
      <table className="table w-full text-left">
        <thead className="bg-white">
          <tr className="shadow">
            <th className="rounded-l py-3 px-2">Symbol</th>
            <th className="py-3 px-2">Total</th>
            <th className="py-3 px-2">Buy In</th>
            <th className="rounded-r py-3 px-2">Profit/Loss</th>
          </tr>
        </thead>
        {!loading && !error && (
          <tbody>
            {balance?.symbols.map((balanceItem) => (
              <tr
                className="border-b border-gray-200"
                key={balanceItem.symbol.symbolID}
              >
                <td className="py-3 px-2 flex items-center">
                  <div className="flex flex-col">
                    <Link
                      className="font-bold hover:underline"
                      to={"/symbol/" + balanceItem.symbol.symbolID}
                    >
                      {balanceItem.symbol.symbolName ||
                        balanceItem.symbol.symbolID}
                    </Link>
                    <span className="text-sm text-gray-600">
                      {new Intl.NumberFormat("de-DE", {
                        minimumFractionDigits: balanceItem.symbol.precision,
                      }).format(balanceItem.amount)}
                    </span>
                  </div>
                </td>
                <td className="py-3 px-2">
                  <div className="flex flex-col">
                    <span>{formatMoney(balanceItem.fiatAssetPrice)}</span>
                    <span className="text-sm text-gray-600">
                      {formatMoney(balanceItem.fiatValue)}
                    </span>
                  </div>
                </td>
                <td className="py-3 px-2">
                  <span>{formatMoney(balanceItem.costBasis)}</span>
                </td>
                <td>
                  <div className="flex flex-col items-start">
                    <span>{formatMoney(balanceItem.pnl.pnl)}</span>
                    <span
                      className={classNames(
                        "text-sm text-white px-2 rounded-full",
                        {
                          "bg-red-600": balanceItem.pnl.pnl < 0,
                          "bg-green-600": balanceItem.pnl.pnl > 0,
                        }
                      )}
                    >
                      {formatPercent(balanceItem.pnl.pnlPercent)}
                    </span>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        )}
      </table>
      {!loading && error && (
        <div className="flex items-center justify-center py-20 flex-col text-gray-500">
          <FontAwesomeIcon icon={faSadTear} size="2x" className="mb-3" />
          <span>There was an error while loading your assets</span>
          <b className="my-1">{error}</b>
          <span>Please try again later</span>
        </div>
      )}
      {loading && (
        <div className="flex items-center justify-center py-24 flex-col text-gray-500">
          <FontAwesomeIcon icon={faSpinner} spin size="2x" className="mb-3 " />
          <span className="">Loading Assets...</span>
        </div>
      )}
    </div>
  );
}
