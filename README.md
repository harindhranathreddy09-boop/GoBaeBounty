# GoBaeBounty - Comprehensive Bug Bounty Automation Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

GoBaeBounty is a production-grade bug bounty framework implemented in Go that automates the process from passive reconnaissance through dynamic fuzzing and vulnerability detection, designed for authorized security research purposes only.

## Features

- Authorization enforced via HMAC-signed auth files
- Passive domain and URL discovery using popular tools integration
- Polite concurrency-controlled crawling and JS endpoint extraction
- Endpoint fingerprinting and priority scoring
- Endpoint-driven recursive fuzzing with intelligent filtering
- Modular vulnerability check plugins for XSS, SQLi, SSRF, IDOR, etc.
- Fully automated report generation formatted for HackerOne disclosures
- User-configurable rate limiting, concurrency, and scanning depth

## Installation

