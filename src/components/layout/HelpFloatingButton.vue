<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useLanguageStore } from "../../stores/languageStore";
import {
  HelpCircle,
  Phone,
  X,
  MessageSquare,
  BookOpen,
  Keyboard,
  Bot,
  Sparkles,
  ChevronRight,
  ShieldCheck,
  Zap,
} from "lucide-vue-next";

// Import des 4 composants d'aide dédiés
import HelpShortcutsModal from "../help/HelpShortcutsModal.vue";
import HelpUserGuideModal from "../help/HelpUserGuideModal.vue";
import HelpSupportContactModal from "../help/HelpSupportContactModal.vue";
import HelpAiAssistantModal from "../help/HelpAiAssistantModal.vue";

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const isHubOpen = ref(false);

// États d'ouverture des 4 modales dédiées
const isShortcutsOpen = ref(false);
const isGuideOpen = ref(false);
const isSupportOpen = ref(false);
const isAiAssistantOpen = ref(false);

const openOption = (option: "shortcuts" | "guide" | "support" | "ai") => {
  isHubOpen.value = false;
  if (option === "shortcuts") isShortcutsOpen.value = true;
  else if (option === "guide") isGuideOpen.value = true;
  else if (option === "support") isSupportOpen.value = true;
  else if (option === "ai") isAiAssistantOpen.value = true;
};

const handleOpenHelpEvent = () => {
  isHubOpen.value = true;
};

onMounted(() => {
  window.addEventListener("elevendrive:open-help", handleOpenHelpEvent);
});

onUnmounted(() => {
  window.removeEventListener("elevendrive:open-help", handleOpenHelpEvent);
});
</script>

