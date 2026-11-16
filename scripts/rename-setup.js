import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const projectRoot = path.resolve(__dirname, "..");
const nsisDir = path.join(projectRoot, "src-tauri", "target", "release", "bundle", "nsis");
const outDir = path.join(projectRoot, "dist-installer");

if (!fs.existsSync(nsisDir)) {
  console.warn(`[rename-setup] Dossier NSIS introuvable : ${nsisDir}`);
  process.exit(0);
}

const files = fs.readdirSync(nsisDir);
const setupExe = files.find(
  (f) => f.toLowerCase().endsWith(".exe") && !f.toLowerCase().includes("elevendrive.exe")
);

if (!setupExe) {
  console.log("[rename-setup] Aucun installateur NSIS intermédiaire à renommer.");
  process.exit(0);
}

if (!fs.existsSync(outDir)) {
  fs.mkdirSync(outDir, { recursive: true });
}

const sourceFile = path.join(nsisDir, setupExe);
const destFile = path.join(outDir, "ElevenDrive.exe");
const inPlaceDestFile = path.join(nsisDir, "ElevenDrive.exe");

fs.copyFileSync(sourceFile, destFile);
fs.copyFileSync(sourceFile, inPlaceDestFile);

console.log("\n=======================================================");
console.log(" INSTALLATEUR GÉNÉRÉ AVEC SUCCÈS :");
console.log(` -> ${destFile}`);
console.log(` -> ${inPlaceDestFile}`);
console.log("=======================================================\n");

