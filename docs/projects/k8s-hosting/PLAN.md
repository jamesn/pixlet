# Project: Kubernetes Hosting for the Sign Server

Status: Draft
Owner: @jamesn
Last updated: 2026-09-24
Track: Part of the independent [firmware fork](../firmware-fork/PLAN.md) track.
Decides its **FD4** ("where the server runs"): **the owner's existing Kubernetes cluster**.
Estimated effort: **~2–5 days** on top of the server itself (§9).

## 1. Summary

Run the self-hosted sign server (the component that renders pixlet apps and
serves images to the signs) on the owner's existing home Kubernetes cluster.
The project also adds a **pixlet container image**, built and published from this
repo, which the server image builds on.

## 2. Context

```
 ┌──────────────── home Kubernetes cluster ────────────────┐
 │                                                         │
 │  Deployment: sign-server (1 replica)                    │
 │    ├─ container: server  (FROM pixlet image)            │
 │    │     renders .star apps with pixlet → WebP          │
 │    ├─ PVC: /data  (installations, app config, cache)    │
 │    ├─ Secret: app secrets, server API token             │
 │    └─ probes: /healthz (liveness), /readyz (readiness)  │
 │                                                         │
 │  Service: type LoadBalancer (stable LAN IP)             │
 └──────────────────────────┬──────────────────────────────┘
                            │ home Wi-Fi / LAN only
          ┌─────────────────┼─────────────────┐
       sign 1            sign 2            sign N    (firmware we control)
```

- The **server** software itself (adopted/forked) is chosen in the firmware
  project's Phase 0 (FD2). This plan is written to fit any server that can run in
  a container. §4 lists what the server must provide.
- The signs are on the home network and must reach the server by a **stable
  address** that survives pod restarts and rescheduling.

### Assumptions (verify in Phase 0)

| # | Assumption | If false |
|---|------------|----------|
| A1 | Cluster is reachable from the home LAN where the signs are | Add a LAN-facing ingress path first |
| A2 | A LoadBalancer implementation exists (e.g. MetalLB, kube-vip, k3s ServiceLB) **or** a fixed NodePort on a stable node IP is acceptable | Use NodePort + DHCP reservation for the node |
| A3 | A default StorageClass with ReadWriteOnce volumes exists | Use a hostPath/local-path volume pinned to one node |
| A4 | Nodes are amd64 and/or arm64 | Both are built (§6 Phase 1) |
| A5 | The owner deploys with either plain manifests/Kustomize or a GitOps tool (Argo CD / Flux) | Plan supports all three |
| A6 | Cluster can pull from GitHub Container Registry (`ghcr.io`) | Mirror to a local registry |

## 3. Goals / Non-goals

**Goals**
- G1: The sign server runs on the existing cluster with self-healing and easy rollbacks.
- G2: Signs reach the server at a stable LAN address/DNS name.
- G3: All state lives in one persistent volume and is backed up.
- G4: Deploys are declarative and reproducible (versioned images, versioned manifests).
- G5: A pixlet container image is published from this repo for any use (server base, CI, local rendering).

**Non-goals**
- Multi-replica high availability. One replica is enough; signs tolerate short outages (firmware R9).
- Exposing anything to the internet.
- A Kubernetes operator/CRDs for app installations. Possible future work (§10); not needed.
- Hosting the OTA firmware signing key in the cluster (see R8).

## 4. Requirements

### What the server must provide (container-readiness, "12-factor")
- **R1** Configuration via environment variables and/or one mounted config file; no interactive setup.
- **R2** All persistent state under a single directory (default `/data`).
- **R3** HTTP health endpoints: liveness (`/healthz`) and readiness (`/readyz`, ready once state is loaded).
- **R4** Logs to stdout/stderr.
- **R5** Graceful shutdown on `SIGTERM` within 30 s (finish in-flight renders, flush state).
- **R6** Runs as a non-root user with a read-only root filesystem (writes only to `/data` and `/tmp`).

If the adopted server lacks any of R1–R6, add it in the server fork (firmware project Phase 3).

### Platform
- **R7** Signs MUST reach the server via a stable LAN IP or local DNS name that doesn't change across pod restarts, node reboots, or redeploys.
- **R8** Secrets (app secrets, server API token) MUST be Kubernetes Secrets (or sealed/encrypted equivalents if GitOps). The **firmware OTA signing key MUST NOT** be stored in the cluster; it stays offline with the owner.
- **R9** Container images MUST be multi-arch (`linux/amd64`, `linux/arm64`), versioned by git tag, and published to `ghcr.io/jamesn/...`.
- **R10** Images MUST be scanned for vulnerabilities in CI; high/critical findings block release.
- **R11** `/data` MUST be backed up at least daily with a tested restore.
- **R12** Resource requests/limits MUST be set; rendering spikes must not starve other cluster workloads.
- **R13** The server Service MUST NOT be exposed outside the LAN (no Ingress to the internet, no port forwarding).

