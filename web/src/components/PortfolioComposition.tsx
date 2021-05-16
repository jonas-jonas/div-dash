import { faChartBar, faChartPie } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useMemo } from "react";
import {
  Cell,
  Label,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
} from "recharts";
import { useRecoilValue } from "recoil";
import { balancesState } from "../state/balanceState";
import { formatMoney } from "../util/formatter";

const COLORS: { [key: string]: string } = {
  BTC: "#f2a900",
  ETH: "#3c3c3d",
  NANO: "#589ae5",
};

export function PortfolioComposition() {
  const balances = useRecoilValue(balancesState);

  const chartData = useMemo(() => {
    return balances.map((balance) => {
      return {
        symbol: balance.symbol,
        total: balance.costBasis * balance.amount,
      };
    });
  }, [balances]);

  const total = useMemo(() => {
    return chartData.reduce((prev, curr) => {
      return prev + curr.total;
    }, 0);
  }, [chartData]);

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
              {chartData.map((entry) => (
                <Cell
                  key={entry.total}
                  name={entry.symbol}
                  fill={COLORS[entry.symbol]}
                />
              ))}
              <Label
                value={formatMoney(total)}
                position="center"
                fontSize={20}
              />
            </Pie>
            <Legend></Legend>
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
