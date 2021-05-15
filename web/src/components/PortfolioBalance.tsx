import ky from "ky";
import { useEffect } from "react";
import { useRecoilState, useRecoilValue } from "recoil";
import { Balance } from "../models/balance";
import { tokenState } from "../state/authState";
import { balancesState } from "../state/balanceState";
import { formatMoney } from "../util/formatter";

export function PortfolioBalance() {
  const [balances, setBalance] = useRecoilState(balancesState);
  const token = useRecoilValue(tokenState);

  useEffect(() => {
    const loadBalance = async () => {
      try {
        const response = await ky.get("/api/balance", {
          headers: {
            Authorization: "Bearer " + token,
          },
        });
        const balance: Balance[] = await response.json();
        setBalance(balance);
      } catch (error) {
        console.error(error);
      }
    };
    loadBalance();
  }, [token, setBalance]);

  return (
    <div className="col-span-2">
      <table className="table w-full text-left">
        <thead className="bg-white">
          <tr className="shadow">
            <th className="rounded-l py-3 px-2">Symbol</th>
            <th className="py-3 px-2">Amount</th>
            <th className="py-3 px-2">Total</th>
            <th className="py-3 px-2">Cost Basis</th>
            <th className="rounded-r py-3 px-2"></th>
          </tr>
        </thead>
        <tbody>
          {balances?.map((balanceItem) => (
            <tr className="border-b border-gray-200" key={balanceItem.symbol}>
              <td className="py-3 px-2 flex items-center">
                <img
                  src={
                    "https://cryptoicons.org/api/black/" +
                    balanceItem.symbol.toLowerCase() +
                    "/20"
                  }
                  width="20"
                  height="20"
                  alt="BTC icon"
                  className="mr-2"
                />{" "}
                {balanceItem.symbol}
              </td>
              <td className="py-3 px-2 capitalize">
                {new Intl.NumberFormat("de-DE", {
                  minimumFractionDigits: 4,
                  maximumFractionDigits: 4,
                }).format(balanceItem.amount)}
              </td>
              <td className="py-3 px-2">
                {formatMoney(balanceItem.costBasis * balanceItem.amount)}
              </td>
              <td className="py-3 px-2">
                {formatMoney(balanceItem.costBasis)}
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
