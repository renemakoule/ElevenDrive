/**
 * Utilitaire d'ouverture sécurisée des liens externes (Web, PWA, Desktop Tauri).
 * Garantit que les liens s'ouvrent toujours dans le navigateur par défaut de l'utilisateur,
 * avec normalisation automatique du protocole https:// et contournement des bloqueurs de popup.
 */

/**
 * Normalise une URL en s'assurant qu'elle possède un protocole http:// ou https://
 */
export function formatExternalUrl(url?: string): string {
  if (!url) return '';
  const trimmed = url.trim();
  if (!trimmed) return '';

  if (/^https?:\/\//i.test(trimmed)) {
    return trimmed;
  }

  // Si l'utilisateur a saisi 'www.site.com' ou 'site.com'
  return `https://${trimmed}`;
}

/**
 * Extrait le nom de domaine lisible d'une URL (ex: "drive.google.com", "codedelaroute.fr")
 */
export function getDomainFromUrl(url?: string): string {
  if (!url) return 'Lien Web';
  try {
    const formatted = formatExternalUrl(url);
    const parsed = new URL(formatted);
    return parsed.hostname.replace(/^www\./i, '');
  } catch {
    return url.replace(/^https?:\/\//i, '').split('/')[0] || 'Lien Web';
  }
}

/**
 * Ouvre une URL externe de façon universelle :
 * 1. Tente via le plugin Shell Tauri sous environnement Desktop
 * 2. Repli robuste via un lien HTML <a> target="_blank" (évite le blocage des popups)
 * 3. Repli direct window.open
 */
export async function openExternalLink(rawUrl?: string): Promise<void> {
  const url = formatExternalUrl(rawUrl);
  if (!url || typeof window === 'undefined') return;

  // 1. Environnement Desktop Tauri
  if ('__TAURI_INTERNALS__' in window || '__TAURI__' in window) {
    try {
      const { open } = await import('@tauri-apps/plugin-shell');
      await open(url);
      return;
    } catch (tauriErr) {
      console.warn("Échec d'ouverture via Tauri shell plugin, repli web:", tauriErr);
    }
  }

  // 2. Déclenchement naturel via élément <a> pour contourner les bloqueurs de popups
  try {
    const link = document.createElement('a');
    link.href = url;
    link.target = '_blank';
    link.rel = 'noopener noreferrer';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    return;
  } catch (domErr) {
    console.warn("Échec d'ouverture DOM link, repli window.open:", domErr);
  }

  // 3. Repli standard
  window.open(url, '_blank', 'noopener,noreferrer');
}
