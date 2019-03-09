import { getOutcomeActions } from "./fsm";
import data from "./testgraph";

it("getOutcomeActions", () => {
  expect(getOutcomeActions(data, "new", "admin")).toEqual([
    "start_auction",
    "delete"
  ]);
  expect(getOutcomeActions(data, "new", "trading_house")).toEqual([
    "start_auction"
  ]);
});
