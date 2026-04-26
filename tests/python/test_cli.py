import json
import os
import subprocess
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[2]
CMD = ["go", "run", "./src/cmd/gotree"]


def run_cmd(args, env=None):
    full_env = os.environ.copy()
    if env:
        full_env.update(env)
    proc = subprocess.run(
        CMD + args,
        cwd=REPO_ROOT,
        env=full_env,
        capture_output=True,
        text=False,
        check=False,
    )
    stdout = proc.stdout.decode("utf-8", errors="replace")
    stderr = proc.stderr.decode("utf-8", errors="replace")
    proc.stdout = stdout
    proc.stderr = stderr
    return proc


def make_fixture(tmp_path: Path) -> Path:
    root = tmp_path / "fixture"
    (root / "nested").mkdir(parents=True)
    (root / "nested" / "a.go").write_text("package main\n", encoding="utf-8")
    (root / "nested" / "b.txt").write_text("skip\n", encoding="utf-8")
    (root / "z.md").write_text("# z\n", encoding="utf-8")
    return root


def test_version_reports_v2():
    proc = run_cmd(["--version"])
    assert proc.returncode == 0
    assert "gotree version 2.0.0" in proc.stdout


def test_json_output_is_valid(tmp_path):
    root = make_fixture(tmp_path)
    proc = run_cmd(["--json", str(root)])
    assert proc.returncode == 0, proc.stderr
    payload = json.loads(proc.stdout)
    assert "root" in payload
    assert "summary" in payload
    assert payload["root"]["is_dir"] is True


def test_include_overrides_ignore(tmp_path):
    root = make_fixture(tmp_path)
    proc = run_cmd(["--json", "--include", "*.go", "--ignore", "*.go", str(root)])
    assert proc.returncode == 0, proc.stderr
    payload = json.loads(proc.stdout)
    names = collect_node_names(payload["root"])
    assert "a.go" in names
    assert "b.txt" not in names
    assert "z.md" not in names


def test_theme_mono_disables_color_sequences(tmp_path):
    root = make_fixture(tmp_path)
    proc = run_cmd(["--theme", "mono", "--icons", "none", str(root)])
    assert proc.returncode == 0, proc.stderr
    assert "\x1b[" not in proc.stdout


def test_completion_outputs_content():
    proc = run_cmd(["--completion", "powershell"])
    assert proc.returncode == 0, proc.stderr
    assert "Register-ArgumentCompleter" in proc.stdout


def collect_node_names(node):
    names = [node.get("name", "")]
    for child in node.get("children", []):
        names.extend(collect_node_names(child))
    return names
