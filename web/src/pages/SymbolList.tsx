import {
  IconChevronLeft,
  IconChevronRight,
  IconLoader,
  IconMoodSad,
} from "@tabler/icons";
import ky from "ky";
import numeral from "numeral";
import { useQuery } from "react-query";
import { Link, useParams } from "react-router-dom";
import { CopyElement } from "../components/CopyElement";
import { PaginatedSymbols, SymbolTypeLabels } from "../models/symbol";
import * as api from "../util/api";

type SymbolListPageParams = {
  type: string;
};

export function SymbolListPage() {
  const { type } = useParams<SymbolListPageParams>() as SymbolListPageParams;
  const {
    data: symbolResp,
    isLoading,
    error,
  } = useQuery<PaginatedSymbols, ky.HTTPError>(
    ["symbols", type],
    () => api.getSymbols(25, type),
    {
      retry: false,
    }
  );

  const handleActivatePageClick = (page: number) => {};

  return (
    <div className="container mb-24 mx-auto pt-8">
      {isLoading && (
        <div className="flex items-center justify-center py-24 flex-col text-gray-500">
          <IconLoader className="animate-spin-slow mr-3" />
          <span className="">Loading Assets...</span>
        </div>
      )}
      {!isLoading && error && (
        <div className="flex items-center justify-center py-20 flex-col text-gray-500">
          <IconMoodSad className="mb-3" />
          <span>There was an error while loading symbols</span>
          <b className="my-1">{error.message}</b>
          <span>Please try again later</span>
        </div>
      )}
      {!isLoading && symbolResp && (
        <>
          <table className="table w-full text-left">
            <thead className="bg-white">
              <tr className="shadow">
                <th className="rounded-l py-3 px-2">Symbol</th>
                <th className="rounded-l py-3 px-2">Name</th>
                <th className="py-3 px-2">Type</th>
              </tr>
            </thead>
            <tbody>
              {symbolResp.symbols.map((symbol) => {
                return (
                  <tr
                    key={symbol.symbolID}
                    className="border-b border-gray-200"
                  >
                    <td className="py-3 px-2">{symbol.symbolID}</td>
                    <td className="py-3 px-2">
                      <div className="flex flex-col">
                        <div className="text-gray-600 text-xs mb-1">
                          {symbol.isin && (
                            <>
                              <CopyElement value={symbol.isin} />
                            </>
                          )}
                          {symbol.wkn && (
                            <>
                              <span className="mx-2">·</span>
                              <CopyElement value={symbol.wkn} />
                            </>
                          )}
                        </div>
                        <Link
                          className="font-bold hover:underline"
                          to={"/symbols/" + symbol.type + "/" + symbol.symbolID}
                        >
                          {symbol.symbolName}
                        </Link>
                      </div>
                    </td>
                    <td className="py-3 px-2">
                      {SymbolTypeLabels[symbol.type]}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
          <div className="flex justify-end mt-4">
            <Pagination
              totalCount={symbolResp.totalCount}
              activePage={symbolResp.activePage}
              pages={symbolResp.pages}
              onPageSelect={handleActivatePageClick}
            />
          </div>
        </>
      )}
    </div>
  );
}

type PaginationProps = {
  totalCount: number;
  pages: number;
  activePage: number;
  onPageSelect: (page: number) => void;
};

function Pagination({
  totalCount,
  pages,
  activePage,
  onPageSelect,
}: PaginationProps) {
  const handlePrevClick = () => onPageSelect(activePage - 1);
  const handleNextClick = () => onPageSelect(activePage + 1);
  const handleSelectPage = (page: number) => () => onPageSelect(page);

  return (
    <div className="bg-white rounded shadow p-2 text-gray-900 flex">
      {numeral(totalCount).format("0a")} Total Items
      <button
        className="px-2 text-blue-700"
        disabled={activePage === 1}
        onClick={handlePrevClick}
      >
        <IconChevronLeft size={20} stroke={3} />
      </button>
      {activePage > 2 && (
        <>
          <button className="px-2" onClick={handleSelectPage(1)}>
            1
          </button>
          <span className="px-2">...</span>
        </>
      )}
      {activePage > 1 && (
        <button className="px-2" onClick={handleSelectPage(activePage - 1)}>
          {activePage - 1}
        </button>
      )}
      <button className="px-2 border border-blue-700 text-blue-700 rounded font-bold">
        {activePage}
      </button>
      {activePage <= pages - 1 && (
        <button className="px-2" onClick={handleSelectPage(activePage + 1)}>
          {activePage + 1}
        </button>
      )}
      {activePage <= pages - 2 && (
        <>
          <span className="px-2">...</span>
          <button className="px-2" onClick={handleSelectPage(pages)}>
            {pages}
          </button>
        </>
      )}
      <button
        className="px-2 text-blue-700"
        disabled={activePage === pages}
        onClick={handleNextClick}
      >
        <IconChevronRight size={20} stroke={3} />
      </button>
    </div>
  );
}