### Compatibility
- **R14** The server image MUST be built on the pixlet image from this repo at a pinned release tag, so app output follows pixlet's contracts (C1/C2).

## 5. Design

### 5.1 Images

| Image | Repo | Contents |
|-------|------|----------|
| `ghcr.io/jamesn/pixlet:<tag>` | this repo | `pixlet` binary + `libwebp` runtime libs; entrypoint `pixlet` |
| `ghcr.io/jamesn/<server>:<tag>` | server fork | server binary/app + `FROM ghcr.io/jamesn/pixlet:<pinned>` (or copies the pixlet binary from it) |

Pixlet `Dockerfile` (multi-stage):
1. **frontend stage:** `node:22` → `npm ci && npm run build` (fills `dist/`).
2. **build stage:** `golang` (version from `go.mod`) + `libwebp-dev` → `make build`.
   For arm64, either build natively via QEMU/buildx or cross-compile with
   `crossbuild-essential-arm64` + `libwebp-dev:arm64` (as `scripts/setup-linux.sh` already does).
3. **runtime stage:** slim Debian/Ubuntu base + `libwebp7 libwebpdemux2 libwebpmux3`
   (cgo shared libs), non-root user, `ENTRYPOINT ["pixlet"]`.

Add a `.dockerignore` (exclude `node_modules`, `.git`, build output).

### 5.2 Kubernetes objects

