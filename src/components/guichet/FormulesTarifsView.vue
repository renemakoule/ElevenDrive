<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useGuichetStore } from "../../stores/guichetStore";
import { useAuthStore } from "../../stores/authStore";
import { realtimeSocket } from "../../services/realtimeSocketService";
import type { FormulePermis } from "../../types";
import { Tag, ShieldCheck, Lock, PlusCircle, Trash2 } from "lucide-vue-next";
import FormuleCardItem from "./tarifs/FormuleCardItem.vue";
import FormuleCreateModal from "./tarifs/FormuleCreateModal.vue";
import FormuleEditPriceModal from "./tarifs/FormuleEditPriceModal.vue";

const guichetStore = useGuichetStore();
const authStore = useAuthStore();
let unsubscribeWs: (() => void) | null = null;

const filtreCategorie = ref<string>("TOUTES");
const isCreateModalOpen = ref(false);
const isEditPriceModalOpen = ref(false);
const selectedFormuleToEdit = ref<FormulePermis | null>(null);
const isSecretaireNoticeOpen = ref(false);
const noticeFormuleName = ref("");

onMounted(() => {
  guichetStore.loadFormules();
  unsubscribeWs = realtimeSocket.subscribe((msg) => {
    if (msg.type === "FORMULES_TARIFS_UPDATED") {
      guichetStore.loadFormules();
    }
  });
});

onUnmounted(() => {
  if (unsubscribeWs) {
    unsubscribeWs();
    unsubscribeWs = null;
  }
});

const canCreate = computed(
  () =>
    authStore.currentRole === "manager" ||
    authStore.currentRole === "secretaire",
);

const canModify = computed(() => authStore.currentRole === "manager");

const formulesFiltrees = computed(() => {
  if (!guichetStore.formules) return [];
  if (filtreCategorie.value === "TOUTES") return guichetStore.formules;
  return guichetStore.formules.filter(
    (f) => f.categorie_permis === filtreCategorie.value,
  );
});

const openEditPriceModal = (formule: FormulePermis) => {
  if (!canModify.value) {
    noticeFormuleName.value = formule.intitule;
    isSecretaireNoticeOpen.value = true;
    return;
  }
  selectedFormuleToEdit.value = formule;
  isEditPriceModalOpen.value = true;
};

const handleSupprimerUneFormule = async (formule: FormulePermis) => {
  if (
    !confirm(
      `Êtes-vous sûr de vouloir supprimer la formule "${formule.intitule}" du catalogue ?`,
    )
  ) {
    return;
  }
  try {
    await guichetStore.supprimerFormule(formule.id);
  } catch (err: any) {
    alert(err?.message || "Impossible de supprimer cette formule.");
  }
};

const handleViderCatalogue = async () => {
  if (
    !confirm(
      "Êtes-vous sûr de vouloir vider l'ensemble du catalogue ? Les formules par défaut seront supprimées pour vous permettre de définir vos propres tarifs dès l'ouverture.",
    )
  ) {
    return;
  }
  try {
    await guichetStore.viderCatalogueFormules();
  } catch (err: any) {
    alert(err?.message || "Erreur lors du vidage du catalogue.");
  }
};
</script>

