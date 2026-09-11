# Implementation continuation review pickup

Base: 9f4ee85ceabd1785747b084b583bc696d60d730a
Tip: 107d60eabf2fbf83c9d523be038b74bb0a39f290

## Standards

Finding count: 0
Worst issue: none

## Spec

Finding count: 0
Worst issue: none

## Coverage

Finding count: 1
Worst issue: The unavailable-model fallback does not require an available route.

- P1, auto-fix: Require the authorized fallback diagnostic route to be
  available. The approved decision requires a route that is available and
  authorized. The guidance and its anchor predicate permit an unavailable
  authorized route to satisfy the approved unavailable-model edge.

## Resolution

The retained author required an available authorized route in the guidance and
its independent predicate. Removing `available` now makes
`docs-currency-workflow` fail with the expected diagnostic.
