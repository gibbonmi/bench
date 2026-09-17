# Calibrated decisions: research

Consumed by: `decisions/calibrated-decisions.md` and the spec that follows it.
Drift: a new primary source on process-level calibration, or a change to the scorecard contract in `capture/agent-performance/README.md`.
Retire when: the calibrated-decisions spec ships or the map closes without a spec.

## Recommendation, scope, and evidence status

Recommendation: adopt the proper-scoring-rule pattern from RLCR and the abstention reward from behavioral calibration as a process measure. Record a stated confidence with each claim that Bench already labels later. Score the pair with the Brier rule. Feed the score into the scorecard routing loop. Reject the Jev evaluation method and the Jev type-safety claim as Bench inputs.

Scope: the TypeSafe article on Jev, and four arXiv papers on reinforcement learning for calibration. Bench cannot train weights. The transfer target is the process reward that steers the routing tiers in `.bench/lines.env`, the guidance prose in the craft skills, and the scorecard measures.

Evidence status: every formula below comes from a paper's HTML full text, retrieved 2026-09-16. The Jev claims come from the article's own text, retrieved 2026-09-16. No claim was tested in this repository.

## Facts

### Q1. What does the TypeSafe article establish, claim by claim?

Source: https://typesafe.ai/blog/introducing-system-one-models-and-jev (retrieved 2026-09-16).

| claim | article text | status |
| --- | --- | --- |
| Training method | "Reinforcement Learning for Calibrated Decisions (RLCD)" | Named only. No reward, no measure, no procedure. The FAQ defers "Why was a new training algorithm needed?" |
| Optimization target | "answers with epistemically honest probabilities on System One tasks" | Matches the published proper-scoring-rule family. |
| Confidence | "Always communicates confidence and uncertainty with every output. Calibrated: higher confidence means higher accuracy" | The definition of calibration, not evidence for it. No calibration error is reported. |
| Output type | "Possible outputs and structure are defined in advance. The model never makes type errors" | A schema constraint on a classifier. It does not apply to prose output. |
| Hallucination rate | "Our number is not empirical. Schema matching is guaranteed, thus we can confidently add 0% into plots" | Stated by construction. The article says the comparison numbers "almost certainly" carry bias. |
| Sampling | "Parallel. Generates all outputs in a single query" | A classifier property. Cardinality is capped: "Jev supports a cardinality up to 255", and higher cardinality uses "a 2 stage-system of scoring independently then explicit choice". |
| Speed | "70ms-500ms", "40x-200x faster" | Self-measured: "our published evals are generally run from our laptops". |
| Cost | "$0.042 / MTok" input, output "FREE" | Vendor price. Not a technique. |
| Evaluation | "We assume there is a correct compute graph" and grade "against the average of smartest models (in this case, Astra and Fable)" | A preference label from reference models, not ground truth. The authors state "some bias could exist" and that the demos are "on the higher end of real world gains". |
| Use cases | "smart if-statements", "Verify everything. Score, judge, verify, guardrail" | The fit is a fixed-label decision inside a workflow, not open prose. |

### Q2. What reward makes a model both accurate and calibrated?

Source: RLCR, https://arxiv.org/html/2507.16806v2, section 3 "Method" (retrieved 2026-09-16).

The reward adds a Brier term to the binary correctness reward:

```
R(y, q, y*) = 1[y = y*] - (q - 1[y = y*])^2
```

Here `y` is the answer, `q` in [0, 1] is the stated confidence, and `y*` is the ground truth. Theorem 1 holds for any bounded proper scoring rule with `S(p,1) - S(p,0) < lambda`. Under that rule, the expected reward is maximal when `q` equals the true success probability, and among calibrated predictions the most likely answer scores highest. Log-loss fails the bound condition because it is unbounded.

The output format is fixed: `<think>`, `<answer>`, `<analysis>`, then `<confidence>` with "a number between 0 and 1 (inclusive)". The confidence comes after a separate uncertainty analysis, not inside the answer.

Metrics: ECE with 10 bins, Brier score, and AUROC. Result: calibration improves with no accuracy loss, in and out of domain. Remaining failure: "out-of-domain calibration error is often high", and a model "may still assign high confidence to multiple contradictory answers".

Test-time use (section 4.4): selection by highest confidence, and a confidence-weighted majority vote that beats a plain majority vote across 77 datasets.

### Q3. How does a reward handle abstention?

Source: Behaviorally calibrated RL, https://arxiv.org/html/2512.19920v1, sections 3.1 to 3.3 and 4.4 (retrieved 2026-09-16).

Behavioral calibration trains one policy to answer or abstain under a risk threshold `t`. The explicit reward is:

```
R(a, y, t) = +1            if answer and valid
           =  0            if abstain
           = -t / (1 - t)  if answer and invalid
```

The explicit-threshold form failed in practice. The paper reports: "The refusal rate and hallucination rate proved insensitive to the specific prompt input t". The model also "over-reject[s] even at t=0". The fix integrates the reward over a prior `u(t)` on the threshold. That integral yields a strictly proper scoring rule on a stated confidence `p`. With a uniform prior the rule is the Brier form `2 p valid(y) - p^2`.

The output format uses an `<IDK>` token for a whole-response abstention and per-claim confidence marks for a long response. Metrics: smoothed ECE, confidence AUC, a signal-to-noise gain across the risk spectrum, and abstention accuracy.

