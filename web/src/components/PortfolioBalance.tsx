import ky from "ky";
import { useEffect } from "react";
import { useRecoilState, useRecoilValue } from "recoil";
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
              key={balanceItem.asset.assetName}
            >
              <td className="py-3 px-2 flex items-center">
                <img
                  src={
                    "https://cryptoicons.org/api/black/" +
                    balanceItem.asset.assetName.toLowerCase() +
                    "/20"
                  }
                  width="20"
                  height="20"
                  alt="BTC icon"
                  className="mr-2"
                />
                <div className="flex flex-col">
                  <span>{balanceItem.asset.assetName}</span>
                  <span className="text-sm text-gray-600">
                    {new Intl.NumberFormat("de-DE", {
                      minimumFractionDigits: balanceItem.asset.precision,
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
                  <span className="text-sm text-white px-2 rounded-full bg-red-600">
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
