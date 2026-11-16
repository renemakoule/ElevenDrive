<script setup lang="ts">
import { ref, nextTick } from "vue";
import { useLanguageStore } from "../../stores/languageStore";
import { useAuthStore } from "../../stores/authStore";
import {
  Bot,
  Sparkles,
  Send,
  X,
  User,
  RotateCcw,
  Lightbulb,
  CheckCircle2,
} from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const languageStore = useLanguageStore();
const authStore = useAuthStore();
const t = (text: string) => languageStore.t(text);

interface Message {
  id: string;
  sender: "user" | "ai";
  text: string;
  time: string;
}

const inputPrompt = ref("");
const isTyping = ref(false);
const chatContainer = ref<HTMLElement | null>(null);

const messages = ref<Message[]>([
  {
    id: "welcome",
    sender: "ai",
    text: `Bonjour ${
      authStore.currentUser?.prenom || authStore.currentUser?.nom || "Directeur"
    } ! Je suis l'assistant intelligent ElevenDrive IA 🤖. Comment puis-je vous aider aujourd'hui dans l'exploitation de votre auto-école ?`,
    time: "À l'instant",
  },
]);

const quickPrompts = [
  "Comment clôturer la caisse avec le Rapport Z ?",
  "Un élève peut-il passer l'examen s'il n'a pas tout payé ?",
  "Comment réimprimer un reçu thermique perdu ?",
  "Que faire si la connexion Internet est coupée ?",
];

const scrollToBottom = async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

const sendPrompt = (customText?: string) => {
  const text = customText || inputPrompt.value.trim();
  if (!text) return;

  const now = new Date();
  const timeStr = `${String(now.getHours()).padStart(2, "0")}:${String(
    now.getMinutes(),
  ).padStart(2, "0")}`;

  messages.value.push({
    id: `msg-${Date.now()}`,
    sender: "user",
    text,
    time: timeStr,
  });

  inputPrompt.value = "";
  isTyping.value = true;
  scrollToBottom();

  // Générer une réponse intelligente contextuelle ElevenDrive
  setTimeout(() => {
    let reponse = "";
    const lower = text.toLowerCase();

    if (
      lower.includes("clôturer") ||
      lower.includes("caisse") ||
      lower.includes("rapport z")
    ) {
      reponse = `Pour clôturer votre caisse :
1. Appuyez sur la touche raccourci **F4** ou ouvrez l'onglet **Point Caisse (Rapport Z)**.
2. Saisissez le comptage physique de vos espèces réelles en caisse.
3. Le système réconcilie les encaissements du jour et calcule les écarts.
4. Cliquez sur **Valider la clôture** : le ticket thermique Z est imprimé et la session est archivée dans le journal d'audit de la Direction.`;
    } else if (
      lower.includes("examen") ||
      lower.includes("payer") ||
      lower.includes("solde")
    ) {
      reponse = `Concernant les examens et les impayés :
- Par défaut, **seuls les candidats ayant soldé 100% de leurs frais** sont automatiquement injectés dans la liste officielle d'examen.
- Les candidats avec un solde restant sont mis en **liste d'attente**.
- **Forçage manuel :** Le Directeur dispose d'un bouton exclusif **'Forcer l'inscription officielle'** permettant de transférer un élève de la liste d'attente vers la liste officielle avec motif tracé.
- Dès que l'élève solde ses frais, il est automatiquement rebasculé dans la session active suivante.`;
    } else if (
      lower.includes("reçu") ||
      lower.includes("imprimer") ||
      lower.includes("perdu") ||
      lower.includes("ticket")
    ) {
      reponse = `Pour réimprimer un reçu thermique :
1. Rendez-vous dans **Guichet & Encaissement** ou **Historique des Opérations**.
2. Recherchez l'élève par son nom ou son matricule (touche **F3**).
3. Cliquez sur l'icône **Imprimante** à côté de la transaction concernée.
4. Le ticket ESC/POS 80mm est réémis avec la mention sécurisée *'Duplicata certifié'*.`;
    } else if (
      lower.includes("internet") ||
      lower.includes("hors-ligne") ||
      lower.includes("connexion") ||
      lower.includes("offline")
    ) {
      reponse = `ElevenDrive fonctionne en mode **100% Hors-Ligne (Offline First)** :
- Vous pouvez continuer à inscrire des élèves et encaisser les paiements sans interruption.
- Les données sont sécurisées dans la base locale chiffrée de votre appareil.
- Dès que le réseau réapparaît, la synchronisation avec le cloud **Neon Postgres** s'effectue automatiquement en arrière-plan sans perte de données.`;
    } else {
      reponse = `J'ai bien pris en compte votre question : "${text}".
ElevenDrive automatise la gestion de votre auto-école :
- **Guichet & Caisse :** Utilisez les touches F1 à F4 pour les opérations rapides.
- **Pédagogie :** Planifiez vos conduites et présences en salle dans l'onglet Planning.
- **Direction :** Suivez vos marges, flotte et taux de réussite en temps réel.
Besoin d'une assistance humaine ? Vous pouvez également contacter notre support direct par WhatsApp ou Hotline !`;
    }

    isTyping.value = false;
    messages.value.push({
      id: `ai-${Date.now()}`,
      sender: "ai",
      text: reponse,
      time: timeStr,
    });
    scrollToBottom();
  }, 700);
};

