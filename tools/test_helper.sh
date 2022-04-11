#!/usr/bin/env
set -eo pipefail

function run_compat_test() {
    local start_commit="$1"
    local end_commit="$2"

    tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
    repo_path="$tmp_dir/bazel-remote"
    bazel_diff="$(bazel info bazel-bin)/tools/bazel-diff/bazel-diff"
    bazel_differ="$(bazel info bazel-bin)/cli/bazel-differ_/bazel-differ"

    git clone https://github.com/buchgr/bazel-remote.git "$repo_path"

    common_opts="-w $repo_path -b $(which bazelisk)"

    pushd "$repo_path"
    git checkout "$start_commit"

    "$bazel_diff" generate-hashes $common_opts $tmp_dir/bazel-diff-starting.txt
    "$bazel_differ" generate-hashes $common_opts $tmp_dir/bazel-differ-starting.txt

    git checkout "$end_commit"

    "$bazel_diff" generate-hashes $common_opts $tmp_dir/bazel-diff-ending.txt
    "$bazel_differ" generate-hashes $common_opts $tmp_dir/bazel-differ-ending.txt

    "$bazel_diff" $common_opts -sh $tmp_dir/bazel-diff-starting.txt -fh $tmp_dir/bazel-diff-ending.txt -o $tmp_dir/bazel-diff-targets.txt
    "$bazel_differ" diff $common_opts -s $tmp_dir/bazel-differ-starting.txt -f $tmp_dir/bazel-differ-ending.txt -o $tmp_dir/bazel-differ-targets.txt

    cat $tmp_dir/bazel-diff-targets.txt | sort > $tmp_dir/bazel-diff-targets-sorted.txt
    cat $tmp_dir/bazel-differ-targets.txt | sort > $tmp_dir/bazel-differ-targets-sorted.txt

    diff $tmp_dir/bazel-diff-targets-sorted.txt $tmp_dir/bazel-differ-targets-sorted.txt

    popd

    rm -rf "$tmp_dir"
}
