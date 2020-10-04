module.exports = {
    someSidebar: [
        {
            type: "doc",
            id: "fisherman"
        },
        {
            type: "category",
            label: "Getting started",
            items: [
                {
                    type: "doc",
                    id: "getting-started/installation"
                }
            ]
        },
        {
            type: "category",
            label: "Configuration",
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
        }
    ],
}
