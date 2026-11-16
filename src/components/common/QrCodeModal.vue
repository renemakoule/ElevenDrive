<script setup lang="ts">
import { ref, watch } from 'vue';
import { generateQrDataUrl, downloadQrCode } from '../../utils/qrGenerator';
import { 
  X, Download, Copy, CheckCircle2, Share2, Printer, 
  Smartphone, ExternalLink, Sparkles 
} from 'lucide-vue-next';

interface Props {
  isOpen: boolean;
  title: string;
  subtitle?: string;
  url: string;
  badgeText?: string;
  badgeDot?: string;
  schoolName?: string;
}

const props = withDefaults(defineProps<Props>(), {
  subtitle: '',
  badgeText: 'Application Mobile',
  badgeDot: 'bg-brand',
  schoolName: 'AutoSchool',
});

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const qrDataUrl = ref<string>('');
const isCopied = ref(false);
const isGenerating = ref(false);

const loadQrCode = async () => {
  if (!props.url) return;
  isGenerating.value = true;
  try {
    qrDataUrl.value = await generateQrDataUrl(props.url, {
      width: 400,
      margin: 2,
      darkColor: '#064e3b',
      lightColor: '#ffffff',
    });
  } finally {
    isGenerating.value = false;
  }
};

watch(
  () => [props.isOpen, props.url],
  ([open]) => {
    if (open) {
      loadQrCode();
    }
  },
  { immediate: true }
);

const handleCopy = async () => {
  try {
    await navigator.clipboard.writeText(props.url);
    isCopied.value = true;
    setTimeout(() => {
      isCopied.value = false;
    }, 2500);
  } catch {
    // Fallback
    const input = document.createElement('input');
    input.value = props.url;
    document.body.appendChild(input);
    input.select();
    document.execCommand('copy');
    document.body.removeChild(input);
    isCopied.value = true;
    setTimeout(() => {
      isCopied.value = false;
    }, 2500);
  }
};

const handleDownload = () => {
  if (!qrDataUrl.value) return;
  const slug = (props.title || 'qrcode')
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-');
  downloadQrCode(qrDataUrl.value, `qrcode-${slug}.png`);
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
      <title>QR Code - ${props.title}</title>
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
          font-size: 24px;
          margin: 0 0 6px 0;
          color: #064e3b;
        }
        p.school {
          font-size: 14px;
          font-weight: 600;
          color: #64748b;
          margin: 0 0 20px 0;
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
          margin-top: 24px;
          font-size: 14px;
          line-height: 1.5;
          color: #334155;
        }
        .instruction strong {
          color: #059669;
        }
        .url {
          margin-top: 14px;
          font-family: monospace;
          font-size: 11px;
          color: #94a3b8;
          word-break: break-all;
        }
        @media print {
          body { padding: 20px; }
        }
      </style>
    </head>
    <body>
      <div class="card">
        <p class="school">${props.schoolName}</p>
        <h1>${props.title}</h1>
        <div class="qr-wrapper">
          <img src="${qrDataUrl.value}" alt="QR Code" />
        </div>
        <p class="instruction">
          <strong>Scannez avec l'appareil photo</strong> de votre smartphone pour accéder à l'application mobile.
        </p>
        <p class="url">${props.url}</p>
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
    `Bonjour ! Voici votre accès direct à ${props.title} (${props.schoolName}) :\n${props.url}\n\nScannez le QR Code ou cliquez sur ce lien pour démarrer.`
  );
  window.open(`https://wa.me/?text=${text}`, '_blank');
};
</script>

<template>
  <div
    v-if="isOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200"
    @click.self="emit('close')"
  >
    <div
      class="relative w-full max-w-md overflow-hidden rounded-2xl border border-border bg-card shadow-2xl animate-in zoom-in-95 duration-200"
    >
      <!-- En-tête -->
      <div class="flex items-center justify-between border-b border-border bg-muted/30 px-5 py-4">
        <div class="flex items-center gap-2.5">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl border border-brand-border bg-brand-light text-brand">
            <Smartphone class="h-4 w-4" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h3 class="text-sm font-semibold text-foreground">{{ title }}</h3>
              <span
                v-if="badgeText"
                class="inline-flex items-center gap-1 rounded-full border border-border bg-background px-2 py-0.5 text-[10px] font-medium text-foreground"
              >
                <span class="h-1.5 w-1.5 rounded-full" :class="badgeDot"></span>
                {{ badgeText }}
              </span>
            </div>
            <p v-if="subtitle || schoolName" class="text-xs text-muted-foreground truncate max-w-[240px]">
              {{ subtitle || schoolName }}
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

      <!-- Corps : Affichage du QR Code -->
      <div class="flex flex-col items-center p-6 text-center">
        <!-- Cadre QR Code lumineux -->
        <div class="relative group">
          <div class="absolute -inset-1 rounded-3xl bg-gradient-to-tr from-brand to-emerald-400 opacity-20 blur-md group-hover:opacity-40 transition duration-300"></div>
          <div class="relative flex h-64 w-64 items-center justify-center rounded-2xl border-2 border-brand-border bg-white p-3 shadow-lg">
            <img
              v-if="qrDataUrl"
              :src="qrDataUrl"
              :alt="'QR Code ' + title"
              class="h-full w-full object-contain rounded-lg select-none"
            />
            <div v-else class="flex flex-col items-center gap-2 text-muted-foreground">
              <Sparkles class="h-6 w-6 animate-spin text-brand" />
              <span class="text-xs">Génération du QR Code...</span>
            </div>
          </div>
        </div>

        <!-- Instructions de scan -->
        <p class="mt-4 text-xs font-medium text-foreground">
          Visez ce QR code avec l'appareil photo de votre smartphone
        </p>
        <p class="mt-1 max-w-xs text-[11px] text-muted-foreground leading-relaxed">
          Ouvre instantanément la PWA sur iPhone (Safari) ou Android (Chrome) sans téléchargement sur un store.
        </p>

        <!-- Champ URL & Copie -->
        <div class="mt-4 flex w-full items-center gap-2">
          <input
            type="text"
            readonly
            :value="url"
            class="h-9 min-w-0 flex-1 select-all rounded-lg border border-border bg-muted/40 px-3 font-mono text-xs text-foreground focus:outline-none focus:border-brand"
          />
          <button
            type="button"
            @click="handleCopy"
            class="inline-flex h-9 items-center gap-1.5 rounded-lg border border-border px-3 text-xs font-medium text-foreground transition hover:bg-muted"
          >
            <CheckCircle2 v-if="isCopied" class="h-3.5 w-3.5 text-brand" />
            <Copy v-else class="h-3.5 w-3.5 text-muted-foreground" />
            <span>{{ isCopied ? 'Copié' : 'Copier' }}</span>
          </button>
        </div>

        <!-- Actions : Télécharger PNG, Imprimer, WhatsApp -->
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
