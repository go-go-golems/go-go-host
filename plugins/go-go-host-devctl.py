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
                "services.postgres.port": 55432,
                "services.postgres.url": "postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable",
                "services.go-go-hostd.port": 8080,
                "services.go-go-hostd.url": "http://127.0.0.1:8080",
                "services.storybook.port": 6007,
                "services.storybook.url": "http://127.0.0.1:6007",
                "services.keycloak.port": 18080,
                "services.keycloak.url": "http://127.0.0.1:18080/realms/go-go-host",
                "services.web-admin.port": 5173,
                "services.web-admin.url": "http://127.0.0.1:5173",
                "services.web.port": 5173,
                "services.web.url": "http://127.0.0.1:5173",
            }, "unset": []}},
        })
    elif op == "validate.run":
        errors = []
        warnings = []
        if shutil.which("pnpm") is None:
            errors.append({"message": "pnpm is required to run the dashboard and Storybook"})
        if shutil.which("docker") is None:
            errors.append({"message": "docker is required to run the dev Postgres service"})
        if not os.path.exists(os.path.join(web_dir, "package.json")):
            errors.append({"message": "web/admin/package.json is missing"})
        if not os.path.isdir(os.path.join(web_dir, "node_modules")):
            warnings.append({"message": "web/admin/node_modules is missing; run `make web-install` or `cd web/admin && pnpm install`"})
        emit({"type": "response", "request_id": rid, "ok": True, "output": {"valid": len(errors) == 0, "errors": errors, "warnings": warnings}})
    elif op == "launch.plan":
        wait_for_pg = "python3 - <<'PY'\nimport socket, time\nfor _ in range(120):\n    try:\n        with socket.create_connection(('127.0.0.1', 55432), timeout=1):\n            raise SystemExit(0)\n    except OSError:\n        time.sleep(1)\nraise SystemExit('Postgres did not become reachable on 127.0.0.1:55432')\nPY"
        wait_for_keycloak = "python3 - <<'PY'\nimport json, time, urllib.request\nurl = 'http://127.0.0.1:18080/realms/go-go-host/.well-known/openid-configuration'\nfor _ in range(180):\n    try:\n        with urllib.request.urlopen(url, timeout=2) as r:\n            data = json.load(r)\n            if data.get('issuer'):\n                raise SystemExit(0)\n    except Exception:\n        time.sleep(1)\nraise SystemExit('Keycloak realm did not become reachable on 127.0.0.1:18080')\nPY"
        emit({
            "type": "response",
            "request_id": rid,
            "ok": True,
            "output": {"services": [
                {
                    "name": "postgres",
                    "command": ["bash", "--noprofile", "--norc", "-lc", "exec docker compose -f deployments/dev/docker-compose.yaml up postgres"],
                },
                {
                    "name": "keycloak",
                    "command": ["bash", "--noprofile", "--norc", "-lc", "exec docker compose -f deployments/dev/docker-compose.yaml up keycloak"],
                    "health": {"type": "http", "url": "http://127.0.0.1:18080/realms/go-go-host/.well-known/openid-configuration", "interval_ms": 1000, "timeout_ms": 2000},
                },
                {
                    "name": "go-go-hostd",
                    "command": ["bash", "--noprofile", "--norc", "-lc", wait_for_pg + "\n" + wait_for_keycloak + "\nexec go run ./cmd/go-go-hostd --config configs/dev.keycloak.yaml"],
                    "env": {"GO_GO_HOST_CONFIG": "configs/dev.keycloak.yaml"},
                    "health": {"type": "http", "url": "http://127.0.0.1:8080/healthz", "interval_ms": 1000, "timeout_ms": 1000},
                },
                {
                    "name": "web-admin",
                    "command": ["bash", "--noprofile", "--norc", "-lc", "cd web/admin && exec pnpm dev"],
                    "env": {"BROWSER": "none", "VITE_GO_GO_HOST_API_TARGET": "http://127.0.0.1:8080"},
                    "health": {"type": "http", "url": "http://127.0.0.1:5173", "interval_ms": 1000, "timeout_ms": 1000},
                },
                {
                    "name": "storybook",
                    "command": ["bash", "--noprofile", "--norc", "-lc", "cd web/admin && exec pnpm storybook"],
                    "env": {"BROWSER": "none"},
                    "health": {"type": "http", "url": "http://127.0.0.1:6007", "interval_ms": 1000, "timeout_ms": 1000},
                },
            ]},
        })
    else:
        emit({"type": "response", "request_id": rid, "ok": False, "error": {"code": "E_UNSUPPORTED", "message": f"unsupported op: {op}"}})
