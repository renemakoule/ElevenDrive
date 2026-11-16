<script setup lang="ts">
import { ref, reactive, computed, watch } from "vue";
import { useAuthStore } from "../../../stores/authStore";
import { useGuichetStore } from "../../../stores/guichetStore";
import {
  Building2,
  CheckCircle2,
  Eye,
  Globe,
  Copy,
  Save,
} from "lucide-vue-next";

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

const isSaving = ref(false);
const copiedLink = ref<string | null>(null);
const perms = computed(() => authStore.rolePermissions);

const schoolForm = reactive({
  nom_commercial: authStore.currentSchool?.nom_commercial || "ElevenDrive",
  numero_agrement:
    authStore.currentSchool?.numero_agrement || "Agrément en cours",
  code_unique: authStore.currentSchool?.code_unique || "autoecole",
  pays: authStore.currentSchool?.pays || "Cameroun",
  ville: authStore.currentSchool?.ville || "Douala",
  devise: authStore.currentSchool?.devise || "XAF",
  telephone: authStore.currentSchool?.telephone || "",
  email: authStore.currentSchool?.email || "",
});

watch(
  () => authStore.currentSchool,
  (s) => {
    if (s) {
      schoolForm.nom_commercial = s.nom_commercial || "ElevenDrive";
      schoolForm.numero_agrement = s.numero_agrement || "Agrément en cours";
      schoolForm.code_unique = s.code_unique || "autoecole";
      schoolForm.pays = s.pays || "Cameroun";
      schoolForm.ville = s.ville || "Douala";
      schoolForm.devise = s.devise || "XAF";
      schoolForm.telephone = s.telephone || "";
      schoolForm.email = s.email || "";
    }
  },
  { deep: true, immediate: true },
);

const copyToClipboard = async (text: string, key: string) => {
  if (!text) return;
  try {
    await navigator.clipboard.writeText(text);
    copiedLink.value = key;
    setTimeout(() => {
      copiedLink.value = null;
    }, 2000);
    guichetStore.showToast("Lien copié dans le presse-papier !", "info");
  } catch {}
};

