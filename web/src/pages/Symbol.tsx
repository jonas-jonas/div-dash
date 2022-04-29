import {
  IconChevronDown,
  IconChevronRight,
  IconExternalLink,
  IconLink,
  IconLoader,
  IconRefresh,
} from "@tabler/icons";
import numeral from "numeral";
import { useQuery } from "react-query";
import { Link, useParams } from "react-router-dom";
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  XAxis,
  YAxis,
} from "recharts";
import { SymbolTypeLabels } from "../models/symbol";
import * as api from "../util/api";
import { formatMoney } from "../util/formatter";

type SymbolPageParams = {
  symbolId: string;
};

export function SymbolPage() {
  const { symbolId } = useParams<SymbolPageParams>() as SymbolPageParams;

  const { data: symbolDetails, isLoading: loadingSymbol } = useQuery(
    ["symbol", symbolId, "details"],
    () => api.getSymbolDetails(symbolId),
    {
      retry: false,
    }
  );

  const { data: chart, isLoading: loadingChart } = useQuery(
    ["symbol", symbolId, "chart"],
    () => api.getSymbolChart(symbolId),
    {
      retry: false,
    }
  );

  return (
    <div className="container mb-24 mx-auto pt-8">
      {loadingSymbol && (
        <div className="w-full py-14 flex items-center justify-center">
          <IconRefresh className="animate-spin" />
        </div>
      )}
      {!loadingSymbol && symbolDetails && (
        <div>
          <div className="text-sm mb-4 text-gray-500 flex">
            <a href="/assets" className="">
              Assets
            </a>
            <IconChevronRight className="mr-2 ml-2" size={20} />
            <Link to={"/symbols/" + symbolDetails.type} className="">
              {SymbolTypeLabels[symbolDetails.type]}
            </Link>
            <IconChevronRight className="mr-2 ml-2" size={20} />
            <span className="text-gray-900">{symbolDetails.name}</span>
          </div>
          <div className="flex justify-between border-b pb-8">
            <div>
              <h2 className="text-4xl font-bold text-gray-800 mb-4 flex items-center">
                {symbolDetails.images?.thumb && (
                  <img
                    src={symbolDetails.images?.thumb}
                    alt={symbolDetails.name + "Logo"}
                    className="mr-4 rounded-full w-16"
                  />
                )}
                <span>{symbolDetails.name}</span>
              </h2>
              <div className="flex">
                {symbolDetails.tags.map((tag) => {
                  if (tag.type === "CHIP") {
                    return (
                      <span
                        className="rounded-full bg-gray-300 text-gray-700 text-sm shadow px-4 py-1 mr-2"
                        key={tag.label}
                      >
                        {tag.label}
                      </span>
                    );
                  } else if (tag.type === "LINK") {
                    return (
                      <a
                        href={tag.link}
                        target="_blank"
                        rel="noreferrer noopener"
                        className="rounded-full bg-gray-300 text-gray-700 text-sm shadow px-4 py-1 hover:bg-gray-500 hover:text-white transition-colors inline-flex items-center"
                        key={tag.link}
                      >
                        <IconLink className="mr-1" size={16} />
                        {tag.link}
                        <IconExternalLink
                          className="ml-1 text-gray-400"
                          size={16}
                        />
                      </a>
                    );
                  }
                  return null;
                })}
              </div>
            </div>
            <div className="flex items-end flex-col">
              <h2 className="text-4xl font-bold text-gray-800 mb-2">
                54,89 €
                <span className="text-base font-normal ml-2">
                  <IconChevronDown className="inline" />
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
                {symbolDetails.indicators
                  .filter((indicator) => indicator.value > 0)
                  .map((indicator) => {
                    return (
                      <div
                        className="px-6 flex flex-col items-center"
                        key={indicator.label}
                      >
                        <span className="text-3xl font-bold flex items-center">
                          {numeral(indicator.value).format(indicator.format)}
                        </span>
                        <span className="text-gray-400">{indicator.label}</span>
                      </div>
                    );
                  })}
              </div>
              <div className="flex-grow-1 py-8 h-1/2">
                {loadingChart && (
                  <div className="flex items-center justify-center h-full">
                    <IconLoader className="animate-spin-slow" />
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
                <div
                  dangerouslySetInnerHTML={{
                    __html: symbolDetails.description,
                  }}
                ></div>
              </div>
              <div className="bg-white rounded p-8 shadow text-justify">
                <h3 className="text-2xl font-bold mb-2 text-gray-800">
                  Company Environment
                </h3>
                <a className="flex justify-between py-2" href="/sector/">
                  Apple - AAPL
                  <IconChevronRight />
                </a>
                <a className="flex justify-between py-2" href="/sector/">
                  Apple - AAPL
                  <IconChevronRight />
                </a>
                <a className="flex justify-between py-2" href="/sector/">
                  Apple - AAPL
                  <IconChevronRight />
                </a>
              </div>
            </div>
            <div className="bg-white rounded p-8 shadow text-justify">
              <h3 className="text-2xl font-bold mb-4 text-gray-800">Markets</h3>
              <a className="flex justify-between" href="/sector/">
                XETRA
                <IconChevronRight />
              </a>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
