<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from "vue";
import {
  Download,
  X,
  ExternalLink,
  Laptop,
  CheckCircle2,
  Calendar,
  HardDrive,
  RefreshCw,
  AlertCircle,
} from "lucide-vue-next";
import appLogo from "@root/app-icon.png";

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

interface ReleaseAsset {
  name: string;
  size: number;
  browser_download_url: string;
}

interface GitHubRelease {
  id: number;
  tag_name: string;
  name: string;
  published_at: string;
  prerelease: boolean;
  html_url: string;
  assets: ReleaseAsset[];
}

interface FormattedVersion {
  version: string;
  title: string;
  date: string;
  isLatest: boolean;
  isPrerelease: boolean;
  exeAsset: {
    name: string;
    sizeFormatted: string;
    downloadUrl: string;
  } | null;
  releaseUrl: string;
}

const releases = ref<GitHubRelease[]>([]);
const isLoading = ref(false);
const error = ref<string | null>(null);

const REPO_OWNER = "renemakoule";
const REPO_NAME = "ElevenDrive";

const formatBytes = (bytes: number): string => {
  if (!bytes || bytes <= 0) return "—";
  const mb = bytes / (1024 * 1024);
  return `${mb.toFixed(1)} Mo`;
};

const formatDate = (dateStr: string): string => {
  if (!dateStr) return "—";
  try {
    const d = new Date(dateStr);
    return new Intl.DateTimeFormat("fr-FR", {
      day: "numeric",
      month: "short",
      year: "numeric",
    }).format(d);
  } catch {
    return dateStr;
  }
};

const fetchReleases = async () => {
  isLoading.value = true;
  error.value = null;

  try {
    const res = await fetch(
      `https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases?per_page=10`,
      {
        headers: {
          Accept: "application/vnd.github.v3+json",
        },
      }
    );

    if (!res.ok) {
      throw new Error(`Code HTTP ${res.status}`);
    }

    const data = await res.json();
    releases.value = Array.isArray(data) ? data : [];
  } catch (err: any) {
    error.value = err?.message || "Impossible de récupérer les versions";
  } finally {
    isLoading.value = false;
  }
};

const formattedVersions = computed<FormattedVersion[]>(() => {
  if (releases.value.length === 0) {
    // Version par défaut de secours si pas encore de release publiée sur GitHub
    return [
      {
        version: "v1.1.0",
        title: "ElevenDrive Desktop v1.1.0",
        date: "Novembre 2026",
        isLatest: true,
        isPrerelease: false,
        exeAsset: {
          name: "ElevenDrive_1.1.0_x64-setup.exe",
          sizeFormatted: "~8.2 Mo",
          downloadUrl: `https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/latest/download/ElevenDrive.exe`,
        },
        releaseUrl: `https://github.com/${REPO_OWNER}/${REPO_NAME}/releases`,
      },
    ];
  }

  return releases.value.map((rel, index) => {
    // Recherche de l'asset .exe (setup d'installation Windows)
    const exe = rel.assets.find(
      (a) =>
        a.name.toLowerCase().endsWith(".exe") &&
        !a.name.toLowerCase().includes("sig")
    );

    return {
      version: rel.tag_name || rel.name,
      title: rel.name || rel.tag_name,
      date: formatDate(rel.published_at),
      isLatest: index === 0,
      isPrerelease: rel.prerelease,
      exeAsset: exe
        ? {
            name: exe.name,
            sizeFormatted: formatBytes(exe.size),
            downloadUrl: exe.browser_download_url,
          }
        : {
            name: "ElevenDrive-setup.exe",
            sizeFormatted: "GitHub Release",
            downloadUrl: rel.html_url,
          },
      releaseUrl: rel.html_url,
    };
  });
});

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === "Escape" && props.isOpen) {
    emit("close");
  }
};

