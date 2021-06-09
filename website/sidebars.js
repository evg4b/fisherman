module.exports = {
  someSidebar: [
    {
      type: "doc",
      id: "introduction"
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
          type: "doc",
          id: "configuration/hooks-configuration"
        },
        {
          type: "doc",
          id: "configuration/rules"
        },
        {
          type: "doc",
          id: "configuration/variables"
        },
        {
          type: "doc",
          id: "configuration/expressions"
        },
        {
          type: "doc",
          id: "configuration/output"
        },
      ],
    },
    {
      type: "doc",
      id: "cli"
    },
    {
      type: "doc",
      id: "faq"
    },
  ],
}
