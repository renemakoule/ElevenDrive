import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.resolve(__dirname, "..");

const pkgPath = path.join(projectRoot, "package.json");
const tauriConfPath = path.join(projectRoot, "src-tauri", "tauri.conf.json");

if (!fs.existsSync(pkgPath) || !fs.existsSync(tauriConfPath)) {
  console.error("❌ Fichiers package.json ou src-tauri/tauri.conf.json introuvables.");
  process.exit(1);
}

const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf-8"));
const tauriConf = JSON.parse(fs.readFileSync(tauriConfPath, "utf-8"));

const currentVersion = pkg.version || "1.0.0";
const arg = process.argv[2] || "patch";

let newVersion = "";

if (arg === "patch") {
  const parts = currentVersion.split(".").map(Number);
  parts[2] = (parts[2] || 0) + 1;
  newVersion = parts.join(".");
} else if (arg === "minor") {
  const parts = currentVersion.split(".").map(Number);
  parts[1] = (parts[1] || 0) + 1;
  parts[2] = 0;
  newVersion = parts.join(".");
} else if (arg === "major") {
  const parts = currentVersion.split(".").map(Number);
  parts[0] = (parts[0] || 0) + 1;
  parts[1] = 0;
  parts[2] = 0;
  newVersion = parts.join(".");
} else if (/^\d+\.\d+\.\d+/.test(arg)) {
  newVersion = arg.replace(/^v/, "");
} else {
  console.error(`❌ Argument invalide : "${arg}". Utilisez "patch", "minor", "major" ou une version ex: "1.1.2"`);
  process.exit(1);
}

pkg.version = newVersion;
tauriConf.version = newVersion;

fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2) + "\n");
fs.writeFileSync(tauriConfPath, JSON.stringify(tauriConf, null, 2) + "\n");

console.log("\n=======================================================");
console.log(` ✅ VERSION INCRÉMENTÉE : ${currentVersion}  ➡️  ${newVersion}`);
console.log(" Fichiers mis à jour automatiquement :");
console.log("  - package.json");
console.log("  - src-tauri/tauri.conf.json");
console.log("=======================================================");
console.log("\nPour publier cette mise à jour sur GitHub :");
console.log(`  git add package.json src-tauri/tauri.conf.json`);
console.log(`  git commit -m "chore(release): v${newVersion}"`);
console.log(`  git tag v${newVersion}`);
console.log(`  git push origin main`);
console.log(`  git push origin v${newVersion}\n`);
