import { faChartBar, faChartPie } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useMemo } from "react";
import { useQuery } from "react-query";
import {
  Cell,
  Label,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
} from "recharts";
import * as api from "../util/api";
import { formatMoney } from "../util/formatter";

const chartColors = [
  "001219",
  "005f73",
  "0a9396",
  "94d2bd",
  "e9d8a6",
  "ee9b00",
  "ca6702",
  "bb3e03",
  "ae2012",
  "9b2226",
];

export function PortfolioComposition() {
  const { data: balance } = useQuery("balance", api.getBalance);

  const chartData = useMemo(() => {
    return balance?.symbols.map((balanceItem) => {
      return {
        symbol: balanceItem.symbol.symbolName,
        total: balanceItem.fiatAssetPrice * balanceItem.amount,
      };
    });
  }, [balance?.symbols]);

  return (
    <div className="col-span-1 row-span-2 bg-white shadow rounded px-6 py-8 flex flex-col">
      <div className="flex justify-between">
        <h2 className="text-2xl">Composition</h2>
        <div>
          <button className="p-1 mr-2 text-blue-700">
            <FontAwesomeIcon icon={faChartPie} />
          </button>
          <button className="p-1 ">
            <FontAwesomeIcon icon={faChartBar} />
          </button>
        </div>
      </div>
      <div className="h-96">
        <ResponsiveContainer>
          <PieChart width={400} height={400}>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              label={false}
              outerRadius={130}
              innerRadius={90}
              paddingAngle={2}
              dataKey="total"
            >
              {chartData?.map((entry, i) => (
                <Cell
                  key={entry.total}
                  name={entry.symbol}
                  fill={"#" + chartColors[i % chartColors.length]}
                />
              ))}

              <Label
                width={30}
                position="center"
                content={
                  <CustomLabel
                    value={balance?.fiatValue!}
                    costBasis={balance?.costBasis!}
                  />
                }
              ></Label>
            </Pie>
            <Legend align="left"></Legend>
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

type CustomLabelProps = {
  viewBox?: { cx: number; cy: number };
  value: number;
  costBasis: number;
};

function CustomLabel({ viewBox, value, costBasis }: CustomLabelProps) {
  const { cx, cy } = viewBox!;
  const pnl = value - costBasis;
  return (
    <>
      <text
        x={cx}
        y={cy - 5}
        fill="rgba(0, 0, 0, 0.87)"
        className="recharts-text recharts-label"
        textAnchor="middle"
        dominantBaseline="central"
      >
        <tspan alignmentBaseline="middle" fontSize="24px" fontWeight="bold">
          {formatMoney(value)}
        </tspan>
      </text>
      <text
        x={cx}
        y={cy + 16}
        className="recharts-text recharts-label"
        textAnchor="middle"
        dominantBaseline="central"
      >
        <tspan fontSize="14px">
          {pnl > 0 ? "+" + formatMoney(pnl) : "-" + formatMoney(pnl)}
        </tspan>
      </text>
    </>
  );
}
