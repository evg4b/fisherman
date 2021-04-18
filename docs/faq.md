---
id: faq
title: FAQ
slug: /faq
---

## I want check commit prefix based on current brach name

Can be used to check that message starts with issue number.

```yaml
hooks:
  commit-msg:
    extract-variables:
      - 'Extract(BranchName, "^refs/heads/(?P<IssueNumber>PROJ\d+)-.*$")'
    rules:
      - type: commit-message
        condition: Defined("IssueNumber")
        prefix: '[{{IssueNumber}}]'
```

related links:

- [extract-variables](/)
- [commit-message](/)
- ['Defined' function](/)
