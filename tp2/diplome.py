"""Diplôme numérique : texte visible + (nom, moyenne, date, signature) cachés en LSB.

usage :
  python3 diplome.py keygen
  python3 diplome.py generer "Alice Martin" 15.2 [fond.png] -o diplome.png
  python3 diplome.py extraire diplome.png
  python3 diplome.py verifier diplome.png
"""
import argparse, base64, json, os, sys
from datetime import date
from PIL import Image, ImageDraw, ImageFont
import qrcode
import registre, signature, stego

PRIVEE, PUBLIQUE = "cle_privee.pem", "cle_publique.pem"
URL_SITE = os.environ.get("DIPLOME_URL", "http://localhost:8080")   # site de validation
POLICES = ["sans.ttf", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"]


def police(taille):
    for p in POLICES:
        if os.path.exists(p):
            return ImageFont.truetype(p, taille)
    return ImageFont.load_default(taille)


def fond_par_defaut():
    img = Image.new("RGB", (1000, 700), (245, 240, 225))
    ImageDraw.Draw(img).rectangle((20, 20, 980, 680), outline=(120, 90, 30), width=6)
    return img


def ajouter_texte(img, texte, y, taille):
    ImageDraw.Draw(img).text((img.width // 2, y), texte, (30, 30, 30),
                             police(taille), anchor="mm")   # centré


def diplome_simple(nom, moyenne, fond=None):
    """3.1 : texte variable sur fond constant, sans cryptographie."""
    img = Image.open(fond).convert("RGB") if fond else fond_par_defaut()
    h = img.height
    ajouter_texte(img, "Université NIHCAMCURT", int(h * .15), 48)
    ajouter_texte(img, "Diplôme de Master", int(h * .30), 36)
    ajouter_texte(img, nom, int(h * .50), 44)
    ajouter_texte(img, f"Moyenne : {moyenne}/20", int(h * .65), 32)
    return img


def donnees_a_signer(nom, moyenne, jour):
    # sérialisation canonique : mêmes octets à la signature et à la vérification
    return json.dumps({"nom": nom, "moyenne": moyenne, "date": jour},
                      sort_keys=True, ensure_ascii=False).encode()


def ajouter_qr_et_code(img, code):
    """Code lisible + QR-code vers le site de validation, en bas du diplôme."""
    qr = qrcode.make(f"{URL_SITE}/verifier/{code}", box_size=1, border=2).get_image().convert("RGB")
    echelle = max(1, int(img.height * .22) // qr.width)           # grossissement entier : modules nets
    qr = qr.resize((qr.width * echelle, qr.height * echelle), Image.NEAREST)
    marge = 45
    img.paste(qr, (img.width - qr.width - marge, img.height - qr.height - marge))
    d = ImageDraw.Draw(img)
    d.text((marge, img.height - marge - 30), "Validation en ligne : scannez le QR-code ou saisissez le code",
           (30, 30, 30), police(16), anchor="lb")
    d.text((marge, img.height - marge), f"Code : {code}", (30, 30, 30), police(22), anchor="lb")


def generer(nom, moyenne, sortie, fond=None, passphrase=None):
    moyenne = str(moyenne)
    img = diplome_simple(nom, moyenne, fond)          # 1) tout ce qui est visible d'abord...
    jour = date.today().isoformat()
    msg = donnees_a_signer(nom, moyenne, jour)
    sig = signature.signer(msg, PRIVEE, passphrase)
    code = registre.code_diplome(sig)
    ajouter_qr_et_code(img, code)
    sig_b64 = base64.b64encode(sig).decode()
    charge = json.dumps({"msg": msg.decode(), "sig": sig_b64}).encode()
    stego.cacher_msg(img, charge).save(sortie)        # 2) ...la stégano en dernier, en PNG
    registre.enregistrer(code, nom, moyenne, jour, sig_b64)
    return code


def extraire(source):
    """source : chemin ou fichier ouvert."""
    return json.loads(stego.extraire_msg(Image.open(source)))


def inspecter(source):
    """Renvoie (signature valide ?, données signées, code du diplôme)."""
    try:
        c = extraire(source)
        sig = base64.b64decode(c["sig"])
        ok = signature.verifier(c["msg"].encode(), sig, PUBLIQUE)
        return ok, json.loads(c["msg"]), registre.code_diplome(sig)
    except Exception as e:
        print("Erreur de vérification :", repr(e), file=sys.stderr)
        return False, None, None


def verifier(source):
    ok, contenu, _ = inspecter(source)
    return ok, contenu


def main():
    p = argparse.ArgumentParser()
    s = p.add_subparsers(dest="cmd", required=True)
    s.add_parser("keygen")
    g = s.add_parser("generer"); g.add_argument("nom"); g.add_argument("moyenne")
    g.add_argument("fond", nargs="?"); g.add_argument("-o", "--sortie", default="diplome.png")
    s.add_parser("extraire").add_argument("image")
    s.add_parser("verifier").add_argument("image")
    a = p.parse_args()
    pw = os.environ.get("DIPLOME_PASSPHRASE") or None   # passphrase de la clé privée
    if a.cmd == "keygen":
        signature.generer_cles(PRIVEE, PUBLIQUE, pw)
    elif a.cmd == "generer":
        print("Code du diplôme :", generer(a.nom, a.moyenne, a.sortie, a.fond, pw))
    elif a.cmd == "extraire":
        print(json.dumps(extraire(a.image), indent=2, ensure_ascii=False))
    else:
        ok, d = verifier(a.image)
        print("Signature VALIDE" if ok else "Signature INVALIDE ou diplôme falsifié")
        if ok:
            print("Contenu signé :", d, "\n=> comparez avec le texte visible sur l'image")
        sys.exit(0 if ok else 1)


if __name__ == "__main__":
    main()