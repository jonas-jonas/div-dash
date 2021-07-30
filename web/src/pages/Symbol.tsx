import {
  faChevronDown,
  faChevronRight,
  faExternalLinkAlt,
  faLink,
  faSpinner,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import numeral from "numeral";
import { useQuery } from "react-query";
import { useParams } from "react-router-dom";
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  XAxis,
  YAxis,
} from "recharts";
import * as api from "../util/api";
import { formatMoney } from "../util/formatter";

type SymbolPageParams = {
  symbolId: string;
};

export function SymbolPage() {
  const { symbolId } = useParams<SymbolPageParams>();

  const { data: symbolDetails, isLoading: loadingSymbol } = useQuery(
    ["symbol", symbolId, "details"],
    () => api.getSymbolDetails(symbolId)
  );
  const { data: chart, isLoading: loadingChart } = useQuery(
    ["symbol", symbolId, "chart"],
    () => api.getSymbolChart(symbolId)
  );

  return (
    <div className="container mb-24 mx-auto pt-8">
      {loadingSymbol && (
        <div className="w-full py-14 flex items-center justify-center">
          <FontAwesomeIcon icon={faSpinner} spin />
        </div>
      )}
      {!loadingSymbol && symbolDetails && (
        <div>
          <div className="text-sm mb-4 text-gray-500">
            <a href="/assets" className="">
              Assets
            </a>
            <FontAwesomeIcon
              icon={faChevronRight}
              className="mr-2 ml-2"
              size="xs"
            />
            <a href="/assets/common-stock" className="">
              Common Stock
            </a>
            <FontAwesomeIcon
              icon={faChevronRight}
              className="mr-2 ml-2"
              size="xs"
            />
            <span className="text-gray-900">{symbolDetails.name}</span>
          </div>
          <div className="flex justify-between border-b pb-8">
            <div>
              <h2 className="text-4xl font-bold text-gray-800 mb-4">
                {symbolDetails.name}
              </h2>
              {symbolDetails.tags.map((tag) => {
                if (tag.type === "CHIP") {
                  return (
                    <span className="rounded-full bg-gray-300 text-gray-700 text-sm shadow px-4 py-1 mr-2">
                      {tag.label}
                    </span>
                  );
                } else if (tag.type === "LINK") {
                  return (
                    <a
                      href={tag.link}
                      target="_blank"
                      rel="noreferrer noopener"
                      className="rounded-full bg-gray-300 text-gray-700 text-sm shadow px-4 py-1 hover:bg-gray-500 hover:text-white transition-colors"
                    >
                      <FontAwesomeIcon
                        icon={faLink}
                        className="mr-1"
                        size="sm"
                      />
                      {tag.link}
                      <FontAwesomeIcon
                        icon={faExternalLinkAlt}
                        className="ml-1 text-gray-400"
                        size="sm"
                      />
                    </a>
                  );
                }
                return null;
              })}
            </div>
            <div className="flex items-end flex-col">
              <h2 className="text-4xl font-bold text-gray-800 mb-2">
                54,89 €
                <span className="text-base font-normal ml-2">
                  <FontAwesomeIcon icon={faChevronDown} />
                  -4,23%
                </span>
              </h2>
              <span className="rounded-full bg-gray-300 text-gray-700 text-sm shadow px-4 py-1">
                XETRA
              </span>
            </div>
          </div>
          <div className="grid grid-cols-3 gap-6 py-8">
            <div className="col-span-2 row-span-2 flex flex-col">
              <div className="flex bg-white rounded shadow text-gray-700 py-3 justify-evenly">
                <div className="px-6 flex flex-col items-center">
                  <span className="text-3xl font-bold flex items-center">
                    {numeral(symbolDetails.marketCap).format("$0.00 a")}
                  </span>
                  <span className="text-gray-400">Market Cap</span>
                </div>
                <div className="px-6 flex flex-col items-center">
                  <span className="text-3xl font-bold flex items-center">
                    {numeral(symbolDetails.peRatio).format("0.00")}
                  </span>
                  <span className="text-gray-400">PE Ratio</span>
                </div>
                <div className="px-6 flex flex-col items-center">
                  <span className="text-3xl font-bold flex items-center">
                    {numeral(symbolDetails.dividendYield).format("0.00%")}
                  </span>
                  <span className="text-gray-400">Dividend Yield</span>
                </div>
                <div className="px-6 flex flex-col items-center">
                  <span className="text-3xl font-bold flex items-center">
                    {numeral(symbolDetails.eps).format("$0.00")}
                  </span>
                  <span className="text-gray-400">EPS</span>
                </div>
              </div>
              <div className="flex-grow-1 py-8 h-1/2">
                {loadingChart && (
                  <div className="flex items-center justify-center h-full">
                    <FontAwesomeIcon icon={faSpinner} spin />
                  </div>
                )}
                {!loadingChart && chart && chart.length > 0 && (
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart width={300} height={100} data={chart}>
                      <Line
                        type="monotone"
                        dataKey="price"
                        stroke="#121826"
                        strokeWidth={2}
                        dot={false}
                      />
                      <CartesianGrid strokeDasharray="3" vertical={false} />

                      <XAxis dataKey="date" />
                      <YAxis
                        domain={[
                          (dataMin: number) => dataMin * 0.98,
                          (dataMax: number) => dataMax * 1.02,
                        ]}
                        tickFormatter={(x) => {
                          return formatMoney(x);
                        }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                )}
              </div>
              <div>
                <div className="border-t-2 relative h-16 mb-4">
                  <div
                    className="absolute flex items-center flex-col"
                    style={{ left: "15%" }}
                  >
                    <span className="transform -translate-y-2 w-1 h-4 bg-gray-700 block"></span>
                    <span className="font-bold text-gray-800">Ex-Div</span>
                    <span className="text-gray-400">08.02.2019</span>
                  </div>
                  <div
                    className="absolute flex items-center flex-col"
                    style={{ left: "25%" }}
                  >
                    <span className="transform -translate-y-2 w-1 h-4 bg-gray-700 block"></span>
                    <span className="font-bold text-gray-800">Dividend</span>
                    <span className="text-gray-400">01.03.2019</span>
                  </div>

                  <div
                    className="absolute flex items-center flex-col"
                    style={{ left: "0%" }}
                  >
                    <span className="transform -translate-y-2 w-1 h-4 bg-gray-700 block"></span>
                    <span className="font-bold text-gray-800">Earnings</span>
                    <span className="text-gray-400">01.01.2019</span>
                  </div>
                </div>
              </div>
            </div>
            <div>
              <div className="bg-white rounded p-8 shadow text-justify mb-6">
                <h3 className="text-2xl font-bold mb-4 text-gray-800">
                  Your Holdings
                </h3>
                8 @ 123,03 €
              </div>
              <div className="bg-white rounded p-8 shadow text-justify mb-6">
                <h3 className="text-2xl font-bold mb-4 text-gray-800">
                  Company Summary
                </h3>
                {symbolDetails.description}
              </div>
              <div className="bg-white rounded p-8 shadow text-justify">
                <h3 className="text-2xl font-bold mb-2 text-gray-800">
                  Company Environment
                </h3>
                <a className="flex justify-between py-2" href="/sector/">
                  Apple - AAPL
                  <FontAwesomeIcon icon={faChevronRight} />
                </a>
                <a className="flex justify-between py-2" href="/sector/">
                  Apple - AAPL
                  <FontAwesomeIcon icon={faChevronRight} />
                </a>
                <a className="flex justify-between py-2" href="/sector/">
                  Apple - AAPL
                  <FontAwesomeIcon icon={faChevronRight} />
                </a>
              </div>
            </div>
            <div className="bg-white rounded p-8 shadow text-justify">
              <h3 className="text-2xl font-bold mb-4 text-gray-800">Markets</h3>
              <a className="flex justify-between" href="/sector/">
                XETRA
                <FontAwesomeIcon icon={faChevronRight} />
              </a>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