const saveSchool = async () => {
  if (!perms.value.canManageSchool) {
    guichetStore.showToast(
      "Accès refusé : seul le Directeur peut modifier les informations de l'établissement.",
      "error",
    );
    return;
  }
  if (!schoolForm.nom_commercial.trim()) {
    guichetStore.showToast("Le nom de l'auto-école est requis.", "error");
    return;
  }

  isSaving.value = true;
  try {
    await authStore.updateSchool({
      nom_commercial: schoolForm.nom_commercial.trim(),
      numero_agrement: schoolForm.numero_agrement.trim(),
      code_unique: schoolForm.code_unique.trim(),
      pays: schoolForm.pays.trim(),
      ville: schoolForm.ville.trim(),
      devise: schoolForm.devise.trim(),
      telephone: schoolForm.telephone.trim(),
      email: schoolForm.email.trim() || undefined,
    });
    guichetStore.showToast(
      "Les informations de l'établissement ont été mises à jour.",
      "success",
    );
  } catch (err: any) {
    guichetStore.showToast(
      err?.message || "Erreur lors de la mise à jour de l'établissement.",
      "error",
    );
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="space-y-6">
    <div
      class="p-4 rounded-[4px] bg-brand-light/30 border border-brand-border flex items-start justify-between"
    >
      <div class="flex items-center space-x-3">
        <div
          class="w-12 h-12 rounded-[5px] bg-brand-light border border-brand-border text-brand-dark flex items-center justify-center font-black text-xl"
        >
          <Building2 class="w-6 h-6" />
        </div>
        <div>
          <h3 class="text-sm font-bold text-foreground">
            {{ schoolForm.nom_commercial }}
          </h3>
          <p class="text-xs text-muted-foreground">
            Agrément officiel :
            <strong class="text-foreground">{{ schoolForm.numero_agrement }}</strong>
          </p>
          <p class="text-[11px] text-brand-dark dark:text-brand font-mono mt-0.5">
            Identifiant : {{ schoolForm.code_unique }}
          </p>
        </div>
      </div>

      <span
        class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold"
        :class="
          perms.canManageSchool
            ? 'bg-emerald-500/10 text-emerald-600 border border-emerald-500/20'
            : 'bg-muted text-muted-foreground border border-border'
        "
      >
        <CheckCircle2 v-if="perms.canManageSchool" class="w-3.5 h-3.5" />
        <Eye v-else class="w-3.5 h-3.5" />
        {{ perms.canManageSchool ? "Mode Édition Directeur" : "Lecture Seule" }}
      </span>
    </div>

    <form @submit.prevent="saveSchool" class="space-y-4">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Nom Commercial de l'Auto-École *
          </label>
          <input
            v-model="schoolForm.nom_commercial"
            type="text"
            :disabled="!perms.canManageSchool"
            required
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm font-semibold focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Identifiant Web Unique (URL & QR Code) *
          </label>
          <input
            v-model="schoolForm.code_unique"
            type="text"
            :disabled="!perms.canManageSchool"
            required
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm font-mono focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Numéro d'Agrément Ministériel *
          </label>
          <input
            v-model="schoolForm.numero_agrement"
            type="text"
            :disabled="!perms.canManageSchool"
            required
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm font-mono focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Pays du Siège Social
          </label>
          <input
            v-model="schoolForm.pays"
            type="text"
            :disabled="!perms.canManageSchool"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Ville du Siège
          </label>
          <input
            v-model="schoolForm.ville"
            type="text"
            :disabled="!perms.canManageSchool"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Devise Comptable Officielle
          </label>
          <select
            v-model="schoolForm.devise"
            :disabled="!perms.canManageSchool"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          >
            <option value="XAF">FCFA (XAF - Afrique Centrale)</option>
            <option value="XOF">FCFA (XOF - Afrique de l'Ouest)</option>
            <option value="GNF">Franc Guinéen (GNF)</option>
            <option value="EUR">Euro (€)</option>
            <option value="USD">Dollar ($)</option>
          </select>
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Téléphone Standard de l'Établissement
          </label>
          <input
            v-model="schoolForm.telephone"
            type="tel"
            :disabled="!perms.canManageSchool"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label class="block text-xs font-semibold text-muted-foreground mb-1.5">
            Email Officiel de Contact
          </label>
          <input
            v-model="schoolForm.email"
            type="email"
            :disabled="!perms.canManageSchool"
            class="w-full px-3 py-2 rounded-xl bg-card border border-border text-foreground text-sm focus:ring-2 focus:ring-brand/40 focus:border-brand outline-none transition disabled:opacity-70 disabled:cursor-not-allowed"
          />
        </div>
      </div>

      <!-- Lien d'accès Mobiles & PWA -->
      <div class="pt-4 border-t border-border mt-4">
        <h3
          class="text-xs font-bold uppercase text-muted-foreground tracking-wider mb-3 flex items-center gap-1.5"
        >
          <Globe class="w-3.5 h-3.5 text-brand" /> Lien d'accès Web Unique
          (Élèves, Moniteurs & Inscriptions)
        </h3>

        <div
          class="p-3 rounded-xl bg-muted/40 border border-border flex items-center justify-between"
        >
          <div class="min-w-0 pr-2">
            <div class="text-xs font-bold text-foreground">
              Portail Web Officiel de l'Auto-École
            </div>
            <div class="text-[11px] font-mono text-muted-foreground truncate">
              {{
                authStore.schoolAccessLinks?.portal_url ||
                `https://autoecole.app/${schoolForm.code_unique}`
              }}
            </div>
          </div>
          <button
            type="button"
            @click="
              copyToClipboard(
                authStore.schoolAccessLinks?.portal_url ||
                  `https://autoecole.app/${schoolForm.code_unique}`,
                'portal',
              )
            "
            class="p-1.5 rounded-lg bg-card border border-border hover:bg-muted text-muted-foreground hover:text-foreground transition shrink-0"
            title="Copier le lien d'accès unique"
          >
            <CheckCircle2
              v-if="copiedLink === 'portal'"
              class="w-4 h-4 text-emerald-500"
            />
            <Copy v-else class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Bouton Sauvegarder uniquement si Directeur -->
      <div v-if="perms.canManageSchool" class="flex justify-end pt-3">
        <button
          type="submit"
          :disabled="isSaving"
          class="flex items-center space-x-2 px-5 py-2.5 rounded-xl bg-brand hover:bg-brand-dark text-brand-dark font-bold text-xs uppercase tracking-wider transition shadow-md disabled:opacity-50"
        >
          <Save class="w-4 h-4" />
          <span>
            {{ isSaving ? "Enregistrement..." : "Enregistrer l'établissement" }}
          </span>
        </button>
      </div>
    </form>
  </div>
</template>
