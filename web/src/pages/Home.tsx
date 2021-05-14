import { PortfolioBalance } from "../components/PortfolioBalance";
import { PortfolioGraph } from "../components/PortfolioGraph";

export function Home() {
  return (
    <div className="container mb-24 mx-auto">
      <div className="grid grid-cols-3 gap-4 p-8">
        <PortfolioGraph />

        <div className="col-span-1 bg-gray-300">
          <div></div>
        </div>
        <PortfolioBalance />
      </div>
    </div>
  );
}
