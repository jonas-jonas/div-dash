import classNames from "classnames";
import ky from "ky";
import { useEffect } from "react";
import { useRecoilState, useRecoilValue } from "recoil";
import { Symbol } from "../models/symbol";
import { Balance } from "../models/balance";
import { tokenState } from "../state/authState";
import { balancesState } from "../state/balanceState";
import { formatMoney, formatPercent } from "../util/formatter";

export function PortfolioBalance() {
  const [balances, setBalances] = useRecoilState(balancesState);
  const token = useRecoilValue(tokenState);

  useEffect(() => {
    const loadBalance = async () => {
      try {
        const response = await ky.get("/api/balance", {
          headers: {
            Authorization: "Bearer " + token,
          },
        });
        const balances: Balance[] = await response.json();
        setBalances(balances);
      } catch (error) {
        console.error(error);
      }
    };
    loadBalance();
  }, [token, setBalances]);

  const getIconURL = (symbol: Symbol) => {
    switch (symbol.source) {
      case "iex":
        return "";
      case "binance":
        return (
          "https://cryptoicons.org/api/black/" +
          symbol.symbolID.toLowerCase() +
          "/20"
        );
    }
  };

  return (
    <div className="col-span-2">
      <table className="table w-full text-left">
        <thead className="bg-white">
          <tr className="shadow">
            <th className="rounded-l py-3 px-2">Symbol</th>
            <th className="py-3 px-2">Total</th>
            <th className="py-3 px-2">Buy In</th>
            <th className="py-3 px-2">Profit/Loss</th>
            <th className="rounded-r py-3 px-2"></th>
          </tr>
        </thead>
        <tbody>
          {balances?.map((balanceItem) => (
            <tr
              className="border-b border-gray-200"
              key={balanceItem.symbol.symbolID}
            >
              <td className="py-3 px-2 flex items-center">
                {balanceItem.symbol.type === "crypto" && (
                  <img
                    src={getIconURL(balanceItem.symbol)}
                    width="20"
                    height="20"
                    alt="BTC icon"
                    className="mr-2"
                  />
                )}
                <div className="flex flex-col">
                  <span>{balanceItem.symbol.symbolName || balanceItem.symbol.symbolID}</span>
                  <span className="text-sm text-gray-600">
                    {new Intl.NumberFormat("de-DE", {
                      minimumFractionDigits: balanceItem.symbol.precision,
                    }).format(balanceItem.amount)}
                  </span>
                </div>
              </td>
              <td className="py-3 px-2">
                {formatMoney(balanceItem.fiatValue)}
              </td>
              <td className="py-3 px-2">
                <span>{formatMoney(balanceItem.costBasis)}</span>
              </td>
              <td>
                <div className="flex flex-col items-start">
                  <span>{formatMoney(balanceItem.plAbsolute)}</span>
                  <span
                    className={classNames(
                      "text-sm text-white px-2 rounded-full",
                      {
                        "bg-red-600": balanceItem.plAbsolute < 0,
                        "bg-green-600": balanceItem.plAbsolute > 0,
                      }
                    )}
                  >
                    {formatPercent(balanceItem.plPercent)}
                  </span>
                </div>
              </td>
              <td className="py-3 px-2">
                <button className="text-blue-700 font-bold">Details</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
