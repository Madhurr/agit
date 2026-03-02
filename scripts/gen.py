#!/usr/bin/env python3
"""
agit codegen — streams GLM output via curl, byte-by-byte, no buffering.
Usage: python3 scripts/gen.py prompts/01_git_wrapper.txt internal/git/git.go
"""
import sys
import json
import subprocess
import os
import re
import time
import threading

MODEL = "MichelRosselli/GLM-4.5-Air:Q4_K_M"
OLLAMA_URL = "http://localhost:11434/api/generate"


def strip_fences(text):
    text = re.sub(r"^```[a-zA-Z]*\n", "", text, flags=re.MULTILINE)
    text = re.sub(r"^```\s*$", "", text, flags=re.MULTILINE)
    return text.strip()


def stream(prompt, out_path):
    payload = json.dumps({
        "model": MODEL,
        "prompt": prompt,
        "stream": True,
        "options": {
            "temperature": 0.2,
            "num_predict": 8192,
        }
    })

    print(f"\n>>> Generating: {out_path}\n", flush=True)
    print("─" * 60, flush=True)

    proc = subprocess.Popen(
        ["curl", "-s", "-N", "-X", "POST", OLLAMA_URL,
         "-H", "Content-Type: application/json",
         "-d", payload],
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
        bufsize=0  # unbuffered
    )

    tokens = []
    buf = b""
    first_token = threading.Event()
    start = time.time()

    # Spinner: shows elapsed time until first token arrives
    def spinner():
        while not first_token.is_set():
            elapsed = time.time() - start
            sys.stdout.write(f"\r  waiting for model... {elapsed:.1f}s")
            sys.stdout.flush()
            time.sleep(0.2)
        sys.stdout.write("\r" + " " * 40 + "\r")
        sys.stdout.flush()

    t = threading.Thread(target=spinner, daemon=True)
    t.start()

    # Read byte-by-byte — same as terminal, zero buffering delay
    while True:
        byte = proc.stdout.read(1)
        if not byte:
            break
        buf += byte
        if byte == b"\n":
            line = buf.decode("utf-8", errors="replace").strip()
            buf = b""
            if not line:
                continue
            try:
                chunk = json.loads(line)
            except json.JSONDecodeError:
                continue
            token = chunk.get("response", "")
            if token:
                first_token.set()  # stop spinner
                sys.stdout.write(token)
                sys.stdout.flush()
                tokens.append(token)
            if chunk.get("done"):
                break

    proc.wait()
    print("\n" + "─" * 60, flush=True)

    raw = "".join(tokens)
    if not raw.strip():
        print("ERROR: No output from model. Is Ollama running?")
        sys.exit(1)

    code = strip_fences(raw)

    os.makedirs(os.path.dirname(out_path) if os.path.dirname(out_path) else ".", exist_ok=True)
    with open(out_path, "w") as f:
        f.write(code + "\n")

    print(f"\n✓ Written to: {out_path} ({len(code.splitlines())} lines)\n", flush=True)


def main():
    if len(sys.argv) != 3:
        print("Usage: python3 scripts/gen.py <prompt_file> <output_file>")
        sys.exit(1)

    prompt_file, out_file = sys.argv[1], sys.argv[2]

    if not os.path.exists(prompt_file):
        print(f"ERROR: Prompt file not found: {prompt_file}")
        sys.exit(1)

    with open(prompt_file) as f:
        prompt = f.read()

    stream(prompt, out_file)


if __name__ == "__main__":
    main()
