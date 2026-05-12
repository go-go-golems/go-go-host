#!/usr/bin/env python3
import json
import os
import shutil
import sys


def emit(obj):
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()


emit({
    "type": "handshake",
    "protocol_version": "v2",
    "plugin_name": "go-go-host",
    "capabilities": {"ops": ["config.mutate", "validate.run", "launch.plan"]},
})


def repo_root(req):
    return req.get("ctx", {}).get("repo_root") or os.getcwd()


for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        req = json.loads(line)
    except Exception as e:
        print(f"invalid request: {e}", file=sys.stderr)
        continue
    rid = req.get("request_id", "")
    op = req.get("op", "")
    root = repo_root(req)
    web_dir = os.path.join(root, "web", "admin")

    if op == "config.mutate":
        emit({
            "type": "response",
            "request_id": rid,
            "ok": True,
            "output": {"config_patch": {"set": {
                "services.storybook.port": 6007,
                "services.storybook.url": "http://127.0.0.1:6007",
                "services.web.port": 5173,
                "services.web.url": "http://127.0.0.1:5173",
            }, "unset": []}},
        })
    elif op == "validate.run":
        errors = []
        warnings = []
        if shutil.which("pnpm") is None:
            errors.append({"message": "pnpm is required to run the dashboard and Storybook"})
        if not os.path.exists(os.path.join(web_dir, "package.json")):
            errors.append({"message": "web/admin/package.json is missing"})
        if not os.path.isdir(os.path.join(web_dir, "node_modules")):
            warnings.append({"message": "web/admin/node_modules is missing; run `make web-install` or `cd web/admin && pnpm install`"})
        emit({"type": "response", "request_id": rid, "ok": True, "output": {"valid": len(errors) == 0, "errors": errors, "warnings": warnings}})
    elif op == "launch.plan":
        emit({
            "type": "response",
            "request_id": rid,
            "ok": True,
            "output": {"services": [
                {
                    "name": "storybook",
                    "command": ["bash", "--noprofile", "--norc", "-lc", "cd web/admin && exec pnpm storybook"],
                    "env": {"BROWSER": "none"},
                    "health": {"type": "http", "url": "http://127.0.0.1:6007", "interval_ms": 1000, "timeout_ms": 1000},
                },
                {
                    "name": "web-admin",
                    "command": ["bash", "--noprofile", "--norc", "-lc", "cd web/admin && exec pnpm dev"],
                    "env": {"BROWSER": "none", "VITE_GO_GO_HOST_API_TARGET": "http://127.0.0.1:8080"},
                    "health": {"type": "http", "url": "http://127.0.0.1:5173", "interval_ms": 1000, "timeout_ms": 1000},
                },
            ]},
        })
    else:
        emit({"type": "response", "request_id": rid, "ok": False, "error": {"code": "E_UNSUPPORTED", "message": f"unsupported op: {op}"}})
