# Kentik Synthetics → Prometheus via the ktranslate Firehose

This guide covers running `ktranslate` as a **Kentik Data Firehose** collector that receives
**synthetic test results** (kflow with `app_protocol=10`) and exports them as **Prometheus**
metrics — either scraped (`/metrics`) or pushed via **remote-write**. It documents every
required flag (including the easy-to-miss `-udrs`, `-mapping`, and `-prom_seen`), and explains
how to interpret the new `outcome` metric.

---

## 1. How the data flows

```
Kentik Firehose ──HTTPS POST kflow──▶ /chf on -listen (8081)
                                         │  decode capnproto → JCHF
                                         │  app_protocol=10 → EventType = "KSynth"
                                         │  -udrs renames raw cols (INT00 → result_type, …)
                                         ▼
                              prometheus format (fromKSynth)
                                         │  emits kentik:synth:* incl. outcome
                                         ▼
                    ┌───────────────────────────────────────────┐
                    │ -format prometheus        → prom sink /metrics (scrape)   │
                    │ -format prometheus_remote → prom sink remote-write (push) │
                    └───────────────────────────────────────────┘
```

Synthetic results arrive over the **same kflow firehose** as flow. `ktranslate` classifies
records with `app_protocol=10` as `KSynth`, and the Prometheus format turns them into metrics.

---

## 2. Required command-line flags

| Flag | Purpose | Notes |
|------|---------|-------|
| `-listen 0.0.0.0:8081` | HTTP(S) listener that exposes the firehose endpoint `/chf`. | Use `off` to disable. |
| `-ssl_cert_file <pem>` / `-ssl_key_file <key>` | TLS for the `/chf` listener. | Firehose posts over HTTPS. The `/metrics` scrape endpoint stays **plain HTTP**. |
| `-kentik_email <email>` + `KENTIK_API_TOKEN` env | Kentik API auth, used to enrich records (device / **test name** / **agent name**). | Token is an **env var**, never a flag. |
| `-api_root <url>` | Kentik API base. | Default `https://api.kentik.com`; EU: `https://api.kentik.eu`. |
| `-mapping ./config/config.json` | Enum-value mapping (named columns / enums). | Improves label readability. |
| `-udrs ./config/udr.csv` | **Critical.** Maps raw synthetic columns to named fields. | Without it there is no `result_type`, and the outcome is wrong (see §5). |
| `-format prometheus` \| `-format prometheus_remote` | Output encoding. | Scrape vs. remote-write. |
| `-sinks prometheus` | The Prometheus sink (serves `/metrics` **or** pushes remote-write). | Same sink for both formats. |
| `-prom_listen 0.0.0.0:8883` | Scrape endpoint bind address (format `prometheus`). | Plain HTTP. |
| `-prom_remote_write <url>` | Remote-write endpoint (format `prometheus_remote`). | Requires `-compression snappy`. |
| `-compression snappy` | Required for `prometheus_remote`. | Remote-write payload is snappy-compressed protobuf. |
| `-prom_seen N` | Gates when a series is registered/exposed. | **Very important** — see §4. |
| `-filters "..."` | Restrict to synthetics only (optional). | e.g. `"string,eventType,==,KSynth or string,eventType,==,KSynthgest"`. |
| `-log_level debug` | Verbose logging. | Look for `Adding kentik:synth:outcome`. |

### Why `-udrs` and `-mapping` are mandatory for synthetics

Synthetic kflow arrives with **generic** custom columns (`INT00`, `INT64_02`, `STR00`, …).
The UDR file ([config/udr.csv](config/udr.csv)) renames them into the fields the format reads:

```
10,INT00,Result Type,Synthetic Agent      # INT00      → result_type
10,INT64_00,Agent ID,Synthetic Agent       # INT64_00   → agent_id
10,INT64_02,Test ID,Synthetic Agent        # INT64_02   → test_id
10,INT05,Ping Avg RTT,Synthetic Agent      # …          → avg_rtt, etc.
```

`ktranslate` lowercases the display name (`Result Type` → `result_type`). **If you omit
`-udrs`, `result_type` never gets set**, every record looks like `result_type=0`, and the
outcome metric will report `error` for everything. Always pass `-udrs ./config/udr.csv`.

`-mapping ./config/config.json` supplies enum/name mappings used during enrichment.

---

## 3. Example commands

### A. Scrape model (`/metrics`)

```bash
KENTIK_API_TOKEN=<token> bin/ktranslate \
  -listen 0.0.0.0:8081 \
  -ssl_cert_file /path/firehose.pem -ssl_key_file /path/firehose.key \
  -kentik_email you@example.com \
  -mapping ./config/config.json -udrs ./config/udr.csv \
  -format prometheus -sinks prometheus -prom_listen 0.0.0.0:8883 \
  -filters "string,eventType,==,KSynth or string,eventType,==,KSynthgest" \
  -prom_seen 4 -log_level debug

# scrape (plain HTTP, even though the firehose is TLS):
curl -s http://127.0.0.1:8883/metrics | grep '^kentik:synth:'
```

### B. Remote-write model (push)

```bash
KENTIK_API_TOKEN=<token> bin/ktranslate \
  -listen 0.0.0.0:8081 \
  -ssl_cert_file /path/firehose.pem -ssl_key_file /path/firehose.key \
  -kentik_email you@example.com \
  -mapping ./config/config.json -udrs ./config/udr.csv \
  -format prometheus_remote -sinks prometheus \
  -prom_remote_write https://prometheus.example.com/api/v1/write \
  -compression snappy \
  -filters "string,eventType,==,KSynth or string,eventType,==,KSynthgest" \
  -log_level debug
```

