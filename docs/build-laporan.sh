#!/usr/bin/env bash
# Membangun PDF laporan dari docs/LAPORAN_RANCANGAN_SISTEM.md.
# Prasyarat: pandoc, xelatex (TeX Live), dan font DejaVu.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC="$ROOT/docs/LAPORAN_RANCANGAN_SISTEM.md"
OUT="$ROOT/docs/LAPORAN_RANCANGAN_SISTEM.pdf"
BODY="$(mktemp --suffix=.md)"
trap 'rm -f "$BODY"' EXIT

command -v pandoc >/dev/null || { echo "pandoc tidak ditemukan" >&2; exit 1; }
command -v xelatex >/dev/null || { echo "xelatex tidak ditemukan" >&2; exit 1; }

# Badan laporan: mulai dari ABSTRAK, tanpa daftar isi manual.
# Judul, subjudul, dan identitas penulis diberikan melalui metadata pandoc
# agar terbentuk halaman sampul tersendiri.
awk '
/^## ABSTRAK$/ {start=1}
start && /^## DAFTAR ISI$/ {skip=1; next}
skip && /^# BAB I/ {skip=0}
start && !skip && $0 != "---" {print}
' "$SRC" > "$BODY"

pandoc "$BODY" -o "$OUT" \
  --pdf-engine=xelatex \
  --toc \
  -V papersize=a4 \
  -V geometry:left=2.5cm \
  -V geometry:right=2.5cm \
  -V geometry:top=3cm \
  -V geometry:bottom=3cm \
  -V fontsize=11pt \
  -V lang=id \
  -V documentclass=report \
  -V mainfont="DejaVu Serif" \
  -V sansfont="DejaVu Sans" \
  -V monofont="DejaVu Sans Mono" \
  --metadata title="Rancangan Sistem Helpdesk Tiket Layanan Akademik “LaporKampus” Berbasis Go dan Fiber dengan Arsitektur Berlapis" \
  --metadata subtitle="Laporan Rancangan Sistem" \
  --metadata author="Vyerru — 434241021" \
  --metadata date="2026" \
  --include-in-header="$ROOT/docs/laporan-header.tex"

echo "PDF dibuat: $OUT"
