#!/usr/bin/env python3

# Copyright 2026 The cert-manager Authors.
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

"""Print peak memory (RSS) per Prow job as a Markdown table.

Method from https://github.com/cert-manager/testing/issues/1240:
  1. Map pod names to job names from prowjob.json in the public artifacts
     bucket. A Prow pod has the same name as its ProwJob.
  2. Ask Cloud Monitoring for the peak container_memory_rss of every test pod.

Needs `gcloud auth login` with read access to Cloud Monitoring in the
untrusted project. Step 1 needs no credentials; --map-only stops after it.

  hack/prow-job-memory-peaks.py --days 30 > peaks.md
"""
import argparse
import concurrent.futures
import json
import os
import shutil
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

BUCKET = "cert-manager-prow-artifacts"
GCS = f"https://storage.googleapis.com/storage/v1/b/{BUCKET}/o"
RAW = f"https://storage.googleapis.com/{BUCKET}/"
PROM = "https://monitoring.googleapis.com/v1/projects/{}/location/global/prometheus/api/v1/query"
QUERY = 'max by (pod) (max_over_time(container_memory_rss{namespace="test-pods",container="test"}[1d]))'
CACHE = os.path.expanduser("~/.cache/prow-job-memory-peaks.json")
DAY = 86400
GIB = 2**30


def get(url, params=None, token=None):
    if params:
        url += "?" + urllib.parse.urlencode(params)
    req = urllib.request.Request(url)
    if token:
        req.add_header("Authorization", "Bearer " + token)
    for attempt in range(5):
        try:
            with urllib.request.urlopen(req, timeout=120) as r:
                return r.read().decode()
        except urllib.error.HTTPError as e:
            if e.code < 500 and e.code != 429:
                raise
        except OSError:
            pass
        time.sleep(2**attempt)
    raise RuntimeError(f"giving up on {url}")


def gcs_list(prefix, start=None):
    """Return (sub-directories, object names) directly under prefix."""
    dirs, names, page = [], [], None
    while True:
        params = {"prefix": prefix, "delimiter": "/", "fields": "prefixes,items(name),nextPageToken"}
        if start:
            params["startOffset"] = prefix + start
        if page:
            params["pageToken"] = page
        d = json.loads(get(GCS, params))
        dirs += d.get("prefixes", [])
        names += [i["name"] for i in d.get("items", [])]
        page = d.get("nextPageToken")
        if not page:
            return dirs, names


def min_build_id(since):
    # Prow build IDs are Twitter-style snowflakes: milliseconds since
    # 2010-11-04T01:42:54.657Z in the bits above the low 22.
    return (int(since * 1000) - 1288834974657) << 22


def builds_of(job_dir, min_id):
    """Paths of the builds of one job that started after min_id."""
    dirs, names = gcs_list(job_dir, start=str(min_id))
    if job_dir.startswith("logs/"):
        ids = [(d.rstrip("/").rsplit("/", 1)[1], d) for d in dirs]
    else:  # pr-logs/directory/<job>/<id>.txt holds the gs:// path of the build
        ids = [(n.rsplit("/", 1)[1].removesuffix(".txt"), n) for n in names]
    return [path for build_id, path in ids if build_id.isdigit() and len(build_id) == 19 and int(build_id) >= min_id]


def describe(path):
    """Return [pod, job, memory request] for one build, or None."""
    try:
        if path.endswith(".txt"):
            path = get(RAW + path).strip().removeprefix(f"gs://{BUCKET}/") + "/"
        pj = json.loads(get(RAW + path + "prowjob.json"))
        containers = (pj["spec"].get("pod_spec") or {}).get("containers") or [{}]
        request = containers[0].get("resources", {}).get("requests", {}).get("memory", "")
        return [pj["metadata"]["name"], pj["spec"]["job"], request]
    except Exception as e:  # a build with no prowjob.json is not worth stopping for
        print(f"skip {path}: {e}", file=sys.stderr)
        return None


def pod_map(since):
    try:
        with open(CACHE) as f:
            cache = json.load(f)
    except FileNotFoundError:
        cache = {}
    min_id = min_build_id(since)
    job_dirs = gcs_list("logs/")[0] + gcs_list("pr-logs/directory/")[0]
    with concurrent.futures.ThreadPoolExecutor(32) as pool:
        paths = {p for ps in pool.map(lambda d: builds_of(d, min_id), job_dirs) for p in ps}
        todo = [p for p in paths if p not in cache]
        print(f"{len(paths)} builds in window, fetching {len(todo)}", file=sys.stderr)
        for path, result in zip(todo, pool.map(describe, todo)):
            if result:
                cache[path] = result
    os.makedirs(os.path.dirname(CACHE), exist_ok=True)
    with open(CACHE, "w") as f:
        json.dump(cache, f)
    return {v[0]: (v[1], v[2]) for p, v in cache.items() if p in paths}


def peak_rss(project, end, days):
    gcloud = shutil.which("gcloud") or os.path.expanduser("~/google-cloud-sdk/bin/gcloud")
    token = subprocess.check_output([gcloud, "auth", "print-access-token"], text=True).strip()
    peaks = {}
    for i in range(days):  # one query per day: a single 30d range can time out
        d = json.loads(get(PROM.format(project), {"query": QUERY, "time": end - i * DAY}, token))
        for s in d["data"]["result"]:
            pod = s["metric"]["pod"]
            peaks[pod] = max(peaks.get(pod, 0), float(s["value"][1]))
    return peaks


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--days", type=int, default=30)
    ap.add_argument("--project", default="cert-manager-tests-untrusted")
    ap.add_argument("--map-only", action="store_true", help="only build the pod to job map")
    args = ap.parse_args()

    end = int(time.time())
    # One extra day: a pod can outlive the day its build ID was issued.
    pods = pod_map(end - (args.days + 1) * DAY)
    print(f"{len(pods)} pods mapped to {len({j for j, _ in pods.values()})} jobs", file=sys.stderr)
    if args.map_only:
        return

    peaks = peak_rss(args.project, end, args.days)
    jobs, unmatched = {}, 0
    for pod, rss in peaks.items():
        if pod in pods:
            jobs.setdefault(pods[pod], []).append(rss)
        else:
            unmatched += 1
    print(f"{len(peaks)} pods measured, {unmatched} not matched to a job", file=sys.stderr)

    print("| Job | Runs | Max RSS | p95 RSS | Request |\n|---|--:|--:|--:|--:|")
    for (job, request), rss in sorted(jobs.items(), key=lambda kv: -max(kv[1])):
        rss.sort()
        p95 = rss[min(len(rss) - 1, int(len(rss) * 0.95))]
        print(f"| `{job}` | {len(rss)} | {rss[-1] / GIB:.1f} | {p95 / GIB:.1f} | {request} |")


if __name__ == "__main__":
    main()
