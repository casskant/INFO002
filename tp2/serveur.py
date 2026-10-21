"""Site de validation des diplômes NIHCAMCURT (prototype).   python3 serveur.py  ->  http://localhost:8080"""
import base64
import io
import re
from urllib.parse import quote
from bottle import Bottle, abort, redirect, request, run, template

import diplome, registre, signature

app = Bottle()
TAILLE_MAX = 5 * 1024 * 1024
CODE = re.compile(r"[0-9a-f]{16}")

BASE = """<!doctype html><html lang="fr"><meta charset="utf-8">
<title>NIHCAMCURT : validation des diplômes</title>
<style>body{font-family:sans-serif;max-width:640px;margin:2em auto;padding:0 1em}
.ok{color:#0a7a2f}.ko{color:#b00020}td{padding:2px 14px 2px 0}</style>
<h1><a href="/" style="color:inherit;text-decoration:none">Validation des diplômes NIHCAMCURT</a></h1>
{{!corps}}</html>"""

ACCUEIL = """<h2>Par code</h2>
<form action="/verifier" method="get">
<input name="code" size="20" placeholder="code du diplôme" required> <button>Valider</button></form>
<h2>Par image</h2>
<p>Déposez le diplôme (PNG) : le site lit le message caché et vérifie la signature.</p>
<form action="/verifier-image" method="post" enctype="multipart/form-data">
<input type="file" name="image" accept="image/png" required> <button>Vérifier</button></form>"""

RESULTAT = """<h2 class="{{'ok' if ok else 'ko'}}">{{verdict}}</h2>
% if donnees:
<table>
<tr><td>Nom</td><td>{{donnees['nom']}}</td></tr>
<tr><td>Moyenne</td><td>{{donnees['moyenne']}}/20</td></tr>
<tr><td>Date</td><td>{{donnees['date']}}</td></tr>
</table>
<p>Comparez ces informations avec le texte visible sur le diplôme.</p>
% end"""


def page(corps, **variables):
    # {{x}} échappe le HTML (nom = "<script>..." inoffensif) ; {{!corps}} insère la page déjà rendue
    return template(BASE, corps=template(corps, **variables))


@app.get("/")
def accueil():
    return page(ACCUEIL)


@app.get("/verifier")
def verifier_formulaire():
    redirect("/verifier/" + quote(request.query.get("code", "").strip().lower()))


@app.get("/verifier/<code>")
def verifier_code(code):
    enr = registre.chercher(code) if CODE.fullmatch(code) else None
    if enr is None:
        return page(RESULTAT, ok=False, donnees=None,
                    verdict="Code inconnu : ce diplôme n'a pas été délivré par NIHCAMCURT.")
    # on revérifie la signature du registre : une base de données modifiée serait détectée
    msg = diplome.donnees_a_signer(enr["nom"], enr["moyenne"], enr["date"])
    sig_ok = signature.verifier(msg, base64.b64decode(enr["signature"]), diplome.PUBLIQUE)
    if enr["revoque"]:
        verdict, ok = "Diplôme RÉVOQUÉ.", False
    elif not sig_ok:
        verdict, ok = "Registre incohérent : signature invalide.", False
    else:
        verdict, ok = "Diplôme authentique, délivré par NIHCAMCURT.", True
    return page(RESULTAT, ok=ok, verdict=verdict, donnees=enr)


@app.post("/verifier-image")
def verifier_image():
    envoi = request.files.get("image")
    if envoi is None:
        abort(400, "Aucun fichier reçu.")
    data = envoi.file.read(TAILLE_MAX + 1)
    if len(data) > TAILLE_MAX:
        abort(413, "Fichier trop gros (5 Mo maximum).")
    valide, contenu, code = diplome.inspecter(io.BytesIO(data))
    if not valide:
        return page(RESULTAT, ok=False, donnees=None,
                    verdict="Signature invalide, ou aucun message caché dans cette image.")
    enr = registre.chercher(code)
    if enr is None:
        verdict, ok = "Signature valide, mais ce diplôme n'est pas dans le registre.", False
    elif enr["revoque"]:
        verdict, ok = "Diplôme RÉVOQUÉ.", False
    elif (enr["nom"], enr["moyenne"], enr["date"]) != (contenu["nom"], contenu["moyenne"], contenu["date"]):
        verdict, ok = "Les données cachées diffèrent du registre.", False
    else:
        verdict, ok = "Diplôme authentique, délivré par NIHCAMCURT.", True
    return page(RESULTAT, ok=ok, verdict=verdict, donnees=contenu)


if __name__ == "__main__":
    run(app, host="127.0.0.1", port=8080)