<template>
  <div>
    <!-- Bouton Flottant (Desktop uniquement : intégré au Menu Drawer sur mobile) -->
    <button
      type="button"
      @click="isHubOpen = !isHubOpen"
      class="hidden md:flex fixed bottom-4 right-3 z-50 bg-brand-dark hover:bg-brand text-white font-medium text-xs px-3.5 py-2 rounded-[3px] items-center space-x-2 transition-all active:scale-95 border border-brand cursor-pointer"
      :title="t('Ouvrir le centre d’aide et d’assistance')"
    >
      <HelpCircle class="w-4 h-4 text-white shrink-0" />
      <span>{{ t("Besoin d'aide ?") }}</span>
    </button>

    <!-- Modale Hub Central : Choix des 4 options d'assistance -->
    <div
      v-if="isHubOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 select-none animate-in fade-in duration-150"
      @click.self="isHubOpen = false"
    >
      <div
        class="bg-card border border-border rounded-[3px] max-w-md w-full overflow-hidden text-xs flex flex-col"
      >
        <!-- En-tête du Hub -->
        <div
          class="p-4 bg-muted/50 border-b border-border flex items-center justify-between"
        >
          <div class="flex items-center space-x-2.5">
            <div
              class="w-8 h-8 rounded-[4px] bg-brand-light text-brand-dark flex items-center justify-center font-bold"
            >
              <HelpCircle class="w-4 h-4 text-brand" />
            </div>
            <div>
              <div
                class="font-bold text-foreground text-sm flex items-center gap-1.5"
              >
                <span>{{ t("Centre d'Aide & Assistance") }}</span>
                <span
                  class="text-[9px] px-1.5 py-0.2 rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 font-semibold"
                >
                  {{ t("En ligne") }}
                </span>
              </div>
              <div class="text-[10px] text-muted-foreground">
                {{
                  t("Support technique, guides, raccourcis et IA ElevenDrive")
                }}
              </div>
            </div>
          </div>
          <button
            type="button"
            @click="isHubOpen = false"
            class="p-1 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition cursor-pointer"
            :title="t('Fermer')"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- 4 Cartes d'Options Interactives -->
        <div class="p-4 space-y-2.5 overflow-y-auto max-h-[70vh]">
          <!-- Option 1 : Assistant ElevenDrive IA -->
          <button
            type="button"
            @click="openOption('ai')"
            class="w-full text-left p-3 rounded-[4px] border border-brand-500/40 bg-brand-500/5 hover:bg-brand-500/10 hover:border-brand-500/70 transition flex items-center space-x-3 group cursor-pointer"
          >
            <div
              class="w-9 h-9 rounded-[4px] bg-brand text-white flex items-center justify-center shrink-0 relative"
            >
              <Bot class="w-5 h-5" />
              <Sparkles
                class="w-2.5 h-2.5 text-amber-300 absolute -top-1 -right-1"
              />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <span
                  class="font-bold text-xs text-foreground group-hover:text-brand"
                >
                  {{ t("Assistant ElevenDrive IA") }}
                </span>
                <span
                  class="text-[9px] px-1.5 py-0.2 rounded bg-brand text-white font-semibold"
                >
                  {{ t("IA") }}
                </span>
              </div>
              <p class="text-[11px] text-muted-foreground mt-0.5 line-clamp-1">
                {{
                  t(
                    "Posez toutes vos questions sur la gestion et la réglementation.",
                  )
                }}
              </p>
            </div>
            <ChevronRight
              class="w-4 h-4 text-muted-foreground group-hover:text-foreground shrink-0 transition"
            />
          </button>

          <!-- Option 2 : Raccourcis Clavier -->
          <button
            type="button"
            @click="openOption('shortcuts')"
            class="w-full text-left p-3 rounded-[4px] border border-border hover:border-brand-500/50 bg-card hover:bg-muted/40 transition flex items-center space-x-3 group cursor-pointer"
          >
            <div
              class="w-9 h-9 rounded-[4px] bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/20 flex items-center justify-center shrink-0"
            >
              <Keyboard class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <span
                  class="font-bold text-xs text-foreground group-hover:text-brand"
                >
                  {{ t("Raccourcis Clavier") }}
                </span>
                <span
                  class="text-[9px] font-mono px-1.5 py-0.2 rounded bg-muted text-muted-foreground border border-border font-bold"
                >
                  F1 - F9
                </span>
              </div>
              <p class="text-[11px] text-muted-foreground mt-0.5 line-clamp-1">
                {{
                  t(
                    "Accès ultra-rapide aux fonctions guichet, encaissement et caisse.",
                  )
                }}
              </p>
            </div>
            <ChevronRight
              class="w-4 h-4 text-muted-foreground group-hover:text-foreground shrink-0 transition"
            />
          </button>

          <!-- Option 3 : Guide d'Utilisation -->
          <button
            type="button"
            @click="openOption('guide')"
            class="w-full text-left p-3 rounded-[4px] border border-border hover:border-brand-500/50 bg-card hover:bg-muted/40 transition flex items-center space-x-3 group cursor-pointer"
          >
            <div
              class="w-9 h-9 rounded-lg bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/20 flex items-center justify-center shrink-0"
            >
              <BookOpen class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <span
                  class="font-bold text-xs text-foreground group-hover:text-brand"
                >
                  {{ t("Guide d'Utilisation") }}
                </span>
                <span
                  class="text-[9px] px-1.5 py-0.2 rounded bg-muted text-muted-foreground font-medium"
                >
                  {{ t("6 étapes") }}
                </span>
              </div>
              <p class="text-[11px] text-muted-foreground mt-0.5 line-clamp-1">
                {{
                  t(
                    "Manuel interactif : inscriptions, planning, examens et caisse.",
                  )
                }}
              </p>
            </div>
            <ChevronRight
              class="w-4 h-4 text-muted-foreground group-hover:text-foreground shrink-0 transition"
            />
          </button>

          <!-- Option 4 : Contacter l'Équipe Support -->
          <button
            type="button"
            @click="openOption('support')"
            class="w-full text-left p-3 rounded-[4px] border border-border hover:border-brand-500/50 bg-card hover:bg-muted/40 transition flex items-center space-x-3 group cursor-pointer"
          >
            <div
              class="w-9 h-9 rounded-lg bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 flex items-center justify-center shrink-0"
            >
              <Phone class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <span
                  class="font-bold text-xs text-foreground group-hover:text-brand"
                >
                  {{ t("Contacter l'Équipe Support") }}
                </span>
                <span
                  class="text-[9px] px-1.5 py-0.2 rounded bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-semibold"
                >
                  {{ t("Direct") }}
                </span>
              </div>
              <p class="text-[11px] text-muted-foreground mt-0.5 line-clamp-1">
                {{
                  t(
                    "Assistance WhatsApp, hotline téléphonique et tickets prioritaires.",
                  )
                }}
              </p>
            </div>
            <ChevronRight
              class="w-4 h-4 text-muted-foreground group-hover:text-foreground shrink-0 transition"
            />
          </button>
        </div>

        <!-- Pied du Hub -->
        <div
          class="p-3 bg-muted/40 border-t border-border flex items-center justify-between"
        >
          <a
            href="https://wa.me/237651727932?text=Bonjour%20Support%20ElevenDrive,%20j'ai%20besoin%20d'aide%20imm%C3%A9diate."
            target="_blank"
            class="inline-flex items-center gap-1.5 text-brand font-semibold text-[11px] hover:underline"
          >
            <MessageSquare class="w-3.5 h-3.5" />
            <span>{{ t("Urgence WhatsApp Directe") }}</span>
          </a>
          <button
            type="button"
            @click="isHubOpen = false"
            class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground font-medium text-xs transition cursor-pointer"
          >
            {{ t("Fermer") }}
          </button>
        </div>
      </div>
    </div>

    <!-- 4 Modales Spécifiques (Chaque option a son composant unique) -->
    <HelpShortcutsModal
      :is-open="isShortcutsOpen"
      @close="isShortcutsOpen = false"
    />

    <HelpUserGuideModal :is-open="isGuideOpen" @close="isGuideOpen = false" />

    <HelpSupportContactModal
      :is-open="isSupportOpen"
      @close="isSupportOpen = false"
    />

    <HelpAiAssistantModal
      :is-open="isAiAssistantOpen"
      @close="isAiAssistantOpen = false"
    />
  </div>
</template>
