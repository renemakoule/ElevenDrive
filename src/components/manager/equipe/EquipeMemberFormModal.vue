<script setup lang="ts">
import { ref, watch } from "vue";
import type { Agence, RegisterUserPayload, UserTeamMember } from "../../../types";
import {
  UserPlus,
  Pencil,
  X,
  AlertCircle,
  Smartphone,
  Monitor,
  Shield,
  Building2,
  Phone,
  Mail,
  Lock,
  Eye,
  EyeOff,
  Sparkles,
} from "lucide-vue-next";

const props = defineProps<{
  isOpen: boolean;
  mode: "ADD" | "EDIT";
  member: UserTeamMember | null;
  agences: Agence[];
  isSubmitting: boolean;
  errorMsg: string | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (
    e: "save",
    payload: { form: RegisterUserPayload; statut: "ACTIF" | "SUSPENDU" },
  ): void;
}>();

const showPassword = ref(false);
const localError = ref<string | null>(null);
const editStatut = ref<"ACTIF" | "SUSPENDU">("ACTIF");

const form = ref<RegisterUserPayload>({
  nom: "",
  prenoms: "",
  telephone: "",
  email: "",
  mot_de_passe: "",
  role: "MONITEUR",
  agence_id: "",
  code_pin_manager: "",
});

