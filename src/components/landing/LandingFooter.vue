<script setup lang="ts">
import { ref, computed } from "vue";
import appLogo from "@root/app-icon.png";
import { ArrowUpRight, ShieldCheck, Heart, Download } from "lucide-vue-next";
import { useLanguageStore } from "../../stores/languageStore";
import DownloadModal from "./DownloadModal.vue";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
  (e: "scroll-to", id: string): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const currentYear = new Date().getFullYear();
const isDownloadModalOpen = ref(false);

const footerLinks = computed(() => ({
  solutions: [
    { label: t("Guichet & Caisse"), id: "service-guichet" },
    { label: t("Pédagogie & Livret"), id: "service-pedagogie" },
    { label: t("Local-First & Archivage R2"), id: "service-archivage" },
    { label: t("Cockpit Multi-Agences"), id: "service-multi-agences" },
    { label: t("Grille Tarifaire"), id: "pricing" },
  ],
  plateformes: [
    { label: t("Application Desktop (Windows / macOS)"), href: "#" },
    { label: t("Portail Web PWA Moniteurs"), href: "#" },
    { label: t("Espace Candidats Mobile"), href: "#" },
    { label: t("Imprimantes Thermiques POS (80/58mm)"), href: "#" },
    { label: t("API Cloud de Synchronisation"), href: "#" },
  ],
  conformite: [
    { label: t("Conformité Ministérielle"), href: "#" },
    { label: t("Protection des Données (RGPD / APDP)"), href: "#" },
    { label: t("Chiffrement & Sécurité R2"), href: "#" },
    { label: t("Conditions Générales d'Utilisation"), href: "#" },
    { label: t("Mentions Légales"), href: "#" },
  ],
}));

const handleScroll = (id: string) => {
  emit("scroll-to", id);
};
</script>

<template>
  <footer
    class="border-t border-border bg-card/60 backdrop-blur-xs text-foreground"
  >
    <!-- Contenu principal -->
    <div class="mx-auto max-w-6xl px-5 py-12 sm:px-6 sm:py-16">
      <div class="grid gap-10 sm:grid-cols-2 md:grid-cols-12">
        <!-- Colonne Marque -->
        <div class="md:col-span-4">
          <div class="flex items-center gap-2">
            <img
              :src="appLogo"
              alt=""
              class="h-6 w-6 rounded-[4px] object-contain"
            />
            <span
              class="text-[15px] font-semibold tracking-tight text-foreground"
            >
              ElevenDrive<span class="text-brand">.</span>
            </span>
          </div>

          <p
            class="mt-3 max-w-xs text-[12px] leading-relaxed text-muted-foreground"
          >
            {{
              t(
                "Système d'exploitation et logiciel de gestion tout-en-un pour auto-écoles modernes : guichet, caisse, livrets numériques de conduite et archivage certifié."
              )
            }}
          </p>

          <!-- Statut opérationnel -->
          <div
            class="mt-5 inline-flex items-center gap-2 rounded-[4px] border border-border bg-background px-2.5 py-1 text-[11px] text-muted-foreground"
          >
            <span
              class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"
            />
            <span>{{ t("Tous les systèmes sont opérationnels") }}</span>
          </div>

          <!-- Bouton Télécharger ElevenDrive Desktop -->
          <div class="mt-3.5">
            <button
              type="button"
              class="inline-flex items-center gap-1.5 h-7.5 px-3 rounded-[4px] font-medium text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition shadow-xs"
              @click="isDownloadModalOpen = true"
            >
              <Download class="h-3.5 w-3.5" />
              <span>{{ t("Télécharger le logiciel") }}</span>
            </button>
          </div>
        </div>

        <!-- Colonne Solutions -->
        <div class="md:col-span-3">
          <p
            class="text-[12px] font-semibold uppercase tracking-wider text-foreground"
          >
            {{ t("Solutions") }}
          </p>
          <ul class="mt-3 space-y-2">
            <li v-for="item in footerLinks.solutions" :key="item.label">
              <button
                type="button"
                class="text-[12px] text-muted-foreground transition hover:text-foreground"
                @click="handleScroll(item.id)"
              >
                {{ item.label }}
              </button>
            </li>
          </ul>
        </div>

        <!-- Colonne Plateformes -->
        <div class="md:col-span-3">
          <p
            class="text-[12px] font-semibold uppercase tracking-wider text-foreground"
          >
            {{ t("Plateformes") }}
          </p>
          <ul class="mt-3 space-y-2">
            <li v-for="item in footerLinks.plateformes" :key="item.label">
              <button
                v-if="item.label.includes('Desktop') || item.label.includes('Application')"
                type="button"
                class="text-[12px] text-muted-foreground hover:text-foreground cursor-pointer transition text-left flex items-center gap-1.5"
                @click="isDownloadModalOpen = true"
              >
                <span>{{ item.label }}</span>
                <span class="px-1.5 py-0.2 rounded-[3px] text-[9px] font-mono bg-brand/15 text-brand border border-brand/30">
                  .exe
                </span>
              </button>
              <span
                v-else
                class="text-[12px] text-muted-foreground hover:text-foreground cursor-pointer transition"
              >
                {{ item.label }}
              </span>
            </li>
          </ul>
        </div>

        <!-- Colonne Sécurité & Légal -->
        <div class="md:col-span-2">
          <p
            class="text-[12px] font-semibold uppercase tracking-wider text-foreground"
          >
            {{ t("Sécurité") }}
          </p>
          <ul class="mt-3 space-y-2">
            <li v-for="item in footerLinks.conformite" :key="item.label">
              <span
                class="text-[12px] text-muted-foreground hover:text-foreground cursor-pointer transition"
              >
                {{ item.label }}
              </span>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bas de page & Copyright -->
      <div
        class="mt-12 flex flex-col items-center justify-between gap-4 border-t border-border pt-6 text-[11px] text-muted-foreground sm:flex-row"
      >
        <p>© {{ currentYear }} ElevenDrive · {{ t("Propulsé par GAS.") }}</p>

        <div class="flex items-center gap-4">
          <span class="inline-flex items-center gap-1">
            <ShieldCheck class="h-3.5 w-3.5 text-emerald-600" />
            {{ t("Infrastructure souveraine & locale") }}
          </span>
          <button
            type="button"
            class="hover:text-foreground transition underline underline-offset-2"
            @click="emit('navigate-auth')"
          >
            {{ t("Accès Guichet & Direction") }}
          </button>
        </div>
      </div>
    </div>

    <!-- Panneau de téléchargement des versions Desktop -->
    <DownloadModal
      :is-open="isDownloadModalOpen"
      @close="isDownloadModalOpen = false"
    />
  </footer>
</template>
