import ky from "ky";
import { useEffect, useState } from "react";
import { useRecoilValue } from "recoil";
import { tokenState } from "../state/authState";
import { formatMoney } from "../util/formatter";

type Balance = {
  symbol: string;
  total: number;
  costBasis: number;
};

export function PortfolioBalance() {
  const [balances, setBalance] = useState<Balance[]>();
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
  }, [token]);

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
                }).format(balanceItem.total)}
              </td>
              <td className="py-3 px-2">
                {formatMoney(balanceItem.costBasis * balanceItem.total)}
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