| Object | Notes |
|--------|-------|
| `Namespace: signs` | Isolates everything. |
| `Deployment: sign-server` | `replicas: 1`, `strategy: Recreate` (RWO volume can't be shared during rollout), probes (R3), securityContext (R6), requests/limits (R12). |
| `PersistentVolumeClaim: sign-server-data` | RWO, e.g. 1 Gi; mounted at `/data`. |
| `Secret: sign-server-secrets` | Server API token, app secrets (R8). |
| `ConfigMap: sign-server-config` | Non-secret settings (R1). |
| `Service: sign-server` | `type: LoadBalancer` with a **fixed IP** (e.g. MetalLB `loadBalancerIPs` annotation) → R7. Fallback: `NodePort` + node DHCP reservation. |
| `NetworkPolicy` | Allow ingress from the LAN CIDR only (R13). |
| Local DNS | `signs.home.arpa` (or similar) → Service IP, via the home router / Pi-hole / CoreDNS. Firmware points at the name, not the IP. |
| Backup | `CronJob` that snapshots `/data` to NAS/object storage, **or** the cluster's existing backup tool (e.g. Velero/Longhorn snapshots) → R11. |

Packaging: **Kustomize** base + overlay (fits plain `kubectl`, Argo CD, and Flux).
A Helm chart is optional later.

### 5.3 Releases and upgrades

```
pixlet tag v0.4x.y ──► CI builds ghcr.io/jamesn/pixlet:v0.4x.y (amd64+arm64, scanned)
                          │
server fork bumps FROM ───┘──► CI builds ghcr.io/jamesn/<server>:vA.B.C
                                    │
manifests repo/overlay bumps image ─┘──► kubectl apply / GitOps sync ──► rollout
                                                                          (rollback = previous tag)
```

- Dependabot on the server fork watches the pixlet image tag (Docker ecosystem).
- Upgrades go to production directly (one replica); rollback is re-applying the previous tag.
  Firmware-side, the signs keep showing their last image during the brief restart.

## 6. Plan

### Phase 0: Cluster fit check (≈0.5 day)
1. Verify assumptions A1–A6; record answers in §2.
2. Reserve a LAN IP for the Service (MetalLB pool or router DHCP reservation) and pick the local DNS name.
3. Decide the backup target (NAS path / existing backup tool).

**Exit:** A1–A6 answered; IP + DNS name reserved.

### Phase 1: Pixlet container image, in this repo (≈1 day)
1. Add `Dockerfile` + `.dockerignore` (§5.1).
2. Add a CI workflow: on tag push, `docker buildx` for amd64+arm64 → push `ghcr.io/jamesn/pixlet:<tag>` and `:latest`; on PRs, build without pushing.
3. Scan the image (e.g. Trivy or Grype) in CI (R10).
4. Smoke test in CI: `docker run ghcr.io/jamesn/pixlet:<tag> render examples/clock/clock.star` produces a WebP identical to the native build's.

**Exit:** tagged multi-arch pixlet image published; smoke test green.

### Phase 2: Server image (≈0.5–1 day, in the server fork; after firmware project FD2)
1. Add/adjust the server's `Dockerfile` to build on the pinned pixlet image (R14).
2. Close any R1–R6 gaps in the server.
3. CI: multi-arch build, scan, publish on tag.

**Exit:** server image runs locally with `docker run -v data:/data ...` and passes health checks.

### Phase 3: Kubernetes manifests (≈1 day)
1. Kustomize base with the objects in §5.2; overlay for the home cluster (IP, DNS, storage class, resource sizes).
2. Deploy to the `signs` namespace; confirm probes, PVC binding, Service IP.
3. Wire local DNS.
4. Hook into the owner's GitOps tool if used (A5).

**Exit:** `curl http://signs.home.arpa/healthz` from the LAN returns OK.

### Phase 4: Resilience and backup verification (≈0.5–1 day)
1. Kill the pod → rescheduled, same IP/DNS, signs reconnect.
2. Drain/reboot the node running it → same.
3. Roll out a new image tag, then roll back.
4. Run the backup, delete the PVC in a test namespace, restore, and verify installations survive (R11).
5. Confirm NetworkPolicy blocks non-LAN sources (R13).

**Exit:** all pass; results recorded here.

### Phase 5: Hand-off to the firmware pilot
- The firmware project's Phase 1 bench spike and Phase 2 pilot point the pilot sign at `signs.home.arpa`.
- The server on Kubernetes becomes the production server for the rollout (firmware Phase 4).

## 7. Relationship to other work

| Item | Relationship |
|------|--------------|
| [Firmware fork](../firmware-fork/PLAN.md) | Parent track. This plan decides FD4 and supplies the server runtime. Its Phase 2 depends on FD2 (server choice). |
| Pixlet roadmap | Phase 1 (pixlet image) is independent and useful immediately (CI, local rendering), so it can ship in any pixlet release. |
| Spec 011 Phase A | If the server accepts pushes, `pixlet push --target` can point at `http://signs.home.arpa`. Optional. |
| Specs 001/002 | The image scan and Dependabot (Docker ecosystem) extend the same security posture. |

## 8. Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Cluster outage takes the server down | Signs stop updating (show last image) | Low–Med | Firmware R9 graceful degradation; simple redeploy; backups |
| Service IP changes after a cluster change | Signs can't find the server | Low | Fixed LB IP + DNS name; firmware uses the DNS name |
| RWO volume stuck on a dead node | Pod can't start elsewhere | Medium on local storage | Use replicated storage if available (e.g. Longhorn); else documented restore-from-backup runbook |
| cgo/libwebp mismatch in the image | Render failures | Low | Runtime libs from the same distro as the build stage; smoke test compares output with the native build |
| Secrets committed in GitOps repo | Credential leak | Medium | SealedSecrets / SOPS / External Secrets; never plain Secrets in git |
| Cluster maintenance becomes a sign dependency | More moving parts | Medium | Accepted: owner already runs the cluster; the server is a small single workload |

## 9. Effort estimate

| Phase | Effort |
|-------|--------|
| 0: Cluster fit check | ~0.5 day |
| 1: Pixlet image + CI | ~1 day |
| 2: Server image | ~0.5–1 day |
| 3: Manifests + DNS | ~1 day |
| 4: Resilience + backup tests | ~0.5–1 day |
| **Total** | **~3.5–4.5 days** (≈2–5 days range, depending on how much of Phase 2 the chosen server already provides) |

These are rough estimates for one person with an existing, working cluster. They exclude choosing and adapting the server itself (firmware project).

## 10. Future options (not planned)

- Render apps as Kubernetes `CronJob`s per installation, or an operator with a
  `SignApp` CRD (declarative apps in git). Adds 2–4 weeks; only if GitOps-managed
  app configs become valuable.
- Helm chart for sharing.
- Prometheus metrics (renders, failures, per-sign last-seen) + Grafana dashboard.

## 11. Deliverables checklist

- [ ] Phase 0: assumptions verified; LB IP + DNS name reserved; backup target chosen
- [ ] Phase 1: `ghcr.io/jamesn/pixlet` multi-arch image, scanned, smoke-tested
- [ ] Phase 2: server image on pinned pixlet; R1–R6 satisfied
- [ ] Phase 3: Kustomize manifests deployed; DNS resolves; health OK from LAN
- [ ] Phase 4: failover, rollback, backup/restore, and NetworkPolicy verified
- [ ] Phase 5: pilot sign pointed at the cluster-hosted server

## 12. Open questions

1. Cluster details: distribution (k3s / Talos / kubeadm / …), node architectures, LoadBalancer implementation, storage class (replicated?), GitOps tool?
2. Preferred local DNS name (e.g. `signs.home.arpa`) and where local DNS is managed.
3. Backup target: NAS, object storage, or the cluster's existing backup tool?
4. Should the pixlet image (Phase 1) ship now as a standalone pixlet improvement, ahead of the server choice?
