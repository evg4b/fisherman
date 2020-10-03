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
                    type: "doc",
                    id: "configuration/commit-msg-hook"
                },
                
                {
                    type: "doc",
                    id: "configuration/prepare-commit-msg-hook"
                }
            ],
        }
    ],
}
