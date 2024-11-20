# Studio Helm Chart

This Helm chart deploys Supabase Studio along with its required services including Supa-manager API and PostgreSQL.

## Overview

The chart deploys the following components:
- Supabase Studio UI
- Supa-manager API
- PostgreSQL database
- DNS Example Service (optional)
- Version Service (optional)

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- Ingress controller (e.g., nginx-ingress)
- cert-manager (optional, for TLS)
- Docker registry accessible from the cluster (when building in-cluster)

## Installation

There are two ways to deploy Studio: using the in-cluster build process or providing your own pre-built image.

### Option 1: In-Cluster Build (Recommended)

This option will build the Studio UI directly in your cluster using the official Supabase repository.

1. Create a values file (e.g., `my-values.yaml`):