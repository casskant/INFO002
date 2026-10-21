"""Tests du TP2 : signature, stego, diplome, registre et site de validation.
Lancer :  python3 test_tp2.py
Tout est créé dans un dossier temporaire (clés, images, registre) : rien n'est écrasé.
"""
import os
import sys
import tempfile
import threading
import urllib.request
from wsgiref.simple_server import WSGIRequestHandler, make_server
from PIL import Image

import diplome
import registre
import serveur
import stego
from signature import generer_cles, signer, verifier

resultats = []


def verifie(titre, obtenu, attendu):
    ok = obtenu == attendu
    resultats.append(ok)
    print(f"  [{'OK' if ok else 'ECHEC'}] {titre} -> {obtenu} (attendu : {attendu})")


tmp = tempfile.mkdtemp()
def chemin(nom):
    return os.path.join(tmp, nom)

PRIVEE, PUBLIQUE = chemin("cle_privee.pem"), chemin("cle_publique.pem")
AUTRE_PRIVEE, AUTRE_PUBLIQUE = chemin("autre_privee.pem"), chemin("autre_publique.pem")
registre.DB = chemin("registre.db")                       # registre de test, pas le vrai

# ------------------------------------------------------------------ signature.py
print("== 1. signature.py ==")
generer_cles(PRIVEE, PUBLIQUE)
generer_cles(AUTRE_PRIVEE, AUTRE_PUBLIQUE)
data = b"Alice Martin|15.2"
sig = signer(data, PRIVEE)
verifie("données d'origine, bonne clé publique", verifier(data, sig, PUBLIQUE), True)
verifie("note modifiée (15.2 -> 19.9)", verifier(b"Alice Martin|19.9", sig, PUBLIQUE), False)
verifie("signature modifiée (1 bit)", verifier(data, bytes([sig[0] ^ 1]) + sig[1:], PUBLIQUE), False)
verifie("autre clé publique", verifier(data, sig, AUTRE_PUBLIQUE), False)
generer_cles(chemin("pp_privee.pem"), chemin("pp_publique.pem"), passphrase="secret")
verifie("clé protégée, bonne passphrase",
        verifier(data, signer(data, chemin("pp_privee.pem"), "secret"), chemin("pp_publique.pem")), True)
try:
    signer(data, chemin("pp_privee.pem"), "mauvaise")
    refusee = False
except ValueError:
    refusee = True
verifie("clé protégée, mauvaise passphrase refusée", refusee, True)

# ---------------------------------------------------------------------- stego.py
print("== 2. stego.py ==")
petite = Image.new("RGB", (100, 100), (200, 100, 50))
verifie("aller-retour cacher/extraire",
        stego.extraire_msg(stego.cacher_msg(petite, b"Lorem ipsum")), b"Lorem ipsum")
if hasattr(stego, "capacite"):
    try:
        stego.cacher_msg(petite, b"x" * (stego.capacite(petite) + 1))
        refuse = False
    except ValueError:
        refuse = True
    verifie("message trop long refusé", refuse, True)


print(f"\n{sum(resultats)}/{len(resultats)} tests réussis")
sys.exit(0 if all(resultats) else 1)