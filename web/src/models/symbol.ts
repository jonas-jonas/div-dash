export type Symbol = {
  symbolID: string;
  symbolName: string;
  type: SymbolType;
  source: string;
  precision: number;
};

export type SymbolType =
  | "ad"
  | "crypto"
  | "cs"
  | "et"
  | "ps"
  | "rt"
  | "struct"
  | "ut"
  | "wt"
  | "cef"
  | "oef"
  | "wi"
  | "";

export const SymbolTypeLabels: Record<SymbolType, string> = {
  cs: "Common Stock",
  crypto: "Crypto",
  et: "ETF",
  ad: "ADR",
  cef: "Closed End Fund",
  oef: "Open Ended Fund",
  ps: "Preferred Stock",
  rt: "Right",
  struct: "Structured Product",
  ut: "Unit",
  wi: "When Issued",
  wt: "Warrant",
  "": "Other",
};

type SymbolTagChip = {
  label: string;
  type: "CHIP";
};
type SymbolTagLink = {
  label: string;
  type: "LINK";
  link: string;
};
type SymbolTag = SymbolTagChip | SymbolTagLink;

type SymbolDate = {
  label: string;
  date: string;
};

export type SymbolDetails = {
  type: SymbolType;
  name: string;
  tags: SymbolTag[];
  marketCap: number;
  peRatio: number;
  dividendYield: number;
  eps: number;
  description: string;
  dates: SymbolDate[];
};

export type SymbolChartEntry = {
  date: string;
  price: number;
};
