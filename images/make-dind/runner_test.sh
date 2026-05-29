#!/usr/bin/env bash

# Copyright 2023 The cert-manager Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# End-to-end tests for `runner`. Each test invokes the runner script itself
# with a synthetic wrapped command and asserts on the observable behavior:
# exit code, output ordering, and signal handling.

set -o errexit
set -o nounset
set -o pipefail

# Enable job control so backgrounded children get default SIGINT/SIGQUIT
# handling. Without this, non-interactive bash sets SIGINT and SIGQUIT to
# SIG_IGN on any process backgrounded with `&` and that inheritance persists
# across exec — which would cause the SIGINT signal-forwarding test below to
# falsely fail because the wrapped `sleep` would ignore the forwarded signal.
# Prow execs the runner as a session leader, so the production code path is
# unaffected.
set -m

RUNNER="$(cd "$(dirname "$0")" && pwd)/runner"
[[ -x "$RUNNER" ]] || { echo "runner not found or not executable at $RUNNER" >&2; exit 1; }

PASS=0
FAIL=0

pass() { echo "PASS: $1"; PASS=$((PASS + 1)); }
fail() { echo "FAIL: $1"; FAIL=$((FAIL + 1)); }

# Invoke the runner with every `set` option reset to its default, so the
# test's own shell options (errexit, nounset, pipefail, monitor, …) cannot
# leak into the runner's execution and mask or amplify bugs.
#
# Must be called inside a subshell — either explicitly `( run_runner … )` or
# implicitly via `run_runner … &`. The function execs the runner, replacing
# the surrounding (sub)shell; calling it unwrapped would replace the test
# itself. The BASHPID guard catches that mistake.
#
# Wrapping the body in `( … )` would be safer, but bash would then fork an
# extra subshell for backgrounded calls and `$!` would point at that wrapper
# instead of the exec'd runner — breaking `kill "$!"` in the signal tests.
run_runner() {
    if [[ "$BASHPID" == "$$" ]]; then
        echo "run_runner must be called in a subshell or backgrounded" >&2
        exit 1
    fi
    while read -r opt _; do set +o "$opt"; done < <(set -o)
    exec "$RUNNER" "$@"
}

# Test 1: runner exits with the wrapped command's exit code (proves $! is the
# wrapped command, since `wait` returns its exit code).
test_exit_code_propagates() {
    local exit_value=0
    ( run_runner bash -c "exit 42" ) >/dev/null 2>&1 || exit_value=$?
    if [[ "$exit_value" -eq 42 ]]; then
        pass "runner exits with wrapped command's exit code (42)"
    else
        fail "Expected exit 42, got $exit_value"
    fi
}