<template>
  <div class="space-y-6 max-w-6xl mx-auto pb-10">
    <!-- En-tête -->
    <div
      class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 border-b border-border pb-4"
    >
      <div>
        <div class="flex items-center space-x-2">
          <Tag class="w-5 h-5 text-brand" />
          <h1 class="text-xl font-bold tracking-tight text-foreground">
            Catalogue & Tarifs des Permis
          </h1>
          <span
            class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-brand-light text-brand border border-brand/20"
          >
            {{ guichetStore.formules.length }} formule{{
              guichetStore.formules.length > 1 ? "s" : ""
            }}
          </span>
        </div>
        <p class="text-xs text-muted-foreground mt-1">
          Barème tarifaire officiel en vigueur avec échéanciers personnalisables
          par formule.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="canModify && guichetStore.formules.length > 0"
          @click="handleViderCatalogue"
          type="button"
          class="inline-flex items-center space-x-1.5 px-3 py-2 rounded-[4px] border border-rose-300 dark:border-rose-900/50 text-rose-700 dark:text-rose-300 bg-rose-50 dark:bg-rose-950/30 hover:bg-rose-100 dark:hover:bg-rose-900/40 text-xs font-semibold transition active:scale-95 shadow-xs cursor-pointer"
          title="Supprimer toutes les formules pour démarrer avec un catalogue totalement vierge"
        >
          <Trash2 class="w-3.5 h-3.5" />
          <span>Vider le Catalogue</span>
        </button>
        <button
          v-if="canCreate"
          @click="isCreateModalOpen = true"
          type="button"
          class="inline-flex items-center space-x-2 px-3.5 py-2 rounded-[4px] bg-brand text-white text-xs font-semibold hover:bg-brand/90 transition active:scale-95 shadow-xs cursor-pointer"
        >
          <PlusCircle class="w-4 h-4" />
          <span>Nouvelle Formule (Fixer Tarif)</span>
        </button>
      </div>
    </div>

    <!-- Bannière Gouvernance -->
    <div
      class="rounded-[4px] border border-border bg-brand-light/30 p-4 text-xs"
    >
      <div class="flex items-start space-x-3">
        <ShieldCheck class="w-5 h-5 text-brand shrink-0 mt-0.5" />
        <div class="space-y-1">
          <h4 class="font-semibold text-brand-dark">
            Règle de Gouvernance Tarifaire & Échéanciers :
          </h4>
          <p class="text-muted-foreground leading-relaxed">
            La fixation initiale d'une formule, de son tarif officiel et de ses
            tranches de paiement peut être effectuée par le
            <strong>Directeur</strong> ou la <strong>Secrétaire</strong> lors de
            la création d'une offre. Toutefois,
            <strong>seul le Directeur</strong> détient l'habilitation exclusive
            pour modifier un tarif fixé ou son plan d'échelonnement.
          </p>
          <div class="pt-1 flex flex-wrap items-center gap-2 font-medium">
            <span
              class="inline-flex items-center px-2 py-0.5 rounded text-[11px] bg-background border border-border"
            >
              Votre rôle actuel :
              <strong class="ml-1 text-foreground uppercase">{{
                authStore.currentRole === "manager"
                  ? "Directeur / Gérant"
                  : "Secrétaire de Guichet"
              }}</strong>
            </span>
            <span
              v-if="canModify"
              class="inline-flex items-center text-emerald-600 dark:text-emerald-400 text-[11px] font-semibold"
            >
              ✓ Habilité à modifier les tarifs et les tranches
            </span>
            <span
              v-else
              class="inline-flex items-center text-amber-600 dark:text-amber-400 text-[11px] font-medium"
            >
              🔒 Modification des tarifs réservée au Directeur
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Filtres -->
    <div class="flex items-center justify-between gap-2 overflow-x-auto pb-1">
      <div class="flex items-center gap-1.5 text-xs font-medium">
        <span class="text-muted-foreground mr-1 text-[11px]">Catégorie :</span>
        <button
          v-for="cat in ['TOUTES', 'A', 'A1', 'B', 'C', 'D', 'E', 'REC']"
          :key="cat"
          type="button"
          @click="filtreCategorie = cat"
          :class="[
            'px-2.5 py-1 rounded-[4px] text-[11px] font-semibold transition',
            filtreCategorie === cat
              ? 'bg-brand text-white'
              : 'bg-muted hover:bg-muted/80 text-foreground border border-border',
          ]"
        >
          {{
            cat === "TOUTES"
              ? "Toutes les catégories"
              : cat === "REC"
                ? "Recyclage / Performance"
                : `Permis ${cat}`
          }}
        </button>
      </div>
      <span class="text-[11px] text-muted-foreground shrink-0">
        {{ formulesFiltrees.length }} formule{{
          formulesFiltrees.length > 1 ? "s" : ""
        }}
        affichée{{ formulesFiltrees.length > 1 ? "s" : "" }}
      </span>
    </div>

    <!-- Grille des Formules -->
    <div
      v-if="formulesFiltrees.length > 0"
      class="grid grid-cols-1 md:grid-cols-2 gap-4"
    >
      <FormuleCardItem
        v-for="f in formulesFiltrees"
        :key="f.id"
        :formule="f"
        :can-modify="canModify"
        @edit="openEditPriceModal"
        @delete="handleSupprimerUneFormule"
      />
    </div>

    <!-- État vide -->
    <div
      v-else
      class="rounded-[4px] border border-dashed border-border p-12 text-center bg-card space-y-3"
    >
      <Tag class="w-10 h-10 text-muted-foreground mx-auto stroke-[1.5]" />
      <h3 class="text-sm font-bold text-foreground">
        Aucune formule trouvée pour cette sélection
      </h3>
      <p class="text-xs text-muted-foreground max-w-sm mx-auto">
        Aucun tarif n'a encore été défini pour cette catégorie. Créez votre
        première formule pour ouvrir les inscriptions au guichet.
      </p>
      <button
        v-if="canCreate"
        @click="isCreateModalOpen = true"
        type="button"
        class="inline-flex items-center space-x-2 px-3 py-1.5 rounded-[4px] bg-brand text-white text-xs font-semibold hover:bg-brand/90 transition"
      >
        <PlusCircle class="w-4 h-4" />
        <span>Créer une Formule</span>
      </button>
    </div>

    <!-- Modales modulaires -->
    <FormuleCreateModal
      :is-open="isCreateModalOpen"
      @close="isCreateModalOpen = false"
    />

    <FormuleEditPriceModal
      :is-open="isEditPriceModalOpen"
      :formule="selectedFormuleToEdit"
      @close="isEditPriceModalOpen = false"
    />

    <!-- Notice Secrétaire -->
    <div
      v-if="isSecretaireNoticeOpen"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
    >
      <div
        class="w-full max-w-sm overflow-hidden rounded-[4px] border border-amber-300 dark:border-amber-700 bg-card p-5 text-center space-y-3 shadow-xl"
      >
        <div
          class="w-10 h-10 rounded-full bg-amber-100 dark:bg-amber-950 flex items-center justify-center mx-auto text-amber-600"
        >
          <Lock class="w-5 h-5" />
        </div>
        <h3 class="text-sm font-bold text-foreground">
          Modification Réservée au Directeur
        </h3>
        <p class="text-xs text-muted-foreground leading-relaxed">
          Le prix et l'échéancier de la formule
          <strong>{{ noticeFormuleName }}</strong> ont été fixés. En tant que
          Secrétaire, vous pouvez fixer le tarif lors de la création d'une
          formule, mais <strong>seul le Directeur</strong> a l'autorité légale
          pour modifier un tarif ou des tranches déjà fixés.
        </p>
        <div class="pt-2">
          <button
            type="button"
            @click="isSecretaireNoticeOpen = false"
            class="w-full py-2 rounded-[4px] bg-muted hover:bg-muted/80 text-foreground text-xs font-semibold border border-border transition"
          >
            J'ai compris
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
