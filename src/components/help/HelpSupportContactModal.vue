<script setup lang="ts">
import { ref } from "vue";
import { useLanguageStore } from "../../stores/languageStore";
import { useAuthStore } from "../../stores/authStore";
import { supportService } from "../../services/supportService";
import {
  Phone,
  MessageSquare,
  Mail,
  X,
  Send,
  CheckCircle2,
  Clock,
  ShieldCheck,
  AlertCircle,
  HelpCircle,
  Printer,
  Wifi,
  Loader2,
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

const isSent = ref(false);
const isSubmitting = ref(false);
const errorMessage = ref<string | null>(null);

const form = ref({
  nom: authStore.currentUser
    ? `${authStore.currentUser.prenom || ""} ${authStore.currentUser.nom || ""}`.trim()
    : "",
  email: authStore.currentUser?.email || "",
  sujet: "probleme_technique",
  priorite: "normale",
  message: "",
});

const envoyerDemande = async () => {
  if (!form.value.message.trim()) return;

  isSubmitting.value = true;
  errorMessage.value = null;

  try {
    await supportService.createTicket({
      auto_ecole_id: authStore.currentUser?.auto_ecole_id || undefined,
      nom: form.value.nom.trim() || "Utilisateur ElevenDrive",
      email: form.value.email.trim() || "contact@autoschool.local",
      sujet: form.value.sujet,
      priorite: form.value.priorite,
      message: form.value.message.trim(),
    });
    isSent.value = true;
  } catch (err: any) {
    errorMessage.value = err.message || t("Échec de la transmission du ticket");
  } finally {
    isSubmitting.value = false;
  }
};

const reinitialiser = () => {
  isSent.value = false;
  form.value.message = "";
  errorMessage.value = null;
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 select-none animate-in fade-in duration-150"
    @click.self="emit('close')"
  >
    <div
      class="bg-card border border-border rounded-[3px] max-w-xl w-full overflow-hidden flex flex-col max-h-[90vh]"
    >
      <!-- En-tête -->
      <div
        class="p-4 bg-muted/40 border-b border-border flex items-center justify-between shrink-0"
      >
        <div class="flex items-center space-x-2.5">
          <div
            class="w-8 h-8 rounded-[4px] bg-brand-light text-brand-dark flex items-center justify-center font-bold"
          >
            <Phone class="w-4 h-4 text-brand" />
          </div>
          <div>
            <div
              class="font-bold text-foreground text-sm flex items-center gap-2"
            >
              <span>{{ t("Contacter l'Équipe Support") }}</span>
              <span
                class="text-[10px] px-1.5 py-0.2 rounded bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-medium border border-emerald-500/20"
              >
                {{ t("Support 7j/7") }}
              </span>
            </div>
            <div class="text-[11px] text-muted-foreground">
              {{
                t(
                  "Une question, un blocage ou une urgence matérielle ? Nous intervenons sans délai.",
                )
              }}
            </div>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="p-1 rounded-[4px] text-muted-foreground hover:text-foreground hover:bg-muted transition"
          :title="t('Fermer')"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Corps du formulaire / Canaux directs -->
      <div class="p-4 overflow-y-auto space-y-4 flex-1 min-h-0 text-xs">
        <!-- Canaux Rapides WhatsApp & Téléphone -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- WhatsApp -->
          <a
            href="https://wa.me/237651727932?text=Bonjour%20Support%20ElevenDrive,%20j'ai%20besoin%20d'aide."
            target="_blank"
            class="p-3 rounded-[4px] bg-emerald-500/10 border border-emerald-500/30 hover:bg-emerald-500/20 text-emerald-800 dark:text-emerald-300 transition flex items-center space-x-3 group"
          >
            <div
              class="w-9 h-9 rounded-[4px] bg-emerald-600 text-white flex items-center justify-center shrink-0"
            >
              <MessageSquare class="w-5 h-5" />
            </div>
            <div class="min-w-0">
              <div
                class="font-bold text-xs text-foreground group-hover:text-emerald-600"
              >
                {{ t("WhatsApp Support") }}
              </div>
              <div class="text-[10px] text-muted-foreground truncate">
                {{ t("Réponse en moins de 10 min") }}
              </div>
            </div>
          </a>

          <!-- Hotline Téléphone -->
          <a
            href="tel:+2250700000000"
            class="p-3 rounded-[4px] bg-blue-500/10 border border-blue-500/30 hover:bg-blue-500/20 text-blue-800 dark:text-blue-300 transition flex items-center space-x-3 group"
          >
            <div
              class="w-9 h-9 rounded-[4px] bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-xs"
            >
              <Phone class="w-5 h-5" />
            </div>
            <div class="min-w-0">
              <div
                class="font-bold text-xs text-foreground group-hover:text-blue-600"
              >
                {{ t("Hotline Téléphonique") }}
              </div>
              <div class="text-[10px] text-muted-foreground truncate">
                {{ t("+225 07 00 00 00 00 (8h - 19h)") }}
              </div>
            </div>
          </a>
        </div>

        <!-- Formulaire de Ticket Écrit -->
        <div
          class="border border-border rounded-[3px] p-4 bg-muted/20 space-y-3"
        >
          <div
            class="flex items-center justify-between border-b border-border/50 pb-2"
          >
            <span class="font-bold text-foreground text-xs">
              {{ t("Ouvrir un Ticket d'Intervention") }}
            </span>
            <span
              class="text-[10px] text-muted-foreground flex items-center gap-1"
            >
              <Clock class="w-3 h-3 text-muted-foreground" />
              <span>{{ t("Traitement prioritaire") }}</span>
            </span>
          </div>

          <div v-if="isSent" class="py-6 text-center space-y-2">
            <CheckCircle2 class="w-10 h-10 text-emerald-500 mx-auto" />
            <div class="font-bold text-foreground text-sm">
              {{ t("Ticket transmis avec succès !") }}
            </div>
            <p class="text-[11px] text-muted-foreground max-w-sm mx-auto">
              {{
                t(
                  "Notre équipe technique a reçu votre signalement et prend contact avec votre établissement dans les plus brefs délais.",
                )
              }}
            </p>
            <button
              type="button"
              @click="reinitialiser"
              class="mt-2 px-3 py-1.5 rounded-[4px] bg-muted text-foreground hover:bg-muted/80 text-xs font-medium transition"
            >
              {{ t("Envoyer un autre message") }}
            </button>
          </div>

          <form v-else @submit.prevent="envoyerDemande" class="space-y-3">
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label
                  class="block text-[11px] font-medium text-muted-foreground mb-1"
                >
                  {{ t("Type d'intervention") }}
                </label>
                <select
                  v-model="form.sujet"
                  class="w-full text-xs p-2 rounded-[4px] border border-border bg-background"
                >
                  <option value="probleme_technique">
                    {{ t("Panne logicielle / Bug") }}
                  </option>
                  <option value="imprimante">
                    {{ t("Imprimante thermique / Ticket ESC/POS") }}
                  </option>
                  <option value="synchronisation">
                    {{ t("Synchronisation Neon Cloud / Réseau") }}
                  </option>
                  <option value="pedagogie">
                    {{ t("Question réglementation / Tarifs") }}
                  </option>
                  <option value="autre">{{ t("Autre demande") }}</option>
                </select>
              </div>

              <div>
                <label
                  class="block text-[11px] font-medium text-muted-foreground mb-1"
                >
                  {{ t("Degré d'urgence") }}
                </label>
                <select
                  v-model="form.priorite"
                  class="w-full text-xs p-2 rounded-[4px] border border-border bg-background"
                >
                  <option value="normale">
                    {{ t("Normale (dans la journée)") }}
                  </option>
                  <option value="urgente">
                    {{ t("Urgente (dans l'heure)") }}
                  </option>
                  <option value="bloquante">
                    {{ t("Bloquante (Guichet arrêté)") }}
                  </option>
                </select>
              </div>
            </div>

            <div>
              <label
                class="block text-[11px] font-medium text-muted-foreground mb-1"
              >
                {{ t("Description détaillée du problème") }}
              </label>
              <textarea
                v-model="form.message"
                rows="3"
                required
                :placeholder="
                  t(
                    'Décrivez la situation, le numéro de reçu ou le problème rencontré...',
                  )
                "
                class="w-full p-2.5 text-xs rounded-[4px] border border-border bg-background focus:outline-hidden focus:ring-1 focus:ring-brand"
              ></textarea>
            </div>

            <!-- Alerte erreur éventuelle -->
            <div
              v-if="errorMessage"
              class="p-2.5 rounded-[4px] bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-400 text-xs flex items-center gap-2"
            >
              <AlertCircle class="w-4 h-4 shrink-0" />
              <span>{{ errorMessage }}</span>
            </div>

            <div class="flex items-center justify-between pt-1">
              <span
                class="text-[10px] text-muted-foreground flex items-center gap-1"
              >
                <ShieldCheck class="w-3.5 h-3.5 text-emerald-500" />
                <span>{{ t("Télédiagnostic ElevenDrive sécurisé") }}</span>
              </span>
              <button
                type="submit"
                :disabled="isSubmitting"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-[4px] bg-brand hover:bg-brand-dark text-white font-semibold text-xs shadow-xs transition active:scale-95 cursor-pointer disabled:opacity-50"
              >
                <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
                <Send v-else class="w-3.5 h-3.5" />
                <span>{{
                  isSubmitting
                    ? t("Envoi en cours...")
                    : t("Transmettre au Support")
                }}</span>
              </button>
            </div>
          </form>
        </div>

        <!-- Diagnostic Rapide Matériel -->
        <div
          class="p-2.5 rounded-[4px] bg-muted/40 border border-border/70 flex items-center justify-between text-[11px]"
        >
          <div class="flex items-center space-x-2 text-muted-foreground">
            <Wifi class="w-3.5 h-3.5 text-emerald-500" />
            <span>{{
              t("Serveur Neon Cloud : Opérationnel (Latence 28ms)")
            }}</span>
          </div>
          <div class="flex items-center space-x-2 text-muted-foreground">
            <Printer class="w-3.5 h-3.5 text-brand" />
            <span>{{ t("Pilote ESC/POS : Prêt") }}</span>
          </div>
        </div>
      </div>

      <!-- Pied de page -->
      <div
        class="p-3 bg-muted/40 border-t border-border flex items-center justify-between text-xs text-muted-foreground shrink-0"
      >
        <span class="text-[11px]">
          {{ t("Email : support@elevendrive.com") }}
        </span>
        <button
          type="button"
          @click="emit('close')"
          class="px-3 py-1.5 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground font-medium text-xs transition"
        >
          {{ t("Fermer") }}
        </button>
      </div>
    </div>
  </div>
</template>