onMounted(() => {
  window.addEventListener("keydown", handleKeyDown);
  fetchReleases();
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeyDown);
});
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5 select-none font-sans"
  >
    <!-- Voile d'ombrage arrière-plan -->
    <div
      class="fixed inset-0 bg-background/80 backdrop-blur-xs transition-opacity"
      @click="emit('close')"
    ></div>

    <!-- Fenêtre modale centrée -->
    <div
      class="relative w-full max-w-3xl rounded-[6px] border border-border bg-card p-5 sm:p-6 shadow-md text-foreground flex flex-col max-h-[90vh] overflow-hidden"
    >
      <!-- En-tête de la modale -->
      <div
        class="flex items-center justify-between pb-3.5 border-b border-border"
      >
        <div class="flex items-center gap-2.5">
          <div
            class="flex h-7 w-7 items-center justify-center rounded-[4px] border border-border bg-muted/40 p-1 shrink-0"
          >
            <img
              :src="appLogo"
              alt="ElevenDrive Logo"
              class="h-full w-full object-contain"
            />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h2 class="text-xs sm:text-sm font-semibold tracking-tight text-foreground">
                Télécharger ElevenDrive Desktop
              </h2>
              <span
                class="px-1.5 py-0.2 rounded-[3px] text-[10px] font-mono font-medium bg-muted text-muted-foreground border border-border"
              >
                Windows (x64)
              </span>
            </div>
            <p class="text-[11px] text-muted-foreground">
              Choisissez la version d'installation pour votre poste de travail
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click="fetchReleases"
            type="button"
            :disabled="isLoading"
            class="h-7 w-7 flex items-center justify-center rounded-[4px] border border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground transition disabled:opacity-50"
            title="Actualiser la liste"
          >
            <RefreshCw
              class="h-3 w-3"
              :class="{ 'animate-spin': isLoading }"
            />
          </button>
          <button
            @click="emit('close')"
            type="button"
            class="h-7 w-7 flex items-center justify-center rounded-[4px] border border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground transition"
            title="Fermer (Échap)"
          >
            <X class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      <!-- Corps défilant -->
      <div class="flex-1 overflow-y-auto py-4">
        <!-- État de chargement -->
        <div
          v-if="isLoading && formattedVersions.length === 0"
          class="flex flex-col items-center justify-center py-12 text-center"
        >
          <RefreshCw class="h-5 w-5 text-brand animate-spin mb-2" />
          <span class="text-xs text-muted-foreground">
            Recherche des versions sur GitHub...
          </span>
        </div>

        <!-- Avertissement erreur non-bloquant -->
        <div
          v-if="error && releases.length === 0"
          class="mb-3 p-2.5 rounded-[4px] border border-border bg-muted/40 text-[11px] text-muted-foreground flex items-center gap-2"
        >
          <AlertCircle class="h-3.5 w-3.5 text-brand shrink-0" />
          <span>
            Affichage de la version stable par défaut (reconnexion GitHub en cours).
          </span>
        </div>

        <!-- VUE DESKTOP : TABLEAU ÉPURÉ (md:table) -->
        <div class="hidden md:block overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b border-border text-[11px] text-muted-foreground font-medium">
                <th class="py-2 px-2.5">Version</th>
                <th class="py-2 px-2.5">Date</th>
                <th class="py-2 px-2.5">Système</th>
                <th class="py-2 px-2.5">Taille</th>
                <th class="py-2 px-2.5 text-right">Téléchargement</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <tr
                v-for="v in formattedVersions"
                :key="v.version"
                class="hover:bg-muted/20 transition-colors"
              >
                <!-- Version + Badge -->
                <td class="py-2.5 px-2.5">
                  <div class="flex items-center gap-1.5">
                    <span class="font-mono font-semibold text-foreground text-xs">
                      {{ v.version }}
                    </span>
                    <span
                      v-if="v.isLatest"
                      class="px-1.5 py-0.2 rounded-[3px] text-[9px] font-medium uppercase tracking-wide bg-brand/15 text-brand border border-brand/30"
                    >
                      Dernière
                    </span>
                    <span
                      v-else-if="v.isPrerelease"
                      class="px-1.5 py-0.2 rounded-[3px] text-[9px] font-medium uppercase tracking-wide bg-muted text-muted-foreground border border-border"
                    >
                      Bêta
                    </span>
                  </div>
                </td>

                <!-- Date -->
                <td class="py-2.5 px-2.5 text-[11px] text-muted-foreground">
                  {{ v.date }}
                </td>

                <!-- Système -->
                <td class="py-2.5 px-2.5">
                  <span class="inline-flex items-center gap-1 text-[11px] text-muted-foreground">
                    <Laptop class="h-3 w-3" />
                    <span>Windows 10 / 11 (x64)</span>
                  </span>
                </td>

                <!-- Taille -->
                <td class="py-2.5 px-2.5 font-mono text-[11px] text-muted-foreground">
                  {{ v.exeAsset?.sizeFormatted || "—" }}
                </td>

                <!-- Bouton action -->
                <td class="py-2.5 px-2.5 text-right">
                  <a
                    v-if="v.exeAsset"
                    :href="v.exeAsset.downloadUrl"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-1.5 h-7 px-3 rounded-[4px] font-medium text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition shadow-xs"
                  >
                    <Download class="h-3 w-3" />
                    <span>Setup (.exe)</span>
                  </a>
                  <a
                    v-else
                    :href="v.releaseUrl"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition underline underline-offset-2"
                  >
                    <span>Voir release</span>
                    <ExternalLink class="h-3 w-3" />
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- VUE MOBILE : CARTES COMPACTES (< md) -->
        <div class="md:hidden space-y-2.5">
          <div
            v-for="v in formattedVersions"
            :key="v.version"
            class="p-3.5 rounded-[6px] border border-border bg-card shadow-xs flex flex-col gap-2.5"
          >
            <!-- Ligne titre & badge -->
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-foreground text-xs">
                  {{ v.version }}
                </span>
                <span
                  v-if="v.isLatest"
                  class="px-1.5 py-0.2 rounded-[3px] text-[9px] font-medium uppercase tracking-wide bg-brand/15 text-brand border border-brand/30"
                >
                  Dernière
                </span>
                <span
                  v-else-if="v.isPrerelease"
                  class="px-1.5 py-0.2 rounded-[3px] text-[9px] font-medium uppercase tracking-wide bg-muted text-muted-foreground border border-border"
                >
                  Bêta
                </span>
              </div>
              <span class="text-[10px] text-muted-foreground">
                {{ v.date }}
              </span>
            </div>

            <!-- Détails techniques -->
            <div class="grid grid-cols-2 gap-2 text-[10px] text-muted-foreground p-2 rounded-[4px] border border-border bg-muted/20">
              <div class="flex items-center gap-1">
                <Laptop class="h-3 w-3 shrink-0" />
                <span>Windows x64</span>
              </div>
              <div class="flex items-center gap-1 font-mono">
                <HardDrive class="h-3 w-3 shrink-0" />
                <span>{{ v.exeAsset?.sizeFormatted || "—" }}</span>
              </div>
            </div>

            <!-- Bouton téléchargement mobile -->
            <a
              v-if="v.exeAsset"
              :href="v.exeAsset.downloadUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="w-full h-8 flex items-center justify-center gap-1.5 rounded-[4px] font-medium text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition shadow-xs"
            >
              <Download class="h-3.5 w-3.5" />
              <span>Télécharger l'installateur</span>
            </a>
          </div>
        </div>
      </div>

      <!-- Pied de modale : Info technique + Lien GitHub -->
      <div
        class="flex flex-col sm:flex-row items-center justify-between gap-2.5 pt-3.5 border-t border-border text-[11px] text-muted-foreground"
      >
        <div class="flex items-center gap-1.5">
          <CheckCircle2 class="h-3.5 w-3.5 text-brand shrink-0" />
          <span>Installateur certifié NSIS avec signature cryptographique</span>
        </div>

        <a
          :href="`https://github.com/${REPO_OWNER}/${REPO_NAME}/releases`"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-1 text-[11px] hover:text-foreground transition underline underline-offset-2"
        >
          <span>Historique complet sur GitHub</span>
          <ExternalLink class="h-3 w-3" />
        </a>
      </div>
    </div>
  </div>
</template>
