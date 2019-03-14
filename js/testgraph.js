export default [
  {
    status: "new"
  },
  {
    sources: [
      {
        status: "new",
        action: "start_auction",
        roles: ["trading_house", "admin"]
      }
    ],
    status: "auction_in_progress"
  },
  {
    sources: [
      {
        status: "auction_in_progress",
        action: "inactive",
        roles: ["trading_house", "admin"]
      }
    ],
    status: "winning_offer_for_approval",
    outcomes: [
      {
        action: "close_auction",
        status: "inactive",
        roles: ["trading_house", "admin"]
      }
    ]
  },
  {
    sources: [
      {
        status: "winning_offer_for_approval",
        action: "reject_offer",
        roles: ["trading_house", "admin"]
      }
    ],
    status: "offer_rejected"
  },
  {
    sources: [
      {
        status: "winning_offer_for_approval",
        action: "accept_offer",
        roles: ["trading_house", "admin"]
      }
    ],
    status: "offer_accepted"
  },
  {
    sources: [
      {
        status: "new",
        action: "delete",
        roles: ["admin"]
      },
      {
        status: "auction_in_progress",
        action: "cancel",
        roles: ["trading_house", "admin"]
      },
      {
        status: "winning_offer_for_approval",
        action: "cancel",
        roles: ["trading_house", "admin"]
      }
    ],
    status: "inactive"
  }
];