# Test 2: stdout and stderr from the wrapped command appear in chronological
# order in the merged output.
test_log_ordering() {
    local out
    out=$(mktemp)
    ( run_runner bash -c '
        echo "line-1-stdout"
        echo "line-2-stderr" >&2
        echo "line-3-stdout"
        echo "line-4-stderr" >&2
    ' ) >"$out" 2>&1

    # Extract just the wrapped command's lines (runner adds its own banners).
    local wrapped_output
    wrapped_output=$(grep -E '^line-[0-9]+-(stdout|stderr)$' "$out" || true)
    local expected
    expected=$'line-1-stdout\nline-2-stderr\nline-3-stdout\nline-4-stderr'

    if [[ "$wrapped_output" == "$expected" ]]; then
        pass "stdout and stderr appear in chronological order"
    else
        fail "Output order mismatch. Got:"
        echo "$wrapped_output" >&2
    fi
    rm -f "$out"
}

# Verifies the runner streams the wrapped command's output line-by-line
# instead of letting libc's default block buffering hold output back until the
# process exits. awk block-buffers stdout when stdout is a pipe (not a TTY).
# The busy-wait deliberately avoids any operation that would flush the buffer
# as a side effect — system(), close(), pipes, or exit. Without the runner
# forcing line-buffered stdio on the wrapped command (via stdbuf -oL -eL),
# the first print sits in libc's buffer for the full duration of the busy-wait
# — well past our poll window.
test_output_is_streamed() {
    local out
    out=$(mktemp)

    ( run_runner awk 'BEGIN { print "early-line"; t = systime() + 3; while (systime() < t) {} print "late-line" }' ) >"$out" 2>&1 &
    local runner_pid=$!

    # Poll for up to 1.5s — comfortably under the 2-3s busy-wait — for early-line.
    local found=0 i
    for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15; do
        if grep -q "^early-line$" "$out" 2>/dev/null; then
            found=1
            break
        fi
        sleep 0.1
    done

    if [[ "$found" -eq 1 ]]; then
        pass "wrapped command output is streamed line-by-line"
    else
        fail "early-line not seen within 1.5s — output appears to be block-buffered"
    fi

    # Let the busy-wait finish so the runner exits cleanly.
    wait "$runner_pid" 2>/dev/null || true

    if [[ "$found" -eq 0 ]]; then
        echo "--- captured runner output ---" >&2
        cat "$out" >&2
    fi

    rm -f "$out"
}

# Asserts that `kill -s $1` to the runner actually kills the wrapped command,
# not just that the runner itself exits quickly. Capturing the wrapped PID up
# front and probing it after the runner is gone distinguishes a real kill
# from an orphan: if `set -m` is missing in the runner, the wrapped command
# inherits SIG_IGN for SIGINT and survives the forwarded signal — the runner
# still exits promptly (its `wait` is interrupted by the trap), so checking
# only the runner's exit time would silently pass the bug.
assert_signal_reaches_wrapped() {
    local signal=$1
    local start elapsed runner_pid wrapped_pid i
    start=$(date +%s)

    run_runner sleep 30 >/dev/null 2>&1 &
    runner_pid=$!

    # Give the runner time to launch the sleep and install its traps.
    sleep 0.5

    wrapped_pid=$(pgrep -P "$runner_pid" -x sleep || true)
    if [[ -z "$wrapped_pid" ]]; then
        fail "SIG$signal: could not locate wrapped sleep as child of runner $runner_pid"
        kill -KILL "$runner_pid" 2>/dev/null || true
        wait "$runner_pid" 2>/dev/null || true
        return
    fi

    kill -s "$signal" "$runner_pid"
    wait "$runner_pid" 2>/dev/null || true
    elapsed=$(( $(date +%s) - start ))

    # Wrapped death is asynchronous from the runner's exit, so poll briefly.
    for i in 1 2 3 4 5 6 7 8 9 10; do
        kill -0 "$wrapped_pid" 2>/dev/null || break
        sleep 0.1
    done

    if kill -0 "$wrapped_pid" 2>/dev/null; then
        fail "SIG$signal: wrapped sleep (pid $wrapped_pid) still alive after runner exit — orphaned, not killed"
        kill -KILL "$wrapped_pid" 2>/dev/null || true
    elif (( elapsed < 10 )); then
        pass "SIG$signal killed wrapped command (runner exited in ${elapsed}s, sleep dead)"
    else
        fail "SIG$signal not forwarded — runner waited ${elapsed}s for the sleep"
    fi
}

# Verifies the runner's stdout closes promptly after the wrapped command
# exits, even if a leaked descendant inherited the stdout pipe. Mirrors
# what Prow sees on its end of the runner's stdout: prow's `Wait()` blocks
# until every writer to the pipe is gone, so a leaked process (envtest's
# kube-apiserver, the watchdog's in-flight `sleep`, etc.) would keep the
# pipe open until Prow's 2h timeout fires. The pipeline-reader exits as
# soon as it sees EOF, so a fast exit here proves the runner reaped all
# writers before returning.
test_pipe_closes_after_leaked_grandchild() {
    local start elapsed
    start=$(date +%s)

    (
        run_runner bash -c '
            # nohup detaches from the controlling terminal; the spawned sleep
            # inherits the wrapped command'\''s stdout/stderr, which is the
            # mawk FIFO. Without the process-group reap, this sleep would
            # outlive the runner and hold the pipe open.
            nohup sleep 60 >&1 2>&1 &
            echo "leaked grandchild pid=$!"
        '
    ) | (
        while IFS= read -r _; do :; done
    )

    elapsed=$(( $(date +%s) - start ))
    if (( elapsed < 10 )); then
        pass "stdout pipe closes promptly after leaked grandchild (${elapsed}s)"
    else
        fail "stdout pipe still open ${elapsed}s after runner exit — a writer leaked"
    fi

    pkill -KILL -f "^sleep 60" 2>/dev/null || true
}

# Verifies the watchdog actually fires: it dumps the process tree to stderr
# after $WATCHDOG_STALL_SECONDS of no output from the wrapped command. We
# shrink the threshold to 2s and run a silent 5s sleep so the watchdog has
# time for one tick, notice the stall, and emit the banner before the
# wrapped command exits.
test_watchdog_fires_on_stall() {
    local out
    out=$(mktemp)

    ( WATCHDOG_STALL_SECONDS=2 run_runner sleep 5 ) >"$out" 2>&1

    if grep -q "## Watchdog: no output from wrapped command" "$out" \
       && grep -qE "(bash|sleep) " "$out"; then
        pass "watchdog fires and prints process tree when output stalls"
    else
        fail "watchdog test: did not see watchdog banner + process tree"
        echo "--- captured runner output ---" >&2
        tail -20 "$out" >&2
    fi

    rm -f "$out"
}

test_exit_code_propagates
test_log_ordering
test_output_is_streamed
assert_signal_reaches_wrapped TERM
assert_signal_reaches_wrapped INT
test_pipe_closes_after_leaked_grandchild
test_watchdog_fires_on_stall

echo ""
echo "Results: $PASS passed, $FAIL failed"
[[ "$FAIL" -eq 0 ]]
