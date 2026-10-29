# A command run caches its pure git queries

Status: accepted

## Decision

A command run caches each pure git answer that it reads.
A pure answer is one of these answers:

- a location answer, which is the common directory, the administration directory, the top level, or a named administration path
- an object read whose revision operands are all full object ids
- an ancestry check whose two operands are both full object ids

A full object id has 40 or 64 hexadecimal digits of either case.
One predicate grades a full object id, and every cache rule and every identity check uses that predicate.

The cache keys each object answer by the exact git arguments, and these arguments include the checkout root and every flag.
The cache gives back a copy of each stored answer.
A failed read stores nothing, so the next read starts git again.
An ancestry check stores a yes and a no, but an ancestry error stores nothing.

The live set stays outside the cache.
The live set holds each ref, HEAD, a short object id, the index, the status, the configuration, the default branch, and the worktree registrations.
A write inside the run can change each of these answers, so the run starts git for each read of them.
Each read of the common directory, the administration directory, or an administration path checks the stored directory again on the filesystem, so a directory that a write removed or replaced still refuses.
A read of the top level does not check the stored answer again.

The cache scope is one command run.
A command creates one run at its entry and gives that run to each reader in the command.
No package value holds a run, and two runs share no answer.
A reader that receives no run starts git for each read and stores nothing.

## Consequence

A command that reads one immutable object many times starts git once for that object.
A reader that names a ref or a short object id reads the current answer, so a ref move inside the command stays visible.
A reader that needs a cached answer must receive the command's run, because a reader with no run gets no cache.

## Considered options

A cache that lasts across commands stays out.
It needs an invalidation rule for each write, and the cost of one command run does not justify that rule.
