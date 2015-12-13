#!/bin/sh

set -x
cat src/core/core.rs > src/core/mod.rs
cat src/core/gui.rs >> src/core/mod.rs
cat src/core/widgets.rs >> src/core/mod.rs
