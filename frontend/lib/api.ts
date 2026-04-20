import { Scan, type ModuleName } from "@/lib/types";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function startScan(payload: { url: string; modules: ModuleName[] }): Promise<Scan> {
  const response = await fetch(`${apiBase}/api/scan`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    const error = await readError(response);
    throw new Error(error);
  }

  return response.json() as Promise<Scan>;
}

export async function getScan(id: string): Promise<Scan> {
  const response = await fetch(`${apiBase}/api/scan/${id}`, {
    cache: "no-store"
  });
  if (!response.ok) {
    const error = await readError(response);
    throw new Error(error);
  }
  return response.json() as Promise<Scan>;
}

export function streamScan(
  id: string,
  onMessage: (scan: Scan) => void,
  onError: (message: string) => void
) {
  const source = new EventSource(`${apiBase}/api/scan/${id}/stream`);
  source.onmessage = (event) => {
    try {
      onMessage(JSON.parse(event.data) as Scan);
    } catch {
      onError("Canlı veri çözümlenemedi.");
    }
  };
  source.onerror = () => {
    onError("SSE bağlantısı düştü; polling devam ediyor.");
  };
  return source;
}

export function pdfUrl(id: string) {
  return `${apiBase}/api/scan/${id}/report.pdf`;
}

async function readError(response: Response) {
  try {
    const payload = (await response.json()) as { error?: string };
    return payload.error ?? `İstek başarısız: ${response.status}`;
  } catch {
    return `İstek başarısız: ${response.status}`;
  }
}
