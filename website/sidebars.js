module.exports = {
  someSidebar: [
    {
      type: "doc",
      id: "fisherman"
    },
    {
      type: "doc",
      id: "getting-started",
    },
    {
      type: "category",
      label: "Configuration",
      collapsed: false,
      items: [
        {
          type: "doc",
          id: "configuration/configuration-files"
        },
        {
          type: "category",
          label: "Hooks configuration",
          items: [
            {
              type: "doc",
              id: "configuration/hooks/commit-msg-hook"
            },
            {
              type: "doc",
              id: "configuration/hooks/prepare-commit-msg-hook"
            },
          ]
        },
        {
          type: "doc",
          id: "configuration/variables"
        },
        {
          type: "doc",
          id: "configuration/output"
        },
      ],
    },
    {
      type: "doc",
      id: "faq"
    },
  ],
}