const resetChat = () => {
  messages.value = [
    {
      id: "welcome",
      sender: "ai",
      text: `Conversation réinitialisée. En quoi puis-je vous guider ? 💡`,
      time: "À l'instant",
    },
  ];
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-xl w-full overflow-hidden flex flex-col h-[85vh]"
    >
      <!-- En-tête -->
      <div
        class="p-4 bg-muted/40 border-b border-border flex items-center justify-between shrink-0"
      >
        <div class="flex items-center space-x-2.5">
          <div
            class="w-8 h-8 rounded-[3px] bg-brand-light text-brand-dark flex items-center justify-center font-bold relative"
          >
            <Bot class="w-4 h-4 text-brand" />
            <span class="absolute -top-1 -right-1 flex h-2 w-2">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
              ></span>
              <span
                class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"
              ></span>
            </span>
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>{{ t("Assistant ElevenDrive IA") }}</span>
              <span
                class="text-[10px] px-1.5 py-0.2 rounded bg-brand text-white font-medium flex items-center gap-1 shadow-2xs"
              >
                <Sparkles class="w-2.5 h-2.5" />
                <span>{{ t("IA Intégrée") }}</span>
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              {{
                t(
                  "Conseiller virtuel spécialisé dans la réglementation et la gestion auto-école",
                )
              }}
            </div>
          </div>
        </div>
        <div class="flex items-center space-x-1">
          <button
            type="button"
            @click="resetChat"
            class="p-1.5 rounded-[3px] text-muted-foreground hover:text-foreground hover:bg-muted transition"
            :title="t('Réinitialiser la conversation')"
          >
            <RotateCcw class="w-3.5 h-3.5" />
          </button>
          <button
            type="button"
            @click="emit('close')"
            class="p-1.5 rounded-[3px] text-muted-foreground hover:text-foreground hover:bg-muted transition"
            :title="t('Fermer')"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Zone de Messages Chat -->
      <div
        ref="chatContainer"
        class="flex-1 p-4 overflow-y-auto space-y-3.5 min-h-0 bg-background/50"
      >
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="flex items-start space-x-2.5"
          :class="
            msg.sender === 'user' ? 'flex-row-reverse space-x-reverse' : ''
          "
        >
          <!-- Avatar -->
          <div
            class="w-7 h-7 rounded-full flex items-center justify-center shrink-0 text-xs font-bold"
            :class="
              msg.sender === 'user'
                ? 'bg-muted text-foreground border border-border'
                : 'bg-brand text-white'
            "
          >
            <User v-if="msg.sender === 'user'" class="w-3.5 h-3.5" />
            <Bot v-else class="w-3.5 h-3.5" />
          </div>

          <!-- Bulle de message -->
          <div
            class="max-w-[82%] rounded-[4px] px-3.5 py-2.5 text-xs"
            :class="
              msg.sender === 'user'
                ? 'bg-brand text-white rounded-tr-[3px]'
                : 'bg-card border border-border text-foreground rounded-tl-[3px]'
            "
          >
            <div class="whitespace-pre-line leading-relaxed">
              {{ msg.text }}
            </div>
            <div
              class="text-[9px] mt-1 text-right"
              :class="
                msg.sender === 'user'
                  ? 'text-white/70'
                  : 'text-muted-foreground'
              "
            >
              {{ msg.time }}
            </div>
          </div>
        </div>

        <!-- Indicateur de saisie IA -->
        <div
          v-if="isTyping"
          class="flex items-center space-x-2 text-muted-foreground text-xs italic"
        >
          <Bot class="w-3.5 h-3.5 text-brand animate-pulse" />
          <span>{{ t("ElevenDrive IA analyse votre demande...") }}</span>
        </div>
      </div>

      <!-- Suggestions de questions rapides -->
      <div
        class="px-4 py-2 border-t border-border/50 bg-muted/20 overflow-x-auto shrink-0 flex items-center gap-1.5 no-scrollbar"
      >
        <span
          class="text-[10px] font-semibold text-muted-foreground shrink-0 flex items-center gap-1"
        >
          <Lightbulb class="w-3 h-3 text-amber-500" />
          <span>{{ t("Suggestions :") }}</span>
        </span>
        <button
          v-for="(p, idx) in quickPrompts"
          :key="idx"
          type="button"
          @click="sendPrompt(p)"
          class="px-2.5 py-1 rounded-full bg-card hover:bg-muted text-[11px] text-muted-foreground hover:text-foreground border border-border shrink-0 transition"
        >
          {{ p }}
        </button>
      </div>

      <!-- Barre de saisie -->
      <form
        @submit.prevent="sendPrompt()"
        class="p-3 bg-card border-t border-border flex items-center space-x-2 shrink-0"
      >
        <input
          v-model="inputPrompt"
          type="text"
          :placeholder="
            t(
              'Posez une question à ElevenDrive IA (ex: règlement, impayés, matériel)...',
            )
          "
          class="flex-1 px-3.5 py-2 text-xs rounded-[3px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
        />
        <button
          type="submit"
          :disabled="!inputPrompt.trim()"
          class="p-2 rounded-[3px] bg-brand hover:bg-brand-dark text-white disabled:opacity-40 disabled:cursor-not-allowed transition shrink-0 cursor-pointer"
        >
          <Send class="w-4 h-4" />
        </button>
      </form>
    </div>
  </div>
</template>
