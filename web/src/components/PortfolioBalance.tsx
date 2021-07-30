import {
  faCopy,
  faSadTear,
  faSpinner,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import classNames from "classnames";
import ky from "ky";
import { useQuery } from "react-query";
import { Link } from "react-router-dom";
import { Balance } from "../models/balance";
import * as api from "../util/api";
import { formatMoney, formatPercent } from "../util/formatter";

export function PortfolioBalance() {
  const {
    data: balance,
    isLoading,
    error,
  } = useQuery<Balance, ky.HTTPError>("balance", api.getBalance);

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
        {!isLoading && !error && (
          <tbody>
            {balance?.symbols.map((balanceItem) => (
              <tr
                className="border-b border-gray-200"
                key={balanceItem.symbol.symbolID}
              >
                <td className="py-3 px-2 flex items-center">
                  <div className="flex flex-col">
                    <div className="text-gray-600 text-xs mb-1">
                      <CopyElement value={balanceItem.symbol.symbolID} />
                      <span className="mx-2">·</span>
                      <CopyElement value={balanceItem.symbol.isin} />
                      <span className="mx-2">·</span>
                      <CopyElement value={balanceItem.symbol.wkn} />
                    </div>
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
                    <span>
                      {formatMoney(
                        balanceItem.fiatAssetPrice * balanceItem.amount
                      )}
                    </span>
                    <span className="text-sm text-gray-600">
                      {formatMoney(balanceItem.fiatAssetPrice)}
                    </span>
                  </div>
                </td>
                <td className="py-3 px-2">
                  <div className="flex flex-col">
                    <span>
                      {formatMoney(balanceItem.costBasis * balanceItem.amount)}
                    </span>
                    <span className="text-sm text-gray-600">
                      {formatMoney(balanceItem.costBasis)}
                    </span>
                  </div>
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
      {!isLoading && error && (
        <div className="flex items-center justify-center py-20 flex-col text-gray-500">
          <FontAwesomeIcon icon={faSadTear} size="2x" className="mb-3" />
          <span>There was an error while loading your assets</span>
          <b className="my-1">{error.message}</b>
          <span>Please try again later</span>
        </div>
      )}
      {isLoading && (
        <div className="flex items-center justify-center py-24 flex-col text-gray-500">
          <FontAwesomeIcon icon={faSpinner} spin size="2x" className="mb-3 " />
          <span className="">Loading Assets...</span>
        </div>
      )}
    </div>
  );
}

type CopyElementProps = {
  value: string;
};

function CopyElement({ value }: CopyElementProps) {
  const handleClick = () => {
    const textField = document.createElement("textarea");
    textField.innerText = value;
    document.body.appendChild(textField);
    textField.select();
    document.execCommand("copy");
    textField.remove();
  };

  return (
    <button
      className="focus:outline-none hover:text-gray-900"
      onClick={handleClick}
    >
      <FontAwesomeIcon icon={faCopy} className="mr-1" />
      {value}
    </button>
  );
}