const generatePassword = () => {
  const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789";
  let pwd = "";
  for (let i = 0; i < 8; i++) {
    pwd += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  form.value.mot_de_passe = pwd;
};

watch(
  () => [props.isOpen, props.mode, props.member] as const,
  ([open, mode, member]) => {
    if (!open) return;
    localError.value = null;
    if (mode === "EDIT" && member) {
      editStatut.value = (member.statut as "ACTIF" | "SUSPENDU") || "ACTIF";
      form.value = {
        nom: member.nom,
        prenoms: member.prenoms,
        telephone: member.telephone,
        email: member.email || "",
        mot_de_passe: "",
        role: member.role,
        agence_id: member.agence_id || props.agences[0]?.id || "",
        code_pin_manager: "",
      };
    } else {
      editStatut.value = "ACTIF";
      form.value = {
        nom: "",
        prenoms: "",
        telephone: "",
        email: "",
        mot_de_passe: "",
        role: "MONITEUR",
        agence_id: props.agences[0]?.id || "",
        code_pin_manager: "",
      };
      generatePassword();
    }
  },
  { immediate: true },
);

const handleSubmit = () => {
  if (props.isSubmitting) return;
  localError.value = null;

  if (!form.value.nom.trim() || !form.value.prenoms.trim()) {
    localError.value = "Le nom et prénom sont obligatoires.";
    return;
  }
  if (!form.value.telephone.trim()) {
    localError.value = "Le numéro de téléphone est obligatoire.";
    return;
  }
  if (
    props.mode === "ADD" &&
    (!form.value.mot_de_passe || form.value.mot_de_passe.length < 6)
  ) {
    localError.value = "Le mot de passe doit comporter au moins 6 caractères.";
    return;
  }
  if (
    props.mode === "EDIT" &&
    form.value.mot_de_passe &&
    form.value.mot_de_passe.length < 6
  ) {
    localError.value =
      "Le nouveau mot de passe doit comporter au moins 6 caractères.";
    return;
  }

  emit("save", { form: { ...form.value }, statut: editStatut.value });
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200"
    @click.self="emit('close')"
  >
    <div
      class="relative w-full max-w-lg overflow-hidden rounded-[5px] border border-border bg-card shadow-2xl animate-in zoom-in-95 duration-200"
    >
      <div
        class="flex items-center justify-between border-b border-border bg-muted/30 px-6 py-4"
      >
        <div class="flex items-center gap-2.5">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-xl border border-brand-border bg-brand-light text-brand"
          >
            <Pencil v-if="mode === 'EDIT'" class="h-5 w-5" />
            <UserPlus v-else class="h-5 w-5" />
          </div>
          <div>
            <h3 class="text-sm font-semibold text-foreground">
              {{
                mode === "ADD"
                  ? "Créer un nouveau collaborateur"
                  : "Modifier le collaborateur"
              }}
            </h3>
            <p class="text-xs text-muted-foreground">
              {{
                mode === "ADD"
                  ? "Affectation à une agence et génération des accès"
                  : "Mise à jour des coordonnées et agence d'affectation"
              }}
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <form
        @submit.prevent="handleSubmit"
        class="p-6 space-y-4 max-h-[80vh] overflow-y-auto"
      >
        <div
          v-if="localError || errorMsg"
          class="flex items-center gap-2 rounded-[5px] border border-rose-500/20 bg-rose-500/10 p-3 text-xs text-rose-600 dark:text-rose-400"
        >
          <AlertCircle class="h-4 w-4 shrink-0" />
          <span>{{ localError || errorMsg }}</span>
        </div>

        <fieldset
          :disabled="isSubmitting"
          class="space-y-4 border-0 p-0 m-0 disabled:opacity-75"
        >
          <!-- Rôle -->
          <div>
            <label class="mb-1.5 block text-xs font-medium text-foreground">
              Rôle et habilitation *
            </label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                @click="form.role = 'MONITEUR'"
                class="flex flex-col items-start rounded-[4px] border p-3 text-left transition"
                :class="
                  form.role === 'MONITEUR'
                    ? 'border-brand bg-brand-light/50 ring-1 ring-brand text-brand-dark'
                    : 'border-border bg-card hover:bg-muted/40 text-foreground'
                "
              >
                <div class="flex items-center gap-1.5 font-semibold text-xs">
                  <Smartphone class="h-3.5 w-3.5 text-sky-600" />
                  <span>Moniteur</span>
                </div>
                <span class="mt-1 text-[11px] text-muted-foreground leading-tight">
                  Livret tactile mobile PWA, conduite & signature
                </span>
              </button>

              <button
                type="button"
                @click="form.role = 'SECRETAIRE'"
                class="flex flex-col items-start rounded-[4px] border p-3 text-left transition"
                :class="
                  form.role === 'SECRETAIRE'
                    ? 'border-brand bg-brand-light/50 ring-1 ring-brand text-brand-dark'
                    : 'border-border bg-card hover:bg-muted/40 text-foreground'
                "
              >
                <div class="flex items-center gap-1.5 font-semibold text-xs">
                  <Monitor class="h-3.5 w-3.5 text-emerald-600" />
                  <span>Secrétaire</span>
                </div>
                <span class="mt-1 text-[11px] text-muted-foreground leading-tight">
                  Guichet d'agence, inscriptions, caisse & reçus
                </span>
              </button>

              <button
                type="button"
                @click="form.role = 'DIRECTEUR'"
                class="flex flex-col items-start rounded-[4px] border p-3 text-left transition"
                :class="
                  form.role === 'DIRECTEUR'
                    ? 'border-brand bg-brand-light/50 ring-1 ring-brand text-brand-dark'
                    : 'border-border bg-card hover:bg-muted/40 text-foreground'
                "
              >
                <div class="flex items-center gap-1.5 font-semibold text-xs">
                  <Shield class="h-3.5 w-3.5 text-indigo-600" />
                  <span>Directeur Adjoint</span>
                </div>
                <span class="mt-1 text-[11px] text-muted-foreground leading-tight">
                  Supervision globale, finances & dérogations
                </span>
              </button>
            </div>
          </div>

          <!-- Agence d'affectation -->
          <div>
            <label class="mb-1.5 block text-xs font-medium text-foreground">
              Agence d'affectation *
            </label>
            <div class="relative">
              <Building2
                class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground"
              />
              <select
                v-model="form.agence_id"
                required
                class="h-10 w-full rounded-lg border border-border bg-background pl-9 pr-3 text-xs text-foreground focus:border-brand focus:outline-none"
              >
                <option v-for="ag in agences" :key="ag.id" :value="ag.id">
                  {{ ag.nom }} ({{
                    ag.est_siege
                      ? "Siège Principal"
                      : ag.commune || ag.ville || "Succursale"
                  }})
                </option>
                <option v-if="agences.length === 0" value="">
                  Agence Principale (Siège)
                </option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1.5 block text-xs font-medium text-foreground">
                Nom de famille *
              </label>
              <input
                v-model="form.nom"
                type="text"
                required
                placeholder="Ex: Makoule"
                class="h-10 w-full rounded-lg border border-border bg-background px-3 text-xs text-foreground uppercase placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-foreground">
                Prénoms *
              </label>
              <input
                v-model="form.prenoms"
                type="text"
                required
                placeholder="Ex: Rene Boris"
                class="h-10 w-full rounded-lg border border-border bg-background px-3 text-xs text-foreground placeholder:text-muted-foreground/60 focus:border-brand focus:outline-none"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1.5 block text-xs font-medium text-foreground">
                Téléphone (Identifiant) *
              </label>
              <div class="relative">
                <Phone
                  class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground"
                />
                <input
                  v-model="form.telephone"
                  type="tel"
                  required
                  placeholder="+237 651 727 932"
                  class="h-10 w-full rounded-lg border border-border bg-background pl-9 pr-3 text-xs font-mono text-foreground focus:border-brand focus:outline-none"
                />
              </div>
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-foreground">
                Adresse Email (Optionnel)
              </label>
              <div class="relative">
                <Mail
                  class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground"
                />
                <input
                  v-model="form.email"
                  type="email"
                  placeholder="collaborateur@autoecole.ci"
                  class="h-10 w-full rounded-lg border border-border bg-background pl-9 pr-3 text-xs text-foreground focus:border-brand focus:outline-none"
                />
              </div>
            </div>
          </div>

          <div v-if="mode === 'EDIT'">
            <label class="mb-1.5 block text-xs font-medium text-foreground">
              Statut du compte *
            </label>
            <select
              v-model="editStatut"
              class="h-10 w-full rounded-lg border border-border bg-background px-3 text-xs text-foreground focus:border-brand focus:outline-none"
            >
              <option value="ACTIF">Actif (Accès autorisé)</option>
              <option value="SUSPENDU">Suspendu (Accès bloqué)</option>
            </select>
          </div>

          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="text-xs font-medium text-foreground">
                {{
                  mode === "ADD"
                    ? "Mot de passe initial *"
                    : "Nouveau mot de passe (optionnel)"
                }}
              </label>
              <button
                type="button"
                @click="generatePassword"
                class="text-[11px] text-brand hover:underline"
              >
                Générer un mot de passe
              </button>
            </div>
            <div class="relative flex items-center">
              <Lock
                class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground"
              />
              <input
                v-model="form.mot_de_passe"
                :type="showPassword ? 'text' : 'password'"
                :required="mode === 'ADD'"
                :minlength="form.mot_de_passe ? 6 : undefined"
                :placeholder="
                  mode === 'ADD'
                    ? 'Au moins 6 caractères'
                    : 'Laisser vide pour conserver le mot de passe actuel'
                "
                class="h-10 w-full rounded-lg border border-border bg-background pl-9 pr-10 text-xs font-mono text-foreground focus:border-brand focus:outline-none"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 text-muted-foreground hover:text-foreground"
              >
                <EyeOff v-if="showPassword" class="h-4 w-4" />
                <Eye v-else class="h-4 w-4" />
              </button>
            </div>
          </div>

          <div v-if="form.role === 'DIRECTEUR'">
            <label class="mb-1.5 block text-xs font-medium text-foreground">
              Code PIN Dérogations (4 chiffres)
            </label>
            <input
              v-model="form.code_pin_manager"
              type="password"
              maxlength="4"
              placeholder="Ex: 1234"
              class="h-10 w-32 rounded-lg border border-border bg-background px-3 font-mono text-xs text-center tracking-widest text-foreground focus:border-brand focus:outline-none"
            />
          </div>
        </fieldset>

        <div
          v-if="isSubmitting"
          class="rounded-[5px] bg-brand/5 border border-brand/20 p-2.5 text-center text-[11px] text-brand animate-pulse"
        >
          Enregistrement sécurisé en cours sur le serveur cloud...
        </div>

        <div
          class="pt-3 flex items-center justify-end gap-2 border-t border-border"
        >
          <button
            type="button"
            :disabled="isSubmitting"
            @click="emit('close')"
            class="h-10 rounded-lg border border-border px-4 text-xs font-medium text-foreground hover:bg-muted transition disabled:opacity-50"
          >
            Annuler
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-brand px-5 text-xs font-semibold text-white hover:bg-brand-dark transition disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
          >
            <Sparkles v-if="isSubmitting" class="h-4 w-4 animate-spin" />
            <span>
              {{
                isSubmitting
                  ? "Enregistrement en cours..."
                  : mode === "ADD"
                    ? "Enregistrer le collaborateur"
                    : "Mettre à jour le collaborateur"
              }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

