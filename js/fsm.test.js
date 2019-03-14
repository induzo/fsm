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
  expect(
    getOutcomeActions(data, "winning_offer_for_approval", "trading_house")
  ).toEqual(["close_auction", "reject_offer", "accept_offer", "cancel"]);
});
