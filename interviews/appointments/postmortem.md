# Postmortem

Feedback: where could I have done better?
Team Standards: how does my code align with team coding standards and practices? (idiomatic go? linting? formatting? test library?)

- a 2 hour project in Golang that requires an API is a bit of a lift. definitely requires efficient use of time.
- TimeZones are difficult
- I would always prefer to use a third-party framework for things like time scheduling. It's a difficult problem.
- Time constraints didn't allow for showing best practices.
  - interfaces
  - testing strategy
  - development environment
  - robust error handling
  - observability from the start
  - separating application layers (rest, service-layer, repo)
  - data transfer types (DTOs) - db fields shouldn't be exposed directly to API
- The "choose time" endpoint doesn't return an ID for the "time slot" which can be problematic. The server side should always be able to understand time slots by a single identifier. This should also work historically.
- Anything that works with time and scheduling needs to be prepared for adjustments:
  - duration not being 30
  - start/end time not on 00/30
  - partial overlaps of fixed (before hours, after hours)

After research it appears the project could have used _merge intervals_ or _interval trees_.

Interval libraries:

- [intersection of two iterable values](https://github.com/juliangruber/go-intersect)
- [interval search tree](https://github.com/rdleal/intervalst)

Interval code challenge solutions:

- [insert interval](https://github.com/halfrost/LeetCode-Go/blob/master/leetcode/0057.Insert-Interval/57.%20Insert%20Interval.go)
- [merge interval](https://github.com/halfrost/LeetCode-Go/tree/master/leetcode/0056.Merge-Intervals)
