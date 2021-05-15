import { PortfolioBalance } from "../components/PortfolioBalance";
import { PortfolioComposition } from "../components/PortfolioComposition";
import { PortfolioGraph } from "../components/PortfolioGraph";

export function Home() {
  return (
    <div className="container mb-24 mx-auto">
      <div className="grid grid-cols-3 gap-4 p-8">
        <PortfolioGraph />

        <PortfolioComposition />
        <PortfolioBalance />
      </div>
    </div>
  );
}