> Remote-write encodes synthetic series as `kentik.synth.<metric>` (dots), including
> `kentik.synth.outcome`, and mesh/syngest series as `kentik.syngest.<metric>`. (Scrape uses
> colons: `kentik:synth:outcome`, `kentik:syngest:<metric>`.)

---

## 4. `-prom_seen` (do not skip this)

`-prom_seen N` sets `FlowsNeeded` (**default 10**). The Prometheus *scrape* format will not
register/expose a metric name until it has been **seen `N` times**, and it uses those first
`N` sightings to learn the metric's **label set**.

| `-prom_seen` | Behavior |
|--------------|----------|
| `0` | Registers on the **first** record — but with **no labels** (all outcomes collapse into a single series; last write wins). Good for a quick "is anything emitted?" check. |
| `1` | Registers on the **2nd** sighting, capturing labels. Send each series **≥2×** to see it. |
| `4`–`10` | Recommended for steady state: avoids churn while a device warms up, then exposes fully-labeled series. |

Notes:
- This gating applies to the **scrape** format only. `prometheus_remote` emits every record
  immediately with full labels (no `-prom_seen` warm-up).
- Startup log lines to watch: `Seen kentik:synth:outcome -> N`, then `Adding kentik:synth:outcome [labels…]`.

---

## 5. Interpreting the `outcome` metric

`result_type` historically conflated **what test ran** (ping/http/trace/…) with **whether it
failed** (error/timeout). The `outcome` metric decouples the failure state so timeouts and
errors — which used to be **dropped** by the metric sinks — are now trackable. It is emitted
for **every** synthetic result.

### Metric name by sink

| Sink / format | Series | Value type |
|---------------|--------|------------|
| `prometheus` (scrape) | `kentik:synth:outcome` | gauge |
| `prometheus_remote` | `kentik.synth.outcome` | gauge |
| `otel`, `new_relic_metric`, `ddog` | `kentik.synth.outcome` | gauge |
| `influx` | `outcome` field in the `ksynth` measurement | integer |
| `elasticsearch` | `synth_outcome` field | integer |

### Value encoding (higher = worse)

| Value | Meaning |
|-------|---------|
| `0` | **ok** — test completed (any test type: ping/http/trace/…) |
| `1` | **timeout** |
| `2` | **error** |

### Labels

Carried alongside the value: `test_type`, `result_type_str` (`ok`/`timeout`/`error`/test name),
plus `test_name`, `agent_name`, `test_id`, `src_*`/`dst_*`, etc. (the whitelisted synth
attributes). `test_type` tells you *what kind of test* independent of the outcome.

### Companion metrics

For successful results you also get the usual gauges (subset depends on test type):
`avg_rtt`, `jit_rtt`, `sent`, `lost`, `time`, `code`, `port`, `status`, `ttlb`, `size`,
`trx_time`, `validation`.

### PromQL examples

```promql
# Any synthetic test currently failing (timeout or error)
kentik:synth:outcome > 0

# Only timeouts / only errors
kentik:synth:outcome == 1
kentik:synth:outcome == 2

# Failing tests grouped by test name
max by (test_name, test_type) (kentik:synth:outcome) > 0
```

Sample alert rule:

```yaml
groups:
  - name: kentik-synthetics
    rules:
      - alert: SyntheticTestFailing
        expr: kentik:synth:outcome > 0
        for: 5m
        labels: { severity: warning }
        annotations:
          summary: "Synthetic test {{ $labels.test_name }} ({{ $labels.test_type }}) is {{ $labels.result_type_str }}"
```

---

## 6. Validate locally without waiting on the firehose

Enable the built-in JCHF HTTP input and inject synthetic records directly:

```bash
# add -http.source=true to the command, then POST synthetic records
# (use https + -k because -listen is TLS; /metrics stays plain http):
curl -sk https://127.0.0.1:8081/input/ktranslate/jchf -H 'Content-Type: application/json' -d '[
  {"eventType":"KSynth","device_name":"agent","custom_int":{"result_type":1},
   "custom_str":{"result_type_str":"timeout","test_type":"http","test_name":"t1"}},
  {"eventType":"KSynth","device_name":"agent","custom_int":{"result_type":0},
   "custom_str":{"result_type_str":"error","test_type":"http","test_name":"t1"}},
  {"eventType":"KSynth","device_name":"agent","custom_int":{"result_type":2},
   "custom_str":{"result_type_str":"ping","test_type":"ping","test_name":"t1"}}
]'

curl -s http://127.0.0.1:8883/metrics | grep '^kentik:synth:'
```

Unit tests that exercise both output paths:

```bash
# scrape format (prints the exposition text)
go test -tags dynamic -run TestPromSynthOutcomeStdout -v ./pkg/formats/prom/
# remote-write format (decodes the protobuf and checks outcome)
go test -tags dynamic -run TestRemotePromSynthOutcome -v ./pkg/formats/prom/
```

---

## 7. Troubleshooting

| Symptom | Cause / fix |
|---------|-------------|
| `/metrics` shows only `go_*` / `process_*` | No synthetic records reaching the format. Confirm the firehose is POSTing to `https://host:8081/chf`, and grep debug logs for `Adding kentik:synth:outcome`. |
| Connection refused on the scrape port | You bound `127.0.0.1` (scrape locally) or scraped with `https://` — the `/metrics` endpoint is **plain HTTP**. |
| Every outcome reports `error` (2) | `-udrs` missing, so `result_type` is never set. Pass `-udrs ./config/udr.csv`. |
| Metric appears but has no labels | `-prom_seen 0`. Use `-prom_seen 1`+ and send each series ≥2×. |
| Series never appears | Not yet seen `-prom_seen` times; send more records or lower the value. |
| Remote-write rejected | Ensure `-compression snappy` and a valid `-prom_remote_write` URL. |
