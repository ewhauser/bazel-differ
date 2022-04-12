dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
. "${dir}"/test_runner.sh
. "${dir}"/test_helper.sh

runner=$(get_test_runner "${1:-local}")

test_single_file_change() {
  run_compat_test de48add f0d09c4
}

test_workspace_change() {
  run_compat_test 50efa15 60b40a3
}

test_multiple_file_change() {
  run_compat_test b013593 95b4f75
}

test_get_targets_single_file_change() {
  run_get_targets de48add f0d09c4 "${dir}"/snapshots/single_file_change_targets.txt
}

test_get_targets_workspace_change() {
  run_get_targets 50efa15 60b40a3 "${dir}"/snapshots/workspace_change_targets.txt
}

test_get_targets_multiple_file_change() {
  run_get_targets b013593 95b4f75 "${dir}"/snapshots/multiple_file_targets.txt
}

$runner test_single_file_change
$runner test_workspace_change
$runner test_multiple_file_change
$runner test_get_targets_single_file_change
$runner test_get_targets_workspace_change
$runner test_get_targets_multiple_file_change
