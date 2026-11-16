<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useAuthStore } from '../../stores/authStore';
import { useGuichetStore } from '../../stores/guichetStore';
import { generateQrDataUrl, downloadQrCode } from '../../utils/qrGenerator';
import { generateSchoolSlug } from '../../services/authBackendService';
import { getPublicAppOrigin } from '../../config/env';
import { 
  X, Download, Copy, CheckCircle2, Share2, Printer, 
  QrCode, Sparkles, Globe 
} from 'lucide-vue-next';

interface Props {
  isOpen: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'close'): void;
}>();

const authStore = useAuthStore();
const guichetStore = useGuichetStore();

// Code unique 100% dynamique (issu de l'auto-école active ou généré depuis son nom)
const displaySchoolSlug = computed(() => {
  if (authStore.currentSchool?.code_unique) {
    return authStore.currentSchool.code_unique;
  }
  if (authStore.currentSchool?.nom_commercial) {
    return generateSchoolSlug(authStore.currentSchool.nom_commercial);
  }
  return 'auto-ecole';
});

const webPortalUrl = computed(() => {
  if (authStore.schoolAccessLinks?.portal_url) {
    return authStore.schoolAccessLinks.portal_url;
  }
  const slug = displaySchoolSlug.value;
  const base = getPublicAppOrigin();
  return base ? `${base}/${slug}` : `/${slug}`;
});

const qrDataUrl = ref<string>('');
const isGenerating = ref(false);
const isCopied = ref(false);

const loadQrCode = async () => {
  if (!webPortalUrl.value) return;
  isGenerating.value = true;
  try {
    qrDataUrl.value = await generateQrDataUrl(webPortalUrl.value, {
      width: 360,
      margin: 2,
      darkColor: '#064e3b',
      lightColor: '#ffffff',
    });
  } catch (err) {
    console.error('Erreur génération QR code unique:', err);
  } finally {
    isGenerating.value = false;
  }
};

watch(
  () => [props.isOpen, webPortalUrl.value],
  ([open]) => {
    if (open) {
      loadQrCode();
    }
  },
  { immediate: true }
);

const handleCopy = async () => {
  if (!webPortalUrl.value) return;
  try {
    await navigator.clipboard.writeText(webPortalUrl.value);
    isCopied.value = true;
    guichetStore.showToast('Lien web unique copié dans le presse-papiers !', 'success');
    setTimeout(() => {
      isCopied.value = false;
    }, 2500);
  } catch {
    isCopied.value = true;
    setTimeout(() => {
      isCopied.value = false;
    }, 2500);
  }
};

const handleDownload = () => {
  if (!qrDataUrl.value) return;
  const schoolSlug = authStore.currentSchool?.code_unique || 'autoecole';
  downloadQrCode(qrDataUrl.value, `qrcode-portail-web-${schoolSlug}.png`);
  guichetStore.showToast('QR Code d\'accès web téléchargé !', 'success');
};

const handlePrint = () => {
  if (!qrDataUrl.value) return;
  const printWindow = window.open('', '_blank');
  if (!printWindow) return;

  printWindow.document.write(`
    <!DOCTYPE html>
    <html lang="fr">
    <head>
      <meta charset="UTF-8">
      <title>QR Code - Portail Web ${authStore.schoolName}</title>
      <style>
        body {
          font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
          margin: 0;
          padding: 40px;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          text-align: center;
          color: #0f172a;
        }
        .card {
          border: 2px dashed #059669;
          border-radius: 20px;
          padding: 30px;
          max-width: 420px;
          margin: 0 auto;
        }
        h1 {
          font-size: 22px;
          margin: 0 0 6px 0;
          color: #064e3b;
        }
        p.school {
          font-size: 13px;
          font-weight: 700;
          color: #64748b;
          margin: 0 0 18px 0;
          text-transform: uppercase;
          letter-spacing: 0.05em;
        }
        .qr-wrapper {
          background: #ffffff;
          padding: 16px;
          border-radius: 16px;
          display: inline-block;
          box-shadow: 0 4px 12px rgba(0,0,0,0.08);
        }
        img {
          width: 280px;
          height: 280px;
          display: block;
        }
        .instruction {
          margin-top: 20px;
          font-size: 14px;
          line-height: 1.5;
          color: #334155;
        }
        .instruction strong {
          color: #059669;
        }
        .url {
          margin-top: 12px;
          font-family: monospace;
          font-size: 12px;
          font-weight: 600;
          color: #064e3b;
          word-break: break-all;
        }
      </style>
    </head>
    <body>
      <div class="card">
        <p class="school">${authStore.schoolName}</p>
        <h1>Portail Web & Accès Mobile</h1>
        <div class="qr-wrapper">
          <img src="${qrDataUrl.value}" alt="QR Code Portail Web" />
        </div>
        <p class="instruction">
          <strong>Scannez ce QR Code</strong> avec l'appareil photo de votre smartphone pour accéder à la plateforme de l'auto-école (Connexion élèves, moniteurs et inscriptions).
        </p>
        <p class="url">${webPortalUrl.value}</p>
      </div>
      <script>
        window.onload = function() {
          window.print();
        };
      <\/script>
    </body>
    </html>
  `);
  printWindow.document.close();
};

