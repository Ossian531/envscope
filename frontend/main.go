package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type envPayload struct {
	Hostname string            `json:"hostname"`
	ServedAt string            `json:"served_at"`
	Count    int               `json:"count"`
	EnvVars  map[string]string `json:"env_vars"`
}

type kv struct {
	Key   string
	Value string
}

type viewData struct {
	BackendURL string
	Hostname   string
	ServedAt   string
	Count      int
	Vars       []kv
	Err        string
}

func main() {
	addr := getenv("FRONTEND_ADDR", ":3000")
	backendURL := getenv("BACKEND_HOST", "localhost")
	backendPort := getenv("BACKEND_PORT", "8080")
	backendURL = "http://" + backendURL + ":" + backendPort + "/"

	tmpl := template.Must(template.New("page").Parse(pageTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data := viewData{BackendURL: backendURL}
		payload, err := fetchBackend(r.Context(), backendURL)
		if err != nil {
			data.Err = err.Error()
		} else {
			data.Hostname = payload.Hostname
			data.ServedAt = payload.ServedAt
			data.Count = payload.Count
			for k, v := range payload.EnvVars {
				data.Vars = append(data.Vars, kv{Key: k, Value: v})
			}
			sort.Slice(data.Vars, func(i, j int) bool {
				return data.Vars[i].Key < data.Vars[j].Key
			})
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("template error: %v", err)
		}
	})

	log.Printf("frontend listening on %s (backend: %s)", addr, backendURL)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func fetchBackend(ctx context.Context, url string) (*envPayload, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p envPayload
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

const pageTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>envscope · backend environment</title>
<style>
  :root {
    --bg: #0b0f17;
    --panel: #121826;
    --panel-2: #0e1420;
    --line: #1e2a3d;
    --text: #e6edf6;
    --muted: #8aa0bd;
    --accent: #5eead4;
    --accent-2: #818cf8;
    --key: #7dd3fc;
  }
  * { box-sizing: border-box; }
  body {
    margin: 0;
    font: 15px/1.5 ui-monospace, SFMono-Regular, "JetBrains Mono", Menlo, monospace;
    color: var(--text);
    background:
      radial-gradient(1200px 600px at 80% -10%, rgba(129,140,248,.18), transparent 60%),
      radial-gradient(900px 500px at -10% 20%, rgba(94,234,212,.12), transparent 55%),
      var(--bg);
    min-height: 100vh;
  }
  .wrap { max-width: 1000px; margin: 0 auto; padding: 48px 24px 80px; }
  header { display: flex; align-items: baseline; gap: 14px; flex-wrap: wrap; }
  h1 {
    margin: 0; font-size: 28px; letter-spacing: -.5px;
    background: linear-gradient(90deg, var(--accent), var(--accent-2));
    -webkit-background-clip: text; background-clip: text; color: transparent;
  }
  .tag {
    font-size: 12px; color: var(--muted); border: 1px solid var(--line);
    padding: 3px 10px; border-radius: 999px;
  }
  .meta {
    display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 12px; margin: 28px 0;
  }
  .card {
    background: linear-gradient(180deg, var(--panel), var(--panel-2));
    border: 1px solid var(--line); border-radius: 14px; padding: 16px 18px;
  }
  .card .label { font-size: 11px; text-transform: uppercase; letter-spacing: 1px; color: var(--muted); }
  .card .val { font-size: 18px; margin-top: 6px; word-break: break-all; }
  .search {
    width: 100%; margin: 8px 0 20px; padding: 12px 16px;
    background: var(--panel-2); border: 1px solid var(--line);
    border-radius: 12px; color: var(--text); font: inherit; outline: none;
  }
  .search:focus { border-color: var(--accent-2); box-shadow: 0 0 0 3px rgba(129,140,248,.18); }
  table { width: 100%; border-collapse: collapse; }
  thead th {
    text-align: left; font-size: 11px; text-transform: uppercase; letter-spacing: 1px;
    color: var(--muted); padding: 0 14px 10px; border-bottom: 1px solid var(--line);
  }
  tbody tr { border-bottom: 1px solid rgba(30,42,61,.5); }
  tbody tr:hover { background: rgba(129,140,248,.06); }
  td { padding: 11px 14px; vertical-align: top; }
  td.k { color: var(--key); white-space: nowrap; font-weight: 600; }
  td.v { color: var(--text); word-break: break-all; }
  .empty, .err {
    background: var(--panel); border: 1px solid var(--line); border-radius: 14px;
    padding: 24px; color: var(--muted);
  }
  .err { border-color: #5b2330; color: #ffb4b4; }
  .foot { margin-top: 32px; color: var(--muted); font-size: 12px; }
  a { color: var(--accent); }
</style>
</head>
<body>
<div class="wrap">
  <header>
    <h1>envscope</h1>
    <span class="tag">backend environment viewer</span>
  </header>

  {{if .Err}}
    <div class="err">
      <strong>Could not reach the backend.</strong><br>
      {{.Err}}<br><br>
      Target: <code>{{.BackendURL}}</code>
    </div>
  {{else}}
    <div class="meta">
      <div class="card"><div class="label">Hostname</div><div class="val">{{.Hostname}}</div></div>
      <div class="card"><div class="label">Variables</div><div class="val">{{.Count}}</div></div>
      <div class="card"><div class="label">Served at</div><div class="val">{{.ServedAt}}</div></div>
      <div class="card"><div class="label">Source</div><div class="val">{{.BackendURL}}</div></div>
    </div>

    <input id="search" class="search" type="text" placeholder="Filter variables…" autocomplete="off">

    {{if .Vars}}
    <table id="envtable">
      <thead><tr><th>Key</th><th>Value</th></tr></thead>
      <tbody>
        {{range .Vars}}
        <tr><td class="k">{{.Key}}</td><td class="v">{{.Value}}</td></tr>
        {{end}}
      </tbody>
    </table>
    {{else}}
      <div class="empty">The backend reported no environment variables.</div>
    {{end}}
  {{end}}

  <div class="foot">Rendered server-side by the Go frontend · data pulled live from the Go backend.</div>
</div>

<script>
  const box = document.getElementById('search');
  if (box) {
    box.addEventListener('input', () => {
      const q = box.value.toLowerCase();
      document.querySelectorAll('#envtable tbody tr').forEach(tr => {
        tr.style.display = tr.textContent.toLowerCase().includes(q) ? '' : 'none';
      });
    });
  }
</script>
</body>
</html>`
