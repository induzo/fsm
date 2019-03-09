export default [
  {
    status: "new"
  },
  {
    sources: [
      {
        status: "new",
        action: "start_auction",
        roles: ["customer", "admin"]
      }
    ],
    status: "auction_in_progress"
  },
  {
    sources: [
      {
        status: "auction_in_progress",
        action: "close_auction",
        roles: ["customer", "admin"]
      }
    ],
    status: "winning_offer_for_approval"
  },
  {
    sources: [
      {
        status: "winning_offer_for_approval",
        action: "reject_offer",
        roles: ["customer", "admin"]
      }
    ],
    status: "offer_rejected"
  },
  {
    sources: [
      {
        status: "winning_offer_for_approval",
        action: "accept_offer",
        roles: ["customer", "admin"]
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
        roles: ["customer", "admin"]
      },
      {
        status: "winning_offer_for_approval",
        action: "cancel",
        roles: ["customer", "admin"]
      }
    ],
    status: "inactive"
  }
];