const handleShareWhatsApp = () => {
  const text = encodeURIComponent(
    `Bonjour ! Voici l'accès officiel au portail web de ${authStore.schoolName} :\n${webPortalUrl.value}\n\nScannez le QR Code ou cliquez sur ce lien pour vous connecter.`
  );
  window.open(`https://wa.me/?text=${text}`, '_blank');
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200 select-none"
    @click.self="emit('close')"
  >
    <div
      class="relative w-full max-w-md overflow-hidden rounded-[5px] border border-border bg-card animate-in zoom-in-95 duration-200"
    >
      <!-- En-tête de la modale -->
      <div class="flex items-center justify-between border-b border-border bg-muted/30 px-5 py-3.5">
        <div class="flex items-center gap-2.5">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl border border-brand-border bg-brand-light text-brand">
            <QrCode class="h-5 w-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h3 class="text-sm font-semibold text-foreground">Accès Web & QR Code</h3>
              <span class="rounded-md border border-brand-border bg-brand-light px-2 py-0.5 font-mono text-[10px] font-semibold text-brand-dark">
                {{ displaySchoolSlug }}
              </span>
            </div>
            <p class="text-xs text-muted-foreground truncate max-w-[260px]">
              {{ authStore.schoolName }} · Portail Unique
            </p>
          </div>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground transition"
          aria-label="Fermer"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <!-- Contenu central : Un seul QR Code et un seul lien Web -->
      <div class="flex flex-col items-center p-6 text-center">
        <!-- Titre & Descriptif -->
        <div class="mb-4">
          <div class="inline-flex items-center gap-1.5 rounded-full border border-brand-border bg-brand-light px-2.5 py-0.5 text-[11px] font-medium text-brand-dark mb-1.5">
            <Globe class="h-3 w-3 text-brand" />
            Lien d'accès unique
          </div>
          <h4 class="text-base font-semibold text-foreground">Portail Web de l'Établissement</h4>
          <p class="mt-0.5 text-xs text-muted-foreground max-w-xs">
            Tous les utilisateurs (élèves, moniteurs, personnels) utilisent ce portail pour se connecter.
          </p>
        </div>

        <!-- Cadre QR Code Unique -->
        <div class="relative group">
          <div class="absolute -inset-1 rounded-[5px] bg-gradient-to-tr from-brand to-emerald-400 opacity-20 blur-md group-hover:opacity-30 transition duration-300"></div>
          <div class="relative flex h-60 w-60 items-center justify-center rounded-[5px] border-2 border-brand-border bg-white p-3">
            <img
              v-if="qrDataUrl"
              :src="qrDataUrl"
              alt="QR Code d'accès au portail web"
              class="h-full w-full object-contain rounded select-none"
            />
            <div v-else class="flex flex-col items-center gap-2 text-muted-foreground">
              <Sparkles class="h-5 w-5 animate-spin text-brand" />
              <span class="text-xs">Génération du QR Code...</span>
            </div>
          </div>
        </div>

        <p class="mt-3.5 text-xs font-medium text-foreground">
          Visez ce QR Code avec l'appareil photo d'un smartphone 📸
        </p>

        <!-- Champ URL Unique & Bouton Copier -->
        <div class="mt-4 flex w-full items-center gap-2">
          <input
            type="text"
            readonly
            :value="webPortalUrl"
            class="h-9 min-w-0 flex-1 select-all rounded-lg border border-border bg-muted/40 px-3 font-mono text-xs text-foreground focus:outline-none focus:border-brand"
          />
          <button
            type="button"
            @click="handleCopy"
            class="inline-flex h-9 items-center gap-1.5 rounded-lg border border-border px-3 text-xs font-medium text-foreground transition hover:bg-muted shrink-0"
          >
            <CheckCircle2 v-if="isCopied" class="h-3.5 w-3.5 text-brand" />
            <Copy v-else class="h-3.5 w-3.5 text-muted-foreground" />
            <span>{{ isCopied ? 'Copié' : 'Copier' }}</span>
          </button>
        </div>

        <!-- Barre d'outils (Télécharger, Imprimer, WhatsApp) -->
        <div class="mt-4 grid w-full grid-cols-3 gap-2">
          <button
            type="button"
            @click="handleDownload"
            class="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg border border-border bg-background px-3 text-xs font-medium text-foreground transition hover:bg-muted hover:border-foreground/20"
          >
            <Download class="h-3.5 w-3.5 text-brand" />
            <span>Télécharger</span>
          </button>

          <button
            type="button"
            @click="handlePrint"
            class="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg border border-border bg-background px-3 text-xs font-medium text-foreground transition hover:bg-muted hover:border-foreground/20"
          >
            <Printer class="h-3.5 w-3.5 text-muted-foreground" />
            <span>Imprimer</span>
          </button>

          <button
            type="button"
            @click="handleShareWhatsApp"
            class="inline-flex h-9 items-center justify-center gap-1.5 rounded-lg border border-emerald-600/30 bg-emerald-500/10 px-3 text-xs font-medium text-emerald-700 transition hover:bg-emerald-500/20 dark:text-emerald-400"
          >
            <Share2 class="h-3.5 w-3.5" />
            <span>WhatsApp</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
