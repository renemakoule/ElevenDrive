import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const envFile = path.resolve(__dirname, "..", ".env");

// Charge le fichier .env en local s'il existe (sans faire planter GitHub Actions s'il n'existe pas)
if (fs.existsSync(envFile) && typeof process.loadEnvFile === "function") {
  try {
    process.loadEnvFile(envFile);
  } catch {}
}

await import("../node_modules/@tauri-apps/cli/tauri.js");
