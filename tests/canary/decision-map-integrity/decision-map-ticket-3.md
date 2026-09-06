# How does the resolver/hook learn the current harness?

Blocked by: none
Type: Grill

### Question

The harness is implied by the output-shape flags today. Keep that, or make the
harness explicit?

### Answer

An explicit `--harness` flag. Each adapter passes the harness name, and the
output shape follows from that harness column. The old shape flags retire.