Transfer: a 4B model trained on mathematics alone reached frontier-level calibration on SimpleQA with "extremely low prediction accuracy". The authors call calibration "a transferable 'meta-skill'" independent of accuracy.

### Q4. What happens to calibration under verifiable-reward RL alone?

Source: Calibration-aware RL, https://arxiv.org/html/2601.13284v1, sections 3.3, "Diagnosing Calibration Failures in RLVR", and "Proposed Method" (retrieved 2026-09-16).

A decision task presents a question and a finite option set. Confidence is the probability of the decision token itself. Finding 2: "Nearly all trajectories from a base model yield overconfident decision token probabilities... there are no calibrated rollouts to reinforce." Finding 3: the decision token "only conveys the conclusion in the reasoning trace, not the uncertainty".

The numbers below come from CommonsenseQA with Qwen3-1.7B. The base model scores 67.49% accuracy and 29.91 ECE. SFT scores 68.55% and 7.36 ECE. GRPO scores 73.67% and 24.39 ECE. Plain RLVR raises accuracy and keeps the overconfidence.

The fix adds a cross-entropy term on the decision token, with weight `lambda = 0.001`. The target is one-hot when correct and uniform when incorrect. The policy advantage is zero at the decision step. The decision probability is constrained to `[1/|C|, 1]`. ECE drops by up to 9 points with the accuracy gain kept.

### Q5. What does a log-scoring variant add?

Source: Rewarding Doubt, https://arxiv.org/html/2503.02623v6 (retrieved 2026-09-16).

The reward is `log(p)` when correct and `log(1 - p)` when incorrect. The confidence is "an integer between 0 and 10", normalized for the reward. The answer is generated first and frozen. The confidence is generated in a separate step, so "answer correctness is not affected by our confidence calibration training".

Calibration transfers from TriviaQA to CommonsenseQA and MedQA. Accuracy stays stable. The paper reports no confidence collapse. Note the contradiction with RLCR: the log rule is unbounded, and RLCR's theorem excludes it.

## Inferences

- Bench already labels four claim types. The gate labels an acceptance row, and the coordinator labels a delegate done-claim. The tree labels a review finding, and the retro labels the expected repair rounds. What is absent is the stated confidence before the label. Without the pair there is no calibration measure, only an accuracy measure (`capture/agent-performance/README.md`, "Measures" table).
- The Bench "done" claim is the binary reward that RLCR criticizes. Nothing penalizes a confident wrong claim more than a hesitant wrong claim. The scorecard rows record only the correction count.
- Bench already rewards abstention in prose. The Sonnet/high scorecard row praises a routed shortfall (`capture/agent-performance/claude-models.md`). No measure counts it, so the routing loop cannot learn from it.
- The Jev evaluation label, an average over the strongest models, is a preference signal. Invariant 1 in `.bench/BENCH.md` makes the gate the oracle and forbids that label source.
- The Jev type-safety point transfers as a schema rule, not as a guarantee. A fixed enum on each claim status, plus a bounded integer confidence, makes the pairs aggregable. Prose claims cannot be aggregated.
- The 0-to-10 integer scale from Rewarding Doubt fits an agent better than a real number. An agent that writes 0.87 states false precision.
- The behavioral-calibration failure mode transfers: an explicit per-task risk threshold did not steer the policy. The Brier form with no threshold did. Bench has no threshold parameter today, so the Brier form is the smaller change.
- The RLCR out-of-domain caveat transfers: a delegate calibrated on Go tickets may not be calibrated on guidance prose. The scorecard's per-role rows already separate these.

## Contradictions and unknowns

- RLCR excludes log-loss by theorem; Rewarding Doubt uses log-loss and reports stability. Resolution for Bench: use the bounded Brier rule, which both papers accept.
- The RLCD name has no arXiv paper and no published method. Unknown.
- No paper studies calibration of an agent's process claims, such as "this ticket is done". All four study answer correctness. The transfer to a process reward is an inference, not a tested result.
- How many pairs a model/effort/role needs before a Brier mean steers routing is unknown. ECE with 10 bins needs tens of samples per bin. The scorecard caps aggregation at 10 assignments, which is below that.

## Proposals

1. Add a stated confidence, an integer from 0 to 10, to each claim on a decision surface that Bench already labels.
2. Fix the claim status enum in advance: `verified`, `claimed`, `abstained`.
3. Score each labeled pair with the Brier rule on the normalized confidence. An abstention scores zero and leaves the Brier aggregate.
4. Record the pair in the retro and the Brier mean and count in the scorecard, per model, effort, and role.
5. Keep the gate and the coordinator as the only label sources.
6. Keep routing changes under the two-run rule until enough pairs exist for a threshold.

## Validation plan

- Dogfood the pairs on one landing before any verb is built. Count how many pairs one landing yields per role.
- Make sure the enum and integer parse from a delegate return without prose repair.
- After ten pairs for one model/effort/role, compute the Brier mean by hand and compare it with the first-pass measure. If the two agree on every routing decision, the measure adds no information.

## Verification record

- [x] Opens with the recommendation, the scope, and the evidence status.
- [x] Facts, inferences, tested results, and proposals in separate sections. No tested results exist.
- [x] One section per question.
- [x] A claim table for the article comparison.
- [x] Contradictions and unknowns kept.
- [ ] No diagram. The relations are tabular.
- [x] States what could not be verified.
- [x] Ends with a validation plan.
- [x] Every load-bearing claim cites a URL with its retrieval date or a local path.
