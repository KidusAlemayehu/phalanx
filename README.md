    # 🛡 Phalanx

    > **A high-performance, lightweight TLS and mTLS proxy written in Go.**

    Phalanx is an experimental, from-scratch network proxy designed to handle secure routing and interservice communication. It acts as a shield for backend services, offloading complex cryptographic handshakes and enforcing strict identity-based access control.

    ## ✨ Features (Planned)

    *   **TCP Byte Streaming:** Fast, bi-directional I/O routing.
    *   **TLS Termination:** Offload computationally heavy cryptographic handshakes at the edge.
    *   **Dynamic SNI Routing:** Route traffic to multiple distinct backends based on the requested Server Name Indication.
    *   **Zero-Trust mTLS:** Enforce mutual TLS for sidecar-to-sidecar interservice communication, ensuring cryptographic identity verification.

    ## 🧠 Why Phalanx?

    In modern distributed systems, relying solely on perimeter firewalls is no longer sufficient. Phalanx is built to explore the foundations of a **Zero-Trust architecture**, demonstrating how proxies like Envoy and Linkerd handle raw bytes, certificates, and secure connections under the hood.

    ## Getting Started

    *(Instructions for building and running the proxy will go here once Phase 1 is complete!)*

    ---
    *Built as an exploration of modern network security and Go's concurrency model.*
    ──────
