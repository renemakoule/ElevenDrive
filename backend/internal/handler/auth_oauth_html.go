package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// renderOAuthResponse renvoie une page HTML moderne permettant d'ouvrir directement la fenêtre Desktop ou de continuer sur le Web.
func renderOAuthResponse(w http.ResponseWriter, redirectURL string, jsData map[string]interface{}, isDesktop bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	jsonData, _ := json.Marshal(jsData)
	status, _ := jsData["status"].(string)
	isError := status == "error"
	errorMessage := ""
	if isError {
		if e, ok := jsData["error"].(string); ok {
			errorMessage = e
		}
	}

	title := "Connexion Validée"
	desc := "Votre compte a été vérifié avec succès. Votre application ElevenDrive Desktop a été synchronisée automatiquement."
	if isError {
		title = "Erreur d'authentification"
		desc = errorMessage
	} else if status == "needs_registration" {
		title = "Identité Vérifiée"
		desc = "Votre compte a été certifié. Complétez la création de votre établissement sur votre application ElevenDrive."
	}

	bgColor, txtColor, borderColor := "rgba(16, 185, 129, 0.12)", "#10b981", "rgba(16, 185, 129, 0.25)"
	iconSVG := `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>`
	if isError {
		bgColor, txtColor, borderColor = "rgba(239, 68, 68, 0.12)", "#ef4444", "rgba(239, 68, 68, 0.25)"
		iconSVG = `<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>`
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ElevenDrive · Authentification</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            display: flex; align-items: center; justify-content: center;
            min-height: 100vh; background: #090d16; color: #f1f5f9; padding: 1.5rem;
        }
        .card {
            background: #111827; border: 1px solid #1f2937; border-radius: 8px;
            padding: 2.25rem 2rem; max-width: 440px; width: 100%%%%; text-align: center;
            box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.7);
        }
        .icon-circle {
            width: 52px; height: 52px; margin: 0 auto 1.25rem; border-radius: 50%%%%;
            display: flex; align-items: center; justify-content: center;
            background: %s; color: %s; border: 1px solid %s;
        }
        .icon-circle svg { width: 28px; height: 28px; stroke-width: 2.5; }
        h1 { font-size: 1.25rem; font-weight: 600; letter-spacing: -0.02em; margin-bottom: 0.5rem; color: #f8fafc; }
        p.desc { font-size: 0.875rem; color: #94a3b8; line-height: 1.5; margin-bottom: 1.5rem; }
        .actions { display: flex; flex-direction: column; gap: 0.625rem; }
        .btn {
            display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem;
            width: 100%%%%; padding: 0.65rem 1rem; font-size: 0.84rem; font-weight: 600;
            border-radius: 6px; cursor: pointer; transition: all 0.15s ease; text-decoration: none; border: none;
        }
        .btn-primary { background: #10b981; color: #ffffff; }
        .btn-primary:hover { background: #059669; }
        .btn-primary:active { transform: scale(0.99); }
        .btn-secondary { background: transparent; color: #cbd5e1; border: 1px solid #334155; }
        .btn-secondary:hover { background: #1e293b; color: #ffffff; }
        .hint { margin-top: 1.25rem; font-size: 0.75rem; color: #64748b; }
        .brand { margin-top: 1.5rem; font-size: 0.72rem; color: #475569; letter-spacing: 0.05em; text-transform: uppercase; }
    </style>
</head>
<body>
    <div class="card">
        <div class="icon-circle">%s</div>
        <h1>%s</h1>
        <p class="desc">%s</p>
        <div class="actions">
            <button type="button" class="btn btn-primary" onclick="openDesktopApp()">
                <svg width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
                    <line x1="8" y1="21" x2="16" y2="21"></line>
                    <line x1="12" y1="17" x2="12" y2="21"></line>
                </svg>
                Ouvrir la fenêtre ElevenDrive
            </button>
            <button type="button" class="btn btn-secondary" onclick="openWebApp()">
                Continuer dans le navigateur Web
            </button>
        </div>
        <p class="hint">Vous pouvez fermer cet onglet en toute sécurité.</p>
        <div class="brand">ElevenDrive · Système d'exploitation Auto-École</div>
    </div>
    <script>
        const data = %s;
        const redirectUrl = "%s";
        const isDesktop = %t;
        function openDesktopApp() {
            window.location.href = "elevendrive://auth";
            setTimeout(() => { window.close(); }, 400);
        }
        function openWebApp() {
            window.location.replace(redirectUrl);
        }
        if (window.opener && window.opener !== window) {
            try {
                window.opener.postMessage({ type: 'AUTOSCHOOL_OAUTH_RESULT', payload: data }, '*');
                setTimeout(() => window.close(), 500);
            } catch (e) {
                if (!isDesktop) window.location.replace(redirectUrl);
            }
        }
    </script>
</body>
</html>`,
		bgColor, txtColor, borderColor,
		iconSVG, title, desc,
		string(jsonData), redirectURL, isDesktop,
	)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

