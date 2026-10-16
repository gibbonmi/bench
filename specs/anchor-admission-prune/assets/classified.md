# Anchor classification snapshot

This file is the classification of each registry row at the spec base. Each ticket re-derives its slice against the tree at its own base, and the tree wins.

```tsv
registry_file	registry_line	guarded_file	kind	section	doc_line	class	reason
registry_data.go	29	.agents/commands/bench-write-spec.md	Require	Entry orientation	9	DROP	rule-text
registry_data.go	30	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	138	LABEL	override
registry_data.go	31	.agents/commands/bench-write-spec.md	Require	Entry orientation	9	ROUTE	short-route-ref
registry_data.go	32	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	129	DROP	rule-text
registry_data.go	33	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	135	LABEL	override
registry_data.go	34	.agents/commands/bench-write-spec.md	Require	Entry orientation	9	DROP	rule-text
registry_data.go	35	.agents/commands/bench-write-spec.md	Require	Entry orientation	9	DROP	rule-text
registry_data.go	36	.agents/commands/bench-write-spec.md	Require	Entry orientation	9	DROP	rule-text
registry_data.go	37	.agents/commands/bench-write-spec.md	Require	Process	65	DROP	rule-text
registry_data.go	38	.agents/commands/bench-write-spec.md	Require	Process	65	LABEL	override
registry_data.go	39	.agents/commands/bench-write-spec.md	Require	Exit handoff	13	LABEL	override
registry_data.go	40	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	82	DROP	rule-text
registry_data.go	41	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_data.go	42	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	85	DROP	rule-text
registry_data.go	43	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	85	DROP	rule-text
registry_data.go	44	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	85	DROP	rule-text
registry_data.go	45	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	85	DROP	rule-text
registry_data.go	46	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	85	DROP	rule-text
registry_data.go	47	.agents/skills/bench-craft-spec/SKILL.md	Require	The acceptance coverage map	44	LABEL	override
registry_data.go	48	.agents/skills/bench-craft-spec/SKILL.md	Require	Scope cuts	61	DROP	rule-text
registry_data.go	49	.agents/skills/bench-craft-spec/SKILL.md	Require	Slicing a build for delegates	63	DROP	rule-text
registry_data.go	50	.agents/skills/bench-craft-tdd/SKILL.md	Require	Acceptance rows	81	ROUTE	short-route-ref
registry_data.go	51	.agents/skills/bench-craft-tdd/SKILL.md	Require	Acceptance rows	81	DROP	rule-text
registry_data.go	52	.agents/skills/bench-craft-tdd/SKILL.md	Require	Apply it only at pre-agreed seams	37	DROP	rule-text
registry_data.go	53	.agents/skills/bench-craft-tdd/SKILL.md	Require	Acceptance rows	90	DROP	rule-text
registry_data.go	54	.agents/skills/bench-craft-tdd/SKILL.md	Require	Acceptance rows	98	DROP	rule-text
registry_data.go	55	.agents/skills/bench-craft-tdd/SKILL.md	Require	Acceptance rows	81	DROP	sentence-with-route
registry_data.go	56	.agents/skills/bench-craft-tdd/SKILL.md	Require	Apply it only at pre-agreed seams	29	DROP	rule-text
registry_data.go	57	.agents/skills/bench-craft-review/SKILL.md	Require	The axes stay separate	66	ROUTE	short-route-ref
registry_data.go	58	.agents/skills/bench-craft-review/SKILL.md	Require	The axes stay separate	69	DROP	rule-text
registry_data.go	59	.agents/skills/bench-craft-review/SKILL.md	Require	Re-derive, then compare	18	DROP	rule-text
registry_data.go	60	.agents/skills/bench-craft-review/SKILL.md	Require	Re-derive, then compare	23	DROP	rule-text
registry_data.go	61	.agents/skills/bench-craft-review/SKILL.md	Require	The axes stay separate	64	GRAY	sentence-command-with-args
registry_data.go	62	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	57	DROP	rule-text
registry_data.go	63	.agents/skills/bench-craft-delegate/SKILL.md	Require	Verifying the done-claim	114	DROP	rule-text
registry_data.go	64	.agents/skills/bench-craft-delegate/SKILL.md	Require	Isolation	94	GRAY	bare-command-with-args
registry_data.go	65	.agents/skills/bench-craft-delegate/SKILL.md	Require	Isolation	103	DROP	rule-text
registry_data.go	66	.agents/skills/bench-craft-delegate/SKILL.md	Require	Verifying the done-claim	125	DROP	sentence-command-noargs
registry_data.go	67	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	54	DROP	rule-text
registry_data.go	68	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	55	DROP	rule-text
registry_data.go	69	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	67	GRAY	sentence-command-with-args
registry_data.go	70	.agents/skills/bench-craft-seams/SKILL.md	Require	Picking the seam for a feature	61	DROP	rule-text
registry_data.go	71	.agents/skills/bench-craft-seams/SKILL.md	Require	Splitting when the structure gate fires	100	DROP	rule-text
registry_data.go	72	.agents/commands/bench-implement-spec.md	Require	When the build stops short	76	DROP	rule-text
registry_data.go	73	.agents/commands/bench-implement-spec.md	Require	Build	42	GRAY	bare-command-with-args
registry_data.go	74	.agents/commands/bench-implement-spec.md	Require	Build	41	GRAY	sentence-command-with-args
registry_data.go	75	.agents/commands/bench-implement-spec.md	Require	Exit handoff	15	DROP	rule-text
registry_data.go	76	.agents/commands/bench-implement-spec.md	Require	Land	68	LABEL	override
registry_data.go	77	.agents/commands/bench-implement-spec.md	Require	Land	68	DROP	sentence-command-noargs
registry_data.go	78	.agents/commands/bench-review-implementation.md	Require	Process	154	DROP	rule-text
registry_data.go	79	.agents/commands/bench-review-implementation.md	Require	Process	156	DROP	rule-text
registry_data.go	80	.agents/commands/bench-review-implementation.md	Require	Process	89	GRAY	bare-command-with-args
registry_data.go	81	.agents/commands/bench-review-implementation.md	Require	Process	94	GRAY	bare-command-with-args
registry_data.go	82	.agents/commands/bench-review-implementation.md	Require	Process	182	LABEL	override
registry_data.go	83	.agents/commands/bench-review-implementation.md	Require	Process	126	DROP	rule-text
registry_data.go	84	.agents/commands/bench-review-implementation.md	Require	Process	129	ROUTE	short-route-ref
registry_data.go	85	.agents/commands/bench-review-implementation.md	Require	Process	127	ROUTE	short-route-ref
registry_data.go	86	.agents/commands/bench-review-implementation.md	Require	Process	205	DROP	rule-text
registry_data.go	87	.agents/commands/bench-review-implementation.md	Require	Process	237	DROP	rule-text
registry_data.go	88	.agents/commands/bench-review-implementation.md	Require	Process	204	DROP	rule-text
registry_data.go	89	.agents/commands/bench-review-implementation.md	Require	Process	235	DROP	rule-text
registry_data.go	90	.agents/commands/bench-review-implementation.md	Require	Process	205	DROP	rule-text
registry_data.go	91	.agents/commands/bench-review-implementation.md	Require	Process	207	DROP	rule-text
registry_data.go	92	.agents/commands/bench-review-implementation.md	Require	Process	208	DROP	rule-text
registry_data.go	93	.agents/commands/bench-review-implementation.md	Require	Process	237	DROP	rule-text
registry_data.go	94	.agents/commands/bench-review-implementation.md	Require	Exit handoff	30	DROP	rule-text
registry_data.go	95	.agents/commands/bench-review-implementation.md	Require	Process	191	DROP	rule-text
registry_data.go	96	.agents/commands/bench-final-check.md	Require	Run it	179	DROP	rule-text
registry_data.go	97	.agents/commands/bench-final-check.md	Require	Run it	180	DROP	rule-text
registry_data.go	98	.agents/commands/bench-final-check.md	Require	Report	191	DROP	rule-text
registry_data.go	99	.agents/commands/bench-final-check.md	Require	Report	203	ROUTE	short-route-ref
registry_data.go	100	.agents/commands/bench-final-check.md	Require	Exit handoff	52	GRAY	bare-command-with-args
registry_data.go	101	.agents/commands/bench-final-check.md	Require	Exit handoff	48	DROP	rule-text
registry_data.go	102	.agents/commands/bench-setup-repo.md	Require	2. Present findings, then ask one decision at a time	121	DROP	rule-text
registry_data.go	103	.agents/commands/bench-setup-repo.md	Require	2. Present findings, then ask one decision at a time	106	ROUTE	short-route-ref
registry_data.go	104	.agents/commands/bench-debug.md	Require	Finding a retired spec	129	DROP	rule-text
registry_data.go	105	.agents/commands/bench-debug.md	Require	Phase 1 — build the loop (this is the whole skill)	49	DROP	rule-text
registry_data.go	106	.agents/commands/bench-debug.md	Require	Phase 1 — build the loop (this is the whole skill)	37	ROUTE	override
registry_data.go	107	.agents/commands/bench-debug.md	Require	Phase 1 — build the loop (this is the whole skill)	63	DROP	rule-text
registry_data.go	108	.agents/commands/bench-debug.md	Require	Phase 2 — reproduce and minimise	83	DROP	rule-text
registry_data.go	109	.agents/commands/bench-debug.md	Require	Phase 1 — build the loop (this is the whole skill)	49	DROP	override
registry_data.go	110	.agents/commands/bench-debug.md	Require	Phase 2 — reproduce and minimise	69	DROP	rule-text
registry_data.go	111	.agents/commands/bench-debug.md	Require	Phase 1 — build the loop (this is the whole skill)	41	DROP	rule-text
registry_data.go	112	.agents/commands/bench-debug.md	Require	Phase 2 — reproduce and minimise	76	DROP	rule-text
registry_data.go	113	.agents/skills/bench-debug/references/loop-constructions.md	Require		20	DROP	rule-text
registry_data.go	114	.agents/skills/bench-debug/references/loop-constructions.md	Require		18	DROP	rule-text
registry_data.go	115	.agents/skills/bench-debug/SKILL.md	Require		3	DROP	rule-text
registry_data.go	116	.agents/commands/bench-drain.md	Require	Entry orientation	41	DROP	rule-text
registry_data.go	117	.agents/commands/bench-drain.md	Require	5. Verdict the journal	208	DROP	rule-text
registry_data.go	118	.agents/commands/bench-drain.md	Require	3. Drain the inbox	151	DROP	rule-text
registry_data.go	119	.agents/commands/bench-drain.md	Require	5. Verdict the journal	198	DROP	rule-text
registry_data.go	120	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	268	DROP	rule-text
registry_data.go	121	.agents/commands/bench-drain.md	Require	Exit handoff	83	DROP	rule-text
registry_data.go	122	.agents/commands/bench-drain.md	Require	7. Leave the sequence to the commitment	259	LABEL	override
registry_data.go	123	.agents/commands/bench-assess.md	Require	1. Verify the previous assessment landed	40	DROP	rule-text
registry_data.go	124	.agents/commands/bench-assess.md	Require	2. Fan out the area sweeps (mid tier)	50	DROP	rule-text
registry_data.go	125	.agents/commands/bench-assess.md	Require		2	DROP	rule-text
registry_data.go	126	.agents/commands/bench-assess.md	Require	5. Write the assessment (the output contract)	122	DROP	rule-text
registry_data.go	127	.agents/commands/bench-assess.md	Require	Exit handoff	28	ROUTE	short-route-ref
registry_data.go	128	projects/benchkit.md	Require	Hostile-input checklist (shell CLI)	176	DROP	rule-text
registry_data.go	129	projects/benchkit.md	Require	Lines (model + effort routing)	603	DROP	sentence-with-route
registry_data.go	130	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	81	DROP	rule-text
registry_data.go	131	.agents/commands/bench-write-spec.md	Require	Process	68	DROP	rule-text
registry_data.go	132	.agents/skills/bench-craft-grill/SKILL.md	Require	Discipline	57	DROP	rule-text
registry_data.go	133	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	19	GRAY	sentence-command-with-args
registry_data.go	134	.agents/commands/bench-review-implementation.md	Require	Entry orientation	12	DROP	sentence-command-noargs
registry_data.go	135	CHANGELOG.md	Require	[Unreleased]	910	GRAY	sentence-command-with-args
registry_data.go	136	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	33	DROP	rule-text
registry_data.go	137	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	33	DROP	rule-text
registry_data.go	138	README.md	Forbid		0	FORBID	kind
registry_data.go	139	README.md	Require	Reviewer quick start	167	DROP	rule-text
registry_data.go	140	.agents/skills/bench-craft-delegate/SKILL.md	Require	Delegate or retain	13	DROP	rule-text
registry_data.go	141	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_data.go	142	.agents/skills/bench-craft-delegate/SKILL.md	Require	Delegate or retain	17	DROP	rule-text
registry_data.go	143	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_data.go	144	.agents/skills/bench-craft-delegate/SKILL.md	Require	Delegate or retain	17	DROP	rule-text
registry_data.go	145	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_data.go	146	.agents/commands/bench-debug.md	Forbid		0	FORBID	kind
registry_data.go	147	.agents/commands/bench-debug.md	Require	How it meets the rest of Bench	160	DROP	rule-text
registry_data.go	148	.agents/commands/bench-debug.md	Require	How it meets the rest of Bench	161	DROP	rule-text
registry_data.go	149	.agents/commands/bench-debug.md	Require	How it meets the rest of Bench	164	DROP	rule-text
registry_data.go	150	.agents/skills/bench-craft-seams/SKILL.md	Require	Splitting when the structure gate fires	97	DROP	rule-text
registry_data.go	151	.agents/skills/bench-craft-delegate/SKILL.md	Require	Isolation	100	DROP	rule-text
registry_data.go	152	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	69	DROP	rule-text
registry_data.go	153	.agents/skills/bench-craft-delegate/SKILL.md	Require	Verifying the done-claim	122	DROP	rule-text
registry_data.go	154	.agents/commands/bench-write-spec.md	Require	Process	60	DROP	rule-text
registry_data.go	155	.agents/commands/bench-drain.md	Require	1. Reconcile first	96	GRAY	sentence-command-with-args
registry_data.go	156	.bench/BENCH-reference.md	Require	How the pieces fit	25	DROP	rule-text
registry_data.go	157	.bench/BENCH-reference.md	Require	How the pieces fit	25	DROP	sentence-command-noargs
registry_data.go	158	.bench/BENCH-reference.md	Require	Skills index	98	ROUTE	short-route-ref
registry_data.go	159	.bench/BENCH.md	Require	Capture	179	DROP	rule-text
registry_data.go	160	.bench/BENCH-reference.md	Require	Files	43	DROP	rule-text
registry_data.go	161	.bench/BENCH.md	RequireInSection	Workflow	154	DROP	bench-workflow
registry_data.go	162	.bench/BENCH.md	RequireInSection	How to talk to me	86	DROP	rule-text
registry_data.go	163	.bench/BENCH.md	RequireInSection	Roles	11	INVARIANT	bench-invariant-section
registry_data.go	164	.bench/BENCH.md	RequireInSection	How to talk to me	102	DROP	rule-text
registry_data.go	165	.bench/BENCH.md	RequireInSection	Workflow	132	DROP	bench-workflow
registry_data.go	166	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	167	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	168	.agents/commands/bench-write-spec.md	Require	Process	65	DROP	rule-text
registry_data.go	169	.agents/commands/bench-drain.md	Require	1. Reconcile first	91	DROP	rule-text
registry_data.go	170	.agents/commands/bench-drain.md	Require	1. Reconcile first	92	DROP	sentence-with-label
registry_data.go	171	.agents/commands/bench-drain.md	Require	1. Reconcile first	93	DROP	rule-text
registry_data.go	172	.agents/commands/bench-shape-idea.md	Require	Two modes	148	DROP	rule-text
registry_data.go	173	.agents/commands/bench-write-spec.md	Require	Process	73	DROP	sentence-with-route
registry_data.go	174	.agents/commands/bench-write-spec.md	Require	Process	79	DROP	sentence-with-label
registry_data.go	175	.agents/commands/bench-write-spec.md	Require	Process	80	DROP	rule-text
registry_data.go	176	.agents/commands/bench-write-spec.md	Require	Process	73	DROP	rule-text
registry_data.go	177	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	102	LABEL	override
registry_data.go	178	.agents/commands/bench-write-spec.md	RequireInSection	Process	85	DROP	rule-text
registry_data.go	179	.agents/commands/bench-write-spec.md	RequireInSection	Process	86	DROP	rule-text
registry_data.go	180	.agents/commands/bench-write-spec.md	RequireInSection	Process	86	DROP	rule-text
registry_data.go	181	.agents/commands/bench-write-spec.md	RequireInSection	Who runs this phase	51	DROP	rule-text
registry_data.go	182	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_data.go	183	.agents/commands/bench-implement-spec.md	Require	Entry orientation	9	DROP	sentence-with-route
registry_data.go	184	.bench/BENCH.md	Require	Workflow	137	DROP	bench-workflow
registry_data.go	185	.bench/BENCH.md	Require	Workflow	137	DROP	bench-workflow
registry_data.go	186	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	187	.bench/BENCH.md	Require	Workflow	137	DROP	bench-workflow
registry_data.go	188	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	51	DROP	rule-text
registry_data.go	189	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	52	DROP	rule-text
registry_data.go	190	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	53	DROP	rule-text
registry_data.go	191	.agents/skills/bench-craft-line/SKILL.md	Require	Classify reds before the ladder moves	80	DROP	rule-text
registry_data.go	192	.agents/skills/bench-craft-tickets/SKILL.md	Require		10	DROP	rule-text
registry_data.go	193	.agents/skills/bench-craft-tickets/SKILL.md	Require		12	DROP	rule-text
registry_data.go	194	.agents/skills/bench-craft-tickets/SKILL.md	Require	Land the frontier	93	DROP	rule-text
registry_data.go	195	.agents/skills/bench-craft-tickets/SKILL.md	Require	Draft the breakdown	26	LABEL	override
registry_data.go	196	.agents/skills/bench-craft-tickets/SKILL.md	Require	What to build	51	LABEL	override
registry_data.go	197	.agents/skills/bench-craft-tickets/SKILL.md	Require	Acceptance	55	LABEL	override
registry_data.go	198	.agents/skills/bench-craft-spec/SKILL.md	Require	Slicing a build for delegates	71	DROP	override
registry_data.go	199	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Write one file per ticket	47	LABEL	override
registry_data.go	200	.agents/skills/bench-craft-tickets/SKILL.md	Require	Acceptance	68	DROP	rule-text
registry_data.go	201	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	25	DROP	rule-text
registry_data.go	202	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	24	DROP	rule-text
registry_data.go	203	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	22	DROP	rule-text
registry_data.go	204	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	26	DROP	sentence-with-label
registry_data.go	205	.agents/skills/bench-craft-tickets/SKILL.md	ForbidInSection	Draft the breakdown	0	FORBID	kind
registry_data.go	206	.agents/skills/bench-craft-tickets/SKILL.md	Require	Draft the breakdown	35	DROP	sentence-with-label
registry_data.go	207	.agents/skills/bench-craft-tickets/SKILL.md	Require	Draft the breakdown	37	DROP	rule-text
registry_data.go	208	.agents/skills/bench-craft-tickets/SKILL.md	Forbid		0	FORBID	kind
registry_data.go	209	.agents/commands/bench-final-check.md	Require	Exit handoff	81	DROP	rule-text
registry_data.go	210	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	113	DROP	sentence-with-label
registry_data.go	211	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	59	DROP	rule-text
registry_data.go	212	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	60	DROP	rule-text
registry_data.go	213	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	74	DROP	rule-text
registry_data.go	214	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	81	DROP	rule-text
registry_data.go	215	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	50	DROP	rule-text
registry_data.go	216	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Isolation	106	DROP	rule-text
registry_data.go	217	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Isolation	109	DROP	rule-text
registry_data.go	218	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Isolation	109	DROP	rule-text
registry_data.go	219	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	171	DROP	rule-text
registry_data.go	220	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	162	DROP	rule-text
registry_data.go	221	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_data.go	222	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	61	DROP	rule-text
registry_data.go	223	.agents/skills/bench-craft-delegate/SKILL.md	Require		8	DROP	rule-text
registry_data.go	224	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_data.go	225	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Verifying the done-claim	118	DROP	rule-text
registry_data.go	226	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Verifying the done-claim	120	DROP	rule-text
registry_data.go	227	.agents/commands/bench-write-spec.md	Require	Process	83	DROP	rule-text
registry_data.go	228	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Require		3	DROP	sentence-with-route
registry_data.go	229	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Before the landing	198	DROP	rule-text
registry_data.go	230	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Before the landing	200	DROP	rule-text
registry_data.go	231	.agents/skills/bench-craft-spec/references/map-discipline.md	Require		3	DROP	sentence-with-route
registry_data.go	232	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	In the edge inventory	155	DROP	rule-text
registry_data.go	233	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Per row	142	DROP	rule-text
registry_data.go	234	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Per row	143	DROP	rule-text
registry_data.go	235	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	In the edge inventory	156	DROP	rule-text
registry_data.go	236	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	27	DROP	rule-text
registry_data.go	237	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	28	DROP	rule-text
registry_data.go	238	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	29	DROP	rule-text
registry_data.go	239	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	30	DROP	rule-text
registry_data.go	240	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	31	DROP	rule-text
registry_data.go	241	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	32	DROP	rule-text
registry_data.go	242	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	33	DROP	rule-text
registry_data.go	243	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	34	DROP	rule-text
registry_data.go	244	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	35	DROP	rule-text
registry_data.go	245	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	At review	190	DROP	rule-text
registry_data.go	246	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Per row	144	DROP	rule-text
registry_data.go	247	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	44	DROP	rule-text
registry_data.go	248	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Per row	145	DROP	rule-text
registry_data.go	249	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	At ticket slicing	163	DROP	rule-text
registry_data.go	250	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	45	DROP	rule-text
registry_data.go	251	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	46	DROP	rule-text
registry_data.go	252	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	47	DROP	rule-text
registry_data.go	253	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	48	DROP	rule-text
registry_data.go	254	.agents/skills/bench-craft-spec/references/map-discipline.md	Forbid		0	FORBID	kind
registry_data.go	255	.agents/skills/bench-craft-spec/SKILL.md	Require	Authoring process	16	DROP	rule-text
registry_data.go	256	.agents/commands/bench-write-spec.md	Require	Process	58	DROP	sentence-with-route
registry_data.go	257	.agents/commands/bench-write-spec.md	Require	Process	69	DROP	sentence-with-route
registry_data.go	258	.agents/skills/bench-craft-spec/SKILL.md	ForbidInSection	Template	0	FORBID	kind
registry_data.go	259	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	36	DROP	rule-text
registry_data.go	260	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	37	DROP	rule-text
registry_data.go	261	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	39	DROP	rule-text
registry_data.go	262	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	39	DROP	rule-text
registry_data.go	263	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	40	DROP	rule-text
registry_data.go	264	.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md	Require		18	DROP	rule-text
registry_data.go	265	.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md	Require		22	DROP	rule-text
registry_data.go	266	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	22	LABEL	override
registry_data.go	267	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	23	LABEL	override
registry_data.go	268	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	24	LABEL	override
registry_data.go	269	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	25	LABEL	override
registry_data.go	270	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	26	LABEL	override
registry_data.go	271	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	65	DROP	rule-text
registry_data.go	272	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	65	DROP	rule-text
registry_data.go	273	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	71	DROP	rule-text
registry_data.go	274	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	71	DROP	rule-text
registry_data.go	275	.agents/skills/bench-craft-tdd/SKILL.md	RequireInSection	Apply it only at pre-agreed seams	37	DROP	rule-text
registry_data.go	276	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	The acceptance coverage map	44	DROP	rule-text
registry_data.go	277	.agents/skills/bench-craft-spec/SKILL.md	ForbidInSection	The acceptance coverage map	0	FORBID	kind
registry_data.go	278	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	138	LABEL	override
registry_data.go	279	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	140	LABEL	override
registry_data.go	280	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	145	DROP	sentence-with-route
registry_data.go	281	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	155	DROP	rule-text
registry_data.go	282	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	155	DROP	rule-text
registry_data.go	283	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Review rubric	83	DROP	rule-text
registry_data.go	284	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	The acceptance coverage map	48	DROP	rule-text
registry_data.go	285	projects/benchkit.md	RequireInSection	Hostile-input checklist (shell CLI)	288	DROP	rule-text
registry_data.go	286	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	72	DROP	rule-text
registry_data.go	287	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Classify reds before the ladder moves	98	DROP	sentence-with-route
registry_data.go	288	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Verifying the done-claim	120	DROP	rule-text
registry_data.go	289	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Retry stops and aggregate readiness	176	DROP	rule-text
registry_data.go	290	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Retry stops and aggregate readiness	178	DROP	rule-text
registry_data.go	291	.agents/commands/bench-write-spec.md	Require	Entry contract	27	DROP	rule-text
registry_data.go	292	.agents/commands/bench-shape-idea.md	Require	Map content	102	DROP	sentence-command-noargs
registry_data.go	293	.agents/commands/bench-write-spec.md	Require	Entry contract	31	DROP	rule-text
registry_data.go	294	.agents/commands/bench-write-spec.md	Require	Entry contract	31	DROP	rule-text
registry_data.go	295	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_data.go	296	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	133	DROP	rule-text
registry_data.go	297	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	119	DROP	rule-text
registry_data.go	298	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_data.go	299	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_data.go	300	.agents/commands/bench-final-check.md	Require	Exit handoff	79	GRAY	sentence-command-with-args
registry_data.go	301	.agents/commands/bench-final-check.md	Require	Exit handoff	84	DROP	rule-text
registry_data.go	302	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	63	DROP	rule-text
registry_data.go	303	.bench/BENCH-reference.md	Require	Files	39	DROP	rule-text
registry_data.go	304	projects/benchkit.md	Forbid		0	FORBID	kind
registry_data.go	305	.bench/BENCH-reference.md	Require	Files	45	ROUTE	short-route-ref
registry_data.go	306	.agents/commands/bench-drain.md	Require	5. Verdict the journal	219	DROP	rule-text
registry_data.go	307	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	294	DROP	rule-text
registry_data.go	308	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	295	DROP	rule-text
registry_data.go	309	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	296	DROP	rule-text
registry_data.go	310	.agents/commands/bench-drain.md	Require	4. Drain implementation retros	169	DROP	rule-text
registry_data.go	311	.agents/commands/bench-drain.md	Require	4. Drain implementation retros	173	DROP	rule-text
registry_data.go	312	.agents/commands/bench-drain.md	Require	4. Drain implementation retros	190	DROP	rule-text
registry_data.go	313	.agents/commands/bench-drain.md	Require	4. Drain implementation retros	182	DROP	rule-text
registry_data.go	314	.bench/BENCH-reference.md	Require	Files	47	ROUTE	short-route-ref
registry_data.go	317	.agents/commands/bench-drain.md	Require	Exit handoff	82	GRAY	sentence-command-with-args
registry_data.go	318	.agents/commands/bench-drain.md	Require	2. Drain occurrence evidence	123	DROP	sentence-with-label
registry_data.go	319	.agents/commands/bench-drain.md	Require	2. Drain occurrence evidence	124	DROP	rule-text
registry_data.go	320	.agents/commands/bench-drain.md	Require	3. Drain the inbox	154	DROP	sentence-with-label
registry_data.go	321	.agents/commands/bench-drain.md	Require	6. Classify every run; restructure on request	253	DROP	rule-text
registry_data.go	322	.agents/commands/bench-shape-idea.md	ForbidInSection	The exit	0	FORBID	kind
registry_data.go	323	.agents/commands/bench-shape-idea.md	RequireInSection	The exit	162	DROP	sentence-with-route
registry_data.go	324	.agents/commands/bench-shape-idea.md	Require	Decision tickets	77	DROP	rule-text
registry_data.go	325	.agents/commands/bench-review-implementation.md	Require	Process	150	DROP	rule-text
registry_data.go	326	.agents/commands/bench-review-implementation.md	Require	Process	151	DROP	rule-text
registry_data.go	327	.agents/commands/bench-implement-spec.md	RequireInSection	`--full <spec>`	85	DROP	rule-text
registry_data.go	328	.agents/commands/bench-implement-spec.md	RequireInSection	`--full <spec>`	85	DROP	sentence-command-noargs
registry_data.go	329	.agents/commands/bench-implement-spec.md	RequireInSection	`--full <spec>`	87	DROP	rule-text
registry_data.go	330	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Verifying the done-claim	124	GRAY	sentence-command-with-args
registry_data.go	331	.agents/commands/bench-final-check.md	RequireInSection	Exit handoff	75	GRAY	sentence-command-with-args
registry_data.go	332	.agents/commands/bench-final-check.md	RequireInSection	Exit handoff	74	DROP	sentence-command-noargs
registry_data.go	333	.agents/commands/bench-final-check.md	ForbidInSection	Exit handoff	0	FORBID	kind
registry_data.go	334	CHANGELOG.md	Require	[Unreleased]	857	DROP	sentence-command-noargs
registry_data.go	335	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_data.go	336	.agents/commands/bench-write-spec.md	Require	Exit handoff	14	DROP	rule-text
registry_data.go	337	.agents/commands/bench-write-spec.md	Require	Exit handoff	14	DROP	sentence-command-noargs
registry_data.go	338	.agents/commands/bench-implement-spec.md	Require	Land	60	DROP	sentence-with-label
registry_data.go	339	.agents/commands/bench-review-implementation.md	Require	Entry orientation	12	DROP	rule-text
registry_data.go	340	.agents/commands/bench-final-check.md	Require	Entry orientation	10	DROP	sentence-command-noargs
registry_data.go	341	.agents/skills/bench-craft-delegate/SKILL.md	Require	Isolation	100	DROP	rule-text
registry_data.go	342	.agents/skills/bench-craft-tickets/SKILL.md	Require	Land the frontier	97	DROP	sentence-with-label
registry_data.go	343	projects/benchkit.md	Require	Seams (test here; everything else is free to change)	38	DROP	rule-text
registry_data.go	344	.bench/BENCH-reference.md	Require	Command Notes	339	DROP	rule-text
registry_data.go	345	docs/field-guide.html	Require		1147	DROP	rule-text
registry_data.go	346	docs/greenfield-build-sequence.md	Require	Sprint plan	24	DROP	rule-text
registry_data.go	347	internal/usage/worktree.go	Require		23	GRAY	override
registry_data.go	348	.agents/commands/bench-write-spec.md	Require	Entry contract	35	DROP	sentence-with-label
registry_data.go	349	.agents/commands/bench-write-spec.md	Require	Entry contract	35	DROP	rule-text
registry_data.go	350	.agents/commands/bench-write-spec.md	Require	Entry contract	36	DROP	rule-text
registry_data.go	351	.agents/commands/bench-write-spec.md	Require	Entry contract	36	DROP	sentence-with-route
registry_data.go	352	.agents/commands/bench-shape-idea.md	Require	Entry orientation	9	DROP	rule-text
registry_data.go	353	.agents/commands/bench-shape-idea.md	Require	Entry orientation	22	GRAY	sentence-command-with-args
registry_data.go	354	.agents/commands/bench-shape-idea.md	Require	Decision tickets	60	DROP	rule-text
registry_data.go	355	.agents/commands/bench-shape-idea.md	Require	Ownership	37	DROP	rule-text
registry_data.go	356	.agents/commands/bench-write-spec.md	Require	Who runs this phase	53	DROP	rule-text
registry_data.go	357	.agents/commands/bench-shape-idea.md	Forbid		0	FORBID	kind
registry_data.go	358	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_data.go	359	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_data.go	360	.agents/skills/bench-craft-grill/SKILL.md	Require	What it feeds	76	DROP	rule-text
registry_data.go	361	.agents/skills/bench-craft-tdd/SKILL.md	Forbid		0	FORBID	kind
registry_data.go	362	.agents/skills/bench-craft-tdd/SKILL.md	Require	Apply it only at pre-agreed seams	23	DROP	rule-text
registry_data.go	363	.agents/skills/bench-craft-spec/SKILL.md	Require		9	DROP	rule-text
registry_data.go	364	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	43	DROP	rule-text
registry_data.go	365	CONTEXT.md	Require	Core terms	269	DROP	override
registry_data.go	366	CONTEXT.md	Require	Core terms	289	DROP	rule-text
registry_data.go	367	CONTEXT.md	Require	Core terms	291	DROP	rule-text
registry_data.go	368	CONTEXT.md	Require	Core terms	294	DROP	rule-text
registry_data.go	369	README.md	Require	Reviewer quick start	175	DROP	rule-text
registry_data.go	370	projects/benchkit.md	Require	Lines (model + effort routing)	582	DROP	rule-text
registry_data.go	371	CHANGELOG.md	Require	[Unreleased]	1016	DROP	rule-text
registry_data.go	372	CHANGELOG.md	Require	[Unreleased]	808	DROP	rule-text
registry_data.go	373	docs/field-guide.html	Require		709	DROP	rule-text
registry_data.go	374	docs/reporesident-distillation.md	Require	2. Refactor lane with a mechanical exit test	93	DROP	rule-text
registry_data.go	375	.bench/BENCH.md	Require	The four invariants (these override convenience, always)	55	INVARIANT	bench-invariant-section
registry_data.go	376	.bench/BENCH.md	Require	The four invariants (these override convenience, always)	61	INVARIANT	bench-invariant-section
registry_data.go	377	.bench/BENCH.md	Require	The four invariants (these override convenience, always)	65	INVARIANT	bench-invariant-section
registry_data.go	378	.bench/BENCH.md	Require	The four invariants (these override convenience, always)	66	INVARIANT	bench-invariant-section
registry_data.go	379	AGENTS.md	Require	This repo	72	DROP	rule-text
registry_data.go	380	AGENTS.md	Require	This repo	78	DROP	rule-text
registry_data.go	381	AGENTS.md	Require	This repo	80	DROP	rule-text
registry_data.go	382	AGENTS.md	Require	This repo	82	DROP	rule-text
registry_data.go	383	AGENTS.md	Require	This repo	86	GRAY	sentence-command-with-args
registry_data.go	384	AGENTS.md	Forbid		0	FORBID	kind
registry_data.go	385	projects/benchkit.md	Require	Notes for cold sessions	685	GRAY	sentence-command-with-args
registry_data.go	386	.agents/skills/bench-craft-spec/SKILL.md	Require	Slicing a build for delegates	73	DROP	rule-text
registry_data.go	387	.agents/skills/bench-craft-seams/SKILL.md	Require	Picking the seam for a feature	72	DROP	rule-text
registry_data.go	388	.agents/commands/bench-debug.md	Require	Entry orientation	14	DROP	rule-text
registry_data.go	389	.agents/commands/bench-review-implementation.md	Require	Exit handoff	36	DROP	rule-text
registry_data.go	390	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_data.go	391	.bench/BENCH-reference.md	Require	Phase manifest	496	DROP	rule-text
registry_data.go	392	.bench/BENCH.md	RequireInSection	Workflow	150	DROP	bench-workflow
registry_data.go	393	.bench/BENCH.md	RequireInSection	Workflow	151	DROP	bench-workflow
registry_data.go	394	.bench/BENCH-reference.md	Require	Command Notes	351	DROP	rule-text
registry_data.go	395	.bench/BENCH-reference.md	Require	Command Notes	358	DROP	rule-text
registry_data.go	396	.bench/BENCH.md	RequireInSection	The four invariants (these override convenience, always)	52	INVARIANT	bench-invariant-section
registry_data.go	397	.bench/BENCH-reference.md	RequireInSection	Command Notes	342	DROP	rule-text
registry_data.go	398	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_data.go	399	.bench/BENCH-reference.md	Require	Hook Layers	634	GRAY	override
registry_data.go	400	.bench/BENCH-reference.md	Require	Hook Layers	624	DROP	rule-text
registry_data.go	401	CHANGELOG.md	Require	[Unreleased]	559	ROUTE	bare-verb
registry_data.go	402	.bench/hooks/block-dangerous-git.sh	Require		10	DROP	rule-text
registry_data.go	403	.bench/hooks/block-dangerous-git.sh	Require		14	DROP	rule-text
registry_data.go	404	.bench/hooks/block-dangerous-git.sh	Require		16	DROP	rule-text
registry_data.go	405	.agents/skills/bench-craft-spec/references/ste-prose.md	Require		18	DROP	rule-text
registry_data.go	406	.agents/skills/bench-craft-spec/references/ste-prose.md	Require		16	DROP	rule-text
registry_data.go	407	projects/benchkit.md	Forbid		0	FORBID	kind
registry_data.go	408	.bench/BENCH-reference.md	Require	Phase manifest	565	DROP	rule-text
registry_data.go	409	.bench/BENCH.md	RequireInSection	The Bench CLI contract	19	INVARIANT	bench-invariant-section
registry_data.go	410	.bench/BENCH.md	Require	The Bench CLI contract	21	INVARIANT	bench-invariant-section
registry_data.go	411	.bench/BENCH-reference.md	Require	How the pieces fit	19	DROP	rule-text
registry_data.go	412	.bench/BENCH-reference.md	Require	How the pieces fit	27	DROP	rule-text
registry_data.go	413	.bench/BENCH-reference.md	Require		13	DROP	rule-text
registry_data.go	414	.bench/BENCH-reference.md	Require	Command Notes	218	DROP	rule-text
registry_data.go	415	.bench/BENCH-reference.md	Require	Command Notes	219	DROP	rule-text
registry_data.go	416	.bench/BENCH-reference.md	Require	Command Notes	220	DROP	rule-text
registry_data.go	417	.bench/BENCH-reference.md	Require	Command Notes	221	DROP	rule-text
registry_data.go	418	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	419	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	420	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	421	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	422	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	423	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	424	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	425	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	426	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	427	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	428	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	429	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	430	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	431	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_data.go	432	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_data.go	433	.agents/skills/bench-craft-comments/SKILL.md	RequireInSection	The register	44	DROP	rule-text
registry_data.go	434	.agents/skills/bench-craft-comments/SKILL.md	RequireInSection	The register	45	DROP	rule-text
registry_data.go	435	.agents/skills/bench-craft-comments/SKILL.md	RequireInSection	The register	48	DROP	rule-text
registry_data.go	436	.agents/skills/bench-craft-comments/SKILL.md	RequireInSection	The register	42	DROP	rule-text
registry_data.go	437	.agents/skills/bench-craft-comments/SKILL.md	RequireInSection	The register	49	DROP	rule-text
registry_data.go	438	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	The axes stay separate	47	DROP	rule-text
registry_data.go	439	.agents/skills/bench-craft-gate/SKILL.md	RequireInSection	Run the real path	51	DROP	rule-text
registry_data.go	440	.agents/skills/bench-craft-gate/SKILL.md	RequireInSection	Prove it bites	27	DROP	rule-text
registry_data.go	441	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_data.go	442	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	53	DROP	sentence-with-label
registry_data.go	443	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	57	DROP	rule-text
registry_data.go	444	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_data.go	445	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	58	DROP	rule-text
registry_data.go	446	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	59	DROP	rule-text
registry_data.go	447	.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md	Require		33	GRAY	override
registry_data.go	448	.agents/commands/bench-implement-spec.md	RequireInSection	`--full <spec>`	85	DROP	rule-text
registry_data.go	449	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	77	GRAY	sentence-command-with-args
registry_data.go	450	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	71	GRAY	sentence-command-with-args
registry_data.go	451	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Probes	128	GRAY	sentence-command-with-args
registry_data.go	452	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Probes	125	DROP	rule-text
registry_data.go	453	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	69	DROP	rule-text
registry_data.go	454	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	83	DROP	rule-text
registry_data.go	455	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	81	DROP	rule-text
registry_data.go	456	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	37	DROP	rule-text
registry_data.go	457	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	18	DROP	rule-text
registry_data.go	458	.agents/skills/bench-craft-review/references/finding-discipline.md	Require		3	DROP	sentence-with-route
registry_data.go	459	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a string expectation proves	9	DROP	rule-text
registry_data.go	460	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a citation points at	15	DROP	rule-text
registry_data.go	461	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	23	DROP	rule-text
registry_data.go	462	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	25	DROP	rule-text
registry_data.go	463	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	30	DROP	rule-text
registry_data.go	464	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	When a seam cannot reach the state	59	DROP	rule-text
registry_data.go	465	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	What a finding must cite	103	DROP	rule-text
registry_data.go	466	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	65	DROP	sentence-with-route
registry_data.go	467	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	66	DROP	rule-text
registry_data.go	468	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	49	DROP	rule-text
registry_data.go	469	.agents/skills/bench-craft-spec/references/ste-prose.md	Require		31	DROP	rule-text
registry_data.go	470	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Isolation	96	GRAY	sentence-command-with-args
registry_data.go	471	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Isolation	96	DROP	rule-text
registry_data.go	472	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Isolation	97	DROP	rule-text
registry_data.go	473	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	52	DROP	rule-text
registry_data.go	474	.agents/skills/bench-craft-tdd/SKILL.md	RequireInSection	The oracle is the gate, not you	115	DROP	rule-text
registry_data.go	475	.agents/commands/bench-final-check.md	RequireInSection	Exit handoff	77	DROP	sentence-command-noargs
registry_data.go	476	.agents/commands/bench-final-check.md	RequireInSection	Exit handoff	55	DROP	rule-text
registry_data.go	477	.agents/commands/bench-review-implementation.md	RequireInSection	Process	97	DROP	rule-text
registry_front_door.go	4	.agents/commands/bench.md	RequireInSection	Exit handoff	24	DROP	sentence-with-route
registry_decision_maps.go	9	.agents/commands/bench-shape-idea.md	Require	Map content	97	DROP	rule-text
registry_decision_maps.go	10	.agents/commands/bench-shape-idea.md	Require	Map content	98	DROP	sentence-command-noargs
registry_decision_maps.go	11	.agents/commands/bench-shape-idea.md	Forbid		0	FORBID	kind
registry_decision_maps.go	12	.agents/commands/bench-shape-idea.md	Forbid		0	FORBID	kind
registry_decision_maps.go	13	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_decision_maps.go	14	internal/maps/schema.go	Require		70	DROP	rule-text
registry_decision_maps.go	15	.agents/commands/bench-write-spec.md	Require	Entry contract	30	DROP	rule-text
registry_decision_maps.go	16	.agents/commands/bench-write-spec.md	Require	Process	65	DROP	rule-text
registry_decision_maps.go	17	.agents/commands/bench-shape-idea.md	Require	Entry orientation	19	DROP	rule-text
registry_decision_maps.go	18	.agents/commands/bench-shape-idea.md	Require	Entry orientation	20	DROP	rule-text
registry_decision_maps.go	19	.agents/commands/bench-shape-idea.md	Require	Entry orientation	21	DROP	rule-text
registry_decision_maps.go	20	.agents/commands/bench-shape-idea.md	Require	Map content	85	DROP	rule-text
registry_decision_maps.go	21	.agents/commands/bench-shape-idea.md	Require	Map content	86	DROP	rule-text
registry_decision_maps.go	22	.agents/commands/bench-shape-idea.md	Require	Decision tickets	72	DROP	rule-text
registry_decision_maps.go	23	.agents/commands/bench-shape-idea.md	Require	Decision tickets	71	DROP	rule-text
registry_decision_maps.go	24	.agents/commands/bench-shape-idea.md	Require	Decision tickets	79	DROP	rule-text
registry_decision_maps.go	25	README.md	Require	Reviewer quick start	169	DROP	rule-text
registry_decision_maps.go	26	CONTEXT.md	Require	Core terms	288	DROP	rule-text
registry_decision_maps.go	27	CONTEXT.md	Require	Core terms	263	DROP	rule-text
registry_decision_maps.go	28	CONTEXT.md	Require	Core terms	300	DROP	rule-text
registry_decision_maps.go	29	CONTEXT.md	Require	Core terms	275	DROP	rule-text
registry_craft_research.go	9	.agents/skills/bench-craft-research/SKILL.md	Require	The question graph	27	DROP	rule-text
registry_craft_research.go	10	.agents/skills/bench-craft-research/SKILL.md	Require	Rounds and fan-out	40	DROP	rule-text
registry_craft_research.go	11	.agents/skills/bench-craft-research/SKILL.md	Require	Rounds and fan-out	62	DROP	rule-text
registry_craft_research.go	12	.agents/skills/bench-craft-research/SKILL.md	Require	Sources and verification	78	DROP	rule-text
registry_craft_research.go	13	.agents/skills/bench-craft-research/SKILL.md	Require	Sources and verification	88	DROP	rule-text
registry_craft_research.go	14	.agents/skills/bench-craft-research/SKILL.md	Require	The durable output	98	DROP	rule-text
registry_craft_research.go	15	.agents/skills/bench-craft-research/SKILL.md	Require	The report contract	118	DROP	rule-text
registry_craft_research.go	16	.agents/skills/bench-craft-research/SKILL.md	Require	Compatibility evidence	137	DROP	rule-text
registry_craft_research.go	17	.agents/skills/bench-craft-research/SKILL.md	Require	Trigger and ownership	21	DROP	rule-text
registry_craft_research.go	18	.agents/commands/bench-shape-idea.md	Require	Decision tickets	65	DROP	sentence-with-route
registry_craft_research.go	19	.agents/skills/bench-craft-spec/references/bootstrap-authority.md	Require		22	DROP	sentence-with-route
registry_craft_research.go	20	.agents/commands/bench-shape-idea.md	Forbid		0	FORBID	kind
registry_craft_research.go	21	.agents/commands/bench-shape-idea.md	Forbid		0	FORBID	kind
registry_craft_research.go	22	.agents/commands/bench-shape-idea.md	Forbid		0	FORBID	kind
registry_craft_research.go	23	.agents/commands/bench-write-spec.md	Require	Who runs this phase	42	DROP	rule-text
registry_craft_research.go	24	.agents/commands/bench-write-spec.md	Require	Who runs this phase	51	DROP	rule-text
registry_ft311_preparation.go	9	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	23	GRAY	sentence-command-with-args
registry_ft311_preparation.go	10	.agents/commands/bench-implement-spec.md	Require	Build	48	DROP	rule-text
registry_ft311_preparation.go	11	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	48	DROP	rule-text
registry_ft311_preparation.go	12	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Require	Read-only returns	167	DROP	rule-text
registry_ft311_preparation.go	20	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	29	GRAY	sentence-command-with-args
registry_ft311_preparation.go	21	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	31	DROP	rule-text
registry_ft311_preparation.go	22	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	31	GRAY	sentence-command-with-args
registry_ft311_preparation.go	23	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	31	DROP	rule-text
registry_ft311_preparation.go	24	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	31	DROP	rule-text
registry_ft311_preparation.go	25	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_ft311_preparation.go	32	.agents/commands/bench-implement-spec.md	Require	Build	50	DROP	rule-text
registry_ft311_preparation.go	33	.agents/commands/bench-implement-spec.md	Require	Build	50	DROP	rule-text
registry_ft311_preparation.go	39	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	25	DROP	sentence-with-label
registry_ft311_preparation.go	40	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_ft311_preparation.go	41	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_ft311_preparation.go	42	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	33	DROP	rule-text
registry_ft311_preparation.go	43	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	33	DROP	sentence-command-noargs
registry_ft311_preparation.go	49	.agents/commands/bench-implement-spec.md	Require	`--full <spec>`	83	DROP	sentence-with-route
registry_ft311_preparation.go	50	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	10	.agents/commands/bench-review-implementation.md	Require	Review modes	64	DROP	rule-text
registry_ft311_review_dispatch.go	11	.agents/commands/bench-review-implementation.md	Require	Process	80	GRAY	sentence-command-with-args
registry_ft311_review_dispatch.go	12	.agents/commands/bench-review-implementation.md	Require	Process	111	GRAY	sentence-command-with-args
registry_ft311_review_dispatch.go	13	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	14	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	15	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	16	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	17	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Require	Read-only returns	163	DROP	rule-text
registry_ft311_review_dispatch.go	18	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Require	Read-only returns	165	DROP	rule-text
registry_ft311_review_dispatch.go	19	.agents/commands/bench-review-implementation.md	Require	Entry orientation	14	DROP	rule-text
registry_ft311_review_dispatch.go	20	.agents/commands/bench-review-implementation.md	Require	Entry orientation	16	GRAY	sentence-command-with-args
registry_ft311_review_dispatch.go	21	.agents/commands/bench-review-implementation.md	Require	Entry orientation	18	DROP	rule-text
registry_ft311_review_dispatch.go	22	.agents/commands/bench-review-implementation.md	Require	Review modes	48	DROP	rule-text
registry_ft311_review_dispatch.go	23	.agents/commands/bench-review-implementation.md	Require	Process	140	DROP	rule-text
registry_ft311_review_dispatch.go	24	.agents/commands/bench-review-implementation.md	Require	Process	164	DROP	rule-text
registry_ft311_review_dispatch.go	25	.agents/commands/bench-review-implementation.md	Require	Process	167	DROP	rule-text
registry_ft311_review_dispatch.go	26	.agents/commands/bench-review-implementation.md	Require	Process	178	DROP	rule-text
registry_ft311_review_dispatch.go	27	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	28	.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md	Require		3	DROP	sentence-with-route
registry_ft311_review_dispatch.go	29	.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	30	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	41	DROP	rule-text
registry_ft311_review_dispatch.go	31	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	111	LABEL	override
registry_ft311_review_dispatch.go	32	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	112	LABEL	override
registry_ft311_review_dispatch.go	33	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Template	113	LABEL	override
registry_ft311_review_dispatch.go	34	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	60	DROP	rule-text
registry_ft311_review_dispatch.go	35	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	60	DROP	rule-text
registry_ft311_review_dispatch.go	36	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	62	DROP	rule-text
registry_ft311_review_dispatch.go	37	.agents/commands/bench-review-implementation.md	Require	Process	135	DROP	sentence-with-route
registry_ft311_review_dispatch.go	38	projects/benchkit.md	Require	Lines (model + effort routing)	607	ROUTE	short-route-ref
registry_ft311_review_dispatch.go	39	.bench/BENCH.md	Require	Workflow	116	DROP	bench-workflow
registry_ft311_review_dispatch.go	40	.bench/BENCH.md	Require	Workflow	116	DROP	bench-workflow
registry_ft311_review_dispatch.go	41	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	42	.agents/commands/bench-review-implementation.md	Require	Exit handoff	34	DROP	rule-text
registry_ft311_review_dispatch.go	43	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	44	docs/field-guide.html	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	45	.bench/BENCH.md	Require	Workflow	118	DROP	bench-workflow
registry_ft311_review_dispatch.go	46	.bench/BENCH.md	Require	Workflow	118	DROP	bench-workflow
registry_ft311_review_dispatch.go	47	.bench/BENCH.md	Require	Workflow	120	DROP	bench-workflow
registry_ft311_review_dispatch.go	48	.bench/BENCH.md	Require	Workflow	120	DROP	bench-workflow
registry_ft311_review_dispatch.go	49	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	50	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	51	.bench/BENCH.md	Require	Workflow	120	DROP	bench-workflow
registry_ft311_review_dispatch.go	52	.bench/BENCH.md	Require	Workflow	116	DROP	bench-workflow
registry_ft311_review_dispatch.go	53	.agents/commands/bench-drain.md	Require	Delegate the evidence	76	DROP	rule-text
registry_ft311_review_dispatch.go	54	.agents/commands/bench-review-implementation.md	Require	Entry orientation	20	GRAY	sentence-command-with-args
registry_ft311_review_dispatch.go	55	.agents/commands/bench-review-implementation.md	Require	Entry orientation	22	DROP	rule-text
registry_ft311_review_dispatch.go	56	.agents/commands/bench-review-implementation.md	Require	Entry orientation	22	GRAY	sentence-command-with-args
registry_ft311_review_dispatch.go	57	.agents/commands/bench-review-implementation.md	Require	Entry orientation	22	DROP	rule-text
registry_ft311_review_dispatch.go	58	.agents/commands/bench-review-implementation.md	Require	Entry orientation	22	DROP	rule-text
registry_ft311_review_dispatch.go	59	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_ft311_review_dispatch.go	63	.agents/commands/bench-review-implementation.md	Require	Entry orientation	24	DROP	rule-text
registry_ft311_review_dispatch.go	64	.agents/commands/bench-review-implementation.md	Require	Entry orientation	24	DROP	rule-text
registry_ft311_review_dispatch.go	69	.agents/commands/bench-review-implementation.md	Require	Process	171	GRAY	sentence-command-with-args
registry_retained_workflow.go	212	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Slicing a build for delegates	68	DROP	rule-text
registry_retained_workflow.go	213	.bench/BENCH.md	Require	Workflow	124	DROP	bench-workflow
registry_retained_workflow.go	214	.bench/BENCH.md	Require	Workflow	124	DROP	bench-workflow
registry_retained_workflow.go	215	.bench/BENCH.md	Require	Workflow	124	DROP	bench-workflow
registry_retained_workflow.go	216	.bench/BENCH.md	Require	Workflow	124	DROP	bench-workflow
registry_retained_workflow.go	217	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	218	.agents/commands/bench-implement-spec.md	Require	Build	54	DROP	sentence-with-label
registry_retained_workflow.go	219	.agents/skills/bench-craft-gate/SKILL.md	Require	Prove it bites	28	DROP	rule-text
registry_retained_workflow.go	220	.agents/skills/bench-craft-tickets/SKILL.md	Require	Acceptance	67	DROP	sentence-with-label
registry_retained_workflow.go	221	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	237	.bench/BENCH.md	Require	Workflow	128	DROP	bench-workflow
registry_retained_workflow.go	238	.bench/BENCH.md	Require	Workflow	130	DROP	bench-workflow
registry_retained_workflow.go	239	.bench/BENCH.md	Require	Workflow	130	DROP	bench-workflow
registry_retained_workflow.go	240	.bench/BENCH.md	Require	Workflow	130	DROP	bench-workflow
registry_retained_workflow.go	241	.bench/BENCH.md	Require	Workflow	128	DROP	bench-workflow
registry_retained_workflow.go	242	.bench/BENCH.md	Require	Workflow	122	DROP	bench-workflow
registry_retained_workflow.go	243	.bench/BENCH.md	Require	Workflow	148	DROP	bench-workflow
registry_retained_workflow.go	244	.bench/BENCH.md	Require	Workflow	148	DROP	bench-workflow
registry_retained_workflow.go	245	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	246	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	247	.agents/skills/bench-craft-delegate/SKILL.md	Require	Delegate or retain	17	DROP	sentence-with-route
registry_retained_workflow.go	248	.agents/skills/bench-craft-delegate/SKILL.md	Require	The agent type	28	DROP	rule-text
registry_retained_workflow.go	249	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	250	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	251	.agents/skills/bench-craft-delegate/SKILL.md	Require	Isolation	101	DROP	rule-text
registry_retained_workflow.go	252	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	253	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	254	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	255	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	256	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	72	DROP	rule-text
registry_retained_workflow.go	257	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	258	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	259	.agents/skills/bench-craft-delegate/SKILL.md	Require	Isolation	103	DROP	rule-text
registry_retained_workflow.go	260	.claude/agents/bench-writer.md	Require		3	DROP	rule-text
registry_retained_workflow.go	261	.claude/agents/bench-writer.md	Require		21	DROP	sentence-command-noargs
registry_retained_workflow.go	262	.claude/agents/bench-writer.md	Require		23	DROP	rule-text
registry_retained_workflow.go	263	.claude/agents/bench-writer.md	Require		20	DROP	sentence-with-route
registry_retained_workflow.go	264	.claude/agents/bench-writer.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	265	.claude/agents/bench-writer.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	266	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	115	DROP	rule-text
registry_retained_workflow.go	267	.agents/skills/bench-craft-line/SKILL.md	Require	Classify reds before the ladder moves	85	DROP	rule-text
registry_retained_workflow.go	268	.agents/skills/bench-craft-tickets/SKILL.md	Require	Land the frontier	93	DROP	rule-text
registry_retained_workflow.go	269	README.md	Require	Reviewer quick start	68	DROP	rule-text
registry_retained_workflow.go	270	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	271	.bench/BENCH-reference.md	Require	Command Notes	339	DROP	rule-text
registry_retained_workflow.go	272	.agents/skills/bench-craft-line/SKILL.md	Require	Delegated author lines	71	DROP	rule-text
registry_retained_workflow.go	273	.agents/skills/bench-craft-line/SKILL.md	Require	Delegated author lines	72	DROP	rule-text
registry_retained_workflow.go	274	.agents/skills/bench-craft-line/SKILL.md	Require	Delegated author lines	73	DROP	rule-text
registry_retained_workflow.go	275	.agents/skills/bench-craft-line/SKILL.md	Require	Delegated author lines	74	DROP	rule-text
registry_retained_workflow.go	276	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	101	DROP	rule-text
registry_retained_workflow.go	277	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	102	DROP	rule-text
registry_retained_workflow.go	278	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	103	DROP	rule-text
registry_retained_workflow.go	279	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	104	DROP	rule-text
registry_retained_workflow.go	280	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	105	DROP	rule-text
registry_retained_workflow.go	281	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	107	DROP	rule-text
registry_retained_workflow.go	282	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Delegated author transfer	106	DROP	rule-text
registry_retained_workflow.go	283	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	39	DROP	sentence-with-label
registry_retained_workflow.go	284	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	285	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	286	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	287	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	288	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	289	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	290	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	170	DROP	sentence-command-noargs
registry_retained_workflow.go	291	.agents/commands/bench-implement-spec.md	Require	`--full <spec>`	90	DROP	rule-text
registry_retained_workflow.go	292	.agents/commands/bench-implement-spec.md	Require	`--full <spec>`	90	DROP	rule-text
registry_retained_workflow.go	293	.agents/commands/bench-implement-spec.md	Require	`--full <spec>`	90	DROP	rule-text
registry_retained_workflow.go	294	.bench/BENCH.md	Require	Workflow	130	DROP	bench-workflow
registry_retained_workflow.go	295	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	296	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	297	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	298	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	299	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	300	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	301	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	302	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	303	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	44	DROP	rule-text
registry_retained_workflow.go	304	.agents/commands/bench-final-check.md	Require	Report	198	DROP	rule-text
registry_retained_workflow.go	305	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	306	.agents/commands/bench-debug.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	307	.agents/commands/bench-debug.md	Require	How it meets the rest of Bench	157	DROP	rule-text
registry_retained_workflow.go	308	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	309	.agents/commands/bench-final-check.md	Require	Entry orientation	17	DROP	rule-text
registry_retained_workflow.go	310	.agents/commands/bench-final-check.md	Require	Report	197	DROP	rule-text
registry_retained_workflow.go	311	.agents/commands/bench-review-implementation.md	Require	Review modes	62	DROP	rule-text
registry_retained_workflow.go	312	.agents/commands/bench-final-check.md	Require	Entry orientation	39	DROP	rule-text
registry_retained_workflow.go	313	.agents/commands/bench-final-check.md	Require	Entry orientation	39	DROP	rule-text
registry_retained_workflow.go	314	.agents/commands/bench-final-check.md	Require	Entry orientation	39	DROP	rule-text
registry_retained_workflow.go	315	.agents/commands/bench-final-check.md	Require	Entry orientation	39	DROP	rule-text
registry_retained_workflow.go	316	.agents/skills/bench-craft-tickets/SKILL.md	Require	Land the frontier	95	DROP	rule-text
registry_retained_workflow.go	186	.agents/skills/bench-craft-synthesis/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	187	.agents/skills/bench-craft-synthesis/SKILL.md	Require	Three quality loops — all must pass before adoption	65	GRAY	sentence-command-with-args
registry_retained_workflow.go	188	.agents/skills/bench-craft-synthesis/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	189	.agents/skills/bench-craft-synthesis/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	190	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	191	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	192	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	193	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	194	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	195	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	196	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	197	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	198	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	199	.agents/commands/bench-final-check.md	Require	Exit handoff	53	DROP	sentence-command-noargs
registry_retained_workflow.go	200	.agents/commands/bench-final-check.md	Require	Entry orientation	31	GRAY	sentence-command-with-args
registry_retained_workflow.go	201	.agents/commands/bench-final-check.md	Require	Exit handoff	67	GRAY	sentence-command-with-args
registry_retained_workflow.go	202	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	203	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	204	.bench/BENCH-reference.md	Require	Command Notes	353	DROP	sentence-command-noargs
registry_retained_workflow.go	205	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	206	.bench/BENCH-reference.md	Require	Command Notes	406	DROP	rule-text
registry_retained_workflow.go	161	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	101	GRAY	sentence-command-with-args
registry_retained_workflow.go	162	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	102	DROP	rule-text
registry_retained_workflow.go	163	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	159	DROP	rule-text
registry_retained_workflow.go	164	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	161	DROP	rule-text
registry_retained_workflow.go	165	.agents/commands/bench-final-check.md	Require	Report	189	DROP	rule-text
registry_retained_workflow.go	166	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	167	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	168	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	169	.agents/commands/bench-final-check.md	Require	Capture the implementation retro	160	DROP	sentence-command-noargs
registry_retained_workflow.go	170	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	171	.bench/BENCH-reference.md	Require	Files	48	DROP	sentence-with-route
registry_retained_workflow.go	172	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	139	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	62	DROP	sentence-with-label
registry_retained_workflow.go	140	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	141	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	142	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	143	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	144	.agents/skills/bench-craft-line/SKILL.md	Require	Classify reds before the ladder moves	87	DROP	rule-text
registry_retained_workflow.go	145	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	146	.agents/skills/bench-craft-line/SKILL.md	Require	Classify reds before the ladder moves	92	DROP	rule-text
registry_retained_workflow.go	147	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	148	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	149	.agents/skills/bench-craft-spec/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	88	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	89	.bench/BENCH.md	Require	The Bench CLI contract	16	INVARIANT	bench-invariant-section
registry_retained_workflow.go	90	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	91	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	92	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	93	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	94	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	95	AGENTS.md	Require		6	DROP	rule-text
registry_retained_workflow.go	96	AGENTS.md	Require	This repo	48	DROP	sentence-command-noargs
registry_retained_workflow.go	97	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	98	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	99	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	100	AGENTS.md	Require	This repo	47	DROP	rule-text
registry_retained_workflow.go	101	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	102	.bench/BENCH-reference.md	Require	Command Notes	457	DROP	rule-text
registry_retained_workflow.go	103	.bench/BENCH-reference.md	Require	Command Notes	458	DROP	rule-text
registry_retained_workflow.go	104	.bench/BENCH-reference.md	Require	Command Notes	459	DROP	rule-text
registry_retained_workflow.go	105	AGENTS.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	106	.bench/BENCH-reference.md	RequireInSection	Plumbing subcommands	479	DROP	sentence-command-noargs
registry_retained_workflow.go	107	.bench/BENCH-reference.md	Require	Phase manifest	569	DROP	rule-text
registry_retained_workflow.go	108	projects/benchkit.md	Require	Gate (`.bench/gate.sh`)	486	DROP	rule-text
registry_retained_workflow.go	109	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	110	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	111	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	112	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	113	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	114	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	115	projects/benchkit.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	116	.bench/BENCH-reference.md	Require	Command Notes	456	DROP	sentence-command-noargs
registry_retained_workflow.go	117	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	118	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	119	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	120	.bench/BENCH-reference.md	Require	Phase manifest	506	DROP	rule-text
registry_retained_workflow.go	121	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	122	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	123	.bench/BENCH-reference.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	124	.bench/BENCH-reference.md	Require	Files	73	DROP	rule-text
registry_retained_workflow.go	125	.claude/README.md	Require		8	DROP	rule-text
registry_retained_workflow.go	126	.claude/README.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	127	.claude/README.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	128	.claude/README.md	Require		23	DROP	rule-text
registry_retained_workflow.go	27	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	28	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	29	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	291	DROP	rule-text
registry_retained_workflow.go	30	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	31	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	32	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	33	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	34	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	35	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	36	.agents/commands/bench-assess.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	37	.agents/commands/bench-assess.md	Require	Exit handoff	29	DROP	sentence-command-noargs
registry_retained_workflow.go	38	.agents/skills/bench-craft-skills/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	39	.agents/skills/bench-craft-skills/SKILL.md	Require	Invocation: model-invoked vs user-invoked	22	DROP	rule-text
registry_retained_workflow.go	40	.agents/skills/bench-craft-cli/SKILL.md	Require	Bench application	83	DROP	rule-text
registry_retained_workflow.go	41	.agents/skills/bench-craft-cli/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	42	.agents/skills/bench-craft-cli/SKILL.md	Require	Bench application	87	GRAY	sentence-command-with-args
registry_retained_workflow.go	43	.agents/skills/bench-craft-spec/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	44	.agents/skills/bench-craft-spec/SKILL.md	Require	The acceptance coverage map	45	DROP	rule-text
registry_retained_workflow.go	45	.agents/skills/bench-craft-spec/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	46	.agents/skills/bench-craft-spec/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	47	.agents/skills/bench-craft-spec/SKILL.md	Require	Slicing a build for delegates	70	DROP	rule-text
registry_retained_workflow.go	48	.agents/skills/bench-craft-tdd/references/tests.md	Require		16	DROP	rule-text
registry_retained_workflow.go	49	.agents/skills/bench-craft-tdd/references/tests.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	50	.agents/skills/bench-craft-adr/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	51	.agents/skills/bench-craft-adr/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	52	.agents/skills/bench-craft-adr/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	53	.agents/skills/bench-craft-adr/SKILL.md	Require	The rules	43	DROP	rule-text
registry_retained_workflow.go	325	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	The declaration	110	DROP	rule-text
registry_retained_workflow.go	326	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	116	DROP	rule-text
registry_retained_workflow.go	327	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	120	DROP	rule-text
registry_retained_workflow.go	328	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	123	DROP	rule-text
registry_retained_workflow.go	329	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	129	DROP	rule-text
registry_retained_workflow.go	330	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	131	DROP	rule-text
registry_retained_workflow.go	331	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	116	DROP	rule-text
registry_retained_workflow.go	332	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	117	DROP	rule-text
registry_retained_workflow.go	333	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	120	DROP	rule-text
registry_retained_workflow.go	334	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	124	DROP	rule-text
registry_retained_workflow.go	335	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	126	DROP	rule-text
registry_retained_workflow.go	336	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	127	DROP	rule-text
registry_retained_workflow.go	337	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	128	DROP	rule-text
registry_retained_workflow.go	338	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	131	DROP	rule-text
registry_retained_workflow.go	339	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	131	DROP	rule-text
registry_retained_workflow.go	340	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	121	DROP	rule-text
registry_retained_workflow.go	341	.bench/BENCH.md	RequireInSection	Workflow	116	DROP	bench-workflow
registry_retained_workflow.go	342	.agents/commands/bench-implement-spec.md	RequireInSection	Build	44	DROP	sentence-with-route
registry_retained_workflow.go	343	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	135	DROP	rule-text
registry_retained_workflow.go	344	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	124	DROP	sentence-with-route
registry_retained_workflow.go	345	.agents/skills/bench-craft-delegate/SKILL.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	346	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	133	DROP	rule-text
registry_retained_workflow.go	347	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Ticketed-build stage defaults	62	DROP	rule-text
registry_retained_workflow.go	348	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	Delegate or retain	15	DROP	rule-text
registry_retained_workflow.go	349	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	156	DROP	rule-text
registry_retained_workflow.go	350	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	157	DROP	rule-text
registry_retained_workflow.go	351	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	133	DROP	rule-text
registry_retained_workflow.go	352	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	160	DROP	rule-text
registry_retained_workflow.go	353	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Read-only returns	161	DROP	rule-text
registry_retained_workflow.go	354	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	5	DROP	rule-text
registry_retained_workflow.go	355	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	Forbid		0	FORBID	kind
registry_retained_workflow.go	356	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	6	DROP	rule-text
registry_retained_workflow.go	357	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	7	DROP	rule-text
registry_retained_workflow.go	358	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	9	DROP	rule-text
registry_retained_workflow.go	359	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	11	DROP	rule-text
registry_retained_workflow.go	360	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	13	DROP	rule-text
registry_retained_workflow.go	361	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	16	DROP	rule-text
registry_retained_workflow.go	362	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	17	DROP	rule-text
registry_retained_workflow.go	363	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	18	DROP	rule-text
registry_retained_workflow.go	364	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	19	DROP	rule-text
registry_retained_workflow.go	365	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	24	DROP	rule-text
registry_retained_workflow.go	366	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	25	DROP	rule-text
registry_retained_workflow.go	367	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	27	DROP	rule-text
registry_retained_workflow.go	368	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	28	DROP	rule-text
registry_retained_workflow.go	369	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	29	DROP	rule-text
registry_retained_workflow.go	370	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	30	DROP	rule-text
registry_retained_workflow.go	371	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	31	DROP	rule-text
registry_retained_workflow.go	372	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	33	DROP	rule-text
registry_retained_workflow.go	373	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	35	DROP	rule-text
registry_retained_workflow.go	374	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	36	DROP	rule-text
registry_retained_workflow.go	375	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	44	DROP	rule-text
registry_retained_workflow.go	376	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	54	DROP	rule-text
registry_retained_workflow.go	377	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	46	DROP	rule-text
registry_retained_workflow.go	378	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	47	DROP	rule-text
registry_retained_workflow.go	379	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	49	DROP	rule-text
registry_retained_workflow.go	380	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	50	DROP	rule-text
registry_retained_workflow.go	381	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	51	DROP	rule-text
registry_retained_workflow.go	382	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	52	DROP	rule-text
registry_retained_workflow.go	383	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	61	DROP	rule-text
registry_retained_workflow.go	384	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	59	DROP	rule-text
registry_retained_workflow.go	385	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	65	DROP	rule-text
registry_retained_workflow.go	386	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	66	DROP	rule-text
registry_retained_workflow.go	387	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	68	DROP	rule-text
registry_retained_workflow.go	388	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	69	DROP	rule-text
registry_retained_workflow.go	389	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	70	DROP	rule-text
registry_retained_workflow.go	390	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	72	DROP	rule-text
registry_retained_workflow.go	391	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	73	DROP	rule-text
registry_retained_workflow.go	392	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	State and handoff	75	DROP	rule-text
registry_retained_workflow.go	393	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	Retained implementation continuation	114	DROP	rule-text
registry_retained_workflow.go	394	.agents/commands/bench-implement-spec.md	RequireInSection	Build	44	DROP	sentence-with-route
registry_retained_workflow.go	395	.agents/skills/bench-craft-review/SKILL.md	Require		14	DROP	sentence-with-route
registry_retained_workflow.go	396	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	42	DROP	sentence-with-route
registry_retained_workflow.go	397	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	75	DROP	rule-text
registry_ticket_passes.go	4	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	28	DROP	rule-text
registry_ticket_passes.go	5	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	28	DROP	rule-text
registry_ticket_passes.go	6	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	29	DROP	sentence-with-route
registry_ticket_passes.go	7	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Write one file per ticket	65	DROP	sentence-with-label
registry_ticket_passes.go	8	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	20	DROP	rule-text
registry_ticket_passes.go	9	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Write one file per ticket	63	DROP	sentence-with-label
registry_ticket_passes.go	10	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Enforced rules	22	DROP	sentence-with-label
registry_ticket_passes.go	11	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Enforced rules	23	DROP	sentence-with-label
registry_ticket_passes.go	12	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	34	DROP	rule-text
registry_ticket_passes.go	13	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	36	DROP	rule-text
registry_ticket_passes.go	14	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	37	DROP	rule-text
registry_ticket_passes.go	15	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	39	DROP	rule-text
registry_ticket_passes.go	16	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	40	DROP	rule-text
registry_ticket_passes.go	17	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	42	GRAY	sentence-command-with-args
registry_ticket_passes.go	18	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	46	DROP	rule-text
registry_ticket_passes.go	19	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	50	DROP	rule-text
registry_ticket_passes.go	20	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	51	DROP	rule-text
registry_ticket_passes.go	21	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	52	DROP	sentence-with-label
registry_ticket_passes.go	22	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	53	DROP	sentence-with-label
registry_ticket_passes.go	23	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	54	DROP	sentence-with-label
registry_ticket_passes.go	24	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Slicing rules	59	DROP	sentence-with-label
registry_ticket_passes.go	25	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	69	DROP	rule-text
registry_ticket_passes.go	26	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	70	DROP	sentence-with-label
registry_ticket_passes.go	27	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	71	DROP	rule-text
registry_ticket_passes.go	28	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	72	DROP	rule-text
registry_ticket_passes.go	29	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Per row	124	GRAY	sentence-command-with-args
registry_ticket_passes.go	30	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	63	DROP	rule-text
registry_ticket_passes.go	31	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	72	DROP	rule-text
registry_ticket_passes.go	32	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	65	DROP	sentence-command-noargs
registry_ticket_passes.go	33	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	66	DROP	rule-text
registry_ticket_passes.go	34	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	67	DROP	rule-text
registry_ticket_passes.go	35	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	68	DROP	rule-text
registry_ticket_passes.go	36	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	69	DROP	rule-text
registry_ticket_passes.go	37	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	70	DROP	rule-text
registry_ticket_passes.go	38	.agents/commands/bench-write-spec.md	RequireInSection	Process	66	DROP	sentence-with-route
registry_ticket_passes.go	39	.agents/commands/bench-write-spec.md	RequireInSection	Process	67	DROP	rule-text
registry_ticket_passes.go	40	.agents/commands/bench-write-spec.md	RequireInSection	Process	71	DROP	rule-text
registry_ticket_passes.go	41	.agents/commands/bench-write-spec.md	RequireInSection	Process	71	DROP	rule-text
registry_ticket_passes.go	42	.agents/commands/bench-write-spec.md	RequireInSection	Process	71	DROP	rule-text
registry_ticket_passes.go	43	.agents/commands/bench-write-spec.md	RequireInSection	Process	71	DROP	rule-text
registry_ticket_passes.go	44	.agents/commands/bench-write-spec.md	RequireInSection	Process	73	DROP	sentence-with-route
registry_ticket_passes.go	45	.agents/commands/bench-write-spec.md	RequireInSection	Process	73	DROP	sentence-with-label
registry_ticket_passes.go	46	.agents/commands/bench-write-spec.md	RequireInSection	Process	73	GRAY	sentence-command-with-args
registry_ticket_passes.go	47	.agents/commands/bench-write-spec.md	RequireInSection	Process	75	DROP	rule-text
registry_ticket_passes.go	48	.agents/commands/bench-write-spec.md	RequireInSection	Process	75	DROP	rule-text
registry_ticket_passes.go	49	.agents/commands/bench-write-spec.md	RequireInSection	Process	77	DROP	sentence-with-route
registry_ticket_passes.go	50	.agents/commands/bench-write-spec.md	Forbid		0	FORBID	kind
registry_ticket_passes.go	51	.agents/commands/bench-write-spec.md	RequireInSection	Process	79	GRAY	sentence-command-with-args
registry_ticket_passes.go	52	.agents/commands/bench-write-spec.md	RequireInSection	Process	83	DROP	rule-text
registry_ticket_passes.go	53	.agents/commands/bench-write-spec.md	RequireInSection	Process	85	DROP	rule-text
registry_ticket_passes.go	54	.agents/commands/bench-write-spec.md	Require	Who runs this phase	51	DROP	rule-text
registry_ticket_passes.go	55	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	54	DROP	sentence-with-route
registry_ticket_passes.go	56	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	55	DROP	sentence-with-route
registry_ticket_passes.go	57	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	56	DROP	sentence-with-route
registry_ticket_passes.go	58	.agents/skills/bench-craft-line/SKILL.md	Forbid		0	FORBID	kind
registry_ticket_passes.go	59	.agents/skills/bench-craft-line/SKILL.md	Require	Ticketed-build stage defaults	58	DROP	rule-text
registry_ticket_passes.go	60	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	74	DROP	sentence-with-label
registry_ticket_passes.go	61	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	74	DROP	rule-text
registry_ticket_passes.go	62	.agents/skills/bench-craft-tickets/references/slicing-checks.md	RequireInSection	Fence-expansion evidence	74	DROP	sentence-with-label
registry_ticket_passes.go	63	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	25	DROP	rule-text
registry_ticket_passes.go	64	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	25	DROP	rule-text
registry_ticket_passes.go	65	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	25	DROP	sentence-with-label
registry_ticket_passes.go	66	.agents/commands/bench-implement-spec.md	Require	Build	54	DROP	sentence-with-label
registry_ticket_passes.go	67	.agents/commands/bench-implement-spec.md	Require	Build	54	DROP	sentence-with-label
registry_ticket_passes.go	68	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	25	DROP	rule-text
registry_ticket_passes.go	69	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	27	DROP	rule-text
registry_ticket_passes.go	70	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	27	DROP	rule-text
registry_ticket_passes.go	71	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	27	DROP	rule-text
registry_ticket_passes.go	72	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	27	DROP	sentence-command-noargs
registry_ticket_passes.go	73	docs/known-issues.md	Forbid		0	FORBID	kind
registry_ticket_passes.go	74	docs/known-issues.md	Forbid		0	FORBID	kind
registry_ticket_passes.go	75	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Eligibility	5	DROP	rule-text
registry_ticket_passes.go	76	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Eligibility	6	DROP	sentence-with-route
registry_ticket_passes.go	77	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Eligibility	7	DROP	sentence-command-noargs
registry_ticket_passes.go	78	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Eligibility	9	DROP	rule-text
registry_ticket_passes.go	79	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	14	DROP	sentence-with-route
registry_ticket_passes.go	80	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	15	DROP	rule-text
registry_ticket_passes.go	81	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	16	DROP	rule-text
registry_ticket_passes.go	82	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	18	DROP	rule-text
registry_ticket_passes.go	83	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	19	DROP	rule-text
registry_ticket_passes.go	84	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	20	DROP	rule-text
registry_ticket_passes.go	85	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	21	DROP	rule-text
registry_ticket_passes.go	86	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	23	DROP	rule-text
registry_ticket_passes.go	87	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	25	DROP	sentence-with-route
registry_ticket_passes.go	88	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	26	DROP	rule-text
registry_ticket_passes.go	89	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	27	DROP	rule-text
registry_ticket_passes.go	90	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Route steps	29	DROP	rule-text
registry_ticket_passes.go	91	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Stop and return	35	DROP	rule-text
registry_ticket_passes.go	92	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Stop and return	36	DROP	rule-text
registry_ticket_passes.go	93	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Stop and return	37	DROP	rule-text
registry_ticket_passes.go	94	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Thin spec	41	DROP	rule-text
registry_ticket_passes.go	95	.agents/skills/bench-craft-spec/references/map-only-route.md	RequireInSection	Thin spec	42	DROP	rule-text
registry_ticket_passes.go	96	.agents/commands/bench-write-spec.md	RequireInSection	Who runs this phase	49	DROP	sentence-with-route
registry_debug_loop.go	4	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	Refute before you report	115	DROP	rule-text
registry_debug_loop.go	5	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	Refute before you report	116	DROP	rule-text
registry_debug_loop.go	6	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	27	DROP	rule-text
registry_debug_loop.go	7	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	28	DROP	rule-text
registry_debug_loop.go	8	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	29	DROP	rule-text
registry_debug_loop.go	9	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	Refute before you report	119	DROP	rule-text
registry_debug_loop.go	10	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	Refute before you report	120	DROP	sentence-with-route
registry_debug_loop.go	11	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	Refute before you report	117	DROP	rule-text
registry_debug_loop.go	12	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	Refute before you report	118	DROP	rule-text
registry_debug_loop.go	13	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	14	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	15	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	16	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	17	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	18	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	19	.agents/commands/bench-debug.md	ForbidInSection	How it meets the rest of Bench	0	FORBID	kind
registry_debug_loop.go	20	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	66	ROUTE	short-route-ref
registry_debug_loop.go	21	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	138	DROP	rule-text
registry_debug_loop.go	22	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	139	DROP	rule-text
registry_debug_loop.go	23	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Scope and allowance	14	DROP	sentence-with-route
registry_debug_loop.go	24	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	27	DROP	sentence-with-route
registry_debug_loop.go	25	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	150	DROP	rule-text
registry_debug_loop.go	26	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	151	DROP	rule-text
registry_debug_loop.go	27	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	152	DROP	rule-text
registry_debug_loop.go	28	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	153	DROP	rule-text
registry_debug_loop.go	29	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	154	DROP	rule-text
registry_debug_loop.go	30	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	146	DROP	rule-text
registry_debug_loop.go	31	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	147	DROP	rule-text
registry_debug_loop.go	32	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	145	DROP	rule-text
registry_debug_loop.go	33	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	144	DROP	rule-text
registry_debug_loop.go	34	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	142	DROP	rule-text
registry_debug_loop.go	35	.agents/commands/bench-debug.md	RequireInSection	How it meets the rest of Bench	140	DROP	rule-text
registry_debug_loop.go	36	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	22	DROP	rule-text
registry_debug_loop.go	37	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	23	DROP	rule-text
registry_debug_loop.go	38	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	24	DROP	rule-text
registry_debug_loop.go	39	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	25	DROP	rule-text
registry_debug_loop.go	40	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	28	DROP	rule-text
registry_debug_loop.go	41	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	31	DROP	rule-text
registry_debug_loop.go	42	.agents/skills/bench-craft-spec/SKILL.md	RequireInSection	Evidence-led authoring	32	DROP	rule-text
registry_debug_loop.go	43	.agents/skills/bench-craft-spec/SKILL.md	Forbid		0	FORBID	kind
registry_debug_loop.go	44	.agents/skills/bench-craft-spec/SKILL.md	ForbidCaseFoldedEmphasis		0	FORBID	kind
registry_debug_loop.go	45	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	31	DROP	rule-text
registry_debug_loop.go	46	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	31	DROP	rule-text
registry_debug_loop.go	47	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	31	DROP	rule-text
registry_debug_loop.go	48	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	31	DROP	rule-text
registry_debug_loop.go	49	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	33	DROP	rule-text
registry_debug_loop.go	50	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	33	DROP	rule-text
registry_debug_loop.go	51	.agents/skills/bench-craft-tickets/SKILL.md	RequireInSection	Draft the breakdown	33	DROP	rule-text
registry_debug_loop.go	52	.agents/commands/bench-implement-spec.md	RequireInSection	Build	46	DROP	sentence-with-route
registry_debug_loop.go	53	.agents/commands/bench-implement-spec.md	RequireInSection	Build	46	DROP	rule-text
registry_debug_loop.go	54	.agents/commands/bench-implement-spec.md	RequireInSection	Build	37	DROP	rule-text
registry_debug_loop.go	55	.agents/commands/bench-implement-spec.md	RequireInSection	Build	37	DROP	sentence-with-route
registry_debug_loop.go	56	.agents/commands/bench-implement-spec.md	RequireInSection	Build	39	DROP	rule-text
registry_debug_loop.go	57	.agents/commands/bench-implement-spec.md	RequireInSection	Build	39	DROP	rule-text
registry_debug_loop.go	58	.agents/commands/bench-implement-spec.md	RequireInSection	Build	39	DROP	rule-text
registry_debug_loop.go	59	.agents/commands/bench-implement-spec.md	RequireInSection	Build	39	DROP	rule-text
registry_debug_loop.go	60	.agents/commands/bench-implement-spec.md	RequireInSection	When the build stops short	79	DROP	sentence-with-route
registry_calibration.go	6	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	142	DROP	rule-text
registry_calibration.go	7	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	144	DROP	rule-text
registry_calibration.go	8	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	146	DROP	rule-text
registry_calibration.go	9	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	147	DROP	rule-text
registry_calibration.go	10	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	149	DROP	rule-text
registry_calibration.go	11	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	150	DROP	rule-text
registry_calibration.go	12	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	45	DROP	rule-text
registry_calibration.go	13	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a confidence states	64	DROP	rule-text
registry_calibration.go	14	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a confidence states	65	DROP	rule-text
registry_calibration.go	15	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a confidence states	67	DROP	rule-text
registry_calibration.go	16	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a confidence states	69	DROP	rule-text
registry_calibration.go	17	.agents/skills/bench-craft-review/SKILL.md	RequireInSection	What a finding must cite	103	DROP	rule-text
registry_calibration.go	18	.agents/commands/bench-review-implementation.md	RequireInStep	Process	209	DROP	rule-text
registry_calibration.go	19	.agents/skills/bench-craft-line/references/calibration-score.md	RequireInSection	What the declaration states	5	DROP	rule-text
registry_calibration.go	20	.agents/skills/bench-craft-line/references/calibration-score.md	RequireInSection	How one claim scores	11	DROP	rule-text
registry_calibration.go	21	.agents/skills/bench-craft-line/references/calibration-score.md	RequireInSection	How one claim scores	13	DROP	rule-text
registry_calibration.go	22	.agents/skills/bench-craft-line/references/calibration-score.md	RequireInSection	Who writes a label	18	DROP	rule-text
registry_calibration.go	23	.agents/skills/bench-craft-line/references/calibration-score.md	RequireInSection	Who writes a label	19	DROP	rule-text
registry_calibration.go	24	.agents/skills/bench-craft-line/references/calibration-score.md	RequireInSection	Who writes a label	20	DROP	rule-text
registry_calibration.go	25	.agents/skills/bench-craft-line/SKILL.md	RequireInSection	The declaration	107	DROP	rule-text
registry_calibration.go	26	.agents/commands/bench-final-check.md	RequireInSection	Capture the implementation retro	115	DROP	rule-text
registry_calibration.go	27	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_calibration.go	28	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_calibration.go	29	.agents/commands/bench-final-check.md	RequireInSection	Capture the implementation retro	116	DROP	rule-text
registry_calibration.go	30	capture/agent-performance/README.md	RequireInSection	Measures	48	DROP	rule-text
registry_calibration.go	31	capture/agent-performance/README.md	RequireInSection	Update contract	33	DROP	rule-text
registry_calibration.go	32	capture/agent-performance/README.md	RequireInSection	Update contract	32	DROP	rule-text
registry_chunk_chain.go	30	.agents/commands/bench-implement-spec.md	Require	Declare the line, validate the tickets, route the venue	19	DROP	rule-text
registry_chunk_chain.go	31	.agents/commands/bench-implement-spec.md	RequireInSection	Declare the line, validate the tickets, route the venue	21	DROP	rule-text
registry_chunk_chain.go	32	.agents/skills/bench-implement-spec/references/staleness-pass.md	RequireInSection	The charge	10	DROP	rule-text
registry_chunk_chain.go	33	.agents/commands/bench-implement-spec.md	RequireInSection	Build	52	GRAY	sentence-command-with-args
registry_chunk_chain.go	34	.agents/commands/bench-implement-spec.md	RequireInSection	Build	54	DROP	rule-text
registry_chunk_chain.go	35	.agents/commands/bench-implement-spec.md	RequireInSection	Build	54	DROP	rule-text
registry_chunk_chain.go	36	.agents/commands/bench-implement-spec.md	ForbidInSection	Build	0	FORBID	kind
registry_chunk_chain.go	37	.agents/commands/bench-implement-spec.md	RequireInSection	Build	56	DROP	rule-text
registry_chunk_chain.go	38	.agents/commands/bench-implement-spec.md	RequireInSection	Build	56	DROP	rule-text
registry_chunk_chain.go	39	docs/adr/0023-each-ticket-gets-a-fresh-author.md	Require	Consequences	19	DROP	rule-text
registry_chunk_chain.go	40	docs/adr/0023-each-ticket-gets-a-fresh-author.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	41	.agents/commands/bench-implement-spec.md	RequireInSection	Land	60	DROP	rule-text
registry_chunk_chain.go	42	.agents/commands/bench-implement-spec.md	RequireInSection	Land	62	GRAY	sentence-command-with-args
registry_chunk_chain.go	43	.agents/commands/bench-implement-spec.md	RequireInSection	Land	62	DROP	rule-text
registry_chunk_chain.go	44	.agents/commands/bench-implement-spec.md	RequireInSection	Land	62	DROP	sentence-command-noargs
registry_chunk_chain.go	45	.agents/commands/bench-implement-spec.md	RequireInSection	Land	62	DROP	rule-text
registry_chunk_chain.go	46	.agents/commands/bench-implement-spec.md	RequireInSection	Land	62	DROP	rule-text
registry_chunk_chain.go	47	.agents/commands/bench-implement-spec.md	RequireInSection	Land	62	DROP	rule-text
registry_chunk_chain.go	48	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	49	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	50	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	51	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	52	.agents/skills/bench-craft-delegate/SKILL.md	Require	The charge	68	DROP	sentence-with-route
registry_chunk_chain.go	53	.agents/skills/bench-craft-tickets/SKILL.md	ForbidCaseFoldedEmphasis		0	FORBID	kind
registry_chunk_chain.go	54	.bench/BENCH-reference.md	ForbidCaseFoldedEmphasis		0	FORBID	kind
registry_chunk_chain.go	55	.agents/skills/bench-craft-tickets/SKILL.md	Require	Land the frontier	97	DROP	rule-text
registry_chunk_chain.go	56	.bench/BENCH-reference.md	Require	Command Notes	339	DROP	rule-text
registry_chunk_chain.go	57	.agents/commands/bench-implement-spec.md	RequireInSection	Land	64	DROP	sentence-command-noargs
registry_chunk_chain.go	58	.agents/commands/bench-implement-spec.md	RequireInSection	Build	52	DROP	rule-text
registry_chunk_chain.go	59	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	60	.agents/commands/bench-implement-spec.md	RequireInSection	Land	64	GRAY	sentence-command-with-args
registry_chunk_chain.go	61	.agents/commands/bench-implement-spec.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	62	.agents/commands/bench-implement-spec.md	RequireInSection	Land	68	DROP	sentence-with-route
registry_chunk_chain.go	63	.agents/commands/bench-implement-spec.md	RequireInSection	Land	68	DROP	rule-text
registry_chunk_chain.go	64	.agents/commands/bench-implement-spec.md	RequireInSection	Land	72	DROP	rule-text
registry_chunk_chain.go	65	.agents/commands/bench-implement-spec.md	RequireInSection	Land	72	DROP	rule-text
registry_chunk_chain.go	66	.agents/commands/bench-implement-spec.md	RequireInSection	Land	72	DROP	rule-text
registry_chunk_chain.go	67	docs/adr/0023-each-ticket-gets-a-fresh-author.md	Require	Consequences	15	DROP	rule-text
registry_chunk_chain.go	68	README.md	Require	Reviewer quick start	82	DROP	rule-text
registry_chunk_chain.go	69	README.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	70	README.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	71	.agents/commands/bench-implement-spec.md	RequireInSection	Land	64	DROP	sentence-command-noargs
registry_chunk_chain.go	72	.agents/commands/bench-implement-spec.md	RequireInSection	`--full <spec>`	85	DROP	sentence-with-route
registry_chunk_chain.go	73	.agents/commands/bench-implement-spec.md	RequireInSection	Land	70	DROP	sentence-command-noargs
registry_chunk_chain.go	74	.agents/commands/bench-review-implementation.md	RequireInStep	Process	98	DROP	rule-text
registry_chunk_chain.go	75	.agents/commands/bench-review-implementation.md	RequireInStep	Process	100	DROP	sentence-command-noargs
registry_chunk_chain.go	76	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	77	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	78	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	47	DROP	rule-text
registry_chunk_chain.go	79	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	48	DROP	rule-text
registry_chunk_chain.go	80	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	81	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	51	DROP	rule-text
registry_chunk_chain.go	82	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	51	DROP	rule-text
registry_chunk_chain.go	83	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	51	DROP	rule-text
registry_chunk_chain.go	84	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Require	In a review round	216	DROP	rule-text
registry_chunk_chain.go	85	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	86	.agents/commands/bench-review-implementation.md	RequireInStep	Process	122	DROP	rule-text
registry_chunk_chain.go	87	.agents/commands/bench-review-implementation.md	RequireInStep	Process	137	DROP	rule-text
registry_chunk_chain.go	88	.agents/commands/bench-review-implementation.md	RequireInStep	Process	238	DROP	rule-text
registry_chunk_chain.go	89	.agents/commands/bench-review-implementation.md	RequireInStep	Process	215	DROP	sentence-command-noargs
registry_chunk_chain.go	90	.agents/commands/bench-review-implementation.md	RequireInStep	Process	216	DROP	sentence-command-noargs
registry_chunk_chain.go	91	.agents/commands/bench-review-implementation.md	RequireInStep	Process	230	DROP	sentence-command-noargs
registry_chunk_chain.go	92	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	93	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	94	.agents/commands/bench-review-implementation.md	Forbid		0	FORBID	kind
registry_chunk_chain.go	95	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	38	DROP	rule-text
registry_chunk_chain.go	96	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	55	DROP	sentence-command-noargs
registry_chunk_chain.go	97	.agents/skills/bench-craft-line/references/bounded-repair-policy.md	RequireInSection	Classification and completion	56	DROP	rule-text
registry_chunk_chain.go	98	.agents/commands/bench-review-implementation.md	RequireInSection	Review modes	49	DROP	sentence-command-noargs
registry_review_rules.go	5	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What a citation points at	19	DROP	rule-text
registry_review_rules.go	6	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	35	DROP	rule-text
registry_review_rules.go	7	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	Where an axis under-reads	36	DROP	rule-text
registry_review_rules.go	8	.agents/skills/bench-craft-review/references/finding-discipline.md	RequireInSection	What an axis return carries	50	DROP	rule-text
registry_charge_binding.go	7	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	56	DROP	sentence-command-noargs
registry_charge_binding.go	8	.agents/skills/bench-craft-delegate/SKILL.md	RequireInSection	The charge	56	DROP	sentence-with-route
registry_charge_binding.go	9	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	92	DROP	sentence-with-label
registry_charge_binding.go	10	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	93	DROP	rule-text
registry_charge_binding.go	11	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	42	DROP	rule-text
registry_charge_binding.go	12	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	42	DROP	sentence-command-noargs
registry_charge_binding.go	13	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	79	DROP	rule-text
registry_charge_binding.go	14	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_charge_binding.go	15	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	94	DROP	sentence-with-label
registry_charge_binding.go	16	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	95	DROP	rule-text
registry_charge_binding.go	17	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	96	DROP	rule-text
registry_charge_binding.go	18	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	96	DROP	rule-text
registry_charge_binding.go	19	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	27	LABEL	override
registry_charge_binding.go	20	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Repair-charge template	28	LABEL	override
registry_charge_binding.go	21	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Probes	138	DROP	rule-text
registry_charge_binding.go	22	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Claim schema	152	DROP	rule-text
registry_charge_binding.go	23	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Retry stops and aggregate readiness	181	DROP	rule-text
registry_charge_binding.go	24	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Retry stops and aggregate readiness	182	DROP	rule-text
registry_charge_binding.go	25	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	Forbid		0	FORBID	kind
registry_charge_binding.go	26	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Before the landing	208	DROP	sentence-command-noargs
registry_charge_binding.go	27	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Before the landing	209	DROP	rule-text
registry_charge_binding.go	28	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	Before the landing	210	DROP	rule-text
registry_charge_binding.go	29	.agents/skills/bench-craft-delegate/references/charge-evidence-format.md	RequireInSection	Consumers	101	DROP	rule-text
registry_charge_binding.go	30	.agents/skills/bench-craft-delegate/references/charge-evidence-format.md	RequireInSection	Consumers	103	DROP	rule-text
registry_charge_binding.go	31	.agents/skills/bench-craft-delegate/references/charge-evidence-format.md	RequireInSection	Consumers	104	DROP	rule-text
registry_commitment.go	17	.bench/BENCH.md	RequireInSection	Workflow	140	DROP	bench-workflow
registry_commitment.go	18	.bench/BENCH.md	RequireInSection	Workflow	140	DROP	bench-workflow
registry_commitment.go	19	.bench/BENCH.md	RequireInSection	Workflow	140	DROP	bench-workflow
registry_commitment.go	20	.bench/BENCH.md	RequireInSection	Workflow	140	DROP	bench-workflow
registry_commitment.go	21	.bench/BENCH.md	RequireInSection	Workflow	140	DROP	bench-workflow
registry_commitment.go	22	.bench/BENCH.md	RequireInSection	Workflow	140	DROP	bench-workflow
registry_commitment.go	23	.bench/BENCH.md	RequireInSection	Workflow	146	DROP	bench-workflow
registry_commitment.go	24	.bench/BENCH.md	RequireInSection	Workflow	146	DROP	bench-workflow
registry_commitment.go	25	.bench/BENCH.md	RequireInSection	Workflow	146	DROP	bench-workflow
registry_commitment.go	26	.bench/BENCH.md	RequireInSection	Workflow	148	DROP	bench-workflow
registry_commitment.go	27	.bench/BENCH.md	RequireInSection	Workflow	148	DROP	bench-workflow
registry_commitment.go	28	.bench/BENCH.md	RequireInSection	Workflow	148	DROP	bench-workflow
registry_commitment.go	29	.bench/BENCH.md	RequireInSection	Workflow	142	DROP	bench-workflow
registry_commitment.go	30	.bench/BENCH.md	RequireInSection	Workflow	144	DROP	bench-workflow
registry_commitment.go	31	.bench/BENCH.md	RequireInSection	Workflow	144	DROP	bench-workflow
registry_commitment.go	32	.agents/commands/bench-drain.md	Require	Exit handoff	85	DROP	rule-text
registry_commitment.go	33	.agents/commands/bench-drain.md	Require	5. Verdict the journal	215	DROP	rule-text
registry_commitment.go	34	.agents/commands/bench-drain.md	Require	5. Verdict the journal	217	DROP	rule-text
registry_commitment.go	35	.agents/commands/bench-drain.md	Require	7. Leave the sequence to the commitment	259	DROP	rule-text
registry_commitment.go	36	.agents/commands/bench-drain.md	Require	6. Classify every run; restructure on request	237	DROP	rule-text
registry_commitment.go	37	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	312	DROP	rule-text
registry_commitment.go	38	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	321	DROP	rule-text
registry_commitment.go	39	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	324	DROP	rule-text
registry_commitment.go	40	.agents/commands/bench-drain.md	Require	8. Batch-propose, then commit once on green	333	DROP	sentence-with-route
registry_commitment.go	41	.agents/commands/bench-final-check.md	RequireInSection	Exit handoff	87	DROP	rule-text
registry_commitment.go	42	.agents/commands/bench-implement-spec.md	Require	Entry orientation	8	DROP	sentence-command-noargs
registry_commitment.go	43	.agents/commands/bench-implement-spec.md	Require	Entry orientation	11	DROP	rule-text
registry_commitment.go	44	.agents/commands/bench-write-spec.md	Require	Exit handoff	14	GRAY	sentence-command-with-args
registry_commitment.go	45	.agents/commands/bench-write-spec.md	RequireInSection	Exit handoff	16	DROP	rule-text
registry_commitment.go	46	.agents/commands/bench-write-spec.md	RequireInSection	Exit handoff	18	DROP	sentence-command-noargs
registry_commitment.go	47	.agents/commands/bench-write-spec.md	RequireInSection	Exit handoff	22	GRAY	sentence-command-with-args
registry_commitment.go	48	.agents/commands/bench-debug.md	Require	Entry orientation	19	DROP	rule-text
registry_commitment.go	49	.agents/commands/bench.md	Require	Entry orientation	13	DROP	rule-text
registry_commitment.go	50	.agents/commands/bench-setup-repo.md	Require	Exit handoff	25	DROP	rule-text
registry_commitment.go	51	DATA_HANDLING.md	Require	Delivery commitment records	228	DROP	rule-text
registry_commitment.go	52	DATA_HANDLING.md	Require	Delivery commitment records	224	DROP	rule-text
registry_commitment.go	54	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	55	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	56	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	57	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	58	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	59	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	60	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	61	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	62	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	63	.agents/commands/bench-drain.md	Forbid		0	FORBID	kind
registry_commitment.go	64	.agents/commands/bench-final-check.md	Forbid		0	FORBID	kind
registry_commitment.go	65	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_commitment.go	66	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_commitment.go	67	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_commitment.go	68	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_commitment.go	69	.bench/BENCH.md	Forbid		0	FORBID	kind
registry_commitment.go	70	.agents/commands/bench.md	Forbid		0	FORBID	kind
registry_spec_trace.go	6	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	37	DROP	rule-text
registry_spec_trace.go	7	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	38	DROP	rule-text
registry_spec_trace.go	8	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	39	DROP	rule-text
registry_spec_trace.go	9	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	39	DROP	rule-text
registry_spec_trace.go	10	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	40	DROP	rule-text
registry_spec_trace.go	11	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	40	DROP	rule-text
registry_spec_trace.go	12	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	41	DROP	rule-text
registry_spec_trace.go	13	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	43	DROP	sentence-command-noargs
registry_spec_trace.go	14	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Before the map locks	42	DROP	rule-text
registry_check_floor.go	7	.agents/commands/bench-implement-spec.md	RequireInSection	Build	48	DROP	sentence-command-noargs
registry_check_floor.go	8	.agents/commands/bench-implement-spec.md	RequireInSection	Build	52	DROP	sentence-command-noargs
registry_check_floor.go	9	.agents/skills/bench-craft-spec/references/map-discipline.md	RequireInSection	Per row	126	DROP	rule-text
registry_check_floor.go	10	.agents/skills/bench-craft-delegate/references/delegation-discipline.md	RequireInSection	In the charge	35	GRAY	sentence-command-with-args
```
