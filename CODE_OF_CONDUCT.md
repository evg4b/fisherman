# Code of Conduct

Fisherman is a small project maintained by one person in their spare time.
There is no committee here — just a few expectations that keep it pleasant to
work on.

## Be decent

* Be respectful in issues, pull requests, and reviews. Disagree about code, not
  about people.
* Assume good faith. Everyone here is doing this voluntarily.
* No harassment, personal attacks, discriminatory remarks, or posting other
  people's private information.
* Keep it on topic. This is a git hook manager, not a forum.

## No AI slop

AI tools are welcome — this project is developed with them. What is not welcome
is unreviewed machine output dumped on the maintainer.

Before you open an issue or a pull request that an AI helped produce:

* **Read it.** If you have not read every line you are submitting, do not
  submit it.
* **Run it.** Code must compile and pass `make` (lint, unit tests, build)
  locally. "The AI said it works" is not a test result.
* **Understand it.** You should be able to explain why the change is written
  the way it is and answer review questions about it.
* **Keep it small.** Do not send sweeping refactors, mass "improvements", or
  drive-by rewrites nobody asked for.
* **No invented facts.** Bug reports must describe behavior you actually
  observed, with real commands and real output — not a plausible-sounding
  scenario. Do not cite APIs, flags, or files that do not exist.

Reviewing a change costs the maintainer far more time than generating it costs
you. Submissions that ignore the above will be closed without detailed feedback.

## Enforcement

The maintainer may edit, lock, or close any issue, pull request, or comment, and
may block repeat offenders. Report problems by opening an issue or emailing
**evg.abramovitch@gmail.com**. Reports are handled privately.

That is the whole thing. Be decent, do your own thinking, and contributions are
very welcome.
