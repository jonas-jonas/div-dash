import { IconChevronDown, IconChevronUp } from "@tabler/icons";
import classNames from "classnames";
import { useMemo } from "react";
import { useQuery } from "react-query";
import { Cell, Label, Pie, PieChart, ResponsiveContainer } from "recharts";
import { SymbolType, SymbolTypeLabels } from "../models/symbol";
import * as api from "../util/api";
import { formatMoney, formatPercent } from "../util/formatter";

const chartColors = ["ef476f", "ffd166", "06d6a0", "118ab2", "073b4c"];

export function PortfolioComposition() {
  const { data: balance } = useQuery("balance", api.getBalance);

  const chartData = useMemo(() => {
    if (!balance) {
      return [];
    }

    const amountByType = balance.symbols.reduce((acc, balanceItem) => {
      const priceAmount = balanceItem.amount * balanceItem.fiatAssetPrice;
      const costBasis = balanceItem.amount * balanceItem.costBasis;
      if (acc[balanceItem.symbol.type]) {
        acc[balanceItem.symbol.type].total += priceAmount;
        acc[balanceItem.symbol.type].costBasis += costBasis;
      } else {
        acc[balanceItem.symbol.type] = {
          total: priceAmount,
          costBasis: costBasis,
        };
      }
      return acc;
    }, {} as Record<SymbolType, { costBasis: number; total: number }>);

    return Object.entries(amountByType).map(([type, values], index) => {
      const percent = values.total / balance.fiatValue;
      return {
        type: type as SymbolType,
        total: values.total,
        costBasis: values.costBasis,
        color: "#" + chartColors[index % chartColors.length],
        percent,
      };
    });
  }, [balance]);

  return (
    <div className="col-span-1 row-span-2px-6 py-8 flex flex-col">
      <div className="h-96">
        <ResponsiveContainer>
          <PieChart width={400} height={400}>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              label={false}
              outerRadius={130}
              innerRadius={100}
              paddingAngle={1}
              dataKey="total"
            >
              {chartData?.map((entry, i) => (
                <Cell key={entry.type} name={entry.type} fill={entry.color} />
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
          </PieChart>
        </ResponsiveContainer>
      </div>
      <div className="px-8">
        {chartData.map((value) => {
          const borderColor = value.color;
          const pnl = (value.total - value.costBasis) / value.costBasis;
          const isUp = pnl > 0;
          return (
            <div
              className="bg-white shadow rounded mb-4 p-3 flex justify-between border-l-8"
              style={{ borderColor: borderColor }}
              key={value.type}
            >
              <div>
                <h3 className="font-bold">{SymbolTypeLabels[value.type]}</h3>
                <h4>{formatPercent(value.percent)}</h4>
              </div>
              <div className="flex flex-col items-end">
                <h3 className="font-bold">{formatMoney(value.total)}</h3>
                <span
                  className={classNames("text-sm text-white px-2 rounded", {
                    "bg-green-600": isUp,
                    "bg-red-600": isUp,
                  })}
                >
                  {isUp ? (
                    <IconChevronUp className="mr-1" />
                  ) : (
                    <IconChevronDown className="mr-1" />
                  )}
                  {formatPercent(pnl)}
                </span>
              </div>
            </div>
          );
        })}
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
        y={cy - 32}
        className="recharts-text recharts-label"
        textAnchor="middle"
        dominantBaseline="central"
      >
        <tspan fontSize="14px">Total</tspan>
      </text>
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
