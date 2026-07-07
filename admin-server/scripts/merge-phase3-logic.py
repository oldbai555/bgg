#!/usr/bin/env python3
"""Merge logic method bodies from .logic-backup into new nested logic paths."""
import os
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
BACKUP = ROOT / ".logic-backup"
LOGIC = ROOT / "internal" / "logic"

DOMAIN_DIRS = {"iam", "blog", "video", "chat", "sdk", "task", "monitoring", "system", "misc"}


def logic_type_name(content: str):
    m = re.search(r"type (\w+Logic) struct", content)
    return m.group(1) if m else None


def package_name(content: str):
    m = re.search(r"^package (\w+)$", content, re.M)
    return m.group(1) if m else None


def build_index():
    idx = {}
    for dirpath, _, files in os.walk(LOGIC):
        rel = Path(dirpath).relative_to(LOGIC)
        parts = rel.parts
        if not parts or parts[0] not in DOMAIN_DIRS:
            continue
        for fn in files:
            if not fn.endswith("_logic.go") and not fn.endswith("logic.go"):
                continue
            path = Path(dirpath) / fn
            content = path.read_text(encoding="utf-8")
            name = logic_type_name(content)
            if name:
                idx[name] = path
    return idx


def merge():
    index = build_index()
    merged = 0
    missing = []
    for dirpath, _, files in os.walk(BACKUP):
        for fn in files:
            if not fn.endswith(".go"):
                continue
            backup_path = Path(dirpath) / fn
            content = backup_path.read_text(encoding="utf-8")
            ltype = logic_type_name(content)
            if not ltype:
                continue
            target = index.get(ltype)
            if not target:
                missing.append((ltype, backup_path))
                continue
            new_pkg = package_name(target.read_text(encoding="utf-8"))
            old_pkg = package_name(content)
            if new_pkg != old_pkg:
                content = re.sub(rf"^package {old_pkg}$", f"package {new_pkg}", content, count=1, flags=re.M)
            # keep updated goctl header
            content = re.sub(
                r"// goctl [\d.]+",
                "// goctl 1.10.1",
                content,
                count=1,
            )
            target.write_text(content, encoding="utf-8")
            merged += 1
            print(f"merged {ltype} -> {target.relative_to(ROOT)}")
    print(f"merged {merged}, missing {len(missing)}")
    for ltype, p in missing[:20]:
        print(f"  MISSING {ltype} from {p.relative_to(ROOT)}")


def remove_flat_dirs():
    for base in [ROOT / "internal" / "handler", ROOT / "internal" / "logic"]:
        for name in os.listdir(base):
            path = base / name
            if not path.is_dir():
                continue
            if name in DOMAIN_DIRS:
                continue
            # remove flat legacy module dirs
            import shutil
            shutil.rmtree(path)
            print(f"removed {path.relative_to(ROOT)}")


if __name__ == "__main__":
    import sys
    merge()
    if "--remove-flat" in sys.argv:
        remove_flat_dirs()
