# Mission Control

A code exercise.

## Prerequisites

- Ruby (I use [asdf-vm](https://asdf-vm.com) for managing development dependencies)
- `bundle install`

## Running

```sh
rake game
# or
LOG_LEVEL=trace rake game # <- trace level logging
```

## Project Concerns

- Use a Text User Interface for the UI
- Converting between different units can be difficult if there isn't standardization

## Next Steps

- Add more unit tests!
- Ensure all game flow paths work

## Launch Overview

```mermaid
flowchart LR

  pending <--> afterburner
  pending <--> abort_launch
  afterburner --> abort_launch

  afterburner --> disengage_release_structure
  disengage_release_structure --> cross_checks
  cross_checks --> launch
  launch --> complete

  afterburner --> explode
  disengage_release_structure --> explode
  cross_checks --> explode
  launch --> explode
```
