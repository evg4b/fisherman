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
    variables:
      from-branch: '^refs/heads/(?P<IssueNumber>PROJ\d+)-.*$'
    commit-prefix: '[{{IssueNumber}}]'
```
