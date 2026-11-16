<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useCoursCodeStore } from "../../stores/coursCodeStore";
import { useAuthStore } from "../../stores/authStore";
import { api } from "../../db/api";
import { JITSI_BASE_URL } from "../../config/env";
import {
  X,
  Video,
  Calendar,
  Clock,
  User,
  Link,
  Sparkles,
} from "lucide-vue-next";

import type { SeanceCode } from "../../types";

const props = defineProps<{
  isOpen: boolean;
  seanceToReprogram?: SeanceCode | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const coursStore = useCoursCodeStore();
const authStore = useAuthStore();

const THEMES_CODE_POPULAIRES = [
  "Signalisation Routière & Panneaux",
  "Règles de Priorité & Intersections",
  "Croisements, Dépassements & Voies Rapides",
  "Arrêt, Stationnement & Pannes",
  "Visibilité, Éclairage & Intempéries",
  "Vitesse, Distances de Sécurité & Freinage",
  "Alcool, Vigilance & Sécurité Routière",
  "Éco-conduite & Notions Mécaniques",
];

interface MoniteurOption {
  id: string;
  nomComplet: string;
  role: string;
  telephone?: string;
}

const moniteursDisponibles = ref<MoniteurOption[]>([]);
const selectedMoniteurId = ref<string>("");
const isCustomMoniteur = ref(false);
const isLoadingMoniteurs = ref(false);

const formTheme = ref(THEMES_CODE_POPULAIRES[0]);
const formType = ref<"PRESENTIEL" | "EN_LIGNE">("EN_LIGNE");
const formDate = ref(new Date().toISOString().slice(0, 10));
const formHeureDebut = ref("18:30");
const formHeureFin = ref("20:30");
const formFormateur = ref(
  authStore.currentUser?.nom
    ? `${authStore.currentUser.prenom} ${authStore.currentUser.nom}`
    : "Moniteur Référent",
);
const formLienVisio = ref("");
const formModeLien = ref<"JITSI" | "CUSTOM">("JITSI");

const chargerMoniteurs = async () => {
  isLoadingMoniteurs.value = true;
  try {
    const list = await api.getMoniteurs();
    moniteursDisponibles.value = list;

    if (list.length > 0) {
      // Pré-sélectionner l'utilisateur connecté s'il fait partie de l'équipe
      const current = list.find((m) => m.id === authStore.currentUser?.id);
      if (current) {
        selectedMoniteurId.value = current.id;
        formFormateur.value = current.nomComplet;
      } else {
        selectedMoniteurId.value = list[0].id;
        formFormateur.value = list[0].nomComplet;
      }
    } else {
      isCustomMoniteur.value = true;
    }
  } catch (err) {
    console.warn("Erreur chargement moniteurs:", err);
    isCustomMoniteur.value = true;
  } finally {
    isLoadingMoniteurs.value = false;
  }
};

onMounted(() => {
  chargerMoniteurs();
});

watch(selectedMoniteurId, (newId) => {
  if (newId === "__AUTRE__") {
    isCustomMoniteur.value = true;
    formFormateur.value = "";
  } else {
    isCustomMoniteur.value = false;
    const found = moniteursDisponibles.value.find((m) => m.id === newId);
    if (found) {
      formFormateur.value = found.nomComplet;
    }
  }
});

const genererSalonJitsi = () => {
  const roomName =
    "ElevenDrive-Code-" +
    Math.random().toString(36).substring(2, 8).toUpperCase();
  formLienVisio.value = `${JITSI_BASE_URL}/${roomName}`;
};

// Génération par défaut si mode Jitsi
if (!formLienVisio.value) {
  genererSalonJitsi();
}

watch(
  () => [props.isOpen, props.seanceToReprogram] as const,
  ([open, target]) => {
    if (!open) return;
    if (target) {
      formTheme.value = target.theme_traite || THEMES_CODE_POPULAIRES[0];
      formType.value = target.type_seance || "EN_LIGNE";
      formDate.value = new Date().toISOString().slice(0, 10);
      formHeureDebut.value = target.heure_debut || "18:30";
      formHeureFin.value = target.heure_fin || "20:30";
      formFormateur.value = target.formateur_nom || "";
      formLienVisio.value = target.lien_visio || "";
    }
  },
  { immediate: true },
);

const handleSubmit = async () => {
  if (props.seanceToReprogram) {
    await coursStore.reprogrammerSeance({
      id: props.seanceToReprogram.id,
      date_seance: formDate.value,
      heure_debut: formHeureDebut.value,
      heure_fin: formHeureFin.value,
    });
    emit("close");
    return;
  }

  const lien =
    formType.value === "EN_LIGNE"
      ? formLienVisio.value.trim() ||
        `${JITSI_BASE_URL}/ElevenDrive-${Date.now().toString(36).toUpperCase()}`
      : undefined;

  const moniteurIdToSend =
    selectedMoniteurId.value && selectedMoniteurId.value !== "__AUTRE__"
      ? selectedMoniteurId.value
      : undefined;

  await coursStore.planifierSeance({
    theme_traite: formTheme.value,
    formateur_nom: formFormateur.value.trim() || "Moniteur Référent",
    moniteur_id: moniteurIdToSend,
    date_seance: formDate.value,
    heure_debut: formHeureDebut.value,
    heure_fin: formHeureFin.value,
    type_seance: formType.value,
    lien_visio: lien,
  });

  emit("close");
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-black/70 backdrop-blur-xs animate-in fade-in"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-card border border-border w-full max-w-xl max-h-[92vh] rounded-[5px] shadow-2xl flex flex-col overflow-hidden text-foreground"
    >
      <!-- En-tête -->
      <header
        class="px-5 py-3.5 border-b border-border flex items-center justify-between bg-muted/20"
      >
        <div class="flex items-center space-x-2">
          <div
            class="w-7 h-7 rounded-[4px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center"
          >
            <Video class="w-4 h-4 text-brand" />
          </div>
          <div>
            <h2 class="text-sm font-bold text-foreground">
              {{
                seanceToReprogram
                  ? "Reprogrammer le cours bloqué"
                  : "Planifier une séance de code"
              }}
            </h2>
            <p class="text-[11px] text-muted-foreground">
              {{
                seanceToReprogram
                  ? "Définissez la nouvelle date et les nouveaux horaires pour réactiver le cours"
                  : "Définition du créneau, du thème et de la visioconférence"
              }}
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="$emit('close')"
          class="text-muted-foreground hover:text-foreground cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </header>

      <!-- Formulaire -->
      <form
        @submit.prevent="handleSubmit"
        class="p-5 overflow-y-auto space-y-4 text-xs"
      >
        <!-- Type de séance : En ligne vs Salle -->
        <div>
          <label class="block font-semibold text-foreground mb-1.5"
            >Type de session *</label
          >
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="formType = 'EN_LIGNE'"
              class="flex items-center justify-center gap-2 p-2.5 rounded-[4px] border text-xs font-semibold transition cursor-pointer"
              :class="
                formType === 'EN_LIGNE'
                  ? 'border-brand bg-brand-light text-primary'
                  : 'border-border bg-card text-muted-foreground hover:bg-muted'
              "
            >
              <Video class="w-4 h-4 text-brand" />
              <span>💻 Cours en ligne (Visioconférence)</span>
            </button>

            <button
              type="button"
              @click="formType = 'PRESENTIEL'"
              class="flex items-center justify-center gap-2 p-2.5 rounded-[4px] border text-xs font-semibold transition cursor-pointer"
              :class="
                formType === 'PRESENTIEL'
                  ? 'border-brand bg-brand-light text-primary'
                  : 'border-border bg-card text-muted-foreground hover:bg-muted'
              "
            >
              <span>🏫 En Salle (Présentiel à l'agence)</span>
            </button>
          </div>
        </div>

        <!-- Thème traité -->
        <div>
          <label class="block font-semibold text-foreground mb-1"
            >Thème pédagogique *</label
          >
          <select
            v-model="formTheme"
            required
            class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark font-medium"
          >
            <option v-for="t in THEMES_CODE_POPULAIRES" :key="t" :value="t">
              {{ t }}
            </option>
          </select>
        </div>

        <!-- Date & Horaires -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div>
            <label class="block font-semibold text-foreground mb-1"
              >Date de la séance *</label
            >
            <input
              v-model="formDate"
              type="date"
              required
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>

          <div>
            <label class="block font-semibold text-foreground mb-1"
              >Heure de début *</label
            >
            <input
              v-model="formHeureDebut"
              type="time"
              required
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>

          <div>
            <label class="block font-semibold text-foreground mb-1"
              >Heure de fin *</label
            >
            <input
              v-model="formHeureFin"
              type="time"
              required
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>
        </div>

        <!-- Formateur / Moniteur avec sélection dynamique en base de données -->
        <div>
          <div class="flex items-center justify-between mb-1">
            <label class="block font-semibold text-foreground">
              Moniteur / Formateur référent *
            </label>
            <span
              v-if="isLoadingMoniteurs"
              class="text-[10px] text-muted-foreground animate-pulse"
            >
              Chargement des moniteurs...
            </span>
          </div>

          <div class="space-y-2">
            <!-- Sélecteur des moniteurs enregistrés dans la base de données -->
            <select
              v-if="moniteursDisponibles.length > 0"
              v-model="selectedMoniteurId"
              class="w-full h-8 px-2 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark font-medium"
            >
              <option
                v-for="m in moniteursDisponibles"
                :key="m.id"
                :value="m.id"
              >
                👤 {{ m.nomComplet }} ({{
                  m.role === "MONITEUR" ? "Moniteur" : "Directeur"
                }}) {{ m.telephone ? "· " + m.telephone : "" }}
              </option>
              <option value="__AUTRE__">
                ✏️ Autre intervenant / Saisie manuelle...
              </option>
            </select>

            <!-- Champ texte si pas de moniteur ou si choix "Autre intervenant" -->
            <input
              v-if="isCustomMoniteur || moniteursDisponibles.length === 0"
              v-model="formFormateur"
              type="text"
              required
              placeholder="Précisez le nom du moniteur ou formateur"
              class="w-full h-8 px-2.5 rounded-[4px] border border-border bg-background text-foreground text-xs focus:outline-none focus:border-brand-dark"
            />
          </div>
        </div>

        <!-- Configuration Visioconférence (si mode EN_LIGNE) -->
        <div
          v-if="formType === 'EN_LIGNE'"
          class="p-3.5 rounded-[4px] bg-brand-light/30 border border-brand-border space-y-3"
        >
          <div class="flex items-center justify-between">
            <span
              class="font-bold text-brand-dark text-xs flex items-center gap-1.5"
            >
              <Video class="w-4 h-4 text-brand" /> Salon Virtuel de la classe
            </span>
            <span class="text-[10px] text-muted-foreground"
              >Accès direct pour élèves et moniteur</span
            >
          </div>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="
                formModeLien = 'JITSI';
                genererSalonJitsi();
              "
              class="px-2.5 py-1 rounded-[4px] border text-[11px] font-semibold transition cursor-pointer"
              :class="
                formModeLien === 'JITSI'
                  ? 'bg-brand text-primary border-brand font-bold'
                  : 'bg-card border-border text-muted-foreground'
              "
            >
              Jitsi Meet Intégré (Recommandé)
            </button>
            <button
              type="button"
              @click="
                formModeLien = 'CUSTOM';
                formLienVisio = '';
              "
              class="px-2.5 py-1 rounded-[4px] border text-[11px] font-semibold transition cursor-pointer"
              :class="
                formModeLien === 'CUSTOM'
                  ? 'bg-brand text-primary border-brand font-bold'
                  : 'bg-card border-border text-muted-foreground'
              "
            >
              Lien personnalisé (Google Meet / Zoom)
            </button>
          </div>

          <div>
            <label class="block font-semibold text-foreground mb-1 text-[11px]">
              {{
                formModeLien === "JITSI"
                  ? "Lien de la salle sécurisée (Généré automatiquement)"
                  : "Collez votre lien de visioconférence (Google Meet, Zoom, Teams...)"
              }}
            </label>
            <div class="flex items-center gap-2">
              <input
                v-model="formLienVisio"
                type="url"
                required
                :placeholder="
                  formModeLien === 'JITSI'
                    ? 'https://meet.jit.si/...'
                    : 'https://meet.google.com/...'
                "
                class="flex-1 h-8 px-2.5 rounded-[4px] border border-border bg-card text-foreground text-xs focus:outline-none focus:border-brand-dark font-mono"
              />
              <button
                v-if="formModeLien === 'JITSI'"
                type="button"
                @click="genererSalonJitsi"
                class="px-2.5 py-1.5 rounded-[4px] border border-border bg-card hover:bg-muted text-foreground text-xs font-semibold flex items-center gap-1 transition cursor-pointer shrink-0"
                title="Générer un nouveau nom de salon"
              >
                <Sparkles class="w-3.5 h-3.5 text-brand" />
                <span>Régénérer</span>
              </button>
            </div>
            <p class="text-[10px] text-muted-foreground mt-1">
              Les apprenants inscrits en mode <strong>En Ligne</strong> ou
              <strong>Hybride</strong> verront ce lien s'activer 15 min avant le
              début du cours.
            </p>
          </div>
        </div>

        <!-- Pied de page -->
        <div
          class="pt-3 border-t border-border flex items-center justify-end space-x-2"
        >
          <button
            type="button"
            @click="$emit('close')"
            class="px-3 py-1.5 rounded-[4px] border border-border text-muted-foreground hover:text-foreground text-xs font-semibold transition cursor-pointer"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="coursStore.loading"
            class="px-4 py-1.5 rounded-[4px] bg-brand hover:bg-brand-dark text-primary font-bold text-xs uppercase tracking-wider transition cursor-pointer disabled:opacity-50"
          >
            {{
              coursStore.loading ? "Enregistrement..." : "Programmer la séance"
            }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
