import QRCode from 'qrcode';

export interface QrOptions {
  width?: number;
  margin?: number;
  darkColor?: string;
  lightColor?: string;
}

/**
 * Génère une DataURL (image/png en base64) représentant un QR Code scannable.
 */
export const generateQrDataUrl = async (
  text: string,
  options?: QrOptions
): Promise<string> => {
  if (!text) return '';
  try {
    return await QRCode.toDataURL(text, {
      width: options?.width ?? 320,
      margin: options?.margin ?? 2,
      color: {
        dark: options?.darkColor ?? '#064e3b', // Teinte Brand Emerald foncée
        light: options?.lightColor ?? '#ffffff',
      },
      errorCorrectionLevel: 'M',
    });
  } catch (err) {
    console.error('Erreur génération QR Code:', err);
    return '';
  }
};

/**
 * Déclenche le téléchargement du QR Code sous forme d'image PNG
 */
export const downloadQrCode = (dataUrl: string, fileName: string) => {
  if (!dataUrl) return;
  const link = document.createElement('a');
  link.href = dataUrl;
  link.download = fileName.endsWith('.png') ? fileName : `${fileName}.png`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